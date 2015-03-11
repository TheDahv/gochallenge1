package drum

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"testing"
)

func TestDecodeFile(t *testing.T) {
	tData := []struct {
		path   string
		output string
	}{
		{"pattern_1.splice",
			`Saved with HW Version: 0.808-alpha
Tempo: 120
(0) kick	|x---|x---|x---|x---|
(1) snare	|----|x---|----|x---|
(2) clap	|----|x-x-|----|----|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(4) hh-close	|x---|x---|----|x--x|
(5) cowbell	|----|----|--x-|----|
`,
		},
		{"pattern_2.splice",
			`Saved with HW Version: 0.808-alpha
Tempo: 98.4
(0) kick	|x---|----|x---|----|
(1) snare	|----|x---|----|x---|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(5) cowbell	|----|----|x---|----|
`,
		},
		{"pattern_3.splice",
			`Saved with HW Version: 0.808-alpha
Tempo: 118
(40) kick	|x---|----|x---|----|
(1) clap	|----|x---|----|x---|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(5) low-tom	|----|---x|----|----|
(12) mid-tom	|----|----|x---|----|
(9) hi-tom	|----|----|-x--|----|
`,
		},
		{"pattern_4.splice",
			`Saved with HW Version: 0.909
Tempo: 240
(0) SubKick	|----|----|----|----|
(1) Kick	|x---|----|x---|----|
(99) Maracas	|x-x-|x-x-|x-x-|x-x-|
(255) Low Conga	|----|x---|----|x---|
`,
		},
		{"pattern_5.splice",
			`Saved with HW Version: 0.708-alpha
Tempo: 999
(1) Kick	|x---|----|x---|----|
(2) HiHat	|x-x-|x-x-|x-x-|x-x-|
`,
		},
	}

	for _, exp := range tData {
		decoded, err := DecodeFile(path.Join("..", "data", "fixtures", exp.path))
		if err != nil {
			t.Fatalf("something went wrong decoding %s - %v", exp.path, err)
		}
		if fmt.Sprint(decoded) != exp.output {
			t.Logf("decoded:\n%#v\n", fmt.Sprint(decoded))
			t.Logf("expected:\n%#v\n", exp.output)
			t.Fatalf("%s wasn't decoded as expected.\nGot:\n%s\nExpected:\n%s",
				exp.path, decoded, exp.output)
		}
	}
}

func TestInvalidPathsReturnsError(t *testing.T) {
	_, err := DecodeFile("meow")
	if err == nil {
		t.Error("Expected an error on file open, and didn't get one")
	}
}

func loadBytesFromFile(fn string) ([]byte, error) {
	f, err := os.Open(path.Join("..", "data", "fixtures", fn))

	if err != nil {
		return nil, err
	}

	scn := bufio.NewScanner(f)
	scn.Scan()

	return scn.Bytes(), nil
}

func TestLoadsHWVersion(t *testing.T) {
	tData := []struct {
		path      string
		HWVersion string
	}{
		{
			"pattern_1.splice",
			"0.808-alpha",
		},
		{
			"pattern_2.splice",
			"0.808-alpha",
		},
		{
			"pattern_3.splice",
			"0.808-alpha",
		},
		{
			"pattern_4.splice",
			"0.909",
		},
		{
			"pattern_5.splice",
			"0.708-alpha",
		},
	}

	for _, tCase := range tData {
		b, err := loadBytesFromFile(tCase.path)
		if err != nil {
			t.Errorf("Failure to open file %s\n", tCase.path)
			t.FailNow()
		}

		result := readHwVersion(b[14:46])
		if result != tCase.HWVersion {
			t.Logf("decoded:\n%#v\n", result)
			t.Logf("expected:\n%#v\n", tCase.HWVersion)
			t.Fatalf("%s wasn't decoded as expected.\nGot:\n%s\nExpected:\n%s",
				tCase.path, result, tCase.HWVersion)
		}
	}
}

