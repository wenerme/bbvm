package parser

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/wenerme/bbvm/bbasm"
	"math"
	"strings"
)

type parser struct {
	line       int
	stack      []interface{}
	assemblies []Assembly
	symbols    map[string]interface{}
}

func (p *parser) Push(v interface{}) {
	fmt.Printf("PUSH %#v\n", v)
	if v == nil {
		panic(errors.Errorf("Can not push nil"))
	}
	p.stack = append(p.stack, v)
}
func (p *parser) Pop() interface{} {
	if len(p.stack) == 0 {
		panic(errors.Errorf("Stack underflow"))
	}
	v := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return v
}

func (p *parser) AddAssembly() Assembly {
	for i := len(p.stack) - 1; i >= 0; i-- {
		v := p.stack[i]

		if a, ok := v.(Assembly); ok {
			var e error
			switch a.(type) {
			case *bbasm.Inst:
				a := a.(*bbasm.Inst)
				//a.Line = p.line
				e = buildInst(a, p.stack[i+1:]...)
			case *Label:
				a := a.(*Label)
				//a.Line = p.line
				e = buildLabel(a, p.stack[i+1:]...)
			case *Comment:
				a := a.(*Comment)
				//a.Line = p.line
				e = buildComment(a, p.stack[i+1:]...)
			case *PseudoBlock:
				a := a.(*PseudoBlock)
				//a.Line = p.line
				e = buildPseudoBlock(a, p.stack[i+1:]...)
			case *PseudoData:
				a := a.(*PseudoData)
				e = buildPseudoData(a, p.stack[i+1:]...)
			}
			//e := build(a, p.stack[i + 1:]...)
			if e != nil {
				panic(errors.Trace(e))
			}

			if m, ok := a.(interface {
				SetLine(int)
			}); ok {
				m.SetLine(p.line)
			}

			p.stack = p.stack[0:i]
			p.assemblies = append(p.assemblies, a)
			return a
		}
	}
	panic(errors.Errorf("No assembly avalible line %v", p.line))
}

func (p *parser) AddComment() {
	a := p.assemblies[len(p.assemblies)-2]
	c := p.assemblies[len(p.assemblies)-1]
	// Pop Comment out
	p.assemblies = p.assemblies[:len(p.assemblies)-1]
	if m, ok := a.(interface {
		SetComment(string)
	}); ok {
		m.SetComment(c.(*Comment).Content)
	} else {
		panic(errors.Errorf("Can not add comment try add %#v to %#v", c, a))
	}
}
func (p *parser) PushInst(op bbasm.Opcode) {
	p.Push(&bbasm.Inst{Opcode: op})
}

func (p *parser) AddOperand(direct bool) {
	v := p.Pop()
	op := &bbasm.Operand{}
	var am bbasm.AddressMode = math.MaxUint8

	switch v.(type) {
	case string:
		r := Lookup(bbasm.R0, v.(string))
		if r != nil {
			op.V = int32(r.(bbasm.RegisterType))
			if direct {
				am = bbasm.AddressRegister
			} else {
				am = bbasm.AddressRegisterDeferred
			}
		} else {
			op.Symbol = v.(string)
		}
	case int:
		op.V = int32(v.(int))
	default:
		panic(errors.Errorf("Can not add operand with %#v", v))
	}

	if am == math.MaxUint8 {
		if direct {
			am = bbasm.AddressImmediate
		} else {
			am = bbasm.AddressDirect
		}
	}
	op.AddressMode = am
	p.Push(op)
}

func (p *parser) AddPseudoDataValue() {
	v := p.Pop()
	v, e := createPseudoDataValue(v)
	if e != nil {
		panic(errors.Trace(e))
	}
	p.Push(v)
}

func (p *parser) AddInteger() {
	v := p.Pop()
	i, e := parseInt(v.(string))
	if e != nil {
		panic(errors.Trace(e))
	}
	p.Push(int(i))
}

func lookup(t interface{}, v string) interface{} {
	ret := Lookup(t, strings.ToUpper(strings.Trim(v, " \t\r")))
	if ret == nil {
		panic(errors.Errorf("Can not lookup '%v' for %T", v, t))
	}
	return ret
}
