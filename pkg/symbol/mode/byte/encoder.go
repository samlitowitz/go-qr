package byte

import (
	"fmt"
	"io"

	"github.com/samlitowitz/go-qr/pkg/bits"

	"github.com/samlitowitz/go-qr/pkg/symbol/mode"
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
	charCount := len(v)
	maxDataLen := 2<<enc.cfg.CharacterCountLength - 1
	if charCount > maxDataLen {
		return 0, &mode.EncodingError{
			Err: &mode.OutOfBoundsError{
				Given:  fmt.Sprintf("%d", charCount),
				Bounds: fmt.Sprintf("[1, %d]", maxDataLen),
			},
		}
	}
	bitsInStream := enc.cfg.ModeIndicatorLength + enc.cfg.CharacterCountLength + bits.BitsPerByte*charCount
	buf := bits.NewBuffer(make([]byte, 0, bitsInStream/bits.BitsPerByte+1))
	// Mode
	_, err := buf.Write([]byte{enc.cfg.ModeIndicator << (bits.BitsPerByte - enc.cfg.ModeIndicatorLength)}, enc.cfg.ModeIndicatorLength)
	if err != nil {
		return 0, &mode.EncodingError{
			Pos: 0,
			Err: err,
		}
	}

	// Character Count
	_, err = buf.Write([]byte{byte(charCount) << (bits.BitsPerByte - enc.cfg.CharacterCountLength)}, enc.cfg.CharacterCountLength)
	if err != nil {
		return 0, &mode.EncodingError{
			Pos: 0,
			Err: err,
		}
	}

	// Data
	_, err = buf.Write(v, bits.BitsPerByte*charCount)
	if err != nil {
		return 0, &mode.EncodingError{
			Pos: 0,
			Err: err,
		}
	}

	return buf.Len(), nil
}
