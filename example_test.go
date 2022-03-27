package gojpegturbo_test

import (
	"fmt"
	"io/ioutil"
	"log"

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
