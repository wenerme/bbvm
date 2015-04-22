package bbvm

type DataType uint8
const (
	T_DWORD DataType = iota
	T_WORD
	T_BYTE
	T_FLOAT
	T_INT
)

func (t DataType)String() string {
	switch t{
		case T_DWORD: return "DWORD"
		case T_WORD: return "WORD"
		case T_BYTE: return "BYTE"
		case T_FLOAT: return "FLOAT"
		case T_INT: return "INT"
	}
	panic("UNKNOWN DataType")
}

type CompareType uint8
const (
	COM_EQ CompareType = iota+1
	COM_LT
	COM_LE
	COM_GT
	COM_GE
	COM_NE
)

func (a CompareType)IsMatch(b CompareType) bool {
	if a == b {
		return true
	}
	switch a{
		case COM_GT: {
			switch b{case COM_GE, COM_NE: return true}
		}
		case COM_LT: {
			switch b{case COM_LE, COM_NE: return true}
		}
		case COM_EQ: {
			switch b{case COM_LE, COM_GE: return true}
		}
		case COM_NE: {
			switch b{case COM_LT, COM_GT: return true}
		}
		case COM_GE: {
			switch b{case COM_GT, COM_EQ: return true}
		}
		case COM_LE: {
			switch b{case COM_LT, COM_EQ: return true}
		}
	}
	return false
}
func (c CompareType)String() string {
	switch c{
		case COM_EQ: return "Z"
		case COM_LT: return "B"
		case COM_LE: return "BE"
		case COM_GT: return "A"
		case COM_GE: return "AE"
		case COM_NE: return "NZ"
	}
	return "UNKNOWN CompareType"
}

type CalculateType uint8
const (
	CAL_ADD CalculateType = iota
	CAL_SUB
	CAL_MUL
	CAL_DIV
	CAL_MOD
)
func (c CalculateType)String() string {
	switch c{
		case CAL_ADD: return "ADD"
		case CAL_SUB: return "SUB"
		case CAL_MUL: return "MUL"
		case CAL_DIV: return "DIV"
		case CAL_MOD: return "MOD"
	}
	return "UNKNOWN CalculateType"
}
type Opcode uint8
const (
	OP_NOP  Opcode = iota
	OP_LD
	OP_PUSH
	OP_POP
	OP_IN
	OP_OUT
	OP_JMP
	OP_JPC
	OP_CALL
	OP_RET
	OP_CMP
	OP_CAL
	OP_EXIT Opcode = 0xF
)
func (o Opcode)Len() int {
	switch o {
		case OP_NOP, OP_RET, OP_EXIT: return 1
		case OP_PUSH, OP_POP, OP_JMP, OP_CALL: return 5
		case OP_JPC: return 6
		case OP_LD, OP_IN, OP_OUT, OP_CMP, OP_CAL: return 10
	}
	panic("Unexcepted")
}
func (o Opcode)String() string {
	switch o {
		case OP_NOP: return "NOP"
		case OP_LD: return "LD"
		case OP_PUSH: return "PUSH"
		case OP_POP: return "POP"
		case OP_IN: return "IN"
		case OP_OUT: return "OUT"
		case OP_JMP: return "JMP"
		case OP_JPC: return "JPC"
		case OP_CALL: return "CALL"
		case OP_RET: return "RET"
		case OP_CMP: return "CMP"
		case OP_CAL: return "CAL"
		case OP_EXIT: return "EXIT"
	}
	return "UNKNOWN Opcode"
}

type AddressMode uint8
const (
	AM_REGISTER AddressMode = iota
	AM_REGISTER_DEFERRED
	AM_IMMEDIATE
	AM_DIRECT
)

type RegisterType uint8
const (
// 程序计数器,指令寻址寄存器
	REG_RP RegisterType = iota
// 标志寄存器,存储比较操作结果
	REG_RF
// 栈寄存器	
// 空栈顶地址，指向的是下一个准备要压入数据的位置
	REG_RS
// 辅助栈寄存器
// 栈开始的地址（文件长度+2）
	REG_RB
// #0 寄存器
	REG_R0
// #1 寄存器
	REG_R1
// #2 寄存器
	REG_R2
// #3 寄存器
	REG_R3
)

func (r RegisterType)String() string {
	switch r{
		case REG_RP: return "RP"
		case REG_RF: return "RF"
		case REG_RS: return "RS"
		case REG_RB: return "RB"
		case REG_R0: return "R0"
		case REG_R1: return "R1"
		case REG_R2: return "R2"
		case REG_R3: return "R3"
	}
	return "UNKNOW RegisterType"
}


