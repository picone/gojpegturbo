package gojpegturbo

import (
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageAttr(t *testing.T) {
	type want struct {
		colorModel color.Model
		bounds     image.Rectangle
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
			},
		},
		//{
		//	name:     "case 2-gray",
		//	filename: "./testdata/gray.jpg",
		//	want: want{
		//		colorModel: color.GrayModel,
		//		bounds: image.Rectangle{
		//			Max: image.Point{
		//				X: 600,
		//				Y: 800,
		//			},
		//		},
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := ioutil.ReadFile(tt.filename)
			require.NoError(t, err)
			img, err := Decode(buf, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.want.colorModel, img.ColorModel())
			assert.Equal(t, tt.want.bounds, img.Bounds())
			err = jpeg.Encode(ioutil.Discard, img, nil)
			assert.NoError(t, err)
		})
	}
}
