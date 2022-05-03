package gojpegturbo_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/picone/gojpegturbo"
)

func ExampleDecode() {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	options := gojpegturbo.NewDecodeOptions()
	options.DctMethod = gojpegturbo.DctMethodIntFast
	options.DitherMode = gojpegturbo.DitherOrdered
	options.TwoPassQuantize = false
	options.DoFancyUpSampling = false
	options.DesiredNumberOfColors = 216
	img, err := gojpegturbo.Decode(buf, options) // options can be nil and all options set to default.
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("width=%d,height=%d\n", img.Bounds().Max.X, img.Bounds().Max.Y)
	// Output:
	// width=600,height=800
}

func ExampleResizeArea() {
	fp, err := os.Open("./testdata/test.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	srcImg, err := gojpegturbo.DecodeReader(fp, nil) // here you can cut your image with options.
	if err != nil {
		log.Fatalln(err)
	}
	// if you use Area algo, it SHOULD NOT pass width or height max than origin.
	dstImg, err := srcImg.ResizeArea(200, 400)
	if err != nil {
		log.Fatalln(err)
	}
	encOptions := gojpegturbo.NewEncodeOptions()
	encOptions.Quality = 60
	buf, err := gojpegturbo.Encode(dstImg, encOptions)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%v", len(buf) > 0) // here you can write `buf` to file
	// Output:
	// true
}
