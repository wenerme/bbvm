package image
import (
	"testing"
	"io/ioutil"
	"bytes"
	"os"
	_ "image"
	_ "image/color"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/bmp"
	"image"
	"image/draw"
	"log"
)



func TestReadLibConfig(t *testing.T) {
	b, err := ioutil.ReadFile("../tests/9688-wener.lib")
	//	b, err := ioutil.ReadFile("../tests/case/9288-wener.lib")
	if err != nil {panic(err)}
	r := bytes.NewReader(b)
	configs, err := DecodeConfig(r)
	_ = configs
	if err != nil {panic(err)}
}

func TestDetectImageFormat(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]imageFormat{
		"../tests/9688-wener.lib": libRGB565Format,
		"../tests/9188-wener.lib": libGray2BEFormat,
		"../tests/9288-wener.lib": libGray2LEFormat,
		"../tests/wener.rlb": rlbFormat,
	}


	for fn, fe := range tests {
		fp, err := os.Open(fn)
		if err != nil {panic(err)}
		f, err := detectImageFormat(fp)
		if err != nil {panic(err)}
		assert.EqualValues(fe, f)
		log.Printf("%s -> %s", fn, f)
	}
}

func saveTemp(i image.Image) {
	p, err := os.Create("temp.bmp")
	if err != nil {panic(err)}
	err = bmp.Encode(p, i)
	if err != nil {panic(err)}
}

func imageConvert(src image.Image, dest draw.Image) draw.Image {
	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			dest.Set(x, y, dest.ColorModel().Convert(src.At(x, y)))
		}
	}
	return dest
}