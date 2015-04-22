package main

import (
	. "../."
	"log"
)

type A struct {
	V int
}
func (a *A)Get() int {
	return a.V
}

func (a *A)Inc() int {
	return a.Get()+1
}
type VA interface {
	Get() int
	Inc() int
}
type B struct {
	A
	V2 int
}
func (a *B)Get() int {
	return a.V2
}
func main() {
	vm := NewVM()
	log.Println(vm)

	//	b := make([]byte, 4)
	//	i := -1
	//	Codec.PutInt(b, i)
	//	log.Print(Codec.Int(b))
	//	log.Print(len(int(1)))
	b := B{}
	b.V = 10
	b.V2 = 20
	log.Print(b.Inc())
}


