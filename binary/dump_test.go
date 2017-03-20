package binary

import (
	"strings"
	"testing"
)

func TestLoadDumpInto(t *testing.T) {
	dump := strings.NewReader("00000000 ab cd ef")
	b := make([]byte, 32, 32)
	if err := LoadDumpInto(dump, b); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := byte(0xcd)
	actual := b[1]
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestLoadDumpIntoMultiline(t *testing.T) {
	dump := strings.NewReader(`
		00000000 ab cd ef
		00000010 12 34 56
	`)
	b := make([]byte, 32, 32)
	if err := LoadDumpInto(dump, b); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := byte(0x34)
	actual := b[0x11]
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestLoadDumpIntoInvalidPos(t *testing.T) {
	dump := strings.NewReader(`
		00000000 ab cd ef
		0000001x 12 34 56
	`)
	b := make([]byte, 32, 32)
	err := LoadDumpInto(dump, b)
	if err == nil {
		t.Errorf("Expected error: %v", err)
	}
	expected := "line 3: strconv.ParseUint: parsing \"0000001x\": invalid syntax"
	actual := err.Error()
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestLoadDumpIntoOutOfRangePos(t *testing.T) {
	dump := strings.NewReader(`
		00000000 ab cd ef
		00001000 12 34 56
	`)
	b := make([]byte, 32, 32)
	err := LoadDumpInto(dump, b)
	if err == nil {
		t.Errorf("Expected error: %v", err)
	}
	expected := "line 3: no such position: 0x1000"
	actual := err.Error()
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestLoadDumpIntoInvalidValue(t *testing.T) {
	dump := strings.NewReader(`
		00000000 ab xx ef
	`)
	b := make([]byte, 32, 32)
	err := LoadDumpInto(dump, b)
	if err == nil {
		t.Errorf("Expected error: %v", err)
	}
	expected := "line 2: strconv.ParseUint: parsing \"xx\": invalid syntax"
	actual := err.Error()
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestLoadDumpIntoRepeatingLine(t *testing.T) {
	dump := strings.NewReader(`
		00000000 00 01 02 03 04 05 06 07
		*
		00000020 ff ff ff ff ff ff ff ff
	`)
	b := make([]byte, 0x30, 0x30)
	if err := LoadDumpInto(dump, b); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := byte(0x04)
	actual := b[0x14]
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestLoadDumpIntoIgnoreTrailing(t *testing.T) {
	dump := `
00000000  00 ff 82 11 00 46 49 4c  45 20 31 a0 a0 a0 a0 a0  |.....FILE 1.....|
00000010  ee a0 a0 a0 a0 00 00 00  00 00 00 00 00 00 01 00  |................|
`
	b := make([]byte, 0x30, 0x30)
	if err := LoadStringDumpInto(dump, b); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := byte(0xee)
	actual := b[0x10]
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}
