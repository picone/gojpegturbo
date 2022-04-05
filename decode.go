package gojpegturbo

/*
#cgo linux LDFLAGS: -lturbojpeg
#cgo darwin LDFLAGS: -L/usr/local/opt/libjpeg-turbo/lib -lturbojpeg
#cgo darwin CFLAGS: -I/usr/local/opt/libjpeg-turbo/include

#include "goturbo.h"
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"unsafe"
)

// DecodeOptions 解码图片时的选项
type DecodeOptions struct {
	// CropRect 图片剪裁区域，默认不剪裁
	CropRect *image.Rectangle
	// DctMethod 解码的时候使用的方法。
	// 现在的计算机上有AVX2，JDCT_IFAST和JDCT_ISLOW有相似的性能。如果JPEG图像使用85质量一下的等级压缩的，那么这两种算法
	// 应该没有差别。如果高于85，实践中如quality=97，JDCT_IFAST通常会比JDCT_ISLOW的PSNR低大约4-6dB。因此一般不是用
	// JDCT_IFAST。对于JDCT_FLOAT并不一定质量就好，因为每个机器的四舍五入情况不一致。
	DctMethod DctMethod
	// TwoPassQuantize 默认是true。
	TwoPassQuantize bool
	// DitherMode 抖动方法，默认是DitherFs。
	DitherMode DitherMode
	// DesiredNumberOfColors 颜色表使用的颜色数量，默认是256。
	DesiredNumberOfColors int
	// DoFancyUpSampling 升采样是否使用精确的，默认是true。
	DoFancyUpSampling bool
	// ScaleNum 按比例缩放图片，在MCU升降采样的时候就能生效，一般为1
	ScaleNum uint
	// ScaleDenom 缩放比例的分母，为了保证图片效果，一般取值为1/2,1/4,1/8...
	ScaleDenom uint
}

// NewDecodeOptions 创建一个默认的解码图片选项
func NewDecodeOptions() *DecodeOptions {
	return &DecodeOptions{
		DctMethod:             DctMethodIntSlow,
		TwoPassQuantize:       true,
		DitherMode:            DitherFs,
		DesiredNumberOfColors: 256,
		DoFancyUpSampling:     true,
	}
}

// 转成cgo传入需要的options对象
func (options *DecodeOptions) toCOptions() (*C.jpeg_decode_options, error) {
	if options == nil {
		return nil, nil
	}
	// 暂时不支持scale & crop同时，会有panic。
	if options.CropRect != nil && (options.ScaleNum > 0 || options.ScaleDenom > 0) {
		return nil, ErrOptionsUnsupported
	}
	co := &C.jpeg_decode_options{
		dct_method:               C.J_DCT_METHOD(options.DctMethod),
		dither_mode:              C.J_DITHER_MODE(options.DitherMode),
		desired_number_of_colors: C.int(options.DesiredNumberOfColors),
		scale_num:                C.uint(options.ScaleNum),
		scale_denom:              C.uint(options.ScaleDenom),
	}
	if options.TwoPassQuantize {
		co.two_pass_quantize = C.int(1)
	}
	if options.DoFancyUpSampling {
		co.do_fancy_upsampling = C.int(1)
	}
	if options.CropRect != nil {
		co.crop.left = C.uint(uint(options.CropRect.Min.X))
		co.crop.top = C.uint(uint(options.CropRect.Min.Y))
		co.crop.width = C.uint(uint(options.CropRect.Max.X - options.CropRect.Min.X))
		co.crop.height = C.uint(uint(options.CropRect.Max.Y - options.CropRect.Min.Y))
	}
	return co, nil
}

var (
	// ErrEmptyImage 传入图片的空的
	ErrEmptyImage = errors.New("empty image")
	// ErrEmptyDecode 解码结果是空的
	ErrEmptyDecode = errors.New("decode image empty")
	// ErrOptionsUnsupported 当前的选项并不支持
	ErrOptionsUnsupported = errors.New("options now unsupported")
)

// Decode 解码JPEG图片
func Decode(img []byte, options *DecodeOptions) (*ImageAttr, error) {
	if len(img) == 0 {
		return nil, ErrEmptyImage
	}
	jres := C.jpeg_decode_result{}
	co, err := options.toCOptions()
	if err != nil {
		return nil, err
	}
	C.jpeg_decode((*C.uchar)(unsafe.Pointer(&img[0])), C.uint(uint(len(img))), co, &jres)
	if jres.img != nil {
		defer C.free(unsafe.Pointer(jres.img))
	}
	if jres.err != nil {
		defer C.free(unsafe.Pointer(jres.err))
		return nil, fmt.Errorf("jpeg_decode failed, err = %s", C.GoString(jres.err))
	}
	if jres.img == nil || int(jres.img_size) == 0 {
		return nil, ErrEmptyDecode
	}
	imgAttr := &ImageAttr{
		Img:           C.GoBytes(unsafe.Pointer(jres.img), C.int(int(jres.img_size))),
		ImageWidth:    int(jres.image_width),
		ImageHeight:   int(jres.image_height),
		OriginWidth:   int(jres.origin_width),
		OriginHeight:  int(jres.origin_height),
		ColorSpace:    ColorSpace(jres.color_space),
		ComponentsNum: int(jres.num_components),
	}
	return imgAttr, nil
}

// DecodeReader 解码reader过来的图片
func DecodeReader(r io.Reader, options *DecodeOptions) (*ImageAttr, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, err
	}
	return Decode(buf.Bytes(), options)
}
