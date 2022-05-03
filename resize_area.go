package gojpegturbo

import (
	"errors"
	"math"
)

type areaTableItem struct {
	srcIdx int
	dstIdx int
	alpha  float32
}

var (
	// ErrWrongDstSize 输入的目标尺寸又唔
	ErrWrongDstSize = errors.New("error input dst image width or height")
)

// ResizeArea 参考opencv的INTER_AREA算法，目前只能用于图片缩小，放大场景下效果不佳。
func ResizeArea(src *ImageAttr, dstWidth, dstHeight int) (*ImageAttr, error) {
	if dstWidth > src.ImageWidth || dstHeight > src.ImageHeight || dstWidth == 0 || dstHeight == 0 {
		return nil, ErrWrongDstSize
	}
	dst := &ImageAttr{
		Img:           make([]byte, dstWidth*dstHeight*src.ComponentsNum),
		ImageWidth:    dstWidth,
		ImageHeight:   dstHeight,
		OriginWidth:   dstWidth,
		OriginHeight:  dstHeight,
		ColorSpace:    src.ColorSpace,
		ComponentsNum: src.ComponentsNum,
	}
	// 计算各个缩放cell的index和对应权重
	hTb := calcAreaTable(src.ImageWidth, dst.ImageWidth, src.ComponentsNum)
	vTb := calcAreaTable(src.ImageHeight, dst.ImageHeight, src.ComponentsNum)
	dstRowBuf := make([]float32, dstWidth*src.ComponentsNum)
	prevDstIdx := 0
	for _, vItem := range vTb {
		// 统计每一行各个pixel加权结果
		srcRowBuf := make([]float32, dstWidth*src.ComponentsNum)
		srcRowIdx := src.ImageWidth * vItem.srcIdx // 计算当前行src.Img的下标开始
		if src.ComponentsNum == 3 {
			for _, hItem := range hTb {
				srcRowBuf[hItem.dstIdx] += float32(src.Img[srcRowIdx+hItem.srcIdx]) * hItem.alpha
				srcRowBuf[hItem.dstIdx+1] += float32(src.Img[srcRowIdx+hItem.srcIdx+1]) * hItem.alpha
				srcRowBuf[hItem.dstIdx+2] += float32(src.Img[srcRowIdx+hItem.srcIdx+2]) * hItem.alpha
			}
		} else if src.ComponentsNum == 1 {
			for _, hItem := range hTb {
				srcRowBuf[hItem.dstIdx] += float32(src.Img[srcRowIdx+hItem.srcIdx]) * hItem.alpha
			}
		} else {
			for _, hItem := range hTb {
				for i := 0; i < src.ComponentsNum; i++ {
					srcRowBuf[hItem.dstIdx+i] += float32(src.Img[srcRowIdx+hItem.srcIdx+i]) * hItem.alpha
				}
			}
		}
		// 统计这行并加上vItem.alpha。若是新的dstIdx则输出结果到dst.Img
		if vItem.dstIdx != prevDstIdx {
			for i := 0; i < dstWidth*src.ComponentsNum; i++ {
				dst.Img[prevDstIdx*dstWidth+i] = byte(dstRowBuf[i])
				dstRowBuf[i] = vItem.alpha * srcRowBuf[i]
			}
			prevDstIdx = vItem.dstIdx
		} else {
			for i := 0; i < dstWidth*src.ComponentsNum; i++ {
				dstRowBuf[i] += vItem.alpha * srcRowBuf[i]
			}
		}
	}
	return dst, nil
}

// calcAreaTable 计算每个像素的权重。
//
//	各个变量的位置分布长这样：
//	0      0.25               1            ...            10              10.75
//	| ----- | --------------- | ------------------------- | --------------- |
//	        ↑                 ↑                           ↑                 ↑
//	     srcStart         srcStartInt      ...        srcEndInt          srcEnd
//	<---------------------------- cell_width ------------------------------>
func calcAreaTable(srcSize, dstSize, pixel int) []areaTableItem {
	factor := float32(srcSize) / float32(dstSize)
	tb := make([]areaTableItem, 0, 2*dstSize) // 2倍的dstSize肯定够用，多申请点避免重新申请内存
	for i := 0; i < dstSize; i++ {
		srcStart := float32(i) * factor
		srcEnd := srcStart + factor // (i+1) * factor
		cellWidth := factor
		if lastCellWidth := float32(srcSize) - srcStart; lastCellWidth < cellWidth {
			cellWidth = lastCellWidth
		}
		srcStartInt := int(math.Ceil(float64(srcStart))) // 向上取整
		srcEndInt := int(math.Floor(float64(srcEnd)))    // 向下取整
		if lastBlock := srcSize - 1; lastBlock < srcEndInt {
			srcEndInt = lastBlock
		}
		if srcStartInt < srcEndInt {
			srcStartInt = srcEndInt
		}
		// 为什么是1e-3? 因为权重太小能忽略影响了。
		if float32(srcStartInt)-srcStart > 1e-3 {
			tb = append(tb, areaTableItem{
				dstIdx: i * pixel,
				srcIdx: (srcStartInt - 1) * pixel,
				alpha:  (float32(srcStartInt) - srcStart) / cellWidth,
			})
		}
		for j := srcStartInt; j < srcEndInt; j++ {
			tb = append(tb, areaTableItem{
				dstIdx: i * pixel,
				srcIdx: j * pixel,
				alpha:  1 / cellWidth,
			})
		}
		// srcEndInt右边的pixel，不足一个，同理需要加上权重alpha
		if srcEnd-float32(srcEndInt) > 1e-3 {
			tb = append(tb, areaTableItem{
				dstIdx: i * pixel,
				srcIdx: srcEndInt * pixel,
				alpha:  (srcEnd - float32(srcEndInt)) / cellWidth,
			})
		}
	}
	return tb
}
