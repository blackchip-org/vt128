package d71

import "bytes"

type Editor struct {
	disk Disk
	Pos  int
}

func (e *Editor) Mark() *Editor {
	return &Editor{disk: e.disk, Pos: e.Pos}
}

func (e *Editor) Move(delta int) *Editor {
	e.Pos += delta
	return e
}

func (e *Editor) Seek(track int, sector int, offset int) {
	toff := Geom[track].Offset
	e.Pos = toff + (sector * SectorLen) + offset
}

func (e *Editor) Peek() int {
	return int(e.disk[e.Pos])
}

func (e *Editor) Poke(val int) {
	e.disk[e.Pos] = byte(val)
}

func (e *Editor) Write(val int) {
	e.disk[e.Pos] = byte(val)
	e.Pos++
}

func (e *Editor) Read() int {
	v := int(e.disk[e.Pos])
	e.Pos++
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
	for i := 0; i < n; i++ {
		e.disk[e.Pos+i] = byte(val[i])
	}
	e.Pos += n
}

func (e *Editor) ReadString(length int) string {
	var buf bytes.Buffer
	for i := 0; i < length; i++ {
		buf.WriteByte(e.disk[e.Pos+i])
	}
	e.Pos += length
	return buf.String()
}

func (e *Editor) Fill(val int, length int) {
	for i := 0; i < length; i++ {
		e.disk[e.Pos+i] = byte(val)
	}
	e.Pos += length
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
