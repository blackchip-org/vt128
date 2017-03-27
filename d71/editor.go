package d71

import "bytes"

type Editor struct {
	disk Disk
	Pos  Pos
}

func (e *Editor) Mark() *Editor {
	return &Editor{disk: e.disk, Pos: e.Pos}
}

func (e *Editor) Move(delta int) *Editor {
	e.Pos.Move(delta)
	return e
}

func (e *Editor) Seek(track int, sector int, at int) {
	e.Pos.Seek(track, sector, at)
}

func (e *Editor) Peek() int {
	return int(e.disk[e.Pos.Offset()])
}

func (e *Editor) Poke(val int) {
	e.disk[e.Pos.Offset()] = byte(val)
}

func (e *Editor) Write(val int) {
	e.disk[e.Pos.Offset()] = byte(val)
	e.Pos.Move(1)
}

func (e *Editor) Read() int {
	v := int(e.disk[e.Pos.Offset()])
	e.Pos.Move(1)
	return v
}

func (e *Editor) ReadWord() int {
	return e.Read() + (e.Read() << 8)
}

func (e *Editor) WriteWord(val int) {
	e.Write(val & 0xff)
	e.Write(val >> 8)
}

func (e *Editor) WriteString(val string) {
	n := len(val)
	off := e.Pos.Offset()
	for i := 0; i < n; i++ {
		e.disk[off+i] = byte(val[i])
	}
	e.Pos.Move(n)
}

func (e *Editor) ReadString(length int) string {
	var buf bytes.Buffer
	off := e.Pos.Offset()
	for i := 0; i < length; i++ {
		buf.WriteByte(e.disk[off+i])
	}
	e.Pos.Move(length)
	return buf.String()
}

func (e *Editor) Fill(val int, length int) {
	off := e.Pos.Offset()
	for i := 0; i < length; i++ {
		e.disk[off+i] = byte(val)
	}
	e.Pos.Move(length)
}

func (e *Editor) WriteStringN(val string, pad int, length int) {
	e.WriteString(val)
	n := length - len(val)
	if n <= 0 {
		return
	}
	for i := 0; i < n; i++ {
		e.Write(pad)
	}
}

func (e *Editor) Track() int {
	return e.Pos.Track
}

func (e *Editor) Sector() int {
	return e.Pos.Sector
}

func (e *Editor) At() int {
	return e.Pos.At
}
