package d71

import "testing"

func TestDirInterleaveStart(t *testing.T) {
	d := NewDisk("", "")
	want := 4
	if got, _ := freeDirSector(d); want != got {
		t.Fatalf("want %v ; got %v", want, got)
	}
}

func TestDirInterleave(t *testing.T) {
	d := NewDisk("", "")
	for i := 0; i < 5; i++ {
		s, ok := freeDirSector(d)
		if !ok {
			t.Fatalf("wanted ok")
		}
		d.BamWrite(DirTrack, s, false)
	}
	want := 2
	if got, _ := freeDirSector(d); want != got {
		t.Fatalf("want %v ; got %v", want, got)
	}
}

func TestDirInterleave2(t *testing.T) {
	d := NewDisk("", "")
	for i := 0; i < 11; i++ {
		s, ok := freeDirSector(d)
		if !ok {
			t.Fatalf("wanted ok")
		}
		d.BamWrite(DirTrack, s, false)
	}
	want := 3
	if got, _ := freeDirSector(d); want != got {
		t.Fatalf("want %v ; got %v", want, got)
	}
}

func TestDirInterleaveFull(t *testing.T) {
	d := NewDisk("", "")
	// Note: first sector BAM, next sector already allocated by format
	for i := 0; i < Geom[DirTrack].Sectors-2; i++ {
		s, ok := freeDirSector(d)
		if !ok {
			t.Fatalf("wanted ok")
		}
		d.BamWrite(DirTrack, s, false)
	}
	if _, ok := freeDirSector(d); ok {
		t.Fatalf("wanted not ok")
	}
}

func TestFindStart(t *testing.T) {
	d := NewDisk("", "")
	twant := 17
	swant := 0
	tgot, sgot, ok := freeBlockFirst(d)
	if !ok {
		t.Fatalf("expected ok")
	}
	if twant != tgot || swant != sgot {
		t.Errorf("wanted %v %v ; got %v %v", twant, swant, tgot, sgot)
	}
}

func TestFindFirstTrack(t *testing.T) {
	d := NewDisk("", "")
	twant := 17
	swant := 3
	d.BamWrite(17, 0, false)
	d.BamWrite(17, 1, false)
	d.BamWrite(17, 2, false)
	tgot, sgot, ok := freeBlockFirst(d)
	if !ok {
		t.Fatalf("expected ok")
	}
	if twant != tgot || swant != sgot {
		t.Errorf("wanted %v %v ; got %v %v", twant, swant, tgot, sgot)
	}
}

func TestFindNextTrack(t *testing.T) {
	d := NewDisk("", "")
	for s := 0; s < Geom[17].Sectors; s++ {
		d.BamWrite(17, s, false)
	}
	d.BamWrite(16, 0, false)
	d.BamWrite(16, 1, false)
	twant := 16
	swant := 2
	tgot, sgot, ok := freeBlockFirst(d)
	if !ok {
		t.Fatalf("expected ok")
	}
	if twant != tgot || swant != sgot {
		t.Errorf("wanted %v %v ; got %v %v", twant, swant, tgot, sgot)
	}
}

func TestFindInterleave(t *testing.T) {
	d := NewDisk("", "")
	d.BamWrite(17, 0, false)
	twant := 17
	swant := 6
	tgot, sgot, ok := freeBlockNext(d, 17, 0)
	if !ok {
		t.Fatalf("expected ok")
	}
	if twant != tgot || swant != sgot {
		t.Errorf("wanted %v %v ; got %v %v", twant, swant, tgot, sgot)
	}
}

func TestFindInterleaveWrap(t *testing.T) {
	d := NewDisk("", "")
	track, sector, ok := freeBlockFirst(d)
	d.BamWrite(track, sector, false)
	for i := 0; i < 3; i++ {
		track, sector, _ = freeBlockNext(d, track, sector)
		d.BamWrite(track, sector, false)
	}
	tgot, sgot, ok := freeBlockNext(d, track, sector)
	twant := 17
	swant := 3
	if !ok {
		t.Fatalf("expected ok")
	}
	if twant != tgot || swant != sgot {
		t.Errorf("wanted %v %v ; got %v %v", twant, swant, tgot, sgot)
	}
}

func TestFindInterleaveNextTrack(t *testing.T) {
	d := NewDisk("", "")
	track, sector, ok := freeBlockFirst(d)
	d.BamWrite(track, sector, false)
	for i := 1; i < Geom[track].Sectors; i++ {
		track, sector, _ = freeBlockNext(d, track, sector)
		d.BamWrite(track, sector, false)
	}
	tgot, sgot, ok := freeBlockNext(d, track, sector)
	twant := 16
	swant := 0
	if !ok {
		t.Fatalf("expected ok")
	}
	if twant != tgot || swant != sgot {
		t.Errorf("wanted %v %v ; got %v %v", twant, swant, tgot, sgot)
	}
}

func TestFindInterleaveFull(t *testing.T) {
	d := NewDisk("", "")
	track, sector, ok := freeBlockFirst(d)
	d.BamWrite(track, sector, false)
	total := d.Info().Free
	for i := 0; i < total; i++ {
		track, sector, ok = freeBlockNext(d, track, sector)
		if !ok {
			t.Fatalf("expected ok")
		}
		d.BamWrite(track, sector, false)
	}
	_, _, ok = freeBlockNext(d, track, sector)
	if ok {
		t.Fatalf("expected not ok")
	}
}
