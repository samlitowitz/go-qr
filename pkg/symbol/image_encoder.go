package symbol

import "image"

type ImageEncoder struct {
	*Encoder
}

func NewImageEncoder() *ImageEncoder {
	return &ImageEncoder{
		Encoder: NewEncoder(),
	}
}

func (enc *ImageEncoder) Encode() (image.Image, error) {
	return nil, nil
}
