package drum

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
// TODO: implement
func DecodeFile(path string) (*Pattern, error) {
	p := &Pattern{}
	return p, nil
}

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

// Track is each line representing the sample and which
// beats the sample should be played on
type Track struct {
	SampleID   int
	SampleName string
	Pattern    []byte
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	HWVersion string
	BPM       float32
	Tracks    []Track
}
