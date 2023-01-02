package byte_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/go-qr/pkg/symbol/mode"
	bytem "github.com/samlitowitz/go-qr/pkg/symbol/mode/byte"
)

func TestEncoder_Encode(t *testing.T) {
	testCases := map[string]struct {
		cfg          *mode.Config
		v            []byte
		expected     []byte
		bitsInStream int
	}{
		"M3 - 1": {
			cfg: &mode.Config{
				ModeIndicatorLength:  2,
				ModeIndicator:        2,
				CharacterCountLength: 4,
			},
			v:            []byte{0xd1},
			expected:     []byte{0x87, 0x44},
			bitsInStream: 14,
		},
		"M3 - 2": {
			cfg: &mode.Config{
				ModeIndicatorLength:  2,
				ModeIndicator:        2,
				CharacterCountLength: 4,
			},
			v:            []byte{0xd1, 0xc1},
			expected:     []byte{0x8B, 0x47, 0x04},
			bitsInStream: 22,
		},
		"M4": {
			cfg: &mode.Config{
				ModeIndicatorLength:  3,
				ModeIndicator:        2,
				CharacterCountLength: 5,
			},
			v:            []byte("0123456789012345"),
			expected:     append([]byte{0x50}, []byte("0123456789012345")...),
			bitsInStream: 136,
		},
		"1": {
			cfg: &mode.Config{
				ModeIndicatorLength:  4,
				ModeIndicator:        4,
				CharacterCountLength: 8,
			},
			v:            []byte("01234"),
			expected:     []byte{0x40, 0x53, 0x03, 0x13, 0x23, 0x33, 0x40},
			bitsInStream: 52,
		},
	}

	for testDesc, testCase := range testCases {
		buf := &bytes.Buffer{}
		encoder := bytem.NewEncoder(testCase.cfg, buf)
		bitsInStream, err := encoder.Encode(testCase.v)
		if err != nil {
			t.Fatalf("%s: Encode failed: %s", testDesc, err)
		}
		if bitsInStream != testCase.bitsInStream {
			t.Fatalf(
				"%s: Encode Failed: expected %d bits in stream, got %d",
				testDesc,
				testCase.bitsInStream,
				bitsInStream,
			)
		}
		actual := buf.Bytes()
		if !cmp.Equal(testCase.expected, actual) {
			t.Fatalf(
				"%s: Encode failed:\n%s",
				testDesc,
				cmp.Diff(testCase.expected, actual),
			)
		}
	}
}
