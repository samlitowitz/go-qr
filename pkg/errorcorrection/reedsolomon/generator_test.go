package reedsolomon_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/go-qr/pkg/errorcorrection/reedsolomon"
)

func TestGenerator_Generate(t *testing.T) {
	testCases := map[string]struct {
		cfg           *reedsolomon.Config
		dataCodewords []byte
		expected      []byte
	}{
		"1-M \"HELLO WORLD\"": {
			cfg: &reedsolomon.Config{
				ErrorCorrectionCodewordCount: 10,
				ErrorCorrections: []*reedsolomon.ErrorCorrection{
					{
						Blocks: 1,
						K:      16,
						N:      26,
					},
				},
				Gx: []byte{45, 32, 94, 64, 70, 118, 61, 46, 67, 251, 0},
			},
			dataCodewords: []byte{32, 91, 11, 120, 209, 114, 220, 77, 67, 64, 236, 17, 236, 17, 236, 17},
			expected:      []byte{196, 35, 39, 119, 235, 215, 231, 226, 93, 23},
		},
	}

	for testDesc, testCase := range testCases {
		gen := reedsolomon.NewGenerator(testCase.cfg)
		actual, err := gen.Generate(testCase.dataCodewords)
		if err != nil {
			t.Fatalf("%s: Generate failed: %s", testDesc, err)
		}
		if !cmp.Equal(testCase.expected, actual) {
			t.Fatalf(
				"%s: Generate failed:\n%s",
				testDesc,
				cmp.Diff(testCase.expected, actual),
			)
		}
	}
}
