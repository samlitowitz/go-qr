package numeric

import (
	"fmt"
	"io"
	"math"
)

const bitsPerByte = 8
const bitsPerGroup = 10
const bufLen = 1024

type EncodingError struct {
	Pos int
	Err error
}

func (err *EncodingError) Error() string {
	return fmt.Sprintf(
		"Error at position %d: %s",
		err.Pos,
		err.Err,
	)
}

func (err *EncodingError) Unwrap() error {
	return err.Err
}

type OutOfBoundsError struct {
	given, bounds string
}

func (err *OutOfBoundsError) Error() string {
	return fmt.Sprintf(
		"Out of bounds: `%s` given, expected in %s",
		err.given,
		err.bounds,
	)
}

type Config struct {
	ModeIndicatorLength  int // The length of the mode indicator in bits
	ModeIndicator        byte
	CharacterCountLength int // The length of the character count in bits
}

type Encoder struct {
	w   io.Writer
	cfg *Config
}

func NewEncoder(cfg *Config, w io.Writer) *Encoder {
	return &Encoder{cfg: cfg, w: w}
}

func (enc *Encoder) Encode(v []byte) error {
	if len(v) > 2<<14-1 {
		return &EncodingError{
			Err: &OutOfBoundsError{
				given:  fmt.Sprintf("%d", len(v)),
				bounds: fmt.Sprintf("[1, %d]", 2<<14-1),
			},
		}
	}

	for i := 0; i < len(v); i++ {
		if v[i] < 0x30 || v[i] > 0x39 {
			return &EncodingError{
				Err: &OutOfBoundsError{
					given:  fmt.Sprintf("%x", v[i]),
					bounds: "[0x30, 0x39]",
				},
			}
		}
		v[i] = v[i] - 0x30
	}

	charCount := len(v)
	buf := make([]byte, bufLen)
	var byteInBuf int
	unusedBitsInByte := bitsPerByte
	var numberOfBitsToPack, numberOfBitsUnpacked int
	var err error

	// Mode
	numberOfBitsToPack = enc.cfg.ModeIndicatorLength
	if numberOfBitsToPack/8+1+byteInBuf >= bufLen {
		_, err = enc.w.Write(buf[:byteInBuf])
		if err != nil {
			return &EncodingError{
				Pos: byteInBuf,
				Err: err,
			}
		}
	}
	for numberOfBitsUnpacked = numberOfBitsToPack; numberOfBitsUnpacked > 0; {
		numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, err = packInt(
			int(enc.cfg.ModeIndicator),
			numberOfBitsToPack,
			unusedBitsInByte,
			byteInBuf,
			&buf,
		)
		if err != nil {
			return &EncodingError{
				Pos: byteInBuf,
				Err: err,
			}
		}
	}

	// Char Count
	numberOfBitsToPack = enc.cfg.CharacterCountLength
	if numberOfBitsToPack/8+1+byteInBuf >= bufLen {
		_, err = enc.w.Write(buf[:byteInBuf])
		if err != nil {
			return &EncodingError{
				Pos: byteInBuf,
				Err: err,
			}
		}
	}
	for numberOfBitsUnpacked = numberOfBitsToPack; numberOfBitsUnpacked > 0; {
		numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, err = packInt(
			charCount,
			numberOfBitsToPack,
			unusedBitsInByte,
			byteInBuf,
			&buf,
		)
		if err != nil {
			return &EncodingError{
				Pos: byteInBuf,
				Err: err,
			}
		}
	}

	// Data
	// -- Process in groups of three digits
	var groupVal int
	for k := 0; k < charCount; k += 3 {
		switch true {
		case charCount-k >= 3:
			groupVal = 100*int(v[k]) + 10*int(v[k+1]) + int(v[k+2])
			numberOfBitsToPack = bitsPerGroup
		case charCount-k >= 2:
			groupVal = 10*int(v[k]) + int(v[k+1])
			numberOfBitsToPack = mostSignificantBit(groupVal)
		case charCount-k >= 1:
			groupVal = int(v[k])
			numberOfBitsToPack = mostSignificantBit(groupVal)
		}
		if numberOfBitsToPack/8+1+byteInBuf >= bufLen {
			_, err = enc.w.Write(buf[:byteInBuf])
			if err != nil {
				return &EncodingError{
					Pos: byteInBuf,
					Err: err,
				}
			}
		}
		for numberOfBitsUnpacked = numberOfBitsToPack; numberOfBitsUnpacked > 0; {
			numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, err = packInt(
				int(enc.cfg.ModeIndicator),
				numberOfBitsToPack,
				unusedBitsInByte,
				byteInBuf,
				&buf,
			)
			if err != nil {
				return &EncodingError{
					Pos: byteInBuf,
					Err: err,
				}
			}
		}
	}

	if byteInBuf > 0 {
		_, err = enc.w.Write(buf[:byteInBuf])
		if err != nil {
			return &EncodingError{
				Pos: byteInBuf,
				Err: err,
			}
		}
	}

	return nil
}

func packInt(v, numberOfBitsToPack, unusedBitsInByte, byteInBuf int, buf *[]byte) (int, int, int, error) {
	var toCopy int
	toCopy = v
	numberOfBitsUnpacked := numberOfBitsToPack
	for numberOfBitsPacked := 0; numberOfBitsPacked < numberOfBitsToPack; {
		switch true {
		case numberOfBitsUnpacked == unusedBitsInByte:
			// copy
			(*buf)[byteInBuf] |= byte(toCopy)

			// bookkeeping
			byteInBuf++
			unusedBitsInByte = bitsPerByte
			numberOfBitsPacked += numberOfBitsUnpacked
			numberOfBitsUnpacked = 0

		case numberOfBitsUnpacked < unusedBitsInByte:
			// copy
			(*buf)[byteInBuf] |= byte(toCopy) << (unusedBitsInByte - numberOfBitsUnpacked)

			// bookkeeping
			unusedBitsInByte -= numberOfBitsUnpacked
			numberOfBitsPacked += numberOfBitsUnpacked
			numberOfBitsUnpacked = 0

		case numberOfBitsUnpacked > unusedBitsInByte:
			// copy
			(*buf)[byteInBuf] |= byte(toCopy >> (numberOfBitsUnpacked - unusedBitsInByte))

			// bookkeeping
			numberOfBitsPacked += unusedBitsInByte
			numberOfBitsUnpacked -= unusedBitsInByte
			byteInBuf++
			unusedBitsInByte = bitsPerByte
		}
	}
	return numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, nil
}

func mostSignificantBit(v int) int {
	return int(math.Log2(float64(v)))
}
