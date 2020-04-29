package convert

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
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

// Image は変換対象画像ファイル情報です
type Image struct {
	File        *os.File
	FilePath    string
	Image       image.Image
	ImageFormat string
}

// ConvertImages は変換対象画像ファイル情報スライスです
type ConvertImages struct {
	cImages []convertImage
}

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

	cis, err := getTargetImages(dir, bExt, aExt)
	if err != nil {
		return err
	}
	err = cis.ConvertImages()
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
func (ci *Image) getFileNameWithoutExt() string {
	p := ci.FilePath
	return p[:len(p)-len(filepath.Ext(p))]
}

type convertImage interface {
	Type() string
	getConvertFn() func(w io.Writer, m image.Image) error
	getImage() *Image
}
type cGif struct {
	image *Image
}

func (c *cGif) Type() string {
	return ExtGif
}
func (c *cGif) getImage() *Image {
	return c.image
}
func (c *cGif) getConvertFn() func(io.Writer, image.Image) error {
	return func(nf io.Writer, ci image.Image) error {
		return gif.Encode(nf, ci, nil)
	}
}

type cJepg struct {
	image *Image
}

func (c *cJepg) Type() string {
	return ExtJpeg
}
func (c *cJepg) getConvertFn() func(io.Writer, image.Image) error {
	return func(nf io.Writer, ci image.Image) error {
		return jpeg.Encode(nf, ci, nil)
	}
}
func (c *cJepg) getImage() *Image {
	return c.image
}

type cPng struct {
	image *Image
}

func (c *cPng) Type() string {
	return ExtPng
}
func (c *cPng) getConvertFn() func(io.Writer, image.Image) error {
	return func(nf io.Writer, ci image.Image) error {
		return png.Encode(nf, ci)
	}
}
func (c *cPng) getImage() *Image {
	return c.image
}

func convert(ce convertImage) error {
	ci := ce.getImage()
	np := ci.getFileNameWithoutExt() + "." + ce.Type()
	nf, err := os.Create(np)
	if err != nil {
		return err
	}

	fn := ce.getConvertFn()
	err = fn(nf, ci.Image)
	if err != nil {
		return err
	}

	bp := ci.FilePath

	ci.File, ci.FilePath, ci.ImageFormat = nf, np, ce.Type()

	if err := os.Remove(bp); err != nil {
		return err
	}

	fmt.Printf("%sの画像形式を%s(%s)に変更しました。", bp, ce.Type(), np)

	return nil
}

func (cis ConvertImages) ConvertImages() (err error) {
	for _, v := range cis.cImages {
		err = convert(v)
		if err != nil {
			return err
		}
	}

	return
}

func NewConvertImage(p, aExt string) (ci convertImage, err error) {

	f, err := os.Open(p)
	if err != nil {
		return
	}
	defer f.Close()

	img, fmt, err := image.Decode(f)
	if err != nil {
		return
	}

	i := &Image{File: f, FilePath: p, Image: img, ImageFormat: fmt}
	switch aExt {
	case ExtGif:
		return &cGif{image: i}, nil
	case ExtJpeg:
		return &cJepg{image: i}, nil
	case ExtPng:
		return &cPng{image: i}, nil
	default:
		err = errors.New("指定されたフォーマットは対応していません")
	}

	return
}

func getTargetImages(dir, tExt, aExt string) (*ConvertImages, error) {
	cis := &ConvertImages{}

	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if ok, _ := isTargetImage(p, tExt); ok {
				ci, err := NewConvertImage(p, aExt)
				if err != nil {
					return err
				}
				cis.cImages = append(cis.cImages, ci)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return cis, nil
}

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
