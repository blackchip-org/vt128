package d71

import "fmt"

type Writer struct {
	d      Disk
	e      *Editor
	track  int
	sector int
	offset int
}

func newWriter(d Disk, track int, sector int) *Writer {
	w := &Writer{d: d, track: track, sector: sector, offset: 2}
	w.e = d.Editor()
	return w
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (w *Writer) write(b byte) error {
	if w.offset >= SectorLen {
		// Find a free block
		newT, newS, ok := freeBlockNext(w.d, w.track, w.sector)
		if !ok {
			return fmt.Errorf("Disk full")
		}
		// Add link to next block
		w.e.Seek(w.track, w.sector, 0)
		w.e.Write(newT)
		w.e.Write(newS)

		w.track = newT
		w.sector = newS
		w.offset = 2 // Skip link marker

		// Goto next block and write EOF marker
		w.e.Seek(w.track, w.sector, 0)
		w.e.Write(0) // End of chain
		w.e.Write(0) // Zero bytes used in this block
	}
	return nil
}
