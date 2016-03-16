package asm

import (
	"fmt"
	"github.com/spacemonkeygo/errors"
)

var ErrDataToShort = errors.NewClass("ErrDataToShort")
var ErrWrongInst = errors.NewClass("ErrWrongInst")

type Operand struct {
	// Address or Immediate value
	Val         int32
	AddressMode AddressMode
	// Symbol associated to this operand
	Symbol string
}

//func (o *Operand)Get() int {
//	switch (o.AddrMode) {
//	case AM_REGISTER:
//		return o.Register.Get();
//	case AM_REGISTER_DEFERRED:
//		return o.VM.GetInt(o.VM.Register(RegisterType(o.Val)).Get());
//	case AM_IMMEDIATE:
//		return int(int32(o.Val));// must convert to int32 first
//	case AM_DIRECT:
//		return o.VM.GetInt(int(o.Val));
//	}
//	panic("Unexcepted")
//}
//func (o *Operand)Float32() float32 {
//	return math.Float32frombits(uint32(o.Get()))
//}
//func (o *Operand)SetFloat32(v float32) {
//	o.Set(int(int32(math.Float32bits(v))))// Must conver to int32 first
//}
//func (o *Operand)Str() string {
//	//if s, err := o.VM.GetStr(o.Get()); err == nil {
//	//	return s
//	//}else {
//	//	//log.Error("Operand string res %d not exists: %s", o.Get(), err.Error())
//	//	return ""
//	//}
//	// FIXME
//	return ""
//}
//func (o *Operand)SetStr(v string) {
//	//if r := o.StrRes(); r != nil {
//	//	o.StrRes().Set(v)
//	//}else {
//	//	//log.Error("Operand string res %d not exists", o.Get())
//	//}
//	// FIXME
//}
//
//func (o *Operand)Set(i int) {
//	switch (o.AddrMode) {
//	case AM_REGISTER:
//		o.VM.Register(RegisterType(o.Val)).Set(i);
//	case AM_REGISTER_DEFERRED:
//		o.VM.SetInt(o.VM.Register(RegisterType(o.Val)).Get(), i);
//	case AM_IMMEDIATE:
//		panic("ERR Set a IMMEDIATE operand")
//	case AM_DIRECT:
//		o.VM.SetInt(int(o.Val), i);
//	default:
//		panic("ERR Unknown address mode when set operand")
//	}
//}

func (o Operand) String() string {
	switch o.AddressMode {
	case AM_REGISTER:
		return fmt.Sprintf("%s", RegisterType(o.Val))
	case AM_REGISTER_DEFERRED:
		return fmt.Sprintf("[ %s ]", RegisterType(o.Val))
	case AM_IMMEDIATE:
		if o.Symbol == "" {
			return fmt.Sprintf("%v", o.Val)
		} else {
			return fmt.Sprintf("%v", o.Symbol)
		}
	case AM_DIRECT:
		if o.Symbol == "" {
			return fmt.Sprintf("[ %v ]", o.Val)
		} else {
			return fmt.Sprintf("[ %v ]", o.Symbol)
		}
	}
	return "ERR Unknown address mode " + o.AddressMode.String()
}
