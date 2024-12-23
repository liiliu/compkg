package img

import (
	"errors"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"image/gif"
	"image/jpeg"
	"image/png"

	"os"
)

// Scale 压缩图片
func (imf *ImageFile) Scale(outPath string, outWith uint, outHeight uint) error {
	outFile, err := os.Create(outPath)
	defer outFile.Close()
	if err != nil {
		return err
	}
	canvas := resize.Resize(outWith, outHeight, imf.Image, resize.Lanczos3)
	switch imf.Type {
	case "jpeg", "jpg", "jfif":
		err = jpeg.Encode(outFile, canvas, nil)
	case "png":
		err = png.Encode(outFile, canvas)
	case "gif":
		err = gif.Encode(outFile, canvas, nil)
	case "bmp":
		err = bmp.Encode(outFile, canvas)
	default:
		return errors.New("no support image scale type")
	}
	return err
}
