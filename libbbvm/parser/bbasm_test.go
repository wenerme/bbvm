package parser

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/juju/errors"
	"github.com/wenerme/bbvm/libbbvm/asm"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestTypeMatch(t *testing.T) {
	v := reflect.ValueOf(buildInst)
	spew.Dump(v, v.Kind() == reflect.Func)
	spew.Dump(reflect.TypeOf(buildInst))
}

func TestBBAsm(t *testing.T) {
	b, e := ioutil.ReadFile("testdata/exp.txt")
	if e != nil {
		panic(e)
	}
	code := string(b)
	p := &BBAsm{Buffer: code}
	p.Init()
	if err := p.Parse(); err != nil {
		log.Fatal(err)
	}
	p.PrintSyntaxTree()

	func() {
		defer func() {
			if e := recover(); e != nil {
				spew.Dump(p.stack)
				fmt.Println("----------------------------------------------")
				for _, a := range p.assemblies {
					fmt.Println(a.Assembly())
				}
				if e, ok := e.(error); ok {
					panic(errors.ErrorStack(e.(error)))
				} else {
					panic(e)
				}
			}

		}()
		p.Execute()
	}()

	spew.Dump(p.stack)
	for _, a := range p.assemblies {
		fmt.Println(a.Assembly())
	}
}

func TestLK(t *testing.T) {
	fmt.Println(asm.Lookup(asm.T_INT, "int"))
}
