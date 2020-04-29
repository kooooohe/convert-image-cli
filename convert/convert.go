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

var exts = []string{"gif", "jpeg", "png"}

// ConvertImage は変換対象画像ファイル情報です
type ConvertImage struct {
	File        *os.File
	FilePath    string
	Image       image.Image
	ImageFormat string
}

// ConvertImages は変換対象画像ファイル情報スライスです
type ConvertImages []*ConvertImage

func Convert(dir, bExt, aExt string) error {
	if dir == "" {
		return errors.New("no dir name")
	}
	if bExt == "" {
		return errors.New("no b ext name")
	}
	if aExt == "" {
		return errors.New("no a ext name")
	}
	if !contains(exts, aExt) {
		return errors.New("no a ext name")
	}
	if !contains(exts, bExt) {
		return errors.New("no b bxt name")

	}

	cis, err := getTargetImages(dir, bExt)
	if err != nil {
		return err
	}
	err = cis.ConvertImagesFromTo(bExt, aExt)
	if err != nil {
		return err
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

// getFileNameWithoutExt は拡張子なしのファイル名を返却します
func (ci *ConvertImage) getFileNameWithoutExt() string {
	p := ci.FilePath
	return p[:len(p)-len(filepath.Ext(p))]
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

	err := cis.ConvertImagesTo(a)

	return err
}

// NewConvertImage は指定されたファイルから生成したConvertImageを返却します。
func NewConvertImage(p string) (ci *ConvertImage, err error) {

	f, err := os.Open(p)
	if err != nil {
		return
	}
	defer f.Close()

	img, fmt, err := image.Decode(f)
	if err != nil {
		return
	}

	ci = &ConvertImage{File: f, FilePath: p, Image: img, ImageFormat: fmt}

	return
}

// NewConvertImagesByDir は指定されたディレクトリに含まれる画像ファイルから生成したImageFileのスライスを返却します。
func getTargetImages(dir, tExt string) (ConvertImages, error) {
	cis := ConvertImages{}

	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if ok, _ := isTargetImage(p, tExt); ok {
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
		return nil, err
	}

	return cis, nil
}

// isImage はファイルパスからそのファイルが画像か判定します
func isTargetImage(p, tExt string) (ok bool, err error) {
	f, err := os.Open(p)
	if err != nil {
		return
	}
	defer f.Close()

	_, fmt, err := image.Decode(f)
	if err != nil {
		return
	}

	ok = fmt == tExt
	return
}
