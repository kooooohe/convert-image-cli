package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"convert_image/convert"
)

var (
	b string
	a string
)

func init() {
	flag.StringVar(&b, "b", convert.ExtJpeg, "変換前画像形式")
	flag.StringVar(&a, "a", convert.ExtPng, "変換後画像形式")
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("対象ディレクトリを指定してください")
		os.Exit(1)
	}

	for _, dir := range flag.Args() {
		cis, err := convert.NewConvertImagesByDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		err = cis.ConvertImagesFromTo(b, a)
		if err != nil {
			log.Fatal(err)
		}
	}
}
