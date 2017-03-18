package d71

import "testing"

func TestWriteString(t *testing.T) {
	d, _ := NewDisk("", "")
	expected := "XY"
	d.Editor().Move(0x160a2).WriteString(expected)
	actual := string(d[0x160a2]) + string(d[0x160a3])
	if expected != actual {
		t.Errorf("expected %v ; actual %v", expected, actual)
	}
}
