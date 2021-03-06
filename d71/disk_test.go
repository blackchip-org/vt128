package d71

import (
	"testing"

	"github.com/blackchip-org/vt128/binary"
)

func TestNewDisk(t *testing.T) {
	expected := "MY DISK\xa0"
	d := NewDisk("MY DISK", "MD")
	actual := d.Editor().Move(0x16590).ReadString(len(expected))
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestNewDiskLongName(t *testing.T) {
	d := NewDisk("123456789ABCDEFXX", "AB")
	expected := "123456789ABCDEF"
	actual := d.Editor().Move(0x16590).ReadString(len(expected))
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestNewDiskLongID(t *testing.T) {
	d := NewDisk("MY DISK", "ABCD")
	expected := "AB"
	actual := d.Editor().Move(0x165a2).ReadString(len(expected))
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestBamPosFront(t *testing.T) {
	d := NewDisk("", "")
	e := d.Editor()
	at, mask := bamPos(e, 10, 15)

	expOffset := Offset(DirTrack, 0, 0x04+(10-1)*4)
	expAt := 2
	expMask := 1 << 7

	if e.Pos.Offset() != expOffset {
		t.Errorf("expected offset $%x ; actual $%x", e.Pos.Offset(), expOffset)
	}
	if expAt != at {
		t.Errorf("expected at %v ; actual %v", expAt, at)
	}
	if mask != expMask {
		t.Errorf("expected mask %b ; actual %b", expMask, mask)
	}
}

func TestBamPosBack(t *testing.T) {
	d := NewDisk("", "")
	e := d.Editor()
	at, mask := bamPos(e, 46, 18)

	expOffset := Offset(BamTrack, 0, (46-Flip)*3)
	expAt := 2
	expMask := 1 << 2

	if e.Pos.Offset() != expOffset {
		t.Errorf("expected offset $%x ; actual $%x", e.Pos.Offset(), expOffset)
	}
	if expAt != at {
		t.Errorf("expected at %v ; actual %v", expAt, at)
	}
	if mask != expMask {
		t.Errorf("expected mask %b ; actual %b", expMask, mask)
	}
}

func TestBamAlloc(t *testing.T) {
	d := NewDisk("", "")
	e := d.Editor()
	d.BamWrite(1, 1, false)
	e.Seek(DirTrack, 0, 0x4)
	expFree := 21 - 1
	free := e.Read()
	if expFree != free {
		t.Errorf("expected free %v ; actual %v", expFree, free)
	}
	expMap := 0xff - (1 << 1)
	bmap := e.Read()
	if expMap != bmap {
		t.Errorf("expected map %b ; actual %b", expMap, bmap)
	}
}

func TestBamFree(t *testing.T) {
	d := NewDisk("", "")
	e := d.Editor()
	d.BamWrite(1, 1, false)
	d.BamWrite(1, 1, true)
	e.Seek(DirTrack, 0, 0x4)
	expFree := 21
	free := e.Read()
	if expFree != free {
		t.Errorf("expected free %v ; actual %v", expFree, free)
	}
	expMap := 0xff
	bmap := e.Read()
	if expMap != bmap {
		t.Errorf("expected map %b ; actual %b", expMap, bmap)
	}
}

func TestBamAllocBack(t *testing.T) {
	d := NewDisk("", "")
	e := d.Editor()
	d.BamWrite(Flip, 1, false)
	e.Seek(DirTrack, 0, 0xdd)
	expFree := 21 - 1
	free := e.Read()
	if expFree != free {
		t.Errorf("expected free %v ; actual %v", expFree, free)
	}
	e.Seek(BamTrack, 0, 0)
	expMap := 0xff - (1 << 1)
	bmap := e.Read()
	if expMap != bmap {
		t.Errorf("expected map %b ; actual %b", expMap, bmap)
	}
	e.Seek(DirTrack, 0, 0xdd)
	expFree2 := 21 - 1
	free2 := e.Read()
	if expFree2 != free2 {
		t.Errorf("expected free2 %v ; actual %v", expFree2, free2)
	}
}

func TestBlankDisk(t *testing.T) {
	expected := NewDisk("BLANK", "")
	actual := NewDisk("", "")
	dump := `
0000000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
*
0016500 12 01 41 80 15 ff ff 1f 15 ff ff 1f 15 ff ff 1f
0016510 15 ff ff 1f 15 ff ff 1f 15 ff ff 1f 15 ff ff 1f
*
0016540 15 ff ff 1f 15 ff ff 1f 11 fc ff 07 13 ff ff 07
0016550 13 ff ff 07 13 ff ff 07 13 ff ff 07 13 ff ff 07
0016560 13 ff ff 07 12 ff ff 03 12 ff ff 03 12 ff ff 03
0016570 12 ff ff 03 12 ff ff 03 12 ff ff 03 11 ff ff 01
0016580 11 ff ff 01 11 ff ff 01 11 ff ff 01 11 ff ff 01
0016590 42 4c 41 4e 4b a0 a0 a0 a0 a0 a0 a0 a0 a0 a0 a0
00165a0 a0 a0 20 20 a0 32 41 a0 a0 a0 a0 00 00 00 00 00
00165b0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
*
00165d0 00 00 00 00 00 00 00 00 00 00 00 00 00 15 15 15
00165e0 15 15 15 15 15 15 15 15 15 15 15 15 15 15 00 13
00165f0 13 13 13 13 13 12 12 12 12 12 12 11 11 11 11 11
0016600 00 ff 00 00 00 00 00 00 00 00 00 00 00 00 00 00
0016610 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
*
0041000 ff ff 1f ff ff 1f ff ff 1f ff ff 1f ff ff 1f ff
0041010 ff 1f ff ff 1f ff ff 1f ff ff 1f ff ff 1f ff ff
0041020 1f ff ff 1f ff ff 1f ff ff 1f ff ff 1f ff ff 1f
0041030 ff ff 1f 00 00 00 ff ff 07 ff ff 07 ff ff 07 ff
0041040 ff 07 ff ff 07 ff ff 07 ff ff 03 ff ff 03 ff ff
0041050 03 ff ff 03 ff ff 03 ff ff 03 ff ff 01 ff ff 01
0041060 ff ff 01 ff ff 01 ff ff 01 00 00 00 00 00 00 00
0041070 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
*
0055600
`
	binary.LoadStringDumpInto(dump, actual)
	if len(expected) != len(actual) {
		t.Errorf("size mismatch: %x %x", len(expected), len(actual))
		return
	}
	diff, same := binary.Compare(expected, actual)
	if !same {
		t.Errorf("diff report:\n%v", diff)
	}
}

func TestBlankFreeSectors(t *testing.T) {
	d := NewDisk("", "")
	info := d.Info()
	expected := 1328
	actual := info.Free
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestTrackInfoFront(t *testing.T) {
	d := NewDisk("", "")
	d.BamWrite(4, 9, false)
	d.BamWrite(4, 19, false)
	ti := d.TrackInfo(4)
	expected := 21 - 2
	actual := ti.Free
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestTrackInfoBack(t *testing.T) {
	d := NewDisk("", "")
	d.BamWrite(44, 9, false)
	d.BamWrite(44, 19, false)
	ti := d.TrackInfo(44)
	expected := 21 - 2
	actual := ti.Free
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestList(t *testing.T) {
	d := NewDisk("", "")
	dump := `
00016600  00 ff 82 11 00 46 49 4c  45 20 31 a0 a0 a0 a0 a0  |.....FILE 1.....|
00016610  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016620  00 00 82 11 01 46 49 4c  45 20 32 a0 a0 a0 a0 a0  |.....FILE 2.....|
00016630  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016640  00 00 82 11 02 46 49 4c  45 20 33 a0 a0 a0 a0 a0  |.....FILE 3.....|
00016650  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016660  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
`
	if err := binary.LoadStringDumpInto(dump, d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	infos := d.List()
	want := "FILE 1"
	got := infos[0].Name
	if want != got {
		t.Errorf("wanted %v ; got %v", want, got)
	}
	want = "FILE 3"
	got = infos[2].Name
	if want != got {
		t.Errorf("wanted %v ; got %v", want, got)
	}
}

func TestFind(t *testing.T) {
	d := NewDisk("", "")
	dump := `
00016600  00 ff 82 11 00 46 49 4c  45 20 31 a0 a0 a0 a0 a0  |.....FILE 1.....|
00016610  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016620  00 00 82 11 01 46 49 4c  45 20 32 a0 a0 a0 a0 a0  |.....FILE 2.....|
00016630  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016640  00 00 82 11 02 46 49 4c  45 20 33 a0 a0 a0 a0 a0  |.....FILE 3.....|
00016650  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016660  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
`
	if err := binary.LoadStringDumpInto(dump, d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	info, found := d.Find("FILE 2")
	if !found {
		t.Fatalf("wanted found ; got not found")
	}
	want := "FILE 2"
	got := info.Name
	if want != got {
		t.Errorf("wanted %v ; got %v", want, got)
	}
}

func TestFindNotFound(t *testing.T) {
	d := NewDisk("", "")
	dump := `
00016600  00 ff 82 11 00 46 49 4c  45 20 31 a0 a0 a0 a0 a0  |.....FILE 1.....|
00016610  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016620  00 00 82 11 01 46 49 4c  45 20 32 a0 a0 a0 a0 a0  |.....FILE 2.....|
00016630  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016640  00 00 82 11 02 46 49 4c  45 20 33 a0 a0 a0 a0 a0  |.....FILE 3.....|
00016650  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016660  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
`
	if err := binary.LoadStringDumpInto(dump, d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	_, found := d.Find("FILE X")
	if found {
		t.Fatalf("wanted not found ; got found")
	}
}
