package bits_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/samlitowitz/go-qr/pkg/bits"
)

type bufNBits struct {
	n    int
	data []byte
}

type testCase struct {
	init   []byte
	writes []bufNBits
	reads  []bufNBits
}

func TestBuffer_Write(t *testing.T) {
	testCases := map[string]testCase{
		"write n bits where n = 8": {
			init: []byte{},
			writes: []bufNBits{
				{
					n:    8,
					data: []byte{0x01},
				},
			},
			reads: []bufNBits{
				{
					n:    8,
					data: []byte{0x01},
				},
			},
		},
		"write n bits where n % 8 = 0, n > 8": {
			init: []byte{},
			writes: []bufNBits{
				{
					n:    16,
					data: []byte{0x01, 0x02},
				},
			},
			reads: []bufNBits{
				{
					n:    16,
					data: []byte{0x01, 0x02},
				},
			},
		},
		"write n bits where n < 8": {
			init: []byte{},
			writes: []bufNBits{
				{
					n:    8,
					data: []byte{0xf5},
				},
			},
			reads: []bufNBits{
				{
					n:    6,
					data: []byte{0xf4},
				},
			},
		},
		"write n bits where n % 8 != 0, n > 8": {
			init: []byte{},
			writes: []bufNBits{
				{
					n:    16,
					data: []byte{0x01, 0xd0},
				},
			},
			reads: []bufNBits{
				{
					n:    10,
					data: []byte{0x01, 0xc0},
				},
			},
		},
		"write 8 bits with empty init": {
			init: []byte{},
			writes: []bufNBits{
				{
					n:    8,
					data: []byte{0x01},
				},
			},
			reads: []bufNBits{
				{
					n:    8,
					data: []byte{0x01},
				},
			},
		},
		"write 16 bits with 2 byte init": {
			init: make([]byte, 0, 2),
			writes: []bufNBits{
				{
					n:    32,
					data: []byte{0x01, 0x02, 0x03, 0x04},
				},
			},
			reads: []bufNBits{
				{
					n:    16,
					data: []byte{0x01, 0x02},
				},
			},
		},
		"write 2 in 2 bits and 6 in 4 bits": {
			init: []byte{},
			writes: []bufNBits{
				{
					n:    2,
					data: []byte{0x02 << 6},
				},
				{
					n:    4,
					data: []byte{0x06 << 4},
				},
			},
			reads: []bufNBits{
				{
					n:    6,
					data: []byte{0x98},
				},
			},
		},
	}

	for desc, tc := range testCases {
		buf := bits.NewBuffer(tc.init)

		for _, write := range tc.writes {
			m, err := buf.Write(write.data, write.n)
			if err != nil {
				t.Fatalf(
					"%s: Write failed: %s",
					desc,
					err,
				)
			}
			if !cmp.Equal(m, write.n) {
				t.Fatalf(
					"%s:\n%s",
					desc,
					cmp.Diff(m, write.n),
				)
			}
		}

		for _, read := range tc.reads {
			l := noLossBitsToBytes(read.n)
			actual := make([]byte, l)
			m, err := buf.Read(actual, read.n)
			if err != nil {
				t.Fatalf(
					"%s: Read failed: %s",
					desc,
					err,
				)
			}
			if !cmp.Equal(m, read.n) {
				t.Fatalf(
					"%s:\n%s",
					desc,
					cmp.Diff(m, read.n),
				)
			}
			if !cmp.Equal(read.data, actual) {
				t.Fatalf(
					"%s: Read failed:\n%s",
					desc,
					cmp.Diff(read.data, actual),
				)
			}
		}

	}
}

func noLossBitsToBytes(n int) int {
	m := n / bits.BitsPerByte
	if n%bits.BitsPerByte > 0 {
		m++
	}
	return m
}
