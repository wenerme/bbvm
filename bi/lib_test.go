package bi
import (
	"testing"
	"io/ioutil"
	"bytes"
	"os"
	_ "image"
	_ "image/color"
)



func TestReadLibConfig(t *testing.T) {
	b, err := ioutil.ReadFile("../tests/case/9688-wener.lib")
	//	b, err := ioutil.ReadFile("../tests/case/9288-wener.lib")
	if err != nil {panic(err)}
	r := bytes.NewReader(b)
	r.Seek(8, os.SEEK_SET)
	i, err := DecodeLibRGB565One(r)
	if err != nil {panic(err)}
	//	i = imageConvert(i, image.NewRGBA(i.Bounds()))
	//		log.Info("%+v", i.At(10, 10))
	//		log.Info("%+v", i.At(50, 50))
	//	log.Info("%+v", i.At(100, 100))

	//		log.Info("%+v", color.RGBAModel.Convert(i.At(10, 10)))
	//		log.Info("%+v", color.RGBAModel.Convert(i.At(50, 50)))
	//	log.Info("%+v", color.RGBAModel.Convert(i.At(100, 100)))

	saveTemp(i)
}
