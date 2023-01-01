package codewords

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/go-qr/pkg/codewords/errorcorrection/reedsolomon"
)

func TestInterleave(t *testing.T) {
	testCases := map[string]struct {
		cfg           *reedsolomon.Config
		dataCodewords []byte
		ecCodewords   []byte
		expected      []byte
	}{
		"5-Q": {
			cfg: &reedsolomon.Config{
				ErrorCorrectionCodewordCount: 72,
				ErrorCorrections: []*reedsolomon.ErrorCorrection{
					{
						Blocks: 2,
						N:      33,
						K:      15,
					},
					{
						Blocks: 2,
						N:      34,
						K:      16,
					},
				},
				Gx: reedsolomon.GenerateGeneratorPolynomial(52),
			},
			dataCodewords: []byte{
				// Block 1
				67, 85, 70, 134, 87, 38, 85, 194, 119, 50, 6, 18, 6, 103, 38,
				// Block 2
				246, 246, 66, 7, 118, 134, 242, 7, 38, 86, 22, 198, 199, 146, 6,
				// Block 3
				182, 230, 247, 119, 50, 7, 118, 134, 87, 38, 82, 6, 134, 151, 50, 7,
				// Block 4
				70, 247, 118, 86, 194, 6, 151, 50, 16, 236, 17, 236, 17, 236, 17, 236,
			},
			ecCodewords: []byte{
				// Block 1
				213, 199, 11, 45, 115, 247, 241, 223, 229, 248, 154, 117, 154, 111, 86, 161, 111, 39,
				// Block 2
				87, 204, 96, 60, 202, 182, 124, 157, 200, 134, 27, 129, 209, 17, 163, 163, 120, 133,
				// Block 3
				148, 116, 177, 212, 76, 133, 75, 242, 238, 76, 195, 230, 189, 10, 108, 240, 192, 141,
				// Block 4
				235, 159, 5, 173, 24, 147, 59, 33, 106, 40, 255, 172, 82, 2, 131, 32, 178, 236,
			},
			expected: []byte{
				// Data Codewords
				67, 246, 182, 70,
				85, 246, 230, 247,
				70, 66, 247, 118,
				134, 7, 119, 86,
				87, 118, 50, 194,
				38, 134, 7, 6,
				85, 242, 118, 151,
				194, 7, 134, 50,
				119, 38, 87, 16,
				50, 86, 38, 236,
				6, 22, 82, 17,
				18, 198, 6, 236,
				6, 199, 134, 17,
				103, 146, 151, 236,
				38, 6, 50, 17,
				7, 236,
				// EC Codewords
				213, 87, 148, 235,
				199, 204, 116, 159,
				11, 96, 177, 5,
				45, 60, 212, 173,
				115, 202, 76, 24,
				247, 182, 133, 147,
				241, 124, 75, 59,
				223, 157, 242, 33,
				229, 200, 238, 106,
				248, 134, 76, 40,
				154, 27, 195, 255,
				117, 129, 230, 172,
				154, 209, 189, 82,
				111, 17, 10, 2,
				86, 163, 108, 131,
				161, 163, 240, 32,
				111, 120, 192, 178,
				39, 133, 141, 236,
			},
		},
	}

	for testDesc, testCase := range testCases {
		actual, err := Interleave(testCase.cfg, testCase.dataCodewords, testCase.ecCodewords)
		if err != nil {
			t.Fatalf(
				"%s: Interleave failed: %s",
				testDesc,
				err,
			)
		}
		if !cmp.Equal(testCase.expected, actual) {
			t.Fatalf(
				"%s: Interleave failed:\n%s",
				testDesc,
				cmp.Diff(testCase.expected, actual),
			)
		}
	}
}
