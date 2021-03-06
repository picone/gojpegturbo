[![Build](https://github.com/picone/gojpegturbo/actions/workflows/main.yml/badge.svg)](https://github.com/picone/gojpegturbo/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/picone/gojpegturbo/branch/main/graph/badge.svg)](https://codecov.io/gh/picone/gojpegturbo)
[![License](https://img.shields.io/badge/license-Apache%202-blue.svg)](LICENSE)
[![rcard](https://goreportcard.com/badge/github.com/picone/gojpegturbo)](https://goreportcard.com/report/github.com/picone/gojpegturbo)

# gojpegturbo

一个go封装 [libjpeg-turbo](https://github.com/libjpeg-turbo/libjpeg-turbo) 的库。
和其他封装不同的是，封装了特别多的接口，方便某些场景能够达到极致的性能，以最大程度利用libjpeg-turbo的特性。

## Getting Start

### 局部图片解码

JPEG图片是由一个个MCU组成的，从左到右，从上到下排列。我们在哈夫曼编码解码后就能知道接下来的多少字节是一个MCU了。在某种特殊场景，需要剪裁其中一小块
区域，这时候局部解码就很有用，我们只需要找到对应区域所在的MCU并解码，能大大减少IDCT的次数，当然内存也减少很多。一个demo：

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/picone/gojpegturbo"
)

func main() {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	options := gojpegturbo.NewDecodeOptions()
	options.CropRect = &image.Rectangle{
		Min: image.Point{X: 100, Y: 200},
		Max: image.Point{X: 300, Y: 400},
	}
	outImg, err := gojpegturbo.Decode(buf, options)
	fmt.Println(outImg.Bounds().Max) // (200,200)
}
```

### 图片解码时缩放

图片解码时，如常用的4:2:2采样Cb和Cr通道是会进行升采样的，如果图片解码后边长需要等比缩放到比原来的1/2还小，可以使CbCr不进行生采样，Y通道使用降采样，
内存节省1/4，CPU也会使用更少。 如图片原来是600×800，如业务需要缩放至300×300，则可以使用`options.ScaleNum = 1`,`options.ScaleDenom = 2`
，解码出来的图片将会是300×400，业务再使用别的缩放算法进行图片变换。jpeg解码过程中的图片缩放性能消耗极少，并且节省更多内存。完整demo如下：

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/picone/gojpegturbo"
)

func main() {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	options := gojpegturbo.NewDecodeOptions()
	options.ScaleNum = 1
	options.ScaleDenom = 2
	outImg, err := gojpegturbo.Decode(buf, options)
	fmt.Println(outImg.Bounds().Max) // (300,400)
}
```

### 更高级的解码参数

通过调整解码的参数，在接受图片质量稍微变差的同时能提供更快的速度，部分场景下适用（如生成较小的缩略图，图片质量并不那么重要了）。

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/picone/gojpegturbo"
)

func main() {
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
	outImg, err := gojpegturbo.Decode(buf, options)
}
```

## Performance

### Decode性能测试

以下仅使用单张图片对比。

```text
goos: darwin
goarch: amd64
pkg: github.com/picone/gojpegturbo
cpu: Intel(R) Core(TM) i7-5650U CPU @ 2.20GHz
BenchmarkDecodeC
BenchmarkDecodeC-4                    	     469	   2585729 ns/op	  33.79 MB/s
BenchmarkDecodeDctFast
BenchmarkDecodeDctFast-4              	     465	   2500954 ns/op	  34.93 MB/s
BenchmarkDecodeDoSloppierSampling
BenchmarkDecodeDoSloppierSampling-4   	     484	   2428223 ns/op	  35.98 MB/s
BenchmarkDecodeDitherOrdered
BenchmarkDecodeDitherOrdered-4        	     448	   3118717 ns/op	  28.01 MB/s
BenchmarkDecodeLessColors
BenchmarkDecodeLessColors-4           	     460	   2749298 ns/op	  31.78 MB/s
BenchmarkDecodeGo
BenchmarkDecodeGo-4                   	     135	   8480772 ns/op	  10.30 MB/s
PASS
```

### Encode性能测试

以下为单张图片测试

```text
goos: darwin
goarch: amd64
pkg: github.com/picone/gojpegturbo
cpu: Intel(R) Core(TM) i7-5650U CPU @ 2.20GHz
BenchmarkEncodeC
BenchmarkEncodeC-4                    	     508	   2334368 ns/op	  37.42 MB/s
BenchmarkEncodeFast
BenchmarkEncodeFast-4                 	     516	   2317246 ns/op	  37.70 MB/s
BenchmarkEncodeProgress
BenchmarkEncodeProgress-4             	      86	  12894303 ns/op	   6.78 MB/s
BenchmarkEncodeGo
BenchmarkEncodeGo-4                   	      43	  26813204 ns/op	   3.26 MB/s
PASS
```

### 缩放性能测试

由于使用了自己的`ImageAttr`类，图片缩放并不能很好地兼容[nfnt/resize](https://github.com/nfnt/resize)，
因此自己实现了 NearestNeighbor 和 Area 的图片缩放算法。

```text
goos: darwin
goarch: amd64
pkg: github.com/picone/gojpegturbo
cpu: Intel(R) Core(TM) i7-5650U CPU @ 2.20GHz
BenchmarkImageAttr_ResizeArea
BenchmarkImageAttr_ResizeArea-4       	     310	   4008882 ns/op
BenchmarkImageAttr_ResizeNN
BenchmarkImageAttr_ResizeNN-4         	    2302	    478747 ns/op
BenchmarkImageAttr_ResizeBilinear
BenchmarkImageAttr_ResizeBilinear-4   	      48	  25154275 ns/op
PASS
```

## Contributing

- Please create an issue in [issue list](https://github.com/picone/gojpegturbo/issues).
- Contact Committers/Owners for further discussion if needed.
- Following the golang coding standards.

## License

gojpegturbo is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
