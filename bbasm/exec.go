package bbasm

import (
	"context"
	"fmt"
	"math"
)

func Execute(ctx context.Context, rt Runtime, inst *Inst) {
	a := inst.A
	b := inst.B

	switch inst.Opcode {
	case NOP:
	// Tick
	case LD:
		a.Set(b.Get())
	case PUSH:
		rt.Push(a.Get())
	case POP:
		a.Set(rt.Pop())
	case IN:
		rt.In(ctx,a.Get(), b.Get())
	case OUT:
		rt.Out(ctx,a.Get(), b.Get())
	case JMP:
		rt.Jump(a.Get())
	case JPC:
		if inst.CompareType.IsMatch(CompareType(rt.Register(RF).Get())) {
			rt.Jump(a.Get())
		}
	case CALL:
		// 压入下一条指令的位置
		rt.Push(rt.Register(RP).Get() + inst.Opcode.Len())
		rt.Jump(a.Get())
	case RET:
		rt.Jump(rt.Pop())
	case CMP:
		rt.Register(RF).Set(Compare(a, b, inst.DataType))
	case CAL:
		a.Set(Calculate(a.Get(), b.Get(), inst.CalculateType, inst.DataType))
	case EXIT:
		rt.Exit()
	}
}

func Compare(a Operand, b Operand, d DataType) int {
	v := float32(0)
	if d == FLOAT {
		v = a.Float() - b.Float()
	} else {
		v = float32(a.Get() - b.Get())
	}
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	default:
		return 0
	}
}

func Calculate(a int, b int, t CalculateType, d DataType) int {
	var oa, ob, oc float64
	if d == FLOAT {
		oa = float64(math.Float32frombits(uint32(a)))
		ob = float64(math.Float32frombits(uint32(b)))
	} else {
		oa = float64(a)
		ob = float64(b)
	}
	switch t {
	case ADD:
		oc = oa + ob
	case SUB:
		oc = oa - ob
	case MUL:
		oc = oa * ob
	case DIV:
		oc = oa / ob
	case MOD:
		oc = float64(int(oa) % int(ob))
	default:
		panic(fmt.Errorf("invalid calc %v", t.String()))
	}
	var ret int
	switch d {
	case DWORD:
		ret = int(oc)
	case WORD:
		ret = int(int16(oc))
	case BYTE:
		ret = int(byte(oc))
	case FLOAT:
		ret = int(math.Float32bits(float32(oc)))
	case INT:
		ret = int(oc)
	default:
		panic(fmt.Errorf("invalid datatype %v", t.String()))
	}
	return ret
}
