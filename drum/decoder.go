package drum

import (
	"bufio"
	"os"
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
//    4         3      - padding
//    8         1      - name length
//    9         n      - pattern name
//    9 + n    16      - pattern (01 for note, 00 for none)

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

	return p, nil
}

// readHwVersion extracts the HWVersion from the input data
func readHwVersion(data []byte) string {
	hwVersion := []byte{}
	for hwIdx := 14; hwIdx < 46; hwIdx++ {
		if data[hwIdx] != 0x00 {
			hwVersion = append(hwVersion, data[hwIdx])
		}
	}

	return string(hwVersion)
}
