/*Package gojpegturbo 封装libjpeg-turbo库，优化图片的缩放和剪裁。

简单的图片解码编码：
	img, err := gojpegturbo.Decode(buf, nil)
	// modify img
	buf, err = gojpegturbo.Encode(img, nil)

图片缩放：
	options := gojpegturbo.NewDecodeOptions()
	options.ExpectWidth = 50
	options.ExpectHeight = 50
	img, err := gojpegturbo.Decode(buf, options)
	// 注意，这里输出的 img 的长和宽总是大于或等于 expect 的尺寸，会在图片解码阶段尽量逼近需要，但一般不会刚好等于，需要使用 resize 函数继续缩小。
	img = img.ResizeNN(50, 50)
	buf, err = gojpegturbo.Encode(img, nil)

图片剪裁：
	options := gojpegturbo.NewDecodeOptions()
	options.CropRect = &image.Rectangle{
		Min: image.Point{X: 100, Y: 200},
		Max: image.Point{X: 300, Y: 621},
	}
	// 这里输出的图片是已经剪裁好的。相比起其他图片解码库，优势在于解码的时候只解码感兴趣的 MCU 区域，大大节省 CPU 和内存。
	img, err := gojpegturbo.Decode(buf, options)
	buf, err = gojpegturbo.Encode(img, nil)
*/
package gojpegturbo
