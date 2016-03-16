package parser

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/wenerme/bbvm/libbbvm/asm"
	"math"
	"strings"
)

type parser struct {
	line       int
	stack      []interface{}
	assemblies []asm.Assembly
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
	v := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return v
}

func (p *parser) PopAssembly() asm.Assembly {
	for i := len(p.stack) - 1; i >= 0; i-- {
		v := p.stack[i]

		if a, ok := v.(asm.Assembly); ok {
			var e error
			switch a.(type) {
			case *asm.Inst:
				a := a.(*asm.Inst)
				a.Line = p.line
				e = buildInst(a, p.stack[i+1:]...)
			case *asm.Label:
				a := a.(*asm.Label)
				a.Line = p.line
				e = buildLabel(a, p.stack[i+1:]...)
			case *asm.Comment:
				a := a.(*asm.Comment)
				a.Line = p.line
				e = buildComment(a, p.stack[i+1:]...)
			case *asm.PseudoBlock:
				a := a.(*asm.PseudoBlock)
				a.Line = p.line
				e = buildPseudoBlock(a, p.stack[i+1:]...)
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
		m.SetComment(c.(*asm.Comment).Content)
	} else {
		panic(errors.Errorf("Can not add comment try add %#v to %#v", c, a))
	}
}
func CreateOperand(val string, num bool, direct bool) *asm.Operand {
	op := &asm.Operand{}
	var am asm.AddressMode = math.MaxUint8

	if num {
		i, e := parseInt(val)
		if e != nil {
			panic(e)
		}
		op.Val = i
	} else {
		var r asm.RegisterType = math.MaxUint8
		switch strings.ToUpper(val) {
		case "RP":
			r = asm.REG_RP
		case "RF":
			r = asm.REG_RF
		case "RS":
			r = asm.REG_RS
		case "RB":
			r = asm.REG_RB
		case "R0":
			r = asm.REG_R0
		case "R1":
			r = asm.REG_R1
		case "R2":
			r = asm.REG_R2
		case "R3":
			r = asm.REG_R3
		default:
			// Symbol
			op.Symbol = val
		}
		if r != math.MaxUint8 {
			op.Val = int32(r)
			if direct {
				am = asm.AM_REGISTER
			} else {
				am = asm.AM_REGISTER_DEFERRED
			}
		}
	}
	if am == math.MaxUint8 {
		if direct {
			am = asm.AM_IMMEDIATE
		} else {
			am = asm.AM_DIRECT
		}
	}
	op.AddressMode = am
	return op
}
func (p *parser) PushInst(op asm.Opcode) {
	p.Push(&asm.Inst{Opcode: op})
}

func lookup(t interface{}, v string) interface{} {
	ret := asm.Lookup(t, strings.ToUpper(strings.Trim(v, " \t\r")))
	if ret == nil {
		panic(errors.Errorf("Can not lookup '%v' for %T", t))
	}
	return ret
}
