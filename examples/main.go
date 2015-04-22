package main

import (
	. "../."
	"fmt"
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
	i := uint32(4294967173)
	fmt.Println(int32(i))

}


