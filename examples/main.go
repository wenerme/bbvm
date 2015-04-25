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
func (a *A)Get() int {
	return a.V
}

func (a A)Inc() int {
	return a.Get()+1
}
type VA interface {
	Get() int
	Inc() int
}
type B A
func (a B)Get() int {
	return (A)(a).Inc() + 10
}
//func (b B)Inc()int{
//	return b.Get()+2
//}
func main() {
	t2()
}
func t2() {
	i := image.NewRGBA(image.Rect(0, 0, 101, 101))
	p := &Pen{}
	p.Image = i
	//	r := image.Rect(20, 20, 80, 80)
	p.Set(50, 50, color.RGBA{0xff, 0, 0, 0xff})

	p.Color = color.RGBA{0, 0, 0, 0xff}
	//	p.DrawLine(0,100, 100,0) // 坡度为反
	//	p.DrawLine(80,20,20,80) // 坡度为反
	//	p.DrawLine(0,100,100,0) // 坡度为反
	//	p.DrawLine(80, 0, 0, 80)

	r := image.Rect(20, 20, 80, 80)
	p.Rect(r)

	p.Color = color.RGBA{0xff, 0, 0, 0xff}
	p.Circle(50, 50, 10)
	//	p.Color = color.RGBA{0, 0xff, 0, 0xff}
	//	p.DrawLine(0, 20, 100, 80)
	//	p.DrawLine(0, 50, 100, 50)
	//	p.DrawLine(100, 75, 0, 75)
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

	//	b := make([]byte, 4)
	//	i := -1
	//	Codec.PutInt(b, i)
	//	log.Print(Codec.Int(b))
	//	log.Print(len(int(1)))
	//	v := NewVal()
	//	fmt.Print(v.Get())

	var b B = B(A{})
	fmt.Println(b.Get())
	fmt.Printf("%010.6f\n", 1.234)

	var i interface{}
	i = nil
	ia, ok := i.(*A)
	fmt.Print(ia, ok)
}


