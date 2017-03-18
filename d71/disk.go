package d71

import (
	"fmt"
	"io/ioutil"
)

const (
	// DiskLen is the maximum number of bytes that can be stored on the disk
	DiskLen = 349696

	// SectorLen is the number of bytes in each sector
	SectorLen = 256

	// Flip is the first track found on the flip side of the disk
	Flip = 36

	// LastTrack is the last track on the back side of the disk
	LastTrack = 70

	// DirTrack is where the directory is located
	DirTrack = 18

	// BamTrack is where the extended BAM information is stored on the
	// back side of the disk
	BamTrack = 53
)

// Track contains a number of sectors and the absolute offset in the disk
// of where the tracks starts
type Track struct {
	Sectors int
	Offset  int
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
		switch {
		case i >= 1 && i <= 17:
			sectors = 21
		case i >= 18 && i <= 24:
			sectors = 19
		case i >= 25 && i <= 30:
			sectors = 18
		case i >= 31 && i <= 35:
			sectors = 17
		case i >= 36 && i <= 52:
			sectors = 21
		case i >= 53 && i <= 59:
			sectors = 19
		case i >= 60 && i <= 65:
			sectors = 18
		case i >= 66 && i <= 70:
			sectors = 17
		default:
			panic(fmt.Sprintf("invalid track: %v", i))
		}
		Geom[i] = Track{Sectors: sectors, Offset: offset}
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

func NewDisk(name string, id string) (Disk, error) {
	if len(name) > 0xf {
		return nil, fmt.Errorf("disk name too long")
	}
	if len(id) > 2 {
		return nil, fmt.Errorf("id too long")
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
		e.Write(sectors) // Sectors available
		e.Write(0xff)    // Sectors 0 - 7 free
		e.Write(0xff)    // Sectors 8 - 15 free
		e.Write(0xff)    // Remaining sectors free
	}

	e.WritePadded(name, 0xa0, 0x10) // Disk Name
	e.Fill(0xa0, 2)                 // Fill
	e.WritePadded(id, 0xa0, 2)      // Disk ID
	e.Write(0xa0)                   // Fill
	e.WriteString("2A")             // DOS Type
	e.Fill(0xa0, 0xaa-0xa7+1)       // Fill

	// Free sector count of back side
	e.Seek(18, 0, 0xdd)
	for i := Flip; i <= LastTrack; i++ {
		sectors := Geom[i].Sectors
		e.Write(sectors) // Sectors available
	}

	// BAM, back side
	e.Seek(53, 0, 0)
	for i := Flip; i <= LastTrack; i++ {
		sectors := Geom[i].Sectors
		e.Write(sectors) // Sectors available
		e.Write(0xff)    // Sectors 0 - 7 free
		e.Write(0xff)    // Sectors 8 - 15 free
		e.Write(0xff)    // Remaining sectors free
	}

	return d, nil
}

func (d Disk) Editor() *Editor {
	return &Editor{disk: d}
}

func (d Disk) Save(filename string) error {
	err := ioutil.WriteFile(filename, d, 0644)
	return err
}

// For a given track and sector, compute the location of the BAM entry.
// This function will move the editor position to the start of the BAM
// record. It returns the offset from that position to the byte that
// holds the bitmap and the mask that should be used to modify the entry.
func bamPos(e *Editor, track int, sector int) (off int, mask int) {
	if track < Flip {
		e.Seek(DirTrack, 0, 4)
	} else {
		e.Seek(BamTrack, 0, 0)
		track = track - Flip + 1
	}
	e.Move((track - 1) * 4)
	off = sector/8 + 1
	mask = 1 << byte(sector%8)
	return off, mask
}

// Returns true of the track/sector is free, false if it is used
func (d *Disk) bamRead(track int, sector int) bool {
	e := d.Editor()
	off, mask := bamPos(e, track, sector)
	bmap := e.Move(off).Peek()
	return bmap&mask > 0
}

// Updates the BAM entry for a track/sector, set to true for free and
// false for used
func (d *Disk) bamWrite(track int, sector int, val bool) {
	// Ensure this is a valid alloc or free
	prev := d.bamRead(track, sector)
	if prev == val && val {
		panic(fmt.Sprintf("double free, track %v, sector %v", track, sector))
	}
	if prev == val && !val {
		panic(fmt.Sprintf("double alloc, track %v, sector %v", track, sector))
	}

	// Update the available sector count by +1 or -1
	delta := -1
	if val {
		delta = 1
	}

	e := d.Editor()
	off, mask := bamPos(e, track, sector)

	// Update the number of available sectors for this track
	e.Poke(e.Peek() + delta)

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
