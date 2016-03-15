package main

import (
	"fmt"
	. "github.com/wenerme/bbvm/libbbvm/image"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	t3()
}
func t4() {
	cd, _ := os.Getwd()
	fmt.Println(cd)
	fp, _ := os.Open("tests/WENER.RLB")
	fp.Seek(0x28+4, os.SEEK_SET)
	i, f, err := image.Decode(fp)
	fmt.Print(f, err)
	xi := i.(*image.NRGBA)
	fmt.Println(xi.NRGBAAt(10, 10))
	fmt.Println(xi.At(10, 10))
	xi.Set(10, 10, color.RGBA{0xff, 0xff, 0, 0xff})

	dy, dx := i.Bounds().Dy(), i.Bounds().Dx()
	for y := 0; y < dy; y += 1 {
		for x := 0; x < dx; x += 1 {
			i.(draw.Image).Set(x, y, setAlpha(i.At(x, y), 0xff))
		}
	}

	saveTemp(i)
}
func setAlpha(c color.Color, alpha uint8) (result color.Color) {
	switch c.(type) {
	case color.NRGBA:
		xc := c.(color.NRGBA)
		xc.A = alpha
		result = xc
	default:
		r, g, b, _ := c.RGBA()
		result = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), alpha}
	}
	return
}
func t3() {
	cd, _ := os.Getwd()
	fmt.Println(cd)
	fp, err := os.Open("tests/WENER.RLB")
	if err != nil {
		panic(err)
	}
	img, err := DecodeAt(fp, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("SIZE: %v\n", img.Bounds())
	fmt.Printf("0,0: %v\n", img.At(0, 0))
	fmt.Printf("100,100: %v\n", img.At(0, 0))
	saveTemp(img)
}
func t2() {
	i := image.NewRGBA(image.Rect(0, 0, 101, 101))
	p := NewImageGraphic(i)
	//	p.Image = i
	//	r := image.Rect(20, 20, 80, 80)
	p.Set(50, 50, color.RGBA{0xff, 0, 0, 0xff})

	p.SetColor(color.RGBA{0, 0, 0, 0xff})
	//	p.DrawLine(0,100, 100,0) // 坡度为反
	//	p.DrawLine(80,20,20,80) // 坡度为反
	//	p.DrawLine(0,100,100,0) // 坡度为反
	//	p.DrawLine(80, 0, 0, 80)

	r := image.Rect(20, 20, 80, 80)
	p.Rect(r)

	p.SetColor(color.RGBA{0, 0x88, 0x88, 0xff})
	p.FillRect(image.Rect(20, 20, 50, 50))
	p.SetColor(color.RGBA{0xff, 0, 0, 0xff})
	p.Circle(50, 50, 10)
	//	p.Color = color.RGBA{0, 0xff, 0, 0xff}
	//	p.DrawLine(0, 20, 100, 80)
	//	p.DrawLine(0, 50, 100, 50)
	//	p.DrawLine(100, 75, 0, 75)
	p.SetColor(color.RGBA{0, 0x88, 0x88, 0xff})
	saveTemp(p)
}

func saveTemp(i image.Image) {
	p, err := os.Create("temp.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(p, i)
	if err != nil {
		panic(err)
	}
}
