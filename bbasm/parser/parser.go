//go:generate peg -switch bbasm.peg

package parser

import (
	"fmt"
)

/*
A parser has a context
	contain
		labels
		assemblies
this context is not belong to parser
	asm should hold this context
will used in
	linking
		resolve label reference
	debugging
		resolve label name
when doing REPL, this context is not enough, the running vm hold
	memory
	stack
	register
	system calling handler

asm
	only define the assemblies and types
parser
	parse a line to an assembly
vm
	memory
	stack
	register
	system calling handler

Problem
	who load the file
	who hold the label table
	who group them together

While REPL
	every line will parse
	resolve label reference
	execute by vm

While running
	file will parse
	add to context
*/

func Compile(s string) ([]byte, error) {
	assemblies, err := Parse(s)

	if err != nil {
		return nil, err
	}
	asm := &Assembler{Lines: assemblies}
	return asm.Assemble()
}

func ParseLine(s string) (Assembly, error) {
	assemblies, err := Parse(s)
	if err != nil {
		return nil, err
	}
	if len(assemblies) != 1 {
		return nil, fmt.Errorf("expect one line got: %v", len(assemblies))
	}
	return assemblies[0], nil
}
func Parse(s string) ([]Assembly, error) {
	p := &BBAsm{Buffer: s}
	p.Init()
	if err := p.Parse(); err != nil {
		return nil, err
	}
	// fixme
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("err %v", e)
				}
			}
		}()
		p.Execute()
	}()
	return p.assemblies, err
}
