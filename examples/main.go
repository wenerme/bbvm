package main

import (
	. "../."
	"log"
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
	log.Println(vm)

	//	b := make([]byte, 4)
	//	i := -1
	//	Codec.PutInt(b, i)
	//	log.Print(Codec.Int(b))
	//	log.Print(len(int(1)))
	//	v := NewVal()
	//	fmt.Print(v.Get())

	var b B = B(A{})
	fmt.Print(b.Get())

}


