package d71

import "fmt"

type Writer struct {
	d Disk
	e *Editor
}

func newWriter(d Disk, track int, sector int) *Writer {
	w := &Writer{d: d}
	w.e = d.Editor()
	w.e.Seek(track, sector, 2) // Offset to first data byte after link
	return w
}

func (w *Writer) Write(p []byte) (n int, err error) {
	n = 0
	for _, val := range p {
		err := w.write(val)
		if err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

func (w *Writer) write(b byte) error {
	w.e.Poke(int(b))

	// Add to the number of bytes used in this block
	ptr := w.e.Mark()
	ptr.Seek(ptr.Track(), ptr.Sector(), 1)
	ptr.Write(ptr.Peek() + 1)

	if w.e.At() == SectorLen-1 {
		// Find a free block
		newT, newS, ok := freeBlockNext(w.d, w.e.Track(), w.e.Sector())
		if !ok {
			return fmt.Errorf("Disk full")
		}
		// Add link to next block
		w.e.Seek(w.e.Track(), w.e.Sector(), 0)
		w.e.Write(newT)
		w.e.Write(newS)

		// Goto next block and write EOF marker
		w.e.Seek(newT, newS, 0)
		w.e.Write(0) // End of chain
		w.e.Write(0) // Zero bytes used in this block
	} else {
		w.e.Move(1)
	}

	return nil
}
