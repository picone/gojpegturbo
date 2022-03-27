package gojpegturbo

/*
#cgo linux LDFLAGS: -ljpeg
#cgo darwin LDFLAGS: -L/usr/local/opt/libjpeg-turbo/lib -ljpeg
#cgo darwin CFLAGS: -I/usr/local/opt/libjpeg-turbo/include

#include <stdio.h>
#include "jpeglib.h"
*/
import "C"

// ColorSpace 色彩空间
type ColorSpace C.J_COLOR_SPACE

const (
	// ColorSpaceUnknown 未知
	ColorSpaceUnknown = C.JCS_UNKNOWN
	// ColorSpaceGrayScale 灰色
	ColorSpaceGrayScale = C.JCS_GRAYSCALE
	// ColorSpaceRGB RGB
	ColorSpaceRGB = C.JCS_RGB
	// ColorSpaceYCbCr YCbCr
	ColorSpaceYCbCr = C.JCS_YCbCr
	// ColorSpaceCMYK CMYK
	ColorSpaceCMYK = C.JCS_CMYK
	// ColorSpaceYCCK YCCK
	ColorSpaceYCCK = C.JCS_YCCK
	// ColorSpaceExtRGB ExtRGB
	ColorSpaceExtRGB = C.JCS_EXT_RGB
	// ColorSpaceExtRGBX ExtRGBX
	ColorSpaceExtRGBX = C.JCS_EXT_RGBX
	// ColorSpaceExtBGR ExtBGR
	ColorSpaceExtBGR = C.JCS_EXT_BGR
	// ColorSpaceExtBGRX ExtBGRX
	ColorSpaceExtBGRX = C.JCS_EXT_BGRX
	// ColorSpaceExtXBGR ExtXBGR
	ColorSpaceExtXBGR = C.JCS_EXT_XBGR
	// ColorSpaceExtXRGB = ExtXRGB
	ColorSpaceExtXRGB = C.JCS_EXT_XRGB
	// ColorSpaceExtRGBA ExtRGBA
	ColorSpaceExtRGBA = C.JCS_EXT_RGBA
	// ColorSpaceExtBGRA ExtBGRA
	ColorSpaceExtBGRA = C.JCS_EXT_BGRA
	// ColorSpaceExtABGR ExtABGR
	ColorSpaceExtABGR = C.JCS_EXT_ABGR
	// ColorSpaceExtARGB ARGB
	ColorSpaceExtARGB = C.JCS_EXT_ARGB
	// ColorSpaceExtRGB565 RGB565
	ColorSpaceExtRGB565 = C.JCS_RGB565
)

// DctMethod DCT/IDCT的方法
type DctMethod C.J_DCT_METHOD

const (
	// DctMethodIntSlow 使用int，较慢（精确）的DCT
	DctMethodIntSlow = C.JDCT_ISLOW
	// DctMethodIntFast 使用int且较快（不太精确）的DCT
	DctMethodIntFast = C.JDCT_IFAST
	// DctMethodFloat 使用浮点数的DCT
	DctMethodFloat = C.JDCT_FLOAT
)

// DitherMode 颜色抖动方法
type DitherMode C.J_DITHER_MODE

const (
	// DitherNone 无，快，低质量。
	DitherNone = C.JDITHER_NONE
	// DitherOrdered 有序抖动，中等速度，中等质量。
	DitherOrdered = C.JDITHER_ORDERED
	// DitherFs Floyd-Steinberg算法的抖动，慢，高质量。
	DitherFs = C.JDITHER_FS
)
