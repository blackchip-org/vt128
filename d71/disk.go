package d71

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// DiskLen is the maximum number of bytes that can be stored on the disk
	DiskLen = 349696

	// SectorLen is the number of bytes in each sector
	SectorLen = 256

	// Flip is the first track found on the flip side of the disk
	Flip = 36

	// MaxTrack is the last track on the back side of the disk
	MaxTrack = 70

	// DirTrack is where the directory is located
	DirTrack = 18

	// BamTrack is where the extended BAM information is stored on the
	// back side of the disk
	BamTrack = 53

	// MaxTrackLen is the maximum number of sectors that can be found in a
	// track
	MaxTrackLen = 21
)

// Track contains a number of sectors and the absolute offset in the disk
// of where the tracks starts
type Track struct {
	Sectors  int
	Offset   int
	lastFree int // Bitmap for the final BAM byte when track free
}

type DiskInfo struct {
	Name        string
	ID          string
	DosVersion  string
	DosType     string
	DoubleSided bool
	Free        int
}

// Geom contains an entry for each track describing the number of sectors
// and absolute offset into the disk. Since there is no track zero, that
// index does not contain any useful information.
var Geom []Track

// Create the geometry table
func init() {
	Geom = make([]Track, 71, 71)

	offset := 0
	for i := 1; i <= 70; i++ {
		sectors := 0
		lastFree := 0
		switch {
		case i >= 1 && i <= 17:
			sectors = 21
			lastFree = 1<<5 - 1
		case i >= 18 && i <= 24:
			sectors = 19
			lastFree = 1<<3 - 1
		case i >= 25 && i <= 30:
			sectors = 18
			lastFree = 1<<2 - 1
		case i >= 31 && i <= 35:
			sectors = 17
			lastFree = 1
		case i >= 36 && i <= 52:
			sectors = 21
			lastFree = 1<<5 - 1
		case i >= 53 && i <= 59:
			sectors = 19
			lastFree = 1<<3 - 1
		case i >= 60 && i <= 65:
			sectors = 18
			lastFree = 1<<2 - 1
		case i >= 66 && i <= 70:
			sectors = 17
			lastFree = 1
		default:
			panic(fmt.Sprintf("invalid track: %v", i))
		}
		Geom[i] = Track{Sectors: sectors, Offset: offset, lastFree: lastFree}
		offset += (sectors * SectorLen)
	}
}

// Pos computes the disk byte offset based on a track, sector,
// and sector offset
func Pos(track int, sector int, offset int) int {
	toff := Geom[track].Offset
	return toff + (sector * SectorLen) + offset
}

// A 1571 floppy disk. Use NewDisk for a formatted disk.
type Disk []byte

func NewDisk(name string, id string) Disk {
	if len(name) > 0xf {
		name = name[:0xf]
	}
	if len(id) > 2 {
		id = id[:2]
	}

	d := make(Disk, DiskLen, DiskLen)
	e := d.Editor()

	e.Seek(18, 0, 0)
	e.Write(18)   // Track of first directory sector
	e.Write(1)    // Sector of first directory sector
	e.Write(0x41) // Disk DOS version type. A = 1541
	e.Write(0x80) // Double-sided flag, set to double-sided

	// BAM, front Side
	for i := 1; i < Flip; i++ {
		sectors := Geom[i].Sectors
		if i != DirTrack {
			e.Write(sectors) // Sectors available
			e.Write(0xff)    // Sectors 0 - 7 free
		} else {
			e.Write(sectors - 2) // BAM and first dir sector in use
			e.Write(0xfc)        // First two sectors in use
		}
		e.Write(0xff)             // Sectors 8 - 15 free
		e.Write(Geom[i].lastFree) // Remaining sectors free
	}

	e.WriteStringN(name, 0xa0, 0x10) // Disk Name
	e.Fill(0xa0, 2)                  // Fill
	e.WriteStringN(id, 0x20, 2)      // Disk ID
	e.Write(0xa0)                    // Fill
	e.WriteString("2A")              // DOS Type
	e.Fill(0xa0, 0xaa-0xa7+1)        // Fill

	// Free sector count of back side
	e.Seek(18, 0, 0xdd)
	for i := Flip; i <= MaxTrack; i++ {
		if i != BamTrack {
			sectors := Geom[i].Sectors
			e.Write(sectors) // Sectors available
		} else {
			e.Write(0) // All sectors in use
		}
	}

	// BAM, back side
	e.Seek(53, 0, 0)
	for i := Flip; i <= MaxTrack; i++ {
		if i != BamTrack {
			e.Write(0xff)             // Sectors 0 - 7 free
			e.Write(0xff)             // Sectors 8 - 15 free
			e.Write(Geom[i].lastFree) // Remaining sectors free
		} else {
			e.Fill(0, 3) // All sectors marked as used
		}
	}

	// Blank directory, set link to nothing
	e.Seek(DirTrack, 1, 1)
	e.Write(0xff)

	return d
}

