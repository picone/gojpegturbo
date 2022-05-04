package gojpegturbo

import "math"

// ResizeNN 邻近插值法缩放图片
func ResizeNN(src *ImageAttr, dstWidth, dstHeight int) *ImageAttr {
	hFactor := float32(src.ImageWidth) / float32(dstWidth)
	vFactor := float32(src.ImageHeight) / float32(dstHeight)
	dst := &ImageAttr{
		Img:           make([]byte, dstWidth*dstHeight*src.ComponentsNum),
		ImageWidth:    dstWidth,
		ImageHeight:   dstHeight,
		OriginWidth:   dstWidth,
		OriginHeight:  dstHeight,
		ColorSpace:    src.ColorSpace,
		ComponentsNum: src.ComponentsNum,
	}
	// 先计算映射表，避免后面多次计算。下标代表目标idx，值代表src图片的idx
	hTb := make([]int, dstWidth)
	for i := 0; i < dstWidth; i++ {
		hTb[i] = int(math.Floor(float64(float32(i)*hFactor))) * src.ComponentsNum
		if hTb[i] >= src.ImageWidth*src.ComponentsNum {
			hTb[i] = src.ImageWidth*src.ComponentsNum - 1 // 防止溢出
		}
	}
	// 缩放图片
	dstRowIdx := 0
	for i := 0; i < dstHeight; i++ {
		srcRowIdx := int(math.Floor(float64(float32(i)*vFactor))) * src.ImageWidth * src.ComponentsNum
		if src.ComponentsNum == 3 {
			for j := 0; j < dstWidth; j++ {
				idx := srcRowIdx + hTb[j]
				dst.Img[dstRowIdx] = src.Img[idx]
				dst.Img[dstRowIdx+1] = src.Img[idx+1]
				dst.Img[dstRowIdx+2] = src.Img[idx+2]
				dstRowIdx += 3
			}
		} else if src.ComponentsNum == 1 {
			for j := 0; j < dstWidth; j++ {
				dst.Img[dstRowIdx] = src.Img[srcRowIdx+hTb[j]]
				dstRowIdx++
			}
		} else {
			for j := 0; j < dstWidth; j++ {
				idx := srcRowIdx + hTb[j]
				for k := 0; k < src.ComponentsNum; k++ {
					dst.Img[dstRowIdx+k] = src.Img[idx+k]
				}
				dstRowIdx += src.ComponentsNum
			}
		}
	}
	return dst
}
