package parser

import (
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/wenerme/bbvm/libbbvm/asm"
	"strings"
)

func buildInst(a *asm.Inst, v ...interface{}) error {
	n := 0
	switch a.Opcode {
	case asm.OP_JMP, asm.OP_PUSH, asm.OP_POP, asm.OP_CALL:
		n = 1
	case asm.OP_IN, asm.OP_OUT, asm.OP_JPC:
		n = 2
	case asm.OP_LD, asm.OP_CMP:
		n = 3
	case asm.OP_CAL:
		n = 4
	}
	if len(v) != n {
		return errors.Errorf("Expecte %v for %v got %v", n, a.Opcode, len(v))
	}

	operands := 0
	for _, o := range v {
		switch o.(type) {
		case *asm.Operand:
			switch operands {
			case 0:
				a.A = *(o.(*asm.Operand))
			case 1:
				a.B = *(o.(*asm.Operand))
			default:
				return errors.New("To many operands")
			}
			operands++
		case asm.DataType:
			a.DataType = o.(asm.DataType)
		case asm.CompareType:
			a.CompareType = o.(asm.CompareType)
		case asm.CalculateType:
			a.CalculateType = o.(asm.CalculateType)
		default:
			return errors.Errorf("Unexpected value %#v", o)
		}
	}
	return nil
}

func buildComment(a *asm.Comment, v ...interface{}) error {
	a.Content = strings.Trim(v[0].(string), " \t")
	return nil
}

func buildLabel(a *asm.Label, v ...interface{}) error {
	a.Name = strings.Trim(v[0].(string), " \t")
	return nil
}

func buildPseudoBlock(a *asm.PseudoBlock, v ...interface{}) error {
	if len(v) != 2 {
		return errors.Errorf(".BLOCK size byte got %v", v)
	}
	a.Size = v[0].(int)
	i64 := v[1].(int)
	if int32(i64%0xff) != int32(i64) {
		log.Warnf("Convert %v to byte", i64)
	}
	a.Byte = byte(i64)
	return nil
}

func buildPseudoData(a *asm.PseudoData, v ...interface{}) error {
	a.Label = v[0].(string)
	a.Values = make([]asm.PseudoDataValue, len(v)-1)
	for i, v := range v[1:] {
		if v, ok := v.(asm.PseudoDataValue); ok {
			a.Values[i] = v
		} else {
			return errors.Errorf("Require PseudoDataValue got %#v", v)
		}
	}
	return nil
}
