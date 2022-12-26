package alphanumeric

import (
	"fmt"
	"io"

	"github.com/samlitowitz/go-qr/pkg"

	"github.com/samlitowitz/go-qr/pkg/mode"
)

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
		mapped, ok := encodeMap[v[i]]
		if !ok {
			return 0, &mode.EncodingError{
				Err: &mode.OutOfBoundsError{
					Given:  fmt.Sprintf("%x", v[i]),
					Bounds: "([0-9A-Z$%*+-./:\\s])",
				},
			}
		}
		v[i] = mapped
	}

	charCount := len(v)
	bitsInStream := enc.cfg.ModeIndicatorLength + enc.cfg.CharacterCountLength + 11*(charCount/2) + 6*(charCount%2)

	buf := make([]byte, bufLen)
	var byteInBuf int
	unusedBitsInByte := pkg.BitsPerByte
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
		numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, err = pkg.PackInt(
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
		numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, err = pkg.PackInt(
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
	// -- Process in groups of two
	var groupVal int
	for k := 0; k < charCount; k += 2 {
		switch true {
		case charCount-k >= 2:
			groupVal = 45*int(v[k]) + int(v[k+1])
			numberOfBitsToPack = 11
		case charCount-k >= 1:
			groupVal = int(v[k])
			numberOfBitsToPack = 6
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
			numberOfBitsUnpacked, unusedBitsInByte, byteInBuf, err = pkg.PackInt(
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
		unusedBitsInByte = pkg.BitsPerByte
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
