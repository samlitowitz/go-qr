package reedsolomon_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/go-qr/pkg/errorcorrection/reedsolomon"
)

func TestGeneratorPolynomial(t *testing.T) {
	testCases := map[string]struct {
		n        int
		expected []byte
	}{
		"68 Error Correction Codewords": {
			n:        68,
			expected: []byte{0, 68, 247, 67, 159, 66, 223, 65, 33, 64, 224, 63, 93, 62, 77, 61, 70, 60, 90, 59, 160, 58, 32, 57, 254, 56, 43, 55, 150, 54, 84, 53, 101, 52, 190, 51, 205, 50, 133, 49, 52, 48, 60, 47, 202, 46, 165, 45, 220, 44, 203, 43, 151, 42, 93, 41, 84, 40, 15, 39, 84, 38, 253, 37, 173, 36, 160, 35, 89, 34, 227, 33, 52, 32, 199, 31, 97, 30, 95, 29, 231, 28, 52, 27, 177, 26, 41, 25, 125, 24, 137, 23, 241, 22, 166, 21, 225, 20, 118, 19, 2, 18, 54, 17, 32, 16, 82, 15, 215, 14, 175, 13, 198, 12, 43, 11, 238, 10, 235, 9, 27, 8, 101, 7, 184, 6, 127, 5, 3, 4, 5, 3, 8, 2, 163, 238},
		},
	}

	for testDesc, testCase := range testCases {
		actual := reedsolomon.GeneratorPolynomial(testCase.n)
		if !cmp.Equal(testCase.expected, actual) {
			t.Fatalf(
				"%s: Invalid generator polynomial:\n%s",
				testDesc,
				cmp.Diff(testCase.expected, actual),
			)
		}
	}
}
