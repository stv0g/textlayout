package truetype

import (
	"encoding/binary"
	"errors"
)

type TableTrak struct {
	Horizontal, Vertical TrakData // may be empty
}

func parseTrakTable(data []byte) (out TableTrak, err error) {
	if len(data) < 12 {
		return out, errors.New("invalid trak table (EOF)")
	}
	// ignoring version and format
	horizOffset := binary.BigEndian.Uint16(data[6:])
	vertOffset := binary.BigEndian.Uint16(data[6:])

	if horizOffset != 0 {
		out.Horizontal, err = parseTrakData(data, int(horizOffset))
		if err != nil {
			return out, err
		}
	}
	if vertOffset != 0 {
		out.Vertical, err = parseTrakData(data, int(vertOffset))
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

type TrackEntry struct {
	Track           float32
	NameIndex       NameID
	PerSizeTracking []int16 // in font units
}

type TrakData struct {
	Entries []TrackEntry
	Sizes   []float32
}

func parseTrakData(data []byte, offset int) (out TrakData, err error) {
	if len(data) < int(offset)+8 {
		return out, errors.New("invalid trak data table (EOF)")
	}
	nTracks := binary.BigEndian.Uint16(data[offset:])
	nSizes := binary.BigEndian.Uint16(data[offset+2:])
	sizeTableOffset := int(binary.BigEndian.Uint32(data[offset+4:]))

	if len(data) < offset+8+8*int(nTracks) {
		return out, errors.New("invalid trak data table (EOF)")
	}

	out.Entries = make([]TrackEntry, nTracks)
	for i := range out.Entries {
		out.Entries[i].Track = fixed1616ToFloat(binary.BigEndian.Uint32(data[offset+8+i*8:]))
		out.Entries[i].NameIndex = NameID(binary.BigEndian.Uint16(data[offset+8+i*8+4:]))
		sizeOffset := binary.BigEndian.Uint16(data[offset+8+i*8+6:])
		out.Entries[i].PerSizeTracking, err = parseTrackSizes(data, int(sizeOffset), nSizes)
		if err != nil {
			return out, err
		}
	}

	if len(data) < sizeTableOffset+8+4*int(nSizes) {
		return out, errors.New("invalid trak data table (EOF)")
	}
	out.Sizes = make([]float32, nSizes)
	for i := range out.Sizes {
		out.Sizes[i] = fixed1616ToFloat(binary.BigEndian.Uint32(data[sizeTableOffset+8+4*i:]))
	}

	return out, nil
}

func parseTrackSizes(data []byte, offset int, count uint16) ([]int16, error) {
	if len(data) < offset+int(count)*2 {
		return nil, errors.New("invalid trak table per-sizes values (EOF)")
	}
	out := make([]int16, count)
	for i := range out {
		out[i] = int16(binary.BigEndian.Uint16(data[offset+2*i:]))
	}
	return out, nil
}