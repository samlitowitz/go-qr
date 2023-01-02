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
		"M3": {
			cfg: &mode.Config{
				ModeIndicatorLength:  2,
				ModeIndicator:        2,
				CharacterCountLength: 4,
			},
			v:            []byte{0xff, 0x01, 0x02, 0x03, 0x04, 0x05},
			expected:     []byte{0x98, 0x00, 0x04, 0x08, 0x0C, 0x10, 0x14},
			bitsInStream: 54,
		},
		//"M4": {
		//	cfg: &mode.Config{
		//		ModeIndicatorLength:  3,
		//		ModeIndicator:        2,
		//		CharacterCountLength: 5,
		//	},
		//	v:            []byte("0123456789012345"),
		//	expected:     []byte{0x20, 0x06, 0x2b, 0x35, 0x37, 0x0a, 0x75, 0x28},
		//	bitsInStream: 61,
		//},
		//"1": {
		//	cfg: &mode.Config{
		//		ModeIndicatorLength:  4,
		//		ModeIndicator:        4,
		//		CharacterCountLength: 8,
		//	},
		//	v:            []byte("01234567"),
		//	expected:     []byte{0x10, 0x20, 0x0c, 0x56, 0x61, 0x80},
		//	bitsInStream: 41,
		//},
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
