package main

import (
	. "../."
	"log"
)

func main() {
	vm := NewVM()
	log.Println(vm)

	b := make([]byte, 4)
	i := -1
	Codec.PutInt(b, i)
	log.Print(Codec.Int(b))
//	log.Print(len(int(1)))
}