func (d Disk) Editor() *Editor {
	return &Editor{disk: d}
}

func (d Disk) Save(filename string) error {
	err := ioutil.WriteFile(filename, d, 0644)
	return err
}

func Load(filename string) (Disk, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	if fi.Size() != DiskLen {
		return nil, fmt.Errorf("File is not a D71 disk: %v", filename)
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Disk(data), nil
}

func (d Disk) Info() DiskInfo {
	e := d.Editor()
	di := DiskInfo{}
	free := 0
	e.Seek(DirTrack, 0, 2)
	di.DosVersion = e.ReadString(1)
	di.DoubleSided = e.Read() == 0x80
	// Front side counts in BAM
	for track := 1; track < Flip; track++ {
		// Don't count directory sectors
		if track == DirTrack {
			e.Read()
		} else {
			free += e.Read()
		}
		e.Move(3)
	}
	di.Name = strings.Trim(e.ReadString(16), "\xa0")
	e.Move(2)
	di.ID = e.ReadString(2)
	e.Move(1)
	di.DosType = e.ReadString(2)

	// Back side counts in aux area
	e.Seek(DirTrack, 0, 0xdd)
	for track := Flip; track <= MaxTrack; track++ {
		// Don't count back side BAM track
		if track == BamTrack {
			e.Read()
		} else {
			free += e.Read()
		}
	}
	di.Free = free
	return di
}

func (d Disk) List() []FileInfo {
	w := newDirWalker(d)
	list := make([]FileInfo, 0, 0)
	for {
		fi, more := w.next()
		if !more {
			break
		}
		list = append(list, fi)
	}
	return list
}

// For a given track and sector, compute the location of the BAM entry.
// This function will move the editor position to the start of the BAM
// record. It returns the offset from that position to the byte that
// holds the bitmap and the mask that should be used to modify the entry.
func bamPos(e *Editor, track int, sector int) (off int, mask int) {
	bmapOffset := 1
	bytesPerRecord := 4
	if track < Flip {
		e.Seek(DirTrack, 0, 4)
	} else {
		e.Seek(BamTrack, 0, 0)
		track = track - Flip + 1
		bmapOffset = 0
		bytesPerRecord = 3
	}
	e.Move((track - 1) * bytesPerRecord)
	off = sector/8 + bmapOffset
	mask = 1 << byte(sector%8)
	return off, mask
}

// BamRead returns true if the given track and sector is marked as free
// in the block availability map. Otherwise returns false.
func (d Disk) BamRead(track int, sector int) bool {
	e := d.Editor()
	off, mask := bamPos(e, track, sector)
	bmap := e.Move(off).Peek()
	return bmap&mask > 0
}

// BamWrite updates the block availability map for the given track and
// sector. True markes it as free, false as allocated.
func (d Disk) BamWrite(track int, sector int, val bool) {
	// Do nothing if the value is the same
	prev := d.BamRead(track, sector)
	if prev == val {
		return
	}

	// Update the available sector count by +1 or -1
	delta := -1
	if val {
		delta = 1
	}

	e := d.Editor()
	off, mask := bamPos(e, track, sector)

	// Update the number of available sectors for this track if on the
	// front side
	if track < Flip {
		e.Poke(e.Peek() + delta)
	}

	// Update the bitmap entry
	bmap := e.Move(off).Peek()
	if val {
		bmap = bmap | mask
	} else {
		bmap = bmap & ^mask
	}
	e.Poke(bmap)

	// If the track was on the back side of the disk, we need to update
	// the supplemental sector free count on the from side BAM sector
	if track >= Flip {
		e.Seek(DirTrack, 0, 0xdd)
		off = track - Flip
		e.Move(off).Poke(e.Peek() + delta)
	}
}
