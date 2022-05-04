package gojpegturbo

import (
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageAttr(t *testing.T) {
	type want struct {
		colorModel  color.Model
		bounds      image.Rectangle
		pixelFormat TJPixelFormat
	}
	tests := []struct {
		name     string
		filename string
		want     want
	}{
		{
			name:     "case 1-YCbCr",
			filename: "./testdata/test.jpg",
			want: want{
				colorModel: color.RGBAModel,
				bounds: image.Rectangle{
					Max: image.Point{
						X: 600,
						Y: 800,
					},
				},
				pixelFormat: TJPixelFormatRGB,
			},
		},
		{
			name:     "case 2-gray",
			filename: "./testdata/gray.jpg",
			want: want{
				colorModel: color.GrayModel,
				bounds: image.Rectangle{
					Max: image.Point{
						X: 600,
						Y: 800,
					},
				},
				pixelFormat: TJPixelFormatGray,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := ioutil.ReadFile(tt.filename)
			require.NoError(t, err)
			img, err := Decode(buf, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want.colorModel, img.ColorModel())
			assert.Equal(t, tt.want.bounds, img.Bounds())
			assert.Equal(t, tt.want.pixelFormat, img.PixelFormat())
			err = jpeg.Encode(ioutil.Discard, img, nil)
			assert.NoError(t, err)
		})
	}
}

func TestImageAttr_ResizeArea(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		dstWidth  int
		dstHeight int
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:      "case 1",
			filename:  "./testdata/test.jpg",
			dstWidth:  500,
			dstHeight: 100,
		},
		{
			name:      "case 2-error size",
			filename:  "./testdata/test.jpg",
			dstWidth:  1000,
			dstHeight: 200,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorAs(t, ErrWrongDstSize, &err)
			},
		},
		{
			name:      "case 3",
			filename:  "./testdata/gray.jpg",
			dstWidth:  422,
			dstHeight: 235,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp, err := os.Open(tt.filename)
			require.NoError(t, err)
			img, err := DecodeReader(fp, nil)
			require.NoError(t, err)
			got, err := img.ResizeArea(tt.dstWidth, tt.dstHeight)
			if tt.wantErr == nil {
				require.NoError(t, err)
				assert.Equal(t, got.ImageWidth, tt.dstWidth)
				assert.Equal(t, got.ImageHeight, tt.dstHeight)
				assert.Equal(t, got.OriginWidth, tt.dstWidth)
				assert.Equal(t, got.OriginHeight, tt.dstHeight)
				err = jpeg.Encode(io.Discard, got, nil)
				require.NoError(t, err)
			} else {
				tt.wantErr(t, err)
			}
		})
	}
}

func TestImageAttr_ResizeNN(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		dstWidth  int
		dstHeight int
	}{
		{
			name:      "case 1",
			filename:  "./testdata/test.jpg",
			dstWidth:  500,
			dstHeight: 100,
		},
		{
			name:      "case 2",
			filename:  "./testdata/gray.jpg",
			dstWidth:  422,
			dstHeight: 235,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp, err := os.Open(tt.filename)
			require.NoError(t, err)
			img, err := DecodeReader(fp, nil)
			require.NoError(t, err)
			got := img.ResizeNN(tt.dstWidth, tt.dstHeight)
			require.NoError(t, err)
			assert.Equal(t, got.ImageWidth, tt.dstWidth)
			assert.Equal(t, got.ImageHeight, tt.dstHeight)
			assert.Equal(t, got.OriginWidth, tt.dstWidth)
			assert.Equal(t, got.OriginHeight, tt.dstHeight)
			err = jpeg.Encode(ioutil.Discard, got, nil)
			require.NoError(t, err)
		})
	}
}

func BenchmarkImageAttr_ResizeArea(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	img, err := Decode(buf, nil)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := img.ResizeArea(233, 455)
		require.NoError(b, err)
	}
}

func BenchmarkImageAttr_ResizeNN(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	img, err := Decode(buf, nil)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = img.ResizeNN(233, 455)
	}
}

func BenchmarkImageAttr_ResizeBilinear(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	img, err := Decode(buf, nil)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = img.ResizeBilinear(233, 455)
	}
}
