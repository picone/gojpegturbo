package gojpegturbo

import (
	"bytes"
	"image"
	_ "image/jpeg" // 注册jpeg解码库
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	type args struct {
		filename string
		options  *DecodeOptions
	}
	tests := []struct {
		name     string
		args     args
		wantSize image.Point
		wantErr  bool
	}{
		{
			name: "case 1",
			args: args{
				filename: "./testdata/test.jpg",
			},
		},
		{
			name: "case 2-options with default",
			args: args{
				filename: "./testdata/test.jpg",
				options:  NewDecodeOptions(),
			},
		},
		{
			name: "case 3-dct int fast",
			args: args{
				filename: "./testdata/test.jpg",
				options: &DecodeOptions{
					DctMethod:             DctMethodIntFast,
					TwoPassQuantize:       true,
					DitherMode:            DitherFs,
					DesiredNumberOfColors: 256,
					DoFancyUpSampling:     true,
				},
			},
		},
		{
			name: "case 4-upsampling",
			args: args{
				filename: "./testdata/test.jpg",
				options: &DecodeOptions{
					DctMethod:             DctMethodIntSlow,
					TwoPassQuantize:       true,
					DitherMode:            DitherFs,
					DesiredNumberOfColors: 256,
					DoFancyUpSampling:     false,
				},
			},
		},
		{
			name: "case 5-two pass quantize",
			args: args{
				filename: "./testdata/test.jpg",
				options: &DecodeOptions{
					DctMethod:             DctMethodIntSlow,
					TwoPassQuantize:       false,
					DitherMode:            DitherFs,
					DesiredNumberOfColors: 256,
					DoFancyUpSampling:     true,
				},
			},
		},
		{
			name: "case 6-crop",
			args: args{
				filename: "./testdata/test.jpg",
				options: &DecodeOptions{
					CropRect: &image.Rectangle{
						Min: image.Point{X: 100, Y: 200},
						Max: image.Point{X: 300, Y: 621},
					},
					DctMethod:             DctMethodIntSlow,
					TwoPassQuantize:       false,
					DitherMode:            DitherFs,
					DesiredNumberOfColors: 256,
					DoFancyUpSampling:     true,
				},
			},
			wantSize: image.Point{X: 200, Y: 421},
		},
		{
			name: "case 7-scale",
			args: args{
				filename: "./testdata/test.jpg",
				options: &DecodeOptions{
					DctMethod:             DctMethodIntSlow,
					TwoPassQuantize:       false,
					DitherMode:            DitherFs,
					DesiredNumberOfColors: 256,
					DoFancyUpSampling:     true,
					ScaleNum:              1,
					ScaleDenom:            2,
				},
			},
			wantSize: image.Point{X: 300, Y: 400},
		},
		{
			name: "case 8-crop&scale",
			args: args{
				filename: "./testdata/test.jpg",
				options: &DecodeOptions{
					CropRect: &image.Rectangle{
						Min: image.Point{X: 100, Y: 200},
						Max: image.Point{X: 300, Y: 621},
					},
					ScaleNum:   1,
					ScaleDenom: 2,
				},
			},
			wantErr: true,
		},
		{
			name: "case 9-error",
			args: args{
				filename: "./testdata/error.jpg",
			},
			wantErr: true,
		},
		{
			name: "case 10-gray",
			args: args{
				filename: "./testdata/gray.jpg",
			},
			wantSize: image.Point{X: 600, Y: 800},
		},
		{
			name: "case 11-expect size",
			args: args{
				filename: "./testdata/test.jpg",
				options: &DecodeOptions{
					ExpectWidth:  100,
					ExpectHeight: 200,
				},
			},
			wantSize: image.Point{X: 150, Y: 200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := ioutil.ReadFile(tt.args.filename)
			require.NoError(t, err)
			got, err := Decode(buf, tt.args.options)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, got.ImageWidth*got.ImageHeight*got.ComponentsNum, len(got.Img))
				if tt.wantSize.X > 0 && tt.wantSize.Y > 0 {
					assert.Equal(t, tt.wantSize, image.Point{X: got.ImageWidth, Y: got.ImageHeight})
				}
			}
		})
	}
}

func BenchmarkDecodeC(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decode(buf, nil)
		assert.NoError(b, err)
	}
	b.SetBytes(int64(len(buf)))
}

func BenchmarkDecodeDctFast(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	options := NewDecodeOptions()
	options.DctMethod = DctMethodIntFast
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decode(buf, options)
		assert.NoError(b, err)
	}
	b.SetBytes(int64(len(buf)))
}

func BenchmarkDecodeDoSloppierSampling(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	options := NewDecodeOptions()
	options.DoFancyUpSampling = false
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decode(buf, options)
		assert.NoError(b, err)
	}
	b.SetBytes(int64(len(buf)))
}

func BenchmarkDecodeDitherOrdered(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	options := NewDecodeOptions()
	options.DitherMode = DitherOrdered
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decode(buf, options)
		assert.NoError(b, err)
	}
	b.SetBytes(int64(len(buf)))
}

func BenchmarkDecodeLessColors(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	options := NewDecodeOptions()
	options.DesiredNumberOfColors = 216
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decode(buf, options)
		assert.NoError(b, err)
	}
	b.SetBytes(int64(len(buf)))
}

func BenchmarkDecodeGo(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := image.Decode(bytes.NewBuffer(buf))
		assert.NoError(b, err)
	}
	b.SetBytes(int64(len(buf)))
}
