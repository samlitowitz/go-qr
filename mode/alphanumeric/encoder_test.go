package alphanumeric_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/samlitowitz/go-qr/mode"
	"github.com/samlitowitz/go-qr/mode/alphanumeric"
)

func TestEncoder_Encode_TooLargeInput(t *testing.T) {
	encoder := alphanumeric.NewEncoder(&mode.Config{}, &bytes.Buffer{})
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
				ModeIndicator:        2,
				CharacterCountLength: 9,
			},
			v:            []byte("AC-42"),
			expected:     []byte{0x20, 0x29, 0xce, 0xe7, 0x21, 0x00},
			bitsInStream: 41,
		},
	}

	for _, testCase := range testCases {
		buf := &bytes.Buffer{}
		encoder := alphanumeric.NewEncoder(testCase.cfg, buf)
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
