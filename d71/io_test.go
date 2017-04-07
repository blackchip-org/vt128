package d71

import (
	"io"
	"testing"
)

func TestWriter(t *testing.T) {
	d := NewDisk("", "")
	w := newWriter(d, 17, 0)
	w.Write([]byte{0xaa, 0xbb, 0xcc})
	want := 0xcc
	got := int(d[Offset(17, 0, 4)]) // Plus 2 for block link
	if want != got {
		t.Fatalf("wanted %v ; got %v", want, got)
	}
}

func TestWriterNextBlock(t *testing.T) {
	d := NewDisk("", "")
	w := newWriter(d, 17, 0)
	block := make([]byte, 254) /// Minus 2 for block link
	w.Write(block)
	w.Write([]byte{0xab})

	want := 0xab
	got := int(d[Offset(17, 6, 2)]) // Plus 2 for block link, interleave +6
	if want != got {
		t.Fatalf("wanted %v ; got %v", want, got)
	}
}

func TestReader(t *testing.T) {
	d := NewDisk("", "")
	e := d.Editor()
	e.Seek(17, 0, 0)
	e.Write(0)
	e.Write(4)
	e.Write(0xab)
	e.Write(0xcd)

	r := newReader(d, 17, 0)
	buf := make([]byte, 254, 254)
	n, err := r.Read(buf)
	if err != io.EOF {
		t.Fatalf("wanted eof error ; got %v", err)
	}
	if n != 2 {
		t.Fatalf("wanted n 2 ; got %v", n)
	}
	if buf[0] != 0xab {
		t.Fatalf("wanted 0xab ; got %x", buf[0])
	}
	if buf[1] != 0xcd {
		t.Fatalf("wanted 0xcd ; got %x", buf[1])
	}
}

func TestReaderNextBlock(t *testing.T) {
	d := NewDisk("", "")
	e := d.Editor()
	e.Seek(17, 0, 0)
	e.Write(17)
	e.Write(8)
	e.Seek(17, 8, 0)
	e.Write(0)
	e.Write(4)
	e.Write(0xab)
	e.Write(0xcd)

	r := newReader(d, 17, 0)
	buf := make([]byte, 254, 254)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("wanted no error ; got %v", err)
	}
	if n != 254 {
		t.Fatalf("wanted n 254 ; got %v", n)
	}
	n, err = r.Read(buf)
	if err != io.EOF {
		t.Fatalf("wanted eof error ; got %v", err)
	}
	if n != 2 {
		t.Fatalf("wanted n 2 ; got %v", n)
	}
	if buf[0] != 0xab {
		t.Fatalf("wanted 0xab ; got %x", buf[0])
	}
	if buf[1] != 0xcd {
		t.Fatalf("wanted 0xcd ; got %x", buf[1])
	}
}
