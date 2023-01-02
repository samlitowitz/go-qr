package symbol

import (
	"github.com/samlitowitz/go-qr/pkg/bits"
	"github.com/samlitowitz/go-qr/pkg/symbol/errorcorrection"
	"github.com/samlitowitz/go-qr/pkg/symbol/mode"
)

type unencodedData struct {
	encMode mode.Mode
	data    []byte
}

type Encoder struct {
	unencodedData []*unencodedData
}

func NewEncoder() *Encoder {
	return &Encoder{
		unencodedData: []*unencodedData{},
	}
}

func (enc *Encoder) Write(data []byte, encMode mode.Mode) (int, error) {
	enc.unencodedData = append(enc.unencodedData, &unencodedData{encMode: encMode, data: data})

	return len(data), nil
}

func (enc *Encoder) Encode(ecLevel errorcorrection.Level) (*bits.Buffer, error) {
	return nil, nil
}
