package bi
import (
	"testing"
	"io/ioutil"
	"bytes"
	"os"
	_ "image"
	_ "image/color"
	"image/color"
)



func TestReadLibConfig(t *testing.T) {
	b, err := ioutil.ReadFile("../tests/case/9288-wener.lib")
	if err != nil {panic(err)}
	r := bytes.NewReader(b)
	r.Seek(8, os.SEEK_SET)
	i, err := DecodeLibGray2One(r)
	if err != nil {panic(err)}
	log.Info("Color %+v", i.At(44, 0))
	log.Info("Color %+v", i.At(45, 0))
	log.Info("Color %+v", i.At(46, 0))
	log.Info("Color %+v", i.At(47, 0))
	log.Info("Color %+v", i.At(48, 0))
	log.Info("Color %+v", color.RGBAModel.Convert(i.At(44, 0)))
	log.Info("Bytes %v", i.(*Gray2).Pix[0:20])
	log.Info("Byte 301 %v", i.(*Gray2).Pix[602:620])
	log.Info("Stride %v", i.(*Gray2).Stride)
	//	i = imageConvert(i, image.NewRGBA(i.Bounds()))
	//		log.Info("%+v", i.At(10, 10))
	//		log.Info("%+v", i.At(50, 50))
	//	log.Info("%+v", i.At(100, 100))

	//		log.Info("%+v", color.RGBAModel.Convert(i.At(10, 10)))
	//		log.Info("%+v", color.RGBAModel.Convert(i.At(50, 50)))
	//	log.Info("%+v", color.RGBAModel.Convert(i.At(100, 100)))

	saveTemp(i)
}
