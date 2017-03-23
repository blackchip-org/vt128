package d71

import "testing"

func TestWriteString(t *testing.T) {
	d := NewDisk("", "")
	expected := "XY"
	d.Editor().WriteString(expected)
	actual := string(d[0]) + string(d[1])
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestReadWord(t *testing.T) {
	d := NewDisk("", "")
	d[0] = 0x34
	d[1] = 0x12
	expected := 0x1234
	actual := d.Editor().ReadWord()
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}

func TestWriteWord(t *testing.T) {
	d := NewDisk("", "")
	d.Editor().WriteWord(0x1234)
	expected := 0x34
	actual := int(d[0])
	if expected != actual {
		t.Errorf("expected %x ; actual %x", expected, actual)
	}
	expected2 := 0x12
	actual2 := int(d[1])
	if expected2 != actual2 {
		t.Errorf("expected %x ; actual %x", expected2, actual2)
	}
}
