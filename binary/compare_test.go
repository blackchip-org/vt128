package binary

import (
	"strings"
	"testing"
)

func TestCompareEqual(t *testing.T) {
	a := []byte{1, 2, 3, 4}
	b := []byte{1, 2, 3, 4}
	if _, same := Compare(a, b); !same {
		t.Errorf("expected to be the same")
	}
}

func TestCompareNotEqual(t *testing.T) {
	a := []byte{1, 2, 3, 4, 5}
	b := []byte{1, 22, 3, 44, 5}
	diff, same := Compare(a, b)
	if same {
		t.Errorf("expected to not be the same")
		return
	}
	if len(diff) != 2 {
		t.Errorf("expected two differences ; actual %v", len(diff))
	}
	expected := Diff{Pos: 1, A: 2, B: 22}
	actual := diff[0]
	if expected != actual {
		t.Errorf("expected %+v ; actual %+v", expected, actual)
	}

	expected2 := Diff{Pos: 3, A: 4, B: 44}
	actual2 := diff[1]
	if expected2 != actual2 {
		t.Errorf("expected %+v ; actual %+v", expected2, actual2)
	}
}

func TestCompareNotEqualReport(t *testing.T) {
	a := []byte{1, 2, 3, 4, 5}
	b := []byte{1, 0x22, 3, 0x44, 5}
	diff, same := Compare(a, b)
	if same {
		t.Errorf("expected to not be the same")
		return
	}
	expected := strings.TrimSpace(`
00000001: 02 22
00000003: 04 44
	`)
	actual := strings.TrimSpace(diff.String())
	if expected != actual {
		t.Errorf("\nexpected\n%v\nactual\n%v\n", expected, actual)
	}

}
