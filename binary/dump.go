package binary

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func LoadStringDumpInto(src string, tgt []byte) error {
	return LoadDumpInto(strings.NewReader(src), tgt)
}

func LoadDumpInto(src io.Reader, tgt []byte) error {
	s := bufio.NewScanner(src)
	n := uint64(len(tgt))

	line := 0
	var prev bytes.Buffer
	var pos uint64
	repeating := false

	for {
		if ok := s.Scan(); !ok {
			if s.Err() != nil {
				return s.Err()
			}
			return nil
		}
		str := strings.TrimSpace(s.Text())
		line++
		if len(str) == 0 {
			continue
		}
		if strings.HasPrefix("*", str) {
			repeating = true
			continue
		}
		f := strings.Fields(str)
		newPos, err := strconv.ParseUint(f[0], 16, 0)
		if err != nil {
			return lerror(line, err)
		}
		if len(f) > 17 {
			f = f[:17]
		}
		if repeating {
			vals := prev.Bytes()
			for i := 0; pos < newPos; pos, i = pos+1, i+1 {
				if i >= len(vals) {
					i = 0
				}
				tgt[pos] = byte(vals[i])
			}
			repeating = false
		}
		prev.Reset()
		pos = newPos
		for i := 1; i < len(f); i++ {
			if err := checkPos(pos, n); err != nil {
				return lerror(line, err)
			}
			val, err := strconv.ParseUint(f[i], 16, 8)
			if err != nil {
				return lerror(line, err)
			}
			tgt[pos] = byte(val)
			prev.WriteByte(byte(val))
			pos++
		}
	}
}

func lerrorf(line int, format string, args ...interface{}) error {
	return fmt.Errorf("%v: %v", line, fmt.Sprintf(format, args))
}

func lerror(line int, err error) error {
	return fmt.Errorf("line %v: %v", line, err)
}

func checkPos(pos uint64, n uint64) error {
	if pos >= n {
		return fmt.Errorf("no such position: 0x%x", pos)
	}
	return nil
}
