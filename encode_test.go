package gojpegturbo

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg" // 注册jpeg解码库
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		imgAttr  *ImageAttr
		options  *EncodeOptions
		wantErr  bool
	}{
		{
			name:     "case 1",
			filename: "./testdata/test.jpg",
		},
		{
			name:    "case 2",
			imgAttr: nil,
			wantErr: true,
		},
		{
			name: "case 3",
			imgAttr: &ImageAttr{
				Img:         []byte(`1234`),
				ImageWidth:  100,
				ImageHeight: 100,
			},
			wantErr: true,
		},
		{
			name:     "case 4",
			filename: "./testdata/test.jpg",
			options: &EncodeOptions{
				Quality:   60,
				SubSample: TjSubSampleGray,
			},
		},
		{
			name:     "case 5",
			filename: "./testdata/test.jpg",
			options: &EncodeOptions{
				Quality:     90,
				Progressive: true,
			},
		},
		{
			name:     "case 6",
			filename: "./testdata/test.jpg",
			options: &EncodeOptions{
				Quality: 90,
				FastDct: true,
			},
		},
		{
			name:     "case 7",
			filename: "./testdata/test.jpg",
			options: &EncodeOptions{
				Quality:     90,
				AccurateDCT: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imgAttr := tt.imgAttr
			if tt.filename != "" {
				fp, err := os.Open(tt.filename)
				require.NoError(t, err)
				imgAttr, err = DecodeReader(fp, nil)
				require.NoError(t, err)
			}
			got, err := Encode(imgAttr, tt.options)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Greater(t, len(got), 0)
				// 尝试用原生的去解码是否正常
				_, _, err := image.Decode(bytes.NewBuffer(got))
				assert.NoError(t, err)
			}
		})
	}
}

func BenchmarkEncodeC(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	img, err := Decode(buf, nil)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Encode(img, nil)
		assert.NoError(b, err)
	}
	b.SetBytes(int64(len(buf)))
}

func BenchmarkEncodeFast(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	img, err := Decode(buf, nil)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		options := NewEncodeOptions()
		options.FastDct = true
		_, err := Encode(img, options)
		assert.NoError(b, err)
	}
	b.SetBytes(int64(len(buf)))
}

func BenchmarkEncodeProgress(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	img, err := Decode(buf, nil)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		options := NewEncodeOptions()
		options.Progressive = true
		_, err := Encode(img, options)
		assert.NoError(b, err)
	}
	b.SetBytes(int64(len(buf)))
}

func BenchmarkEncodeGo(b *testing.B) {
	buf, err := ioutil.ReadFile("./testdata/test.jpg")
	require.NoError(b, err)
	img, err := Decode(buf, nil)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := bytes.NewBuffer(nil)
		err := jpeg.Encode(out, img, nil)
		assert.NoError(b, err)
	}
	b.SetBytes(int64(len(buf)))
}
