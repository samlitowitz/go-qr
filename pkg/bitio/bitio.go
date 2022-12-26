package bitio

// Reader is the interface that wraps the basic Read method.
//
// Read reads up to n bits into p. It returns the number of bytes
// read (0 <= m <= n) and any error encountered. Even if Read
// returns m < n, it may use all of p as scratch space during the call.
// If some data is available but not n bits, Read conventionally
// returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after
// successfully reading m > 0 bytes, it returns the number of
// bits read. It may return the (non-nil) error from the same call
// or return the error (and m == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the m > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when m == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
type Reader interface {
	Read(p []byte, n int) (m int, err error)
}

// Writer is the interface that wraps the basic Write method.
//
// Write writes n bits from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= m <= n)
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns m < n.
// Write must not modify the slice data, even temporarily.
//
// Implementations must not retain p.

type Writer interface {
	Write(p []byte, n int) (m int, err error)
}
