package d71

import (
	"testing"

	"github.com/blackchip-org/vt128/binary"
)

func TestDirWalkOneFile(t *testing.T) {
	d := NewDisk("", "")
	dump := `
00016600  00 ff 82 11 00 46 49 4c  45 20 31 a0 a0 a0 a0 a0  |.....FILE 1.....|
00016610  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016620  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
`
	if err := binary.LoadStringDumpInto(dump, d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	w := newDirWalker(d)
	fi, end := w.next()
	if !end {
		t.Errorf("unexpected end")
	}
	expected := "FILE 1"
	actual := fi.Name
	if expected != actual {
		t.Errorf("expected [%v] ; actual [%v]", expected, actual)
	}

	fi, end = w.next()
	if end {
		t.Errorf("expected end")
	}
}

func TestDirWalkTwoFiles(t *testing.T) {
	d := NewDisk("", "")
	dump := `
00016600  00 ff 82 11 00 46 49 4c  45 20 31 a0 a0 a0 a0 a0  |.....FILE 1.....|
00016610  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016620  00 00 82 11 01 46 49 4c  45 20 32 a0 a0 a0 a0 a0  |.....FILE 2.....|
00016630  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016640  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
`
	if err := binary.LoadStringDumpInto(dump, d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	w := newDirWalker(d)
	w.next()
	fi, end := w.next()
	if !end {
		t.Errorf("unexpected end")
		return
	}
	expected := "FILE 2"
	actual := fi.Name
	if expected != actual {
		t.Errorf("expected [%v] ; actual [%v]", expected, actual)
	}

	fi, end = w.next()
	if end {
		t.Errorf("expected end")
	}
}

func TestDirWalkDelFile(t *testing.T) {
	d := NewDisk("", "")
	dump := `
00016600  00 ff 82 11 00 46 49 4c  45 20 31 a0 a0 a0 a0 a0  |.....FILE 1.....|
00016610  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016620  00 00 00 11 01 46 49 4c  45 20 32 a0 a0 a0 a0 a0  |.....FILE 2.....|
00016630  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016640  00 00 82 11 02 46 49 4c  45 20 33 a0 a0 a0 a0 a0  |.....FILE 3.....|
00016650  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016660  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
`
	if err := binary.LoadStringDumpInto(dump, d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	w := newDirWalker(d)
	w.next()
	fi, end := w.next()
	if !end {
		t.Errorf("unexpected end")
		return
	}
	expected := "FILE 3"
	actual := fi.Name
	if expected != actual {
		t.Errorf("expected [%v] ; actual [%v]", expected, actual)
	}

	fi, end = w.next()
	if end {
		t.Errorf("expected end")
	}
}

func TestDirWalkFindDelFile(t *testing.T) {
	d := NewDisk("", "")
	dump := `
00016600  00 ff 82 11 00 46 49 4c  45 20 31 a0 a0 a0 a0 a0  |.....FILE 1.....|
00016610  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016620  00 00 00 11 01 46 49 4c  45 20 32 a0 a0 a0 a0 a0  |.....FILE 2.....|
00016630  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016640  00 00 82 11 02 46 49 4c  45 20 33 a0 a0 a0 a0 a0  |.....FILE 3.....|
00016650  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016660  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
`
	if err := binary.LoadStringDumpInto(dump, d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	w := newDirWalker(d)
	w.skipDeleted = false
	w.next()
	fi, _ := w.next()
	expected := "FILE 2"
	actual := fi.Name
	if expected != actual {
		t.Errorf("expected [%v] ; actual [%v]", expected, actual)
	}
}

func TestDirWalkTwoSectors(t *testing.T) {
	d := NewDisk("", "")
	dump := `
00016600  12 04 82 11 00 46 49 4c  45 20 31 a0 a0 a0 a0 a0  |.....FILE 1.....|
00016610  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016620  00 00 82 11 01 46 49 4c  45 20 32 a0 a0 a0 a0 a0  |.....FILE 2.....|
00016630  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016640  00 00 82 11 02 46 49 4c  45 20 33 a0 a0 a0 a0 a0  |.....FILE 3.....|
00016650  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016660  00 00 82 11 03 46 49 4c  45 20 34 a0 a0 a0 a0 a0  |.....FILE 4.....|
00016670  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016680  00 00 82 11 04 46 49 4c  45 20 35 a0 a0 a0 a0 a0  |.....FILE 5.....|
00016690  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
000166a0  00 00 82 11 05 46 49 4c  45 20 36 a0 a0 a0 a0 a0  |.....FILE 6.....|
000166b0  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
000166c0  00 00 82 11 06 46 49 4c  45 20 37 a0 a0 a0 a0 a0  |.....FILE 7.....|
000166d0  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
000166e0  00 00 82 11 07 46 49 4c  45 20 38 a0 a0 a0 a0 a0  |.....FILE 8.....|
000166f0  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016700  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
00016900  00 ff 82 11 08 46 49 4c  45 20 39 a0 a0 a0 a0 a0  |.....FILE 9.....|
00016910  a0 a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
00016920  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
`
	if err := binary.LoadStringDumpInto(dump, d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	w := newDirWalker(d)
	for i := 1; i <= 8; i++ {
		_, _ = w.next()
	}
	fi, end := w.next()
	if !end {
		t.Errorf("unexpected end")
		return
	}
	expected := "FILE 9"
	actual := fi.Name
	if expected != actual {
		t.Errorf("expected [%v] ; actual [%v]", expected, actual)
	}

}
