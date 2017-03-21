package d71

import (
	"strings"
)

type FileType int

const (
	Del FileType = iota // Deleted
	Seq                 // Sequential
	Prg                 // Program
	Usr                 // User
	Rel                 // Relative
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
	Track  int      // Location of first
	Sector int      // Location of first
}

type dirWalker struct {
	e          *Editor
	nextTrack  int // Track of the next directory block or $00
	nextSector int // Sector of the next directory block
	entry      int // Entry (0-7) in this sector that the walker is at
}

func newDirWalker(d Disk) *dirWalker {
	w := &dirWalker{e: d.Editor()}
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

func (w *dirWalker) next() (FileInfo, bool) {
	for {
		e := w.e.Mark().Move(2)
		ftype := e.Read()
		// Skip over this entry if it has been deleted
		if ftype == 0 {
			ok := w.advance()
			if !ok {
				// Reached the end, no more entries
				return FileInfo{}, false
			}
		} else {
			// If not deleted, this is a valid entry.
			break
		}
	}

	fi := FileInfo{}

	w.e.Move(2)
	ftype := w.e.Read()
	fi.Type = FileType(ftype & 0x7)
	fi.SaveAt = ftype&(1<<5) > 0
	fi.Locked = ftype&(1<<6) > 0
	fi.Splat = ftype&(1<<7) == 0
	fi.Track = w.e.Read()
	fi.Sector = w.e.Read()
	fi.Name = strings.Trim(w.e.ReadString(16), "\xa0")
	w.e.Move(0x1e - 0x15)
	fi.Size = w.e.ReadWord()

	return fi, true
}
