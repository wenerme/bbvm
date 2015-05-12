package bi
import (
	"testing"
	"io/ioutil"
	"bytes"
	"os"
	"image"
	"image/draw"
	"image/color"
	"golang.org/x/image/bmp"
)


func TestReadConfig(t *testing.T) {
	b, err := ioutil.ReadFile("../tests/case/wener.rlb")
	if err != nil {panic(err)}
	is, cfgs, err := Decode(bytes.NewReader(b))
	if err != nil {panic(err)}
	i := is[0]
	di := i.(*image.NRGBA)
	_ = di
	//	for i := 0; i < len(di.Pix); i ++ {
	//		if i % 5 == 0 {
	//			di.Pix[i] = uint8(0xff)
	//		}
	//	}
	d := i.(draw.Image)
	for y := 0; y < i.Bounds().Dy(); y ++ {
		for x := 0; x < i.Bounds().Dx(); x ++ {
			c := d.At(x, y)
			col := c.(color.NRGBA)
			col.A = uint8(0xff)
			d.Set(x, y, col)
		}
	}
	log.Info("%+v", cfgs)
	log.Info("%#v", i.At(10, 10))
	log.Info("%#v", i.At(100, 100))
	log.Info("%+v", i.Bounds())
	log.Info("%T", i)

	if _, ok := i.(draw.Image); ok {
		log.Info("Can edit")
	}else {
		log.Info("Can not edit")
	}

	saveTemp(i)
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