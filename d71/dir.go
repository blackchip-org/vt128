package d71

import (
	"strings"
)

type FileType int

const (
	Del FileType = iota
	Seq
	Prg
	Usr
	Rel
)

type FileInfo struct {
	Type   FileType
	SaveAt bool
	Locked bool
	Splat  bool
	Name   string
	Size   int
	Track  int
	Sector int
}

type dirWalker struct {
	e          *Editor
	nextTrack  int
	nextSector int
	entry      int
}

func newDirWalker(d Disk) *dirWalker {
	w := &dirWalker{e: d.Editor()}
	w.e.Seek(DirTrack, 0, 0)

	firstTrack := w.e.Read()
	firstSector := w.e.Read()
	w.e.Seek(firstTrack, firstSector, 0)
	w.readLink()
	return w
}

func (w *dirWalker) readLink() {
	e := w.e.Mark()
	w.nextTrack = e.Read()
	w.nextSector = e.Read()
}

func (w *dirWalker) advance() bool {
	w.entry++
	if w.entry == 8 {
		if w.nextTrack == 0 {
			return false
		}
		w.entry = 0
		w.e.Seek(w.nextTrack, w.nextSector, 0)
		w.readLink()
	} else {
		w.e.Move(0x20)
	}
	return true
}

func (w *dirWalker) next() (*FileInfo, bool) {
	for {
		e := w.e.Mark().Move(2)
		ftype := e.Read()
		if ftype == 0 {
			ok := w.advance()
			if !ok {
				return nil, false
			}
		} else {
			break
		}
	}

	fi := &FileInfo{}
	w.e.Move(2)
	ftype := w.e.Read()
	fi.Type = FileType(ftype & 0x7)
	fi.SaveAt = ftype&(1<<5) > 0
	fi.Locked = ftype&(1<<6) > 0
	fi.Splat = ftype&(1<<7) == 0
	fi.Track = w.e.Read()
	fi.Sector = w.e.Read()
	fi.Name = strings.Trim(w.e.ReadString(16), "\xa0")
	w.e.Move(0x1e - 0x15)
	fi.Size = w.e.ReadWord()
	return fi, true
}
