package betterapng

import (
	"bytes"
	"image"
	"image/png"
	"image/jpeg"
)

func EncodeJPEG(buf *bytes.Buffer, img image.Image, o *jpeg.Options) error {
	return jpeg.Encode(buf, img, o)
}

func EncodePNG(buf *bytes.Buffer, img image.Image) error {
	return png.Encode(buf, img)
}

func UniversalDecoder(bytes2 []byte, codec string) (image.Image, error) {
	switch codec {
	case "PNG":
		res := new(bytes.Buffer)
		res.Write(bytes2)
		return png.Decode(res)
	case "JPEG":
		res := new(bytes.Buffer)
		res.Write(bytes2)
		return jpeg.Decode(res)
	default:
		return nil, nil
	}
}