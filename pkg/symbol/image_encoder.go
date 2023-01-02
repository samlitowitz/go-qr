package symbol

import (
	"image"

	"github.com/samlitowitz/go-qr/pkg/symbol/errorcorrection"
)

type ImageEncoder struct {
	*Encoder
}

func NewImageEncoder() *ImageEncoder {
	return &ImageEncoder{
		Encoder: NewEncoder(),
	}
}

func (enc *ImageEncoder) Encode(ecLevel errorcorrection.Level) (image.Image, error) {
	return nil, nil
}
