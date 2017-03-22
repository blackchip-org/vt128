package d71

import (
	"testing"

	"github.com/blackchip-org/vt128/binary"
)

func TestNewDisk(t *testing.T) {
	expected := "MY DISK\xa0"
	d, _ := NewDisk("MY DISK", "MD")
	actual := d.Editor().Move(0x16590).ReadString(len(expected))
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestBamPosFront(t *testing.T) {
	d, _ := NewDisk("", "")
	e := d.Editor()
	off, mask := bamPos(e, 10, 15)

	expPos := Pos(DirTrack, 0, 0x04+(10-1)*4)
	expOff := 2
	expMask := 1 << 7

	if e.Pos != expPos {
		t.Errorf("expected pos $%x ; actual $%x", expPos, e.Pos)
	}
	if off != expOff {
		t.Errorf("expected off %v ; actual %v", expOff, off)
	}
	if mask != expMask {
		t.Errorf("expected mask %b ; actual %b", expMask, mask)
	}
}

func TestBamPosBack(t *testing.T) {
	d, _ := NewDisk("", "")
	e := d.Editor()
	off, mask := bamPos(e, 46, 18)

	expPos := Pos(BamTrack, 0, (46-Flip)*3)
	expOff := 2
	expMask := 1 << 2

	if e.Pos != expPos {
		t.Errorf("expected pos $%x ; actual $%x", expPos, e.Pos)
	}
	if off != expOff {
		t.Errorf("expected off %v ; actual %v", expOff, off)
	}
	if mask != expMask {
		t.Errorf("expected mask %b ; actual %b", expMask, mask)
	}
}

func TestBamAlloc(t *testing.T) {
	d, _ := NewDisk("", "")
	e := d.Editor()
	d.bamWrite(1, 1, false)
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
	d, _ := NewDisk("", "")
	e := d.Editor()
	d.bamWrite(1, 1, false)
	d.bamWrite(1, 1, true)
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

func TestBamDoubleAlloc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on double alloc")
		}
	}()
	d, _ := NewDisk("", "")
	d.bamWrite(1, 1, false)
	d.bamWrite(1, 1, false)
}

func TestBamDoubleFree(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on double free")
		}
	}()
	d, _ := NewDisk("", "")
	d.bamWrite(1, 1, true)
}

func TestBamAllocBack(t *testing.T) {
	d, _ := NewDisk("", "")
	e := d.Editor()
	d.bamWrite(Flip, 1, false)
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
	expected, _ := NewDisk("BLANK", "")
	actual, _ := NewDisk("", "")
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
	d, _ := NewDisk("", "")
	info := d.Info()
	expected := 1328
	actual := info.Free
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}
