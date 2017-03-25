package d71

const (
	dirInterleave  = 3
	fileInterleave = 6
)

var trackOrder = make([]int, MaxTrack+1, MaxTrack+1)

func init() {
	i := 0
	for track := DirTrack - 1; track >= 1; track-- {
		trackOrder[i] = track
		i++
	}
	for track := DirTrack + 1; track < Flip; track++ {
		trackOrder[i] = track
		i++
	}
	for track := BamTrack - 1; track >= Flip; track-- {
		trackOrder[i] = track
		i++
	}
	for track := BamTrack + 1; track <= MaxTrack; track++ {
		trackOrder[i] = track
		i++
	}
}

func freeDirSector(d Disk) (sector int, ok bool) {
	rem := Geom[DirTrack].Sectors - 1
	i := 1
	for {
		if d.BamRead(DirTrack, i) {
			return i, true
		}
		rem--
		if rem == 0 {
			return 0, false
		}
		i = (i + dirInterleave)
		if i >= Geom[DirTrack].Sectors {
			i = (i % Geom[DirTrack].Sectors) + 2
		}
	}
}

func freeBlockFirst(d Disk) (track int, sector int, ok bool) {
	for i := 0; i < len(trackOrder); i++ {
		track := trackOrder[i]
		if d.TrackInfo(track).Free == 0 {
			continue
		}
		for sector := 0; sector < Geom[track].Sectors; sector++ {
			if d.BamRead(track, sector) {
				return track, sector, true
			}
		}
	}
	return 0, 0, false
}

func freeBlockNext(d Disk, track int, sector int) (int, int, bool) {
	if d.TrackInfo(track).Free == 0 {
		return freeBlockFirst(d)
	}
	sector = (sector + fileInterleave) % Geom[track].Sectors
	for {
		if d.BamRead(track, sector) {
			return track, sector, true
		}
		sector = (sector + 1) % Geom[track].Sectors
	}
}
