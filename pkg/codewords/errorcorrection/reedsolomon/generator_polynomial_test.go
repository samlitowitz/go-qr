package reedsolomon_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/go-qr/pkg/codewords/errorcorrection/reedsolomon"
)

func TestGenerateGeneratorPolynomial(t *testing.T) {
	testCases := map[string]struct {
		n        int
		expected []byte
	}{
		"7 Error Correction Codewords": {
			n:        7,
			expected: []byte{21, 102, 238, 149, 146, 229, 87, 0},
		},
		"68 Error Correction Codewords": {
			n:        68,
			expected: []byte{238, 163, 8, 5, 3, 127, 184, 101, 27, 235, 238, 43, 198, 175, 215, 82, 32, 54, 2, 118, 225, 166, 241, 137, 125, 41, 177, 52, 231, 95, 97, 199, 52, 227, 89, 160, 173, 253, 84, 15, 84, 93, 151, 203, 220, 165, 202, 60, 52, 133, 205, 190, 101, 84, 150, 43, 254, 32, 160, 90, 70, 77, 93, 224, 33, 223, 159, 247, 0},
		},
	}

	for testDesc, testCase := range testCases {
		actual := reedsolomon.GenerateGeneratorPolynomial(testCase.n)
		if !cmp.Equal(testCase.expected, actual) {
			t.Fatalf(
				"%s: Invalid generator polynomial:\n%s",
				testDesc,
				cmp.Diff(testCase.expected, actual),
			)
		}
	}
}
