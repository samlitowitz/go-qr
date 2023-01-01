package numeric

import (
	"fmt"
	"io"

	"github.com/samlitowitz/go-qr/pkg/bits"

	"github.com/samlitowitz/go-qr/pkg/symbol/mode"
)

const bitsPerGroup = 10
const bufLen = 1024

type Encoder struct {
	w   io.Writer
	cfg *mode.Config
}

func NewEncoder(cfg *mode.Config, w io.Writer) *Encoder {
	return &Encoder{cfg: cfg, w: w}
}

func (enc *Encoder) Encode(v []byte) (int, error) {
	maxDataLen := 2<<enc.cfg.CharacterCountLength - 1
	if len(v) > maxDataLen {
		return 0, &mode.EncodingError{
			Err: &mode.OutOfBoundsError{
				Given:  fmt.Sprintf("%d", len(v)),
				Bounds: fmt.Sprintf("[1, %d]", maxDataLen),
			},
		}
	}

	for i := 0; i < len(v); i++ {
		if v[i] < 0x30 || v[i] > 0x39 {
			return 0, &mode.EncodingError{
				Err: &mode.OutOfBoundsError{
					Given:  fmt.Sprintf("%x", v[i]),
					Bounds: "[0x30, 0x39]",
				},
			}
		}
		v[i] = v[i] - 0x30
	}

	charCount := len(v)
	remainder := 0
	switch charCount % 3 {
	case 1:
		remainder = 4
	case 2:
		remainder = 7
	}
	bitsInStream := enc.cfg.ModeIndicatorLength + enc.cfg.CharacterCountLength + 10*(charCount/3) + remainder

	buf := make([]byte, bufLen)
	var byteInBuf int
	unusedBitsInByte := bits.BitsPerByte
	var numberOfBitsToPack, numberOfBitsUnpacked int
	var err error

	// Mode
	numberOfBitsToPack = enc.cfg.ModeIndicatorLength
	if numberOfBitsToPack/8+1+byteInBuf >= bufLen {
		_, err = enc.w.Write(buf[:byteInBuf])
		if err != nil {
			return 0, &mode.EncodingError{
				Pos: byteInBuf,
				Err: err,
			}
		}
	}
	for numberOfBitsUnpacked = numberOfBitsToPack; numberOfBitsUnpacked > 0; {
		numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, err = bits.PackInt(
			int(enc.cfg.ModeIndicator),
			numberOfBitsToPack,
			unusedBitsInByte,
			byteInBuf,
			&buf,
		)
		if err != nil {
			return 0, &mode.EncodingError{
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
			return 0, &mode.EncodingError{
				Pos: byteInBuf,
				Err: err,
			}
		}
	}
	for numberOfBitsUnpacked = numberOfBitsToPack; numberOfBitsUnpacked > 0; {
		numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, err = bits.PackInt(
			charCount,
			numberOfBitsToPack,
			unusedBitsInByte,
			byteInBuf,
			&buf,
		)
		if err != nil {
			return 0, &mode.EncodingError{
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
			numberOfBitsToPack = 7
		case charCount-k >= 1:
			groupVal = int(v[k])
			numberOfBitsToPack = 4
		}
		if numberOfBitsToPack/8+1+byteInBuf >= bufLen {
			_, err = enc.w.Write(buf[:byteInBuf])
			if err != nil {
				return 0, &mode.EncodingError{
					Pos: byteInBuf,
					Err: err,
				}
			}
		}
		for numberOfBitsUnpacked = numberOfBitsToPack; numberOfBitsUnpacked > 0; {
			numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, err = bits.PackInt(
				groupVal,
				numberOfBitsToPack,
				unusedBitsInByte,
				byteInBuf,
				&buf,
			)
			if err != nil {
				return 0, &mode.EncodingError{
					Pos: byteInBuf,
					Err: err,
				}
			}
		}
	}

	if unusedBitsInByte > 0 {
		byteInBuf++
		unusedBitsInByte = bits.BitsPerByte
	}

	if byteInBuf > 0 {
		_, err = enc.w.Write(buf[:byteInBuf])
		if err != nil {
			return 0, &mode.EncodingError{
				Pos: byteInBuf,
				Err: err,
			}
		}
	}

	return bitsInStream, nil
}
