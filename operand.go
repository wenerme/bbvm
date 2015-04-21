package bbvm
import (
	"fmt"
	"math"
	"errors"
	"encoding/binary"
)



type Inst struct {
	DataType DataType
	CompareType CompareType
	CalculateType CalculateType
	Opcode Opcode
	A Operand
	B Operand

	VM *vm
}
type Operand struct {
	Val uint32
	AddrMode AddressMode
	Vm *vm
}
func (o *Operand)Get() int {
	switch (o.AddrMode) {
		case AM_REGISTER:
		return o.Vm.Register(RegisterType(o.Val)).Get();
		case AM_REGISTER_DEFERRED:
		return o.Vm.GetInt(o.Vm.Register(RegisterType(o.Val)).Get());
		case AM_IMMEDIATE:
		return int(o.Val);
		case AM_DIRECT:
		return o.Vm.GetInt(int(o.Val));
	}
	panic("Unexcepted")
}
func (o *Operand)Uint() uint {
	return uint(o.Get())
}
func (o *Operand)Float() float32 {
	return math.Float32frombits(uint32(o.Get()))
}
func (o *Operand)SetFloat(v float32) {
	o.Set(int(math.Float32bits(v)))
}
func (o *Operand)StrRes() Res {
	return o.Vm.StrPool().Get(o.Get())
}
func (o *Operand)Str() string {
	if s, ok := o.Vm.GetStr(o.Get()); ok {
		return s
	}else {
		log.Error("Operand string res %d not exists", o.Get())
		return ""
	}
}
func (o *Operand)SetStr(v string) {
	if r := o.StrRes(); r!= nil {
		o.StrRes().Set(v)
	}else {
		log.Error("Operand string res %d not exists", o.Get())
	}
}
func (o *Operand)Set(i int) {
	switch (o.AddrMode) {
		case AM_REGISTER:
		o.Vm.Register(RegisterType(o.Val)).Set(i);
		case AM_REGISTER_DEFERRED:
		o.Vm.SetInt(o.Vm.Register(RegisterType(o.Val)).Get(), i);
		case AM_IMMEDIATE:
		log.Error("ERR Set a IMMEDIATE operand")
		case AM_DIRECT:
		o.Vm.SetInt(int(o.Val), i);
		default:
		log.Error("ERR Unknown address mode when set operand")
	}
}
func (o Operand)String() string {
	switch o.AddrMode{
		case AM_REGISTER:
		return fmt.Sprintf("%s", RegisterType(o.Val))
		case AM_REGISTER_DEFERRED:
		return fmt.Sprintf("[ %s ]", RegisterType(o.Val))
		case AM_IMMEDIATE:
		return fmt.Sprintf("%d", o.Val)
		case AM_DIRECT:
		return fmt.Sprintf("[ %d ]", o.Val)
	}
	return "ERR Unknown address mode"
}
var ErrDataToShort = errors.New("Data to short to UnmarshalBinary")

func (i *Inst)MarshalBinary() ([]byte, error) {
	data := make([]byte, i.Opcode.Len())
	data[0] = uint8(i.Opcode) << 4

	switch i.Opcode.Len() {
		case 1: // 无操作数
		case 5: // 一个操作数
		data[0] |= uint8(i.A.AddrMode)
		i.A.AddrMode = AddressMode(data[0] & 0xF)
		binary.LittleEndian.PutUint32(data[1:], i.A.Val)
		case 10: // 两个操作数
		data[0] |= uint8(i.DataType)

		data[1] |= uint8(i.A.AddrMode << 2 | i.B.AddrMode)
		binary.LittleEndian.PutUint32(data[2:], i.A.Val)
		binary.LittleEndian.PutUint32(data[6:], i.B.Val)

		if i.Opcode == OP_CAL {
			data[1] |= uint8(i.CalculateType) << 4
		}

		case 6: // JPC
		data[0] |= uint8(i.CompareType)
		data[1] |= uint8(i.A.AddrMode)
		binary.LittleEndian.PutUint32(data[2:], i.A.Val)
	}
	return data, nil
}
func (i *Inst)UnmarshalBinary(data []byte) error {
	/*
指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数 + 第二个操作数
0x 0       0         0           0        00000000     00000000

无操作数 1byte
指令码 + 无用
0x 0       0
一个操作数 5byte
指令码 + 寻址方式 + 第一个操作数
0x 0       0        00000000
两个操作数 10byte
指令码 + 数据类型 + 保留字节 + 寻址方式 + 第一个操作数 + 第二个操作数
0x 0       0         0        0        00000000     00000000
JPC指令 6byte
指令码 + 比较操作 + 保留字节 + 寻址方式 + 第一个操作数
0x 0       0         0        0        00000000
	*/
	i.Opcode = Opcode(data[0] >> 4)
	if i.Opcode.Len() > len(data) {
		return ErrDataToShort
	}
	switch i.Opcode.Len() {
		case 1: // 无操作数
		case 5: // 一个操作数
		i.A.AddrMode = AddressMode(data[0] & 0xF)
		i.A.Val = binary.LittleEndian.Uint32(data[1:])
		case 10: // 两个操作数
		i.DataType = DataType(data[0] & 0xF)

		addrMode := data[1] & 0xF
		i.A.AddrMode = AddressMode(addrMode/4)
		i.B.AddrMode = AddressMode(addrMode%4)
		i.A.Val = binary.LittleEndian.Uint32(data[2:])
		i.B.Val = binary.LittleEndian.Uint32(data[6:])

		if i.Opcode == OP_CAL {
			i.CalculateType = CalculateType(data[1] >> 4)
		}

		case 6: // JPC
		i.CompareType = CompareType(data[0] & 0xF)
		i.A.AddrMode = AddressMode(data[1] & 0xF)
		i.A.Val = binary.LittleEndian.Uint32(data[2:])
	}
	return nil
}

func (i Inst)String() string {
	switch i.Opcode.Len(){
		case 1:
		return fmt.Sprint(i.Opcode)
		case 5:
		return fmt.Sprintf("%s %s", i.Opcode, i.A)
		case 10:
		switch i.Opcode{
			case OP_CAL:
			return fmt.Sprintf("%s %s %s %s, %s", i.Opcode, i.DataType, i.CalculateType, i.A, i.B)
			case OP_LD:
			return fmt.Sprintf("%s %s %s, %s", i.Opcode, i.DataType, i.A, i.B)
			case OP_CMP:
			return fmt.Sprintf("%s %s %s, %s", i.Opcode, i.CompareType, i.A, i.B)
		}
		return fmt.Sprintf("%s %s %s", i.Opcode, i.A, i.B)
		case 6:
		return fmt.Sprintf("%s %s %s", i.Opcode, i.CompareType, i.A)
	}
	return "ERR Unknown opcode len"
}