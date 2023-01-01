package bits

import (
	"errors"
	"io"
)

// From Go standard library definition of [bytes.Buffer](https://cs.opensource.google/go/go/+/refs/tags/go1.19.4:src/bytes/buffer.go).
// Most of this is copied directly or adapted with minor tweaks to the original to support bit operations

// smallBufferSize is an initial allocation minimal capacity.
const smallBufferSize = 64

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrTooLarge = errors.New("bytes.Buffer: too large")

const maxInt = int(^uint(0) >> 1)

// A Buffer is a variable-sized buffer of bits with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type Buffer struct {
	buf         []byte // contents are the bytes buf[off : len(buf)]
	readOffByte int    // read at &buf[off], write at &buf[len(buf)]
	readOffBit  int    // read at &buf[off] & (1 << readOffBit) - 1

	writeOffBit int // write at &buf[len(buf)] | (((1<< readOffBit) - 1) & input[i])
}

// Bytes returns a slice of length len(b.buf) - b.readOffByte holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like Read, Write, Reset, or Truncate).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
func (b *Buffer) Bytes() []byte { return b.buf[b.readOffByte:] }

// String returns the contents of the unread portion of the buffer
// as a string. If the Buffer is a nil pointer, it returns "<nil>".
//
// To build strings more efficiently, see the strings.Builder type.
func (b *Buffer) String() string {
	if b == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(b.buf[b.readOffBit:])
}

// empty reports whether the unread portion of the buffer is empty.
func (b *Buffer) empty() bool {
	return (len(b.buf) < b.readOffByte) || (len(b.buf) == b.readOffByte && b.readOffBit == 0)
}

// Len returns the number of bits of the unread portion of the buffer.
func (b *Buffer) Len() int { return BitsPerByte*(len(b.buf)-b.readOffByte) + b.writeOffBit }

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
	b.readOffByte = 0
	b.readOffBit = 0
	b.writeOffBit = 0
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value m is the number of bits to write, n; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (b *Buffer) Write(p []byte, n int) (m int, err error) {
	offByte, ok := b.tryGrowByReslice(n)
	if !ok {
		offByte = b.grow(n)
	}
	// Bit offset is zero, just copy it over
	if b.writeOffBit == 0 && n%BitsPerByte == 0 {
		b.writeOffBit = n % BitsPerByte
		m = BitsPerByte * copy(b.buf[offByte:], p)

		if m > n {
			m = n
		}
		return m, nil
	}

	var inputOffByte, inputOffBit int
	var unwrittenBitsBufByte, unreadBitsInInputByte, bytesToWrite int
	var readMask, writeByte byte

	for m = 0; m < n; {
		inputOffByte = m / BitsPerByte
		inputOffBit = m % BitsPerByte

		// how many bits left in current buf byte
		unwrittenBitsBufByte = BitsPerByte - b.writeOffBit
		// how many bits left in current input byte
		unreadBitsInInputByte = n - m
		if BitsPerByte-inputOffBit < unreadBitsInInputByte {
			unreadBitsInInputByte = BitsPerByte - inputOffBit
		}

		bytesToWrite = unwrittenBitsBufByte
		if unreadBitsInInputByte < bytesToWrite {
			bytesToWrite = unreadBitsInInputByte
		}

		// Build read mask
		readMask = (1 << bytesToWrite) - 1
		// Position read mask
		readMask <<= 8 - bytesToWrite - inputOffBit

		// Apply read mask
		writeByte = p[inputOffByte] & readMask
		// Position at MSB
		writeByte <<= inputOffBit
		// Position at write offset
		writeByte >>= b.writeOffBit

		// Write bits
		b.buf[offByte] |= writeByte

		b.writeOffBit += bytesToWrite
		if b.writeOffBit >= 8 {
			offByte++
			b.writeOffBit = b.writeOffBit % 8
		}
		m += bytesToWrite
	}

	return m, nil
}

// Read reads the next n bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless m is zero);
// otherwise it is nil.
func (b *Buffer) Read(p []byte, n int) (m int, err error) {
	if b.empty() {
		// Buffer is empty, reset to recover space.
		b.Reset()
		if len(p) == 0 || n == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	nBytes := noLossBitsToBytes(n)
	if nBytes > len(p) {
		nBytes = len(p)
	}
	if b.readOffBit == 0 {
		m = copy(p, b.buf[b.readOffByte:nBytes])
		b.readOffByte += m
		m *= BitsPerByte
		if m > n {
			m = n
		}
		return m, nil
	}

	var outputOffByte int
	var unreadBitsBufByte, unwrittenBitsInOutputByte int
	var readMask, writeMask byte

	for m = 0; m < n; {
		// how many bits left in current buf byte
		unreadBitsBufByte = BitsPerByte - b.readOffBit
		// how many bits left in current output byte
		unwrittenBitsInOutputByte = m % BitsPerByte

		readMask = (1 << unreadBitsBufByte) - 1
		writeMask = (1 << unwrittenBitsInOutputByte) - 1

		switch true {
		case unreadBitsBufByte == unwrittenBitsInOutputByte:
			p[outputOffByte] |= writeMask & (readMask & b.buf[b.readOffByte])

			m += unreadBitsBufByte
			outputOffByte++
			b.readOffByte++
			b.readOffBit = 0

		case unreadBitsBufByte < unwrittenBitsInOutputByte:
			p[outputOffByte] |= writeMask & ((readMask & b.buf[b.readOffByte]) << (unreadBitsBufByte - unwrittenBitsInOutputByte))

			m += unreadBitsBufByte
			outputOffByte++
			b.readOffBit += unreadBitsBufByte

		case unreadBitsBufByte > unwrittenBitsInOutputByte:
			p[outputOffByte] |= writeMask & ((readMask & b.buf[b.readOffByte]) >> (unwrittenBitsInOutputByte - unreadBitsBufByte))

			m += unreadBitsBufByte - unwrittenBitsInOutputByte
			b.readOffByte++
			b.readOffBit = 0
		}
	}
	return m, nil
}

// tryGrowByReslice is an inlineable version of grow for the fast-case where the
// internal buffer only needs to be re-sliced.
// It returns the index where bytes should be written and whether it succeeded.
func (b *Buffer) tryGrowByReslice(n int) (int, bool) {
	nBytes := noLossBitsToBytes(n)
	if l := len(b.buf); nBytes <= cap(b.buf)-l {
		b.buf = b.buf[:l+nBytes]
		return l, true
	}
	return 0, false
}

// grow grows the buffer to guarantee space for n more bits.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer) grow(n int) int {
	nBytes := noLossBitsToBytes(n)
	m := b.Len()
	mBytes := noLossBitsToBytes(m)
	// If buffer is empty, reset to recover space.
	if m == 0 && b.readOffByte != 0 {
		b.Reset()
	}
	// Try to grow by means of a reslice.
	if i, ok := b.tryGrowByReslice(n); ok {
		return i
	}
	if b.buf == nil && nBytes <= smallBufferSize {
		b.buf = make([]byte, nBytes, smallBufferSize)
		return 0
	}
	c := cap(b.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(b.buf, b.buf[b.readOffByte:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Add b.off to account for b.buf[:b.off] being sliced off the front.
		b.buf = growSlice(b.buf[b.readOffByte:], b.readOffByte+nBytes)
	}
	// Restore b.off and len(b.buf).
	b.readOffByte = 0
	b.buf = b.buf[:mBytes+nBytes]
	return mBytes
}

// growSlice grows b by n, preserving the original content of b.
// If the allocation fails, it panics with ErrTooLarge.
func growSlice(b []byte, n int) []byte {
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	// TODO(http://golang.org/issue/51462): We should rely on the append-make
	// pattern so that the compiler can call runtime.growslice. For example:
	//	return append(b, make([]byte, n)...)
	// This avoids unnecessary zero-ing of the first len(b) bytes of the
	// allocated slice, but this pattern causes b to escape onto the heap.
	//
	// Instead use the append-make pattern with a nil slice to ensure that
	// we allocate buffers rounded up to the closest size class.
	c := len(b) + n // ensure enough space for n elements
	if c < 2*cap(b) {
		// The growth rate has historically always been 2x. In the future,
		// we could rely purely on append to determine the growth rate.
		c = 2 * cap(b)
	}
	b2 := append([]byte(nil), make([]byte, c)...)
	copy(b2, b)
	return b2[:len(b)]
}

func noLossBitsToBytes(n int) int {
	m := n / BitsPerByte
	if n%BitsPerByte > 0 {
		m++
	}
	return m
}

// NewBuffer creates and initializes a new Buffer using buf as its
// initial contents. The new Buffer takes ownership of buf, and the
// caller should not use buf after this call. NewBuffer is intended to
// prepare a Buffer to read existing data. It can also be used to set
// the initial size of the internal buffer for writing. To do that,
// buf should have the desired capacity but a length of zero.
//
// In most cases, new(Buffer) (or just declaring a Buffer variable) is
// sufficient to initialize a Buffer.
func NewBuffer(buf []byte) *Buffer { return &Buffer{buf: buf} }
