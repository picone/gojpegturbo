[![Build](https://github.com/picone/gojpegturbo/actions/workflows/main.yml/badge.svg)](https://github.com/picone/gojpegturbo/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/picone/gojpegturbo/branch/main/graph/badge.svg)](https://codecov.io/gh/picone/gojpegturbo)
[![License](https://img.shields.io/badge/license-Apache%202-blue.svg)](LICENSE)

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
	options.DctMethod = gojpegturbo.DctMethodIntFast
	options.DitherMode = gojpegturbo.DitherOrdered
	options.TwoPassQuantize = false
	options.DoFancyUpSampling = false
	options.DesiredNumberOfColors = 216
	outImg, err := gojpegturbo.Decode(buf, options)
	fmt.Println(outImg.Bounds().Max) // (600,800)
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
	options.ScaleNum = 1
	options.ScaleDenom = 2
	outImg, err := gojpegturbo.Decode(buf, options)
	fmt.Println(outImg.Bounds().Max) // (300,400)
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
BenchmarkDecodeC-4                    	     391	   2700117 ns/op
BenchmarkDecodeDctFast
BenchmarkDecodeDctFast-4              	     477	   2461163 ns/op
BenchmarkDecodeDoSloppierSampling
BenchmarkDecodeDoSloppierSampling-4   	     500	   2656346 ns/op
BenchmarkDecodeDitherOrdered
BenchmarkDecodeDitherOrdered-4        	     464	   2536009 ns/op
BenchmarkDecodeLessColors
BenchmarkDecodeLessColors-4           	     459	   2599797 ns/op
BenchmarkDecodeGo
BenchmarkDecodeGo-4                   	     133	   8788681 ns/op
PASS
```

## Contributing

- Please create an issue in [issue list](https://github.com/picone/gojpegturbo/issues).
- Contact Committers/Owners for further discussion if needed.
- Following the golang coding standards.

## License

gojpegturbo is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
