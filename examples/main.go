package main

import (
	. "../."
	"fmt"
	"image"
	"os"
	"image/png"
	"image/color"
)

type A struct {
	V int
}

type B struct {
	A
}
func main() {
	t3()
}
func t3() {
	a := []byte("ABCDEF")
	b := make([]byte, 0)
	b = append(b, a[1:3]...)
	b[0]=57
	fmt.Println(string(b))
	fmt.Println(string(a))

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
	if err != nil {panic(err)}
	err = png.Encode(p, i)
	if err != nil {panic(err)}
}
func t1() {
	vm := NewVM()
	_ = vm

	var i interface{}
	i = nil
	ia, ok := i.(*A)
	fmt.Print(ia, ok)
}


