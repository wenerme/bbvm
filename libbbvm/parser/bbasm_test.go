package parser

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/wenerme/bbvm/libbbvm/asm"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
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

func TestParseCase(t *testing.T) {
	assert := assert.New(t)
	parseWholeDir("../testdata/case", assert)
}

func parseWholeDir(dir string, assert *assert.Assertions) {
	f, e := os.Open(dir)
	assert.NoError(e)
	fi, e := f.Readdir(-1)
	assert.NoError(e)
	for _, f := range fi {
		if f.IsDir() {
			parseWholeDir(path.Join(dir, f.Name()), assert)
		} else if strings.HasSuffix(f.Name(), ".basm") {
			testParse(path.Join(dir, f.Name()), assert)
		}

	}
}
func testParse(f string, assert *assert.Assertions) {
	fmt.Println("Parse ", f)
	b, e := ioutil.ReadFile(f)
	assert.NoError(e)
	p := &BBAsm{Buffer: string(b)}
	func() {
		defer func() {
			if e := recover(); e != nil {
				spew.Dump(p.stack)
				fmt.Println("-------------------- PARSE FAILED --------------------------")
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

		p.Init()
		if err := p.Parse(); err != nil {
			panic(err)
		}
		p.Execute()
	}()
}
