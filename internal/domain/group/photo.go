package group

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

type GroupPhoto struct {
	Data   []byte `json:"data"`
	Mime   string `json:"mime"`
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
}

func NewGroupPhoto(data []byte) (*GroupPhoto, error) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if format != "jpeg" && format != "png" {
		return nil, ErrPhotoUnsupportedFormat
	}

	if cfg.Width < MinPhotoWidth || cfg.Height < MinPhotoHeight || cfg.Width > MaxPhotoWidth || cfg.Height > MaxPhotoHeight {
		return nil, ErrPhotoInvalidDimensions
	}

	return &GroupPhoto{
		Data:   data,
		Mime:   http.DetectContentType(data),
		Width:  uint32(cfg.Width),
		Height: uint32(cfg.Height),
	}, nil
}
