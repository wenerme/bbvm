package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
	. "github.com/wenerme/bbvm/libbbvm/asm"
)

func (v *vm) Pop() int {
	v.rs.Set(v.rs.Get() - 4)
	return Codec.Int(v.mem[v.rs.Get():])
}
func (v *vm) Push(i int) {
	Codec.PutInt(v.mem[v.rs.Get():], i)
	v.rs.Set(v.rs.Get() + 4)
}
func (v *vm) Call(addr int) {
	v.Push(v.rp.Get() + v.inst.Opcode.Len())
	v.Jump(addr)
}
func (v *vm) Jump(addr int) {
	v.rp.Set(addr)
}

func (v *vm) GetInt(addr int) int {
	return int(int32(binary.LittleEndian.Uint32(v.mem[addr:]))) // must convert to int32 first
}

func (v *vm) SetInt(addr int, i int) {
	binary.LittleEndian.PutUint32(v.mem[addr:], uint32(i))
}

func (v *vm) MustGetStr(addr int) string {
	s, err := v.GetStr(addr)
	if err != nil {
		log.Error(err.Error())
	}
	return s
}
func (v *vm) GetStr(addr int) (string, error) {
	if addr < 0 {
		if r := v.strPool.Get(addr); r != nil {
			if s, ok := r.Get().(string); ok {
				return s, nil
			}
			return "", errors.New("Resource is not string " + fmt.Sprint(r.Get()))
		}
		return "", errors.New("String resource not exits " + fmt.Sprint(addr))
	}
	if len(v.mem) < addr {
		return "", errors.New(fmt.Sprintf("Address out of memory range: %d > %d", addr, len(v.mem)))
	}
	end := addr
	for ; v.mem[end] != 0; end += 1 {
	}
	if end == addr {
		return "", nil
	}
	// Make sure the origin memory is not changed
	return "" + string(v.mem[addr:end]), nil
}
func (v *vm) Proc() {
	i := &v.inst

	log.Debug("%s", v.Report())
	switch v.inst.Opcode {
	case OP_EXIT:
		v.Exit()
	case OP_NOP:
	// NOP
	case OP_CAL:
		i.A.Set(calculate(i.A.Get(), i.B.Get(), i.CalculateType, i.DataType))
	case OP_CALL:
		v.Call(i.A.Get())
	case OP_JMP:
		v.Jump(i.A.Get())
	case OP_IN:
		h := v.In(i.A.Get(), i.B.Get())
		if h == nil {
			log.Error("Can not handle %s", i)
		} else {
			h(i)
		}
	case OP_OUT:
		h := v.Out(i.A.Get(), i.B.Get())
		if h == nil {
			log.Error("Can not handle %s", i)
		} else {
			h(i)
		}
	case OP_RET:
		v.Jump(v.Pop())
	case OP_CMP:
		{
			//var oc float32
			//if i.DataType == T_FLOAT {
			//	oc = i.A.Float32() - i.B.Float32()
			//}else {
			//	oc = float32(i.A.Get() - i.B.Get())
			//}
			//switch {
			//case oc > 0:v.rf.Set(int(COM_GT))
			//case oc < 0:v.rf.Set(int(COM_LT))
			//default:v.rf.Set(int(COM_EQ))
			//}
		}
	case OP_JPC:
		if i.CompareType.IsMatch(CompareType(v.rf.Get())) {
			v.Jump(i.A.Get())
		}
	case OP_PUSH:
		v.Push(i.A.Get())
	case OP_POP:
		i.A.Set(v.Pop())
	case OP_LD:
		// 没考虑类型
		i.A.Set(i.B.Get())
	}
}
