package d71

import (
	"fmt"
	"testing"

	"github.com/blackchip-org/vt128/binary"
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
	fmt.Printf("END OF ZERO BLOCK\n")
	w.Write([]byte{0xab})

	d2 := NewDisk("", "")
	report, _ := binary.Compare(d, d2)
	fmt.Printf("REPORT:\n%v", report)

	want := 0xab
	got := int(d[Offset(17, 6, 2)]) // Plus 2 for block link, interleave +6
	if want != got {
		t.Fatalf("wanted %v ; got %v", want, got)
	}
}
