package d71

import "testing"

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

	expPos := Pos(BamTrack, 0, (46-Flip)*4)
	expOff := 3
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
	e.Seek(BamTrack, 0, 0)
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
	e.Seek(DirTrack, 0, 0xdd)
	expFree2 := 21 - 1
	free2 := e.Read()
	if expFree2 != free2 {
		t.Errorf("expected free2 %v ; actual %v", expFree2, free2)
	}
}
