package drum

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

const (
	tOffset  int = 50 // Bytes offset where track data starts.
	tPadding int = 3  // Number of empty bytes after the track sample ID.
	pLen     int = 16 // Number of bytes used to store the track pattern.
)

// Pattern is the high level representation of the drum pattern contained
// in a .splice file.
type Pattern struct {
	HWVersion string
	BPM       float32
	Tracks    []Track
}

// String returns multiple lines describing the pattern's
// metadata and the beat pattern for each of its tracks.
func (p Pattern) String() string {
	buffer := bytes.NewBufferString("")

	fmt.Fprint(buffer, fmt.Sprintf("Saved with HW Version: %s\n", p.HWVersion))
	fmt.Fprint(buffer, fmt.Sprintf("Tempo: %v\n", p.BPM))
	for _, track := range p.Tracks {
		fmt.Fprint(buffer, fmt.Sprintf("%v\n", track))
	}

	return buffer.String()
}

// Track is a line representing the sample and which beats the sample
// should be played on.
type Track struct {
	SampleID   int
	SampleName string
	Pattern    []byte
}

// String returns the track ID, name, and beat pattern.
//
// TODO: handle errors in missing data
func (t Track) String() string {
	prelude := fmt.Sprintf("(%d) %s\t", t.SampleID, t.SampleName)

	bytesToSteps := func(steps []byte) string {
		var buf bytes.Buffer
		for _, b := range steps {
			if b == 0x00 {
				buf.WriteString("-")
			} else {
				buf.WriteString("x")
			}
		}

		return buf.String()
	}

	track := fmt.Sprintf("|%s|%s|%s|%s|",
		bytesToSteps(t.Pattern[0:4]),
		bytesToSteps(t.Pattern[4:8]),
		bytesToSteps(t.Pattern[8:12]),
		bytesToSteps(t.Pattern[12:16]),
	)

	return prelude + track
}

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	p := &Pattern{}

	f, err := os.Open(path)
	if err != nil {
		return p, err
	}

	scn := bufio.NewScanner(f)
	scn.Scan()
	b := scn.Bytes()

	err = scn.Err()
	if err != nil {
		return p, err
	}

	p.HWVersion = readHwVersion(b)
	p.BPM = readBPM(b)
	p.Tracks = readTracks(b[tOffset:])

	return p, nil
}

// readHwVersion extracts the HWVersion from the input data.
//
// TODO: isolate from the preceding pattern data
// TODO: handle errors for string conversion
func readHwVersion(data []byte) string {
	var str []byte
	for i := 14; i < 46; i++ {
		if data[i] != 0x00 {
			str = append(str, data[i])
		}
	}

	return string(str)
}

// readBPM extracts the BPM bytes and converts them to a float32 integer.
//
// TODO: isolate from the preceding pattern data
// TODO: handle errors for number conversion
func readBPM(data []byte) float32 {
	bpm := data[46:50]
	bits := binary.LittleEndian.Uint32(bpm)
	return math.Float32frombits(bits)
}

// readTracks takes in the byte array representing the data
// with the leading pattern data removed. The first track starts
// at byte 0 of the input byte slice.
//
// Since each track is variable length, it loops its way through the
// data and converts each section to track it represents as it goes.
//
// Extracted tracks are added to a results array and the values returned.
//
// TODO: error handling for incomplete tracks and malformed byte data
func readTracks(data []byte) []Track {
	var tracks []Track
	pos := 0

	for {
		tData := data[pos:]
		if len(tData) == 0 {
			break
		} else {
			t, read := readTrack(data[pos:])
			tracks = append(tracks, t)
			pos += read
		}
	}

	return tracks
}

// readTrack takes in a byte array and seeks through to extract
// one Track object. Note, the track extracted may represent a leading
// subset of the entire input.
//
// It returns the number of bytes used to store the Track so that
// the following Track position can be calculated if one exists.
//
// TODO: return errors if track data is incomplete or malformed
func readTrack(data []byte) (Track, int) {
	t := Track{}
	read := 0

	// Get Sample ID
	t.SampleID = int(data[read])
	read++

	// Skip over padding
	read += tPadding

	// Get Sample Name
	nLen := int(data[read])
	read++
	var name []byte

	for _, b := range data[read : read+nLen] {
		name = append(name, b)
	}
	t.SampleName = string(name)
	read += nLen

	// Get Track Data
	t.Pattern = data[read : read+pLen]
	read += pLen

	return t, read
}
