package reedsolomon

type ErrorCorrection struct {
	Blocks int
	N      int // N is the total number of symbols (RS) in the codeword (RS)
	K      int // K is the number of data symbols (RS) in the codeword (RS)
}
