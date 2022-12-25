package numeric_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/samlitowitz/go-qr/mode/numeric"
)

func TestEncoder_Encode_TooLargeInput(t *testing.T) {
	encoder := numeric.NewEncoder(&numeric.Config{}, &bytes.Buffer{})
	err := encoder.Encode(make([]byte, 2<<14))

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	encodingError := &numeric.EncodingError{}
	if !errors.As(err, &encodingError) {
		t.Fatalf("Expected EncodingError")
	}
	outOfOBoundsError := &numeric.OutOfBoundsError{}
	err = errors.Unwrap(err)
	if !errors.As(err, &outOfOBoundsError) {
		t.Fatalf("Expected OutOfBoundsError")
	}
}

func TestEncoder_Encode(t *testing.T) {
	testCases := []struct {
		cfg      *numeric.Config
		v        []byte
		expected []byte
	}{
		{
			cfg: &numeric.Config{
				ModeIndicatorLength:  4,
				ModeIndicator:        1,
				CharacterCountLength: 10,
			},
			v:        []byte("01234567"),
			expected: []byte{0x10, 0x20, 0x0c, 0x56, 0x61, 0x80},
		},
		{
			cfg: &numeric.Config{
				ModeIndicatorLength:  2,
				ModeIndicator:        0,
				CharacterCountLength: 5,
			},
			v:        []byte("0123456789012345"),
			expected: []byte{0x20, 0x06, 0x2b, 0x35, 0x37, 0x0a, 0x28},
		},
	}

	for _, testCase := range testCases {
		buf := &bytes.Buffer{}
		encoder := numeric.NewEncoder(testCase.cfg, buf)
		err := encoder.Encode(testCase.v)
		if err != nil {
			t.Fatalf("Encode failed: %s", err)
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
