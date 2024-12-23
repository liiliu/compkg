package img

import (
	"errors"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

type ImageFile struct {
	File  *os.File
	Image image.Image
	Path  string
	Type  string
}

// GetImgFile 读取图片文件
func GetImgFile(filePath string) (ImageFile, error) {
	var imf ImageFile
	fType, err := os.Open(filePath)
	defer fType.Close()
	fileTypeData := make([]byte, 2)
	_, _ = fType.Read(fileTypeData)

	if fileTypeData[0] == 0xff && fileTypeData[1] == 0xd8 {
		imf.Type = "jpg"
	} else if fileTypeData[0] == 0x89 && fileTypeData[1] == 0x50 {
		imf.Type = "png"
	} else if fileTypeData[0] == 0x47 && fileTypeData[1] == 0x49 {
		imf.Type = "gif"
	} else if fileTypeData[0] == 0x42 && fileTypeData[1] == 0x4d {
		imf.Type = "bmp"
	} else {
		return imf, errors.New("no support image type")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return imf, err
	}
	imf.File = file
	imf.Path = filePath
	//fileType := strings.ToLower(path.Ext(filePath))
	//imf.Type = fileType[1:]

	var img image.Image
	switch imf.Type {
	case "jpeg", "jpg", "jfif":
		img, err = jpeg.Decode(file)
	case "png":
		img, err = png.Decode(file)
	case "gif":
		img, err = gif.Decode(file)
	case "bmp":
		img, err = bmp.Decode(file)
	default:
		return imf, errors.New("no support image type")
	}
	imf.Image = img
	return imf, err
}

func GetDirSp() string {
	return string(os.PathSeparator)
}
