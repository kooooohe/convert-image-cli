package convert

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

// 対応画像形式
const (
	ExtGif  = "gif"
	ExtJpeg = "jpeg"
	ExtPng  = "png"
)

// ConvertImage は変換対象画像ファイル情報です
type ConvertImage struct {
	File        *os.File
	FilePath    string
	Image       image.Image
	ImageFormat string
}

// ConvertImages は変換対象画像ファイル情報スライスです
type ConvertImages []*ConvertImage

// getFileNameWithoutExt は拡張子なしのファイル名を返却します
func (ci *ConvertImage) getFileNameWithoutExt() string {
	p := ci.FilePath
	return p[:len(p)-len(filepath.Ext(p))]
}

// isImageGif は画像形式がGIFか判定します
func (ci *ConvertImage) isImageGif() bool {
	return ci.ImageFormat == ExtGif
}

// isImageJpeg は画像形式がJPEGか判定します
func (ci *ConvertImage) isImageJpeg() bool {
	return ci.ImageFormat == ExtJpeg
}

// isImagePng は画像形式がPNGか判定します
func (ci *ConvertImage) isImagePng() bool {
	return ci.ImageFormat == ExtPng
}

// convertImageToGif は画像形式をGIFに変換します
func (ci *ConvertImage) convertImageToGif() error {
	np := ci.getFileNameWithoutExt() + "." + ExtGif
	nf, err := os.Create(np)
	if err != nil {
		return err
	}

	err = gif.Encode(nf, ci.Image, nil)
	if err != nil {
		return err
	}

	bp := ci.FilePath

	ci.File, ci.FilePath, ci.ImageFormat = nf, np, ExtGif

	if err := os.Remove(bp); err != nil {
		return err
	}

	fmt.Printf("%sの画像形式を%s(%s)に変更しました。", bp, ExtGif, np)

	return nil
}

// convertImageToJpeg は画像形式をJPEFに変換します
func (ci *ConvertImage) convertImageToJpeg() error {
	np := ci.getFileNameWithoutExt() + "." + ExtJpeg
	nf, err := os.Create(np)
	if err != nil {
		return err
	}

	err = jpeg.Encode(nf, ci.Image, nil)
	if err != nil {
		return err
	}

	bp := ci.FilePath

	ci.File, ci.FilePath, ci.ImageFormat = nf, np, ExtJpeg

	if err := os.Remove(bp); err != nil {
		return err
	}

	fmt.Printf("%sの画像形式を%s(%s)に変更しました。", bp, ExtJpeg, np)

	return nil
}

// convertImageToPng は画像形式をPNGに変換します
func (ci *ConvertImage) convertImageToPng() error {
	np := ci.getFileNameWithoutExt() + "." + ExtPng
	nf, err := os.Create(np)
	if err != nil {
		return err
	}

	err = png.Encode(nf, ci.Image)
	if err != nil {
		return err
	}

	bp := ci.FilePath

	ci.File, ci.FilePath, ci.ImageFormat = nf, np, ExtPng

	if err := os.Remove(bp); err != nil {
		return err
	}

	fmt.Printf("%sの画像形式を%s(%s)に変更しました。", bp, ExtPng, np)

	return nil
}

// ConvertImageTo は画像を指定された画像形式に変換します
func (ci *ConvertImage) ConvertImageTo(fmt string) (err error) {
	switch fmt {
	case ExtGif:
		err = ci.convertImageToGif()
	case ExtJpeg:
		err = ci.convertImageToJpeg()
	case ExtPng:
		err = ci.convertImageToPng()
	default:
		err = errors.New("指定されたフォーマットは対応していません")
	}

	return
}

// getOnlyImageGif はConvertImagesからGIF形式のイメージを保持するConvertImageを抽出して返却します
func (cis ConvertImages) getOnlyImageGif() (rcis ConvertImages) {
	for _, v := range cis {
		if v.isImageGif() {
			rcis = append(rcis, v)
		}
	}

	return rcis
}

// getOnlyImageJpeg はConvertImagesからJPEG形式のイメージを保持するConvertImageを抽出して返却します
func (cis ConvertImages) getOnlyImageJpeg() (rcis ConvertImages) {
	for _, v := range cis {
		if v.isImageJpeg() {
			rcis = append(rcis, v)
		}
	}

	return rcis
}

// GetOnlyImagePng はConvertImagesからPNG形式のイメージを保持するConvertImageを抽出して返却します
func (cis ConvertImages) getOnlyImagePng() (rcis ConvertImages) {
	for _, v := range cis {
		if v.isImagePng() {
			rcis = append(rcis, v)
		}
	}

	return rcis
}

// GetOnly はConvertImagesから指定された画像形式のイメージを保持するConvertImageを抽出して返却します
func (cis ConvertImages) GetOnly(fmt string) (ConvertImages, error) {
	rcis := ConvertImages{}
	switch fmt {
	case ExtGif:
		rcis = cis.getOnlyImageGif()
	case ExtJpeg:
		rcis = cis.getOnlyImageJpeg()
	case ExtPng:
		rcis = cis.getOnlyImagePng()
	default:
		return rcis, errors.New("指定されたフォーマットは対応していません")
	}

	return rcis, nil
}

// ConvertImagesTo はConvertImagesに含まれる画像を指定された画像形式の画像に変換します
func (cis ConvertImages) ConvertImagesTo(fmt string) (err error) {
	for _, v := range cis {
		err = v.ConvertImageTo(fmt)
		if err != nil {
			return err
		}
	}

	return
}

// ConvertImagesFromTo はConvertImagesに含まれる指定された画像形式の画像を指定された画像形式の画像に変換します
func (cis ConvertImages) ConvertImagesFromTo(b string, a string) error {
	rcis, err := cis.GetOnly(b)
	if err != nil {
		return err
	}

	err = rcis.ConvertImagesTo(a)

	return err
}

// NewConvertImage は指定されたファイルから生成したConvertImageを返却します。
func NewConvertImage(p string) (*ConvertImage, error) {
	ci := new(ConvertImage)

	f, err := os.Open(p)
	if err != nil {
		return ci, err
	}
	defer f.Close()

	img, fmt, err := image.Decode(f)
	if err != nil {
		return ci, err
	}

	ci.File, ci.FilePath, ci.Image, ci.ImageFormat = f, p, img, fmt

	return ci, nil
}

// NewConvertImagesByDir は指定されたディレクトリに含まれる画像ファイルから生成したImageFileのスライスを返却します。
func GetTargetImages(dir string) (ConvertImages, error) {
	cis := ConvertImages{}

	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if ok, _ := isImage(p); ok {
				ci, err := NewConvertImage(p)
				if err != nil {
					return err
				}
				cis = append(cis, ci)
			}
		}
		return nil
	})
	if err != nil {
		return cis, err
	}

	return cis, nil
}

// isImage はファイルパスからそのファイルが画像か判定します
func isImage(p string) (bool, error) {
	f, err := os.Open(p)
	if err != nil {
		return false, err
	}
	defer f.Close()

	if _, _, err := image.Decode(f); err != nil {
		return false, err
	}

	return true, nil
}
