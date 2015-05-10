package rlb
import (
	"testing"
	"io/ioutil"
	"bytes"
	"os"
	"image"
	"image/jpeg"
)


func TestReadConfig(t *testing.T) {
	b, err := ioutil.ReadFile("../tests/case/wener.rlb")
	if err != nil {panic(err)}
	is, cfgs, err := Decode(bytes.NewReader(b))
	if err != nil {panic(err)}
	i := is[0]
	log.Info("%+v", cfgs)
	log.Info("%+v", i.Bounds())
	saveTemp(i)
}



func saveTemp(i image.Image) {
	p, err := os.Create("temp.jpeg")
	if err != nil {panic(err)}
	err = jpeg.Encode(p, i, nil)
	if err != nil {panic(err)}
}
