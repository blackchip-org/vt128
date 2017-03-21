package binary

import (
	"bytes"
	"fmt"
)

type Diff struct {
	Pos int
	A   byte
	B   byte
}

type DiffReport []Diff

func (d DiffReport) String() string {
	if len(d) == 0 {
		return "ok"
	}
	var buf bytes.Buffer
	for _, diff := range d {
		buf.WriteString(fmt.Sprintf("%08x: %02x %02x\n", diff.Pos, diff.A, diff.B))
	}
	return buf.String()
}

func Compare(a []byte, b []byte) (DiffReport, bool) {
	same := true
	diff := make([]Diff, 0)
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			same = false
			diff = append(diff, Diff{Pos: i, A: a[i], B: b[i]})
		}
	}
	return diff, same
}
