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
	// The number of empty bytes after the track sample ID
	trackPaddingLength int = 3
	// The number of bytes used to store the track pattern
	patternLength int = 16
	// The number of bytes into a pattern to start looking for
	// track data
	tracksOffset int = 50
)

// Pattern description
//    offset - length  - value
//     0        6      - "SPLICE"
//     7        5      - padding
//    13        1      - ???
//    14       32      - HW Name
//    46        4      - BPM (Little Endian Int32)
// -- loop --
//    0         1      - pattern (Hex -> Int)
//    1         3      - padding
//    4         1      - name length
//    5         n      - pattern name
//    5 + n + 1    16      - pattern (01 for note, 00 for none)

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	HWVersion string
	BPM       float32
	Tracks    []Track
}

// Track is each line representing the sample and which
// beats the sample should be played on
type Track struct {
	SampleID   int
	SampleName string
	Pattern    []byte
}

// String Given a Track, arrange its data into a string representing
// the track information and beat pattern
//
// TODO: handle errors in missing data
func (t Track) String() string {
	prelude := fmt.Sprintf("(%d) %s\t", t.SampleID, t.SampleName)

	bytesToSteps := func(beat []byte) string {
		var buffer bytes.Buffer
		for _, byte := range beat {
			if byte == 0x00 {
				buffer.WriteString("-")
			} else {
				buffer.WriteString("x")
			}
		}

		return buffer.String()
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
func DecodeFile(path string) (p *Pattern, err error) {
	p = &Pattern{}

	file, err := os.Open(path)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	bytes := scanner.Bytes()

	err = scanner.Err()
	if err != nil {
		return
	}

	p.HWVersion = readHwVersion(bytes)
	p.BPM = readBPM(bytes)
	p.Tracks = readTracks(bytes[tracksOffset:])

	return p, nil
}

// readHwVersion extracts the HWVersion from the input data
// TODO: isolate from the preceding pattern data
// TODO: handle errors for string conversion
func readHwVersion(data []byte) string {
	hwVersion := []byte{}
	for hwIdx := 14; hwIdx < 46; hwIdx++ {
		if data[hwIdx] != 0x00 {
			hwVersion = append(hwVersion, data[hwIdx])
		}
	}

	return string(hwVersion)
}

// readBPM extracts the BPM bytes and converts them to a float32 integer
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
// Extracted tracks are added to a results array and the values returned
//
// TODO: error handling for incomplete tracks and malformed byte data
func readTracks(tracksData []byte) (tracks []Track) {
	tracks = []Track{}
	lastReadPos := 0

	for {
		trackData := tracksData[lastReadPos:]
		if len(trackData) == 0 {
			break
		} else {
			track, bytesRead := readTrack(tracksData[lastReadPos:])
			tracks = append(tracks, track)
			lastReadPos += bytesRead
		}
	}

	return
}

// readTrack takes in a byte array and seeks through to extract
// one Track object. It returns the number of bytes used
// to store the Track so that the following Track position can be
// calculated if one exists
//
// TODO: return errors if track data is incomplete or malformed
func readTrack(tracksData []byte) (track Track, bytesRead int) {
	track = Track{}
	bytesRead = 0

	// Get Sample ID
	track.SampleID = int(tracksData[bytesRead])
	bytesRead++

	// Skip over padding
	bytesRead += trackPaddingLength

	// Get Sample Name
	nameLength := int(tracksData[bytesRead])
	bytesRead++
	nameBytes := []byte{}

	for _, byte := range tracksData[bytesRead : bytesRead+nameLength] {
		nameBytes = append(nameBytes, byte)
	}
	track.SampleName = string(nameBytes)
	bytesRead += nameLength

	// Get Track Data
	track.Pattern = tracksData[bytesRead : bytesRead+patternLength]
	bytesRead += patternLength

	return
}
