package parser

import (
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/wenerme/bbvm/bbasm"
	"strings"
)

func buildInst(a *bbasm.Inst, v ...interface{}) error {
	n := 0
	switch a.Opcode {
	case bbasm.JMP, bbasm.PUSH, bbasm.POP, bbasm.CALL:
		n = 1
	case bbasm.IN, bbasm.OUT, bbasm.JPC:
		n = 2
	case bbasm.LD, bbasm.CMP:
		n = 3
	case bbasm.CAL:
		n = 4
	}
	if len(v) != n {
		return errors.Errorf("Expecte %v for %v got %v", n, a.Opcode, len(v))
	}

	operands := 0
	for _, o := range v {
		switch o.(type) {
		case *bbasm.Operand:
			switch operands {
			case 0:
				a.A = *(o.(*bbasm.Operand))
			case 1:
				a.B = *(o.(*bbasm.Operand))
			default:
				return errors.New("To many operands")
			}
			operands++
		case bbasm.DataType:
			a.DataType = o.(bbasm.DataType)
		case bbasm.CompareType:
			a.CompareType = o.(bbasm.CompareType)
		case bbasm.CalculateType:
			a.CalculateType = o.(bbasm.CalculateType)
		default:
			return errors.Errorf("Unexpected value %#v", o)
		}
	}
	return nil
}

func buildComment(a *Comment, v ...interface{}) error {
	a.Content = strings.Trim(v[0].(string), " \t")
	return nil
}

func buildLabel(a *Label, v ...interface{}) error {
	a.Name = strings.Trim(v[0].(string), " \t")
	return nil
}

func buildPseudoBlock(a *PseudoBlock, v ...interface{}) error {
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

func buildPseudoData(a *PseudoData, v ...interface{}) error {
	a.Label = v[0].(string)
	a.Values = make([]PseudoDataValue, len(v)-1)
	for i, v := range v[1:] {
		if v, ok := v.(PseudoDataValue); ok {
			a.Values[i] = v
		} else {
			return errors.Errorf("Require PseudoDataValue got %#v", v)
		}
	}
	return nil
}
