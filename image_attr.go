package gojpegturbo

import "C"
import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
)

var _ image.Image = (*ImageAttr)(nil)

// ImageAttr 图像属性。保存图片原始宽高，目前宽高，颜色空间等基本属性。除此之外，实现了image.Image接口，能够直接使用众多第三方包二次处理。
type ImageAttr struct {
	Img                       []byte
	ImageWidth, ImageHeight   int        // 输出的图片宽高，若没有剪裁，和Origin的宽高一样
	OriginWidth, OriginHeight int        // 原始图片宽高
	ColorSpace                ColorSpace // 色彩空间。目前只有gray和YCbCr。
	ComponentsNum             int        // 颜色分量数，如YCbCr就是3。
}

// ColorModel 色彩空间
func (img *ImageAttr) ColorModel() color.Model {
	if img.ColorSpace == ColorSpaceGrayScale {
		return color.GrayModel
	}
	return color.RGBAModel
}

// Bounds 选择范围。若还需要图片剪裁可以直接在这里返回别的值。目前只返回原始尺寸
func (img *ImageAttr) Bounds() image.Rectangle {
	return image.Rectangle{
		Max: image.Point{
			X: int(img.ImageWidth),
			Y: int(img.ImageHeight),
		},
	}
}

// At 获取指定像素点的RBG颜色
func (img *ImageAttr) At(x, y int) color.Color {
	offset := (x + y*img.ImageWidth) * img.ComponentsNum
	if img.ColorSpace == ColorSpaceGrayScale {
		return &color.Gray{Y: img.Img[offset]}
	}
	return &color.RGBA{
		R: img.Img[offset],
		G: img.Img[offset+1],
		B: img.Img[offset+2],
	}
}

// PixelFormat 像素格式
func (img *ImageAttr) PixelFormat() TJPixelFormat {
	if img.ColorSpace == ColorSpaceGrayScale {
		return TJPixelFormatGray
	}
	return TJPixelFormatRGB
}

// ResizeArea 用 INTER_AREA 方法缩小图片，这个方法不能用于图片放大。
func (img *ImageAttr) ResizeArea(dstWidth, dstHeight int) (*ImageAttr, error) {
	return ResizeArea(img, dstWidth, dstHeight)
}

// ResizeNN 使用 NearestNeighbor 算法缩放图片，速度很快，但是会有锯齿。
func (img *ImageAttr) ResizeNN(dstWidth, dstHeight int) *ImageAttr {
	return ResizeNN(img, dstWidth, dstHeight)
}

// ResizeBilinear 使用 Bilinear （双线插值）算法，速度次之，但是锯齿会少很多。
func (img *ImageAttr) ResizeBilinear(dstWidth, dstHeight uint) image.Image {
	return resize.Resize(dstWidth, dstHeight, img, resize.Bilinear)
}
