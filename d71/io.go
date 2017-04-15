package d71

import (
	"fmt"
	"io"
)

type Writer struct {
	d Disk
	e *Editor
}

func newWriter(d Disk, track int, sector int) *Writer {
	w := &Writer{d: d}
	w.e = d.Editor()
	w.seek(track, sector)
	return w
}

func (w *Writer) seek(track int, sector int) {
	w.e.Seek(track, sector, 0)
	w.e.Write(0) // No next track link yet
	w.e.Write(2) // Two bytes used so far
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
		w.seek(newT, newS)
	} else {
		w.e.Move(1)
	}

	return nil
}

type Reader struct {
	d          Disk
	e          *Editor
	nextTrack  int
	nextSector int
	len        int
}

func newReader(d Disk, track int, sector int) *Reader {
	r := &Reader{d: d}
	r.e = d.Editor()
	r.seek(track, sector)
	return r
}

func (r *Reader) Read(p []byte) (n int, err error) {
	for ; n < len(p); n++ {
		val, err := r.read()
		if err != nil {
			return n, err
		}
		p[n] = val
	}
	return n, nil
}

func (r *Reader) seek(track int, sector int) {
	r.e.Seek(track, sector, 0)
	r.nextTrack = r.e.Read()
	if r.nextTrack == 0 {
		r.len = r.e.Read()
	} else {
		r.nextSector = r.e.Read()
	}
}

func (r *Reader) read() (byte, error) {
	if r.nextTrack == 0 && r.e.At() == r.len {
		return 0, io.EOF
	}
	if r.e.At() == SectorLen {
		r.seek(r.nextTrack, r.nextSector)
	}
	val := byte(r.e.Read())
	// If we moved off the sector, seek to the next block
	if r.e.At() == 0 {
		r.seek(r.nextTrack, r.nextSector)
	}
	return val, nil
}
