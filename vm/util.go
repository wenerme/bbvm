package vm
import (
	"golang.org/x/image/bmp"
	"os"
	"image"
)


func saveImage(i image.Image, fn string) {
	p, err := os.Create(fn)
	if err != nil {panic(err)}
	err = bmp.Encode(p, i)
	if err != nil {panic(err)}
}
