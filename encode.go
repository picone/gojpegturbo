package gojpegturbo

/*
#cgo linux LDFLAGS: -lturbojpeg
#cgo darwin LDFLAGS: -L/usr/local/opt/libjpeg-turbo/lib -lturbojpeg
#cgo darwin CFLAGS: -I/usr/local/opt/libjpeg-turbo/include

#include "goturbo.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	// ErrImgEmpty 图片为空
	ErrImgEmpty = errors.New("image is empty")
	// ErrImgSizeInvalid 图片尺寸不合法
	ErrImgSizeInvalid = errors.New("image size is invalid")
	// ErrQualityOption 图片质量不合法
	ErrQualityOption = errors.New("quality option invalid")
)

// EncodeOptions 图片编码的选项
type EncodeOptions struct {
	// Quality 图片压缩质量
	Quality int
	// FastDct 快速的DCT
	FastDct bool
	// 精确的DCT
	AccurateDCT bool
	// SubSample 采样率
	SubSample TJSubSample
	// Progressive 是否使用渐进式编码
	Progressive bool
}

// NewEncodeOptions 创建一个默认的图片编码选项
func NewEncodeOptions() *EncodeOptions {
	return &EncodeOptions{
		Quality:   95,
		SubSample: TjSubSample420,
	}
}

// toCOptions 转成C的struct
func (options *EncodeOptions) toCOptions() (*C.jpeg_encode_options, error) {
	if options == nil {
		return nil, nil
	}
	// 暂时不支持scale & crop同时，会有panic。
	if options.Quality > 100 || options.Quality < 0 {
		return nil, ErrQualityOption
	}
	tjFlag := 0
	if options.AccurateDCT {
		tjFlag |= TjFlagFastDCT
	}
	if options.AccurateDCT {
		tjFlag |= TjFlagAccurateDCT
	}
	if options.Progressive {
		tjFlag |= TjFlagProgressive
	}
	co := &C.jpeg_encode_options{
		quality:    C.int(options.Quality),
		tj_flag:    C.int(tjFlag),
		sub_sample: C.int(options.SubSample),
	}
	return co, nil
}

// Encode jpeg图片编码
func Encode(img *ImageAttr, options *EncodeOptions) ([]byte, error) {
	if img == nil || len(img.Img) == 0 {
		return nil, ErrImgEmpty
	}
	if img.ImageWidth*img.ImageHeight*img.ComponentsNum != len(img.Img) {
		return nil, ErrImgSizeInvalid
	}
	jres := C.jpeg_encode_result{}
	co, err := options.toCOptions()
	if err != nil {
		return nil, err
	}
	C.jpeg_encode((*C.uchar)(unsafe.Pointer(&img.Img[0])), C.int(img.ImageWidth), C.int(img.ImageHeight),
		C.int(img.PixelFormat()), co, &jres)
	if jres.img == nil {
		defer C.free(unsafe.Pointer(jres.img))
	}
	if jres.err != nil {
		defer C.free(unsafe.Pointer(jres.err))
		return nil, fmt.Errorf("jpeg_encode failed, err = %s", C.GoString(jres.err))
	}
	return C.GoBytes(unsafe.Pointer(jres.img), C.int(int(jres.img_size))), nil
}
