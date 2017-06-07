package d71

import "strings"

type FileType int

const (
	Del FileType = iota // Deleted
	Seq                 // Sequential
	Prg                 // Program
	Usr                 // User
	Rel                 // Relative

	MaxFilenameLen = 16
)

const (
	bitSaveAt = (1 << 5)
	bitLocked = (1 << 6)
	bitSplat  = (1 << 7)
)

var fileTypeStr = map[FileType]string{
	Del: "DEL",
	Seq: "SEQ",
	Prg: "PRG",
	Usr: "USR",
	Rel: "REL",
}

func (f FileType) String() string {
	if str, ok := fileTypeStr[f]; ok {
		return str
	}
	return "???"
}

type FileInfo struct {
	Type   FileType //
	SaveAt bool     // SAVE-@ operation
	Locked bool     //
	Splat  bool     // True if the file wasn't properly closed
	Name   string   //
	Size   int      // Number of sectors
	First  Pos      // Location of first block
	pos    Pos      // Position of this file entry in the directory
}

type dirWalker struct {
	skipDeleted bool // If false, returns deleted entries
	e           *Editor
	nextTrack   int  // Track of the next directory block or $00
	nextSector  int  // Sector of the next directory block
	entry       int  // Entry (0-7) in this sector that the walker is at
	eof         bool //
}

func newDirWalker(d Disk) *dirWalker {
	w := &dirWalker{e: d.Editor()}
	w.skipDeleted = true
	w.e.Seek(DirTrack, 0, 0)

	// BAM sector contains the location of the first directory block
	firstTrack := w.e.Read()
	firstSector := w.e.Read()
	w.e.Seek(firstTrack, firstSector, 0)
	w.readLink()
	return w
}

// Remember the link for the next directory block
func (w *dirWalker) readLink() {
	e := w.e.Mark()
	w.nextTrack = e.Read()
	w.nextSector = e.Read()
}

// Advance the walker to the next directory entry. If at the last entry
// for this sector, move to the next sector or return false indicating
// the end of the listing.
func (w *dirWalker) advance() bool {
	w.entry++
	if w.entry == 8 {
		if w.nextTrack == 0 {
			return false
		}
		w.entry = 0
		w.e.Seek(w.nextTrack, w.nextSector, 0)
		w.readLink()
	} else {
		w.e.Move(0x20)
	}
	return true
}

func (w *dirWalker) next() (*FileInfo, bool) {
	if w.eof {
		return nil, false
	}
	for {
		e := w.e.Mark().Move(2)
		ftype := e.Read()
		// Skip over this entry if it has been deleted
		if w.skipDeleted && ftype == 0 {
			ok := w.advance()
			if !ok {
				// Reached the end, no more entries
				w.eof = true
				return nil, false
			}
		} else {
			// If not deleted, this is a valid entry.
			break
		}
	}

	fi := &FileInfo{}
	fi.pos = w.e.Pos

	e := w.e.Mark().Move(2)
	ftype := e.Read()
	fi.Type = FileType(ftype & 0x7)
	fi.SaveAt = ftype&bitSaveAt > 0
	fi.Locked = ftype&bitLocked > 0
	fi.Splat = ftype&bitSplat == 0
	fi.First.Track = e.Read()
	fi.First.Sector = e.Read()
	fi.Name = strings.Trim(e.ReadString(16), "\xa0")
	e.Move(0x1e - 0x15)
	fi.Size = e.ReadWord()

	ok := w.advance()
	if !ok {
		w.eof = true
	}
	return fi, true
}

func writeFileInfo(d Disk, fi *FileInfo) {
	e := d.Editor()
	e.Pos = fi.pos
	e.Move(2)

	ftype := int(fi.Type)
	if fi.SaveAt {
		ftype = ftype | bitSaveAt
	}
	if fi.Locked {
		ftype = ftype | bitLocked
	}
	if fi.Splat {
		ftype = ftype | bitSplat
	}
	e.Write(ftype)
	e.Write(fi.First.Track)
	e.Write(fi.First.Sector)
	e.WriteStringN(fi.Name, 0xa0, MaxFilenameLen)
	e.Move(2) // Track/Sector location of first side-sector block (REL file only)
	e.Move(1) // REL file record length (REL file only, max. value 254)
	e.Move(6) // $18-$1D: Unused (except with GEOS disks)
	e.Write(fi.Size)
}

func createDirEntry(d Disk) (*FileInfo, error) {
	w := newDirWalker(d)

	// See if we can reuse a delete entry. Also unused entries on a
	// directory sector appear as deleted since the file type is zero.
	w.skipDeleted = false
	for {
		fi, ok := w.next()
		if !ok {
			break
		}
		if fi.Type == Del {
			return fi, nil
		}
	}
	// In this case, the last directory sector was full. Create a new
	// one.
	dirSector, ok := freeDirSector(d)
	if !ok {
		return nil, ErrDirFull
	}
	fi := &FileInfo{
		pos: Pos{
			Track:  DirTrack,
			Sector: dirSector,
		},
	}
	return fi, nil
}
