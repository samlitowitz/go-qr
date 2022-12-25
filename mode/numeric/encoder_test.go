package numeric_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/samlitowitz/go-qr/mode"
	"github.com/samlitowitz/go-qr/mode/numeric"
)

func TestEncoder_Encode_TooLargeInput(t *testing.T) {
	encoder := numeric.NewEncoder(&mode.Config{}, &bytes.Buffer{})
	_, err := encoder.Encode(make([]byte, 2<<14))

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	encodingError := &mode.EncodingError{}
	if !errors.As(err, &encodingError) {
		t.Fatalf("Expected EncodingError")
	}
	outOfOBoundsError := &mode.OutOfBoundsError{}
	err = errors.Unwrap(err)
	if !errors.As(err, &outOfOBoundsError) {
		t.Fatalf("Expected OutOfBoundsError")
	}
}

func TestEncoder_Encode(t *testing.T) {
	testCases := []struct {
		cfg          *mode.Config
		v            []byte
		expected     []byte
		bitsInStream int
	}{
		{
			cfg: &mode.Config{
				ModeIndicatorLength:  4,
				ModeIndicator:        1,
				CharacterCountLength: 10,
			},
			v:            []byte("01234567"),
			expected:     []byte{0x10, 0x20, 0x0c, 0x56, 0x61, 0x80},
			bitsInStream: 41,
		},
		{
			cfg: &mode.Config{
				ModeIndicatorLength:  2,
				ModeIndicator:        0,
				CharacterCountLength: 5,
			},
			v:            []byte("0123456789012345"),
			expected:     []byte{0x20, 0x06, 0x2b, 0x35, 0x37, 0x0a, 0x75, 0x28},
			bitsInStream: 61,
		},
	}

	for _, testCase := range testCases {
		buf := &bytes.Buffer{}
		encoder := numeric.NewEncoder(testCase.cfg, buf)
		bitsInStream, err := encoder.Encode(testCase.v)
		if err != nil {
			t.Fatalf("Encode failed: %s", err)
		}
		if bitsInStream != testCase.bitsInStream {
			t.Fatalf(
				"Expected %d bits in stream, got %d",
				testCase.bitsInStream,
				bitsInStream,
			)
		}
		if buf.String() != string(testCase.expected) {
			t.Fatalf(
				"Expected `%x`, got `%x`",
				testCase.expected,
				buf.String(),
			)
		}
	}
}
