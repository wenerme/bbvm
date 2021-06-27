package bbasm

import (
	"fmt"
	"github.com/spacemonkeygo/errors"
	"math"
)

var ErrDataToShort = errors.NewClass("ErrDataToShort")
var ErrWrongInst = errors.NewClass("ErrWrongInst")

type Operand struct {
	V           int32 // Address or Immediate value
	AddressMode AddressMode
	RT          Runtime
	Symbol      string // Symbol associated to this operand - run from assembly
}

func (o *Operand) Get() int {
	switch o.AddressMode {
	case AddressRegister:
		return o.RT.Register(RegisterType(o.V)).Get()
	case AddressRegisterDeferred:
		return o.RT.GetInt(o.RT.Register(RegisterType(o.V)).Get())
	case AddressImmediate:
		return int(int32(o.V)) // must convert to int32 first
	case AddressDirect:
		return o.RT.GetInt(int(o.V))
	}
	panic(fmt.Errorf("invalid %v", o.AddressMode.String()))
}
func (o *Operand) Float() float32 {
	return math.Float32frombits(uint32(o.Get()))
}
func (o *Operand) SetFloat(v float32) {
	o.Set(int(int32(math.Float32bits(v)))) // Must conver to int32 first
}

//func (o *Operand) Str() string {
//	//if s, err := o.VM.GetStr(o.Get()); err == nil {
//	//	return s
//	//}else {
//	//	//log.Error("Operand string res %d not exists: %s", o.Get(), err.Error())
//	//	return ""
//	//}
//	// FIXME
//	return ""
//}
//func (o *Operand) SetStr(v string) {
//	//if r := o.StrRes(); r != nil {
//	//	o.StrRes().Set(v)
//	//}else {
//	//	//log.Error("Operand string res %d not exists", o.Get())
//	//}
//	// FIXME
//}

func (o *Operand) Set(i int) {
	switch o.AddressMode {
	case AddressRegister:
		o.RT.Register(RegisterType(o.V)).Set(i)
	case AddressRegisterDeferred:
		o.RT.SetInt(o.RT.Register(RegisterType(o.V)).Get(), i)
	case AddressImmediate:
		panic("ERR Set a IMMEDIATE operand")
	case AddressDirect:
		o.RT.SetInt(int(o.V), i)
	default:
		panic("ERR Unknown address mode when set operand")
	}
}

func (o Operand) String() string {
	switch o.AddressMode {
	case AddressRegister:
		return fmt.Sprintf("%s", RegisterType(o.V))
	case AddressRegisterDeferred:
		return fmt.Sprintf("[ %s ]", RegisterType(o.V))
	case AddressImmediate:
		if o.Symbol == "" {
			return fmt.Sprintf("%v", o.V)
		} else {
			return fmt.Sprintf("%v", o.Symbol)
		}
	case AddressDirect:
		if o.Symbol == "" {
			return fmt.Sprintf("[ %v ]", o.V)
		} else {
			return fmt.Sprintf("[ %v ]", o.Symbol)
		}
	}
	return "ERR Unknown address mode " + o.AddressMode.String()
}