func TestLoadsBPM(t *testing.T) {
	tData := []struct {
		path string
		BPM  float32
	}{
		{
			"pattern_1.splice",
			120,
		},
		{
			"pattern_2.splice",
			98.4,
		},
		{
			"pattern_3.splice",
			118,
		},
		{
			"pattern_4.splice",
			240,
		},
		{
			"pattern_5.splice",
			999,
		},
	}

	for _, tCase := range tData {
		b, err := loadBytesFromFile(tCase.path)
		if err != nil {
			t.Errorf("Failure to open file %s\n", tCase.path)
			t.FailNow()
		}

		result := readBPM(b[46:50])
		if result != tCase.BPM {
			t.Logf("decoded:\n%#v\n", result)
			t.Logf("expected:\n%#v\n", tCase.BPM)
			t.Fatalf("%s wasn't decoded as expected.\nGot:\n%f\nExpected:\n%f",
				tCase.path, result, tCase.BPM)
		}
	}
}

// There's probably a better way to enter a byte array literal
// This is the Kick pattern from Sample 3
// (40) kick	|x---|----|x---|----|
var trackData = []byte{0x28, 0x00, 0x00, 0x00, 0x04, 0x6B, 0x69, 0x63, 0x6B, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

func TestLoadsSingleTrack(t *testing.T) {
	tracks := readTracks(trackData)
	if numTracks := len(tracks); numTracks != 1 {
		t.Errorf("Expected to get 1 track back, got %d\n", numTracks)
	}

	track := tracks[0]
	if track.SampleID != 40 {
		t.Errorf("Got SampleID of %d, expected 40\n", track.SampleID)
	}

	if track.SampleName != "kick" {
		t.Errorf("Got sample name '%s', expected 'kick'\n", track.SampleName)
	}
}

func TestMultipleTracks(t *testing.T) {
	newTrackData := []byte{0x01, 0x00, 0x00, 0x00, 0x04, 0x63, 0x6C, 0x61, 0x70, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}
	multipleTracks := []byte{}
	multipleTracks = append(multipleTracks, trackData...)
	multipleTracks = append(multipleTracks, newTrackData...)

	tracks := readTracks(multipleTracks)
	if numTracks := len(tracks); numTracks != 2 {
		t.Errorf("Got %d tracks back, expected 2\n", numTracks)
	}

	track1 := tracks[0]
	track2 := tracks[1]

	if track1.SampleID != 40 {
		t.Errorf("Got SampleID %d, expected 40\n", track1.SampleID)
	}
	if track1.SampleName != "kick" {
		t.Errorf("Got SampleName '%s', expected 'kick'\n", track1.SampleName)
	}

	if track2.SampleID != 1 {
		t.Errorf("Got SampleID %d, expected 1\n", track2.SampleID)
	}
	if track2.SampleName != "clap" {
		t.Errorf("Got SampleName '%s', expected 'clap'\n", track2.SampleName)
	}
}

func TestReadTrack(t *testing.T) {
	trackData := []byte{0x28, 0x00, 0x00, 0x00, 0x04, 0x6B, 0x69, 0x63, 0x6B, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	track, trackSize := readTrack(trackData)
	if trackSize != len(trackData) {
		t.Errorf("Got track %d bytes read, expected %d\n", trackSize, len(trackData))
	}

	if track.SampleID != 40 {
		t.Errorf("Got SampleID %d, expected %d\n", track.SampleID, 40)
	}

	if track.SampleName != "kick" {
		t.Errorf("Got SampleName '%s', expected 'kick'\n", track.SampleName)
	}
}

func TestString(t *testing.T) {
	exp := "(40) kick	|x---|----|x---|----|"
	track, _ := readTrack(trackData)

	if result := track.String(); result != exp {
		t.Errorf("Got tracks tring '%s', expected '%s'\n", result, exp)
	}
}
