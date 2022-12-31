package bits_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/go-qr/pkg"

	"github.com/samlitowitz/go-qr/pkg/bits"
)

func TestBuffer_Write(t *testing.T) {
	testCases := map[string]struct {
		input    []byte
		n        int
		expected []byte
	}{
		"write n bits where n = 8": {
			input:    []byte{0x01},
			n:        8,
			expected: []byte{0x01},
		},
		"write n bits where n % 8 = 0, n > 8": {
			input:    []byte{0x01, 0x02},
			n:        16,
			expected: []byte{0x01, 0x02},
		},
		"write n bits where n < 8": {
			input:    []byte{0xf5},
			n:        6,
			expected: []byte{0xf4},
		},
		"write n bits where n % 8 != 0, n > 8": {
			input:    []byte{0x01, 0xd0},
			n:        10,
			expected: []byte{0x01, 0xc0},
		},
	}

	for testDesc, testCase := range testCases {
		buf := &bits.Buffer{}
		m, err := buf.Write(testCase.input, testCase.n)
		if err != nil {
			t.Fatalf(
				"%s: Write failed: %s",
				testDesc,
				err,
			)
		}
		if m != testCase.n {
			t.Fatalf(
				"%s: Write failed: expected %d bits written, got %d",
				testDesc,
				testCase.n,
				m,
			)
		}
		l := testCase.n / pkg.BitsPerByte
		if testCase.n%pkg.BitsPerByte > 0 {
			l++
		}
		actual := make([]byte, l)
		m, err = buf.Read(actual, testCase.n)
		if err != nil {
			t.Fatalf(
				"%s: Read failed: %s",
				testDesc,
				err,
			)
		}
		if m != testCase.n {
			t.Fatalf(
				"%s: Read failed: expected %d bits read, got %d",
				testDesc,
				testCase.n,
				m,
			)
		}
		if !cmp.Equal(testCase.expected, actual) {
			t.Fatalf(
				"%s: Read failed:\n%s",
				testDesc,
				cmp.Diff(testCase.input, actual),
			)
		}
	}
}