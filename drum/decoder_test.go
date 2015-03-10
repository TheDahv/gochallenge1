package drum

import (
	"fmt"
	"path"
	"testing"
)

func TestDecodeFile(t *testing.T) {
	t.Skip()

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
		decoded, err := DecodeFile(path.Join("..", "fixtures", exp.path))
		if err != nil {
			t.Fatalf("something went wrong decoding %s - %v", exp.path, err)
		}
		if fmt.Sprint(decoded) != exp.output {
			t.Logf("decoded:\n%#v\n", fmt.Sprint(decoded))
			t.Logf("expected:\n%#v\n", exp.output)
			t.Fatalf("%s wasn't decoded as expect.\nGot:\n%s\nExpected:\n%s",
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
		decoded, err := DecodeFile(path.Join("..", "data", "fixtures", tCase.path))
		if err != nil {
			t.Fatalf("something went wrong decoding %s - %v", tCase.path, err)
		}
		if decoded.HWVersion != tCase.HWVersion {
			t.Logf("decoded:\n%#v\n", decoded.HWVersion)
			t.Logf("expected:\n%#v\n", tCase.HWVersion)
			t.Fatalf("%s wasn't decoded as expect.\nGot:\n%s\nExpected:\n%s",
				tCase.path, decoded.HWVersion, tCase.HWVersion)
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
		decoded, err := DecodeFile(path.Join("..", "data", "fixtures", tCase.path))
		if err != nil {
			t.Fatalf("something went wrong decoding %s - %v", tCase.path, err)
		}
		if decoded.BPM != tCase.BPM {
			t.Logf("decoded:\n%#v\n", decoded.BPM)
			t.Logf("expected:\n%#v\n", tCase.BPM)
			t.Fatalf("%s wasn't decoded as expect.\nGot:\n%f\nExpected:\n%f",
				tCase.path, decoded.BPM, tCase.BPM)
		}
	}
}
