package reedsolomon

type Config struct {
	ErrorCorrectionCodewordCount int // Total number of error correction codewords
	ErrorCorrections             []*ErrorCorrection
	Gx                           []byte // Bit-wise representation of the generator polynomial
}
