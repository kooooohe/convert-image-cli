package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"convert_image/convert"
)

var (
	tDir string
	bExt string
	aExt string
)

func init() {
	flag.StringVar(&tDir, "dir", "", "Target Dir")
	flag.StringVar(&bExt, "bExt", convert.ExtJpeg, "変換前画像形式")
	flag.StringVar(&aExt, "aExt", convert.ExtPng, "変換後画像形式")
}

func main() {
	flag.Parse()

	if tDir == "" {
		fmt.Println("対象ディレクトリを指定してください")
		os.Exit(1)
	}

	cis, err := convert.GetTargetImages(tDir)
	if err != nil {
		log.Print(err)
	}
	err = cis.ConvertImagesFromTo(bExt, aExt)
	if err != nil {
		log.Print(err)
	}
}
