package util

import (
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path"
)

type imageUtil struct{}

var ImageUtil imageUtil

func (i *imageUtil) CheckImageExt(filename string) (bool, string) {
	ext := path.Ext(filename)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".bmp" || ext == ".gif", ext[1:]
}

func (i *imageUtil) SaveAsJpg(imageFile multipart.File, ext string, filePath string) error {
	if err := CommonUtil.CheckCreateDir(filePath); err != nil {
		return err
	}

	jpgImgFile, err := os.Create(filePath)
	defer func() {
		jpgImgFile.Close()
	}()

	if err != nil {
		return err
	}
	var decodeImage image.Image
	switch ext {
	case "png":
		decodeImage, err = i.DecodePng(imageFile)
		break
	case "bmp":
		decodeImage, err = bmp.Decode(imageFile)
		break
	case "gif":
		decodeImage, err = gif.Decode(imageFile)
		break
	default:
		decodeImage, err = jpeg.Decode(imageFile)
		break
	}
	if err != nil {
		return err
	}

	err = jpeg.Encode(jpgImgFile, decodeImage, &jpeg.Options{Quality: 100})
	return err
}

func (i *imageUtil) DecodePng(pngImageFile multipart.File) (image.Image, error) {
	pngImage, err := png.Decode(pngImageFile)
	if err != nil {
		return nil, err
	}

	jpgImage := image.NewRGBA(pngImage.Bounds())
	draw.Draw(jpgImage, jpgImage.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(jpgImage, jpgImage.Bounds(), pngImage, pngImage.Bounds().Min, draw.Over)
	return jpgImage, nil
}
