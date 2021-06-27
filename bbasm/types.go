package bbasm

import (
	"fmt"
)

// go:vet

type DataType uint8

const (
	DWORD DataType = iota
	WORD
	BYTE
	FLOAT
	INT
)

func (t DataType) String() string {
	switch t {
	case DWORD:
		return "DWORD"
	case WORD:
		return "WORD"
	case BYTE:
		return "BYTE"
	case FLOAT:
		return "FLOAT"
	case INT:
		return "INT"
	}
	return fmt.Sprintf("DataType(%v)", uint8(t))
}

type CompareType uint8

const (
	Z  CompareType = iota + 1 // Equal
	B                         // Blow
	BE                        // Blow or Equal
	A                         // Above
	AE                        // Above or Equal
	NZ                        // Not Equal
)

func (ct CompareType) IsMatch(b CompareType) bool {
	if ct == b {
		return true
	}
	switch ct {
	case A:
		{
			switch b {
			case AE, NZ:
				return true
			}
		}
	case B:
		{
			switch b {
			case BE, NZ:
				return true
			}
		}
	case Z:
		{
			switch b {
			case BE, AE:
				return true
			}
		}
	case NZ:
		{
			switch b {
			case B, A:
				return true
			}
		}
	case AE:
		{
			switch b {
			case A, Z:
				return true
			}
		}
	case BE:
		{
			switch b {
			case B, Z:
				return true
			}
		}
	}
	return false
}
func (ct CompareType) String() string {
	switch ct {
	case Z:
		return "Z"
	case B:
		return "B"
	case BE:
		return "BE"
	case A:
		return "A"
	case AE:
		return "AE"
	case NZ:
		return "NZ"
	}
	return fmt.Sprintf("CompareType(%v)", uint8(ct))
}

type CalculateType uint8

const (
	ADD CalculateType = iota
	SUB
	MUL
	DIV
	MOD
)

func (c CalculateType) String() string {
	switch c {
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case MUL:
		return "MUL"
	case DIV:
		return "DIV"
	case MOD:
		return "MOD"
	}
	return fmt.Sprintf("CalculateType(%v)", uint8(c))
}

type Opcode uint8

const (
	NOP Opcode = iota
	LD
	PUSH
	POP
	IN
	OUT
	JMP
	JPC
	CALL
	RET
	CMP
	CAL
	EXIT Opcode = 0xF
)

func (o Opcode) Len() int {
	switch o {
	case NOP, RET, EXIT:
		return 1
	case PUSH, POP, JMP, CALL:
		return 5
	case JPC:
		return 6
	case LD, IN, OUT, CMP, CAL:
		return 10
	}
	panic(fmt.Errorf("unexpected op(%v)", uint8(o)))
}
func (o Opcode) String() string {
	switch o {
	case NOP:
		return "NOP"
	case LD:
		return "LD"
	case PUSH:
		return "PUSH"
	case POP:
		return "POP"
	case IN:
		return "IN"
	case OUT:
		return "OUT"
	case JMP:
		return "JMP"
	case JPC:
		return "JPC"
	case CALL:
		return "CALL"
	case RET:
		return "RET"
	case CMP:
		return "CMP"
	case CAL:
		return "CAL"
	case EXIT:
		return "EXIT"
	}
	return fmt.Sprintf("Opcode(%v)", uint8(o))
}

type AddressMode uint8

const (
	AddressRegister         AddressMode = iota // 寄存器寻址
	AddressRegisterDeferred                    // 寄存器间接寻址
	AddressImmediate                           // 立即数
	AddressDirect                              // 直接寻址
)

func (v AddressMode) String() string {
	switch v {
	case AddressRegister:
		return "Register"
	case AddressRegisterDeferred:
		return "Register Deferred"
	case AddressImmediate:
		return "Immediate"
	case AddressDirect:
		return "Direct"
	}
	return fmt.Sprintf("AddressMode(%v)", uint8(v))
}

type RegisterType uint8

const (
	RP RegisterType = iota // 程序计数器,指令寻址寄存器
	RF                     // 标志寄存器,存储比较操作结果
	RS                     // 栈寄存器 - 空栈顶地址，指向的是下一个准备要压入数据的位置
	RB                     // 辅助栈寄存器 - 栈开始的地址（文件长度+2）
	R0                     // #0 寄存器
	R1                     // #1 寄存器
	R2                     // #2 寄存器
	R3                     // #3 寄存器
)

func (r RegisterType) String() string {
	switch r {
	case RP:
		return "RP"
	case RF:
		return "RF"
	case RS:
		return "RS"
	case RB:
		return "RB"
	case R0:
		return "R0"
	case R1:
		return "R1"
	case R2:
		return "R2"
	case R3:
		return "R3"
	}
	return fmt.Sprintf("RegisterType(%v)", uint8(r))
}
