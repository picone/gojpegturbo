package gojpegturbo

/*
#cgo linux LDFLAGS: -ljpeg
#cgo darwin LDFLAGS: -L/usr/local/opt/libjpeg-turbo/lib -ljpeg
#cgo darwin CFLAGS: -I/usr/local/opt/libjpeg-turbo/include

#include <stdio.h>
#include "turbojpeg.h"
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
	DitherNone DitherMode = C.JDITHER_NONE
	// DitherOrdered 有序抖动，中等速度，中等质量。
	DitherOrdered DitherMode = C.JDITHER_ORDERED
	// DitherFs Floyd-Steinberg算法的抖动，慢，高质量。
	DitherFs DitherMode = C.JDITHER_FS
)

// TJPixelFormat 像素存储的格式。
type TJPixelFormat int

const (
	// TJPixelFormatRGB RGB
	TJPixelFormatRGB TJPixelFormat = C.TJPF_RGB
	// TJPixelFormatBGR BGR
	TJPixelFormatBGR TJPixelFormat = C.TJPF_BGR
	// TJPixelFormatRGBX RGBX
	TJPixelFormatRGBX TJPixelFormat = C.TJPF_RGBX
	// TJPixelFormatBGRX BGRX
	TJPixelFormatBGRX TJPixelFormat = C.TJPF_BGRX
	// TJPixelFormatXBGR XBGR
	TJPixelFormatXBGR TJPixelFormat = C.TJPF_XBGR
	// TJPixelFormatXRGB XRGB
	TJPixelFormatXRGB TJPixelFormat = C.TJPF_XRGB
	// TJPixelFormatGray gray
	TJPixelFormatGray TJPixelFormat = C.TJPF_GRAY
	// TJPixelFormatRGBA RGBA
	TJPixelFormatRGBA TJPixelFormat = C.TJPF_RGBA
	// TJPixelFormatBGRA BGRA
	TJPixelFormatBGRA TJPixelFormat = C.TJPF_BGRA
	// TJPixelFormatABGR ABGR
	TJPixelFormatABGR TJPixelFormat = C.TJPF_ABGR
	// TJPixelFormatARGB ARGB
	TJPixelFormatARGB TJPixelFormat = C.TJPF_ARGB
	// TJPixelFormatCMYK CMYK
	TJPixelFormatCMYK TJPixelFormat = C.TJPF_CMYK
	// TJPixelFormatUnknown 未知
	TJPixelFormatUnknown TJPixelFormat = C.TJPF_UNKNOWN
)

// TJSubSample 二次采样方法，一般是4:2:0采样。
type TJSubSample int

const (
	// TjSubSample444 每个颜色分量对应1个像素点。
	TjSubSample444 TJSubSample = C.TJSAMP_444
	// TjSubSample422 每个颜色分量对应包含2x1个像素区块。
	TjSubSample422 TJSubSample = C.TJSAMP_422
	// TjSubSample420 每个颜色分量对应包含2x2个像素区块。
	TjSubSample420 TJSubSample = C.TJSAMP_420
	// TjSubSampleGray 只有明暗分量没有颜色分量
	TjSubSampleGray TJSubSample = C.TJSAMP_GRAY
	// TjSubSample440 每个颜色分量对应包含1x2个像素区块。
	// note: 4:4:0采样在turbojpeg中没有完全加速。
	TjSubSample440 TJSubSample = C.TJSAMP_440
	// TjSubSample411 每个颜色分量对应包含4x1个像素区块，和420大小一样的，但是更好地表现横向特征。
	TjSubSample411 TJSubSample = C.TJSAMP_411
)

const (
	// TjFlagFastDCT 使用更快速的DCT/IDCT算法
	TjFlagFastDCT = C.TJFLAG_FASTDCT
	// TjFlagAccurateDCT 使用更准确的DCT/IDCT算法
	TjFlagAccurateDCT = C.TJFLAG_ACCURATEDCT
	// TjFlagProgressive 使用渐进式编码
	TjFlagProgressive = C.TJFLAG_PROGRESSIVE
)
