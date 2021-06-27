package bbasm

import (
	"encoding/binary"
	"fmt"
)

type Inst struct {
	DataType      DataType
	CompareType   CompareType
	CalculateType CalculateType
	Opcode        Opcode
	A             Operand
	B             Operand

	Line    int
	Comment string
}

func (i *Inst) Len() int {
	return i.Opcode.Len()
}

func (i *Inst) SetComment(v string) {
	i.Comment = v
}
func (i *Inst) GetComment() string {
	return i.Comment
}
func (i *Inst) MarshalBinary() ([]byte, error) {
	data := make([]byte, i.Opcode.Len())
	data[0] = uint8(i.Opcode) << 4

	switch i.Opcode.Len() {
	case 1: // 无操作数
	case 5: // 一个操作数
		data[0] |= uint8(i.A.AddressMode)
		i.A.AddressMode = AddressMode(data[0] & 0xF)
		binary.LittleEndian.PutUint32(data[1:], uint32(i.A.V))
	case 10: // 两个操作数
		data[0] |= uint8(i.DataType)

		data[1] |= uint8(i.A.AddressMode<<2 | i.B.AddressMode)
		binary.LittleEndian.PutUint32(data[2:], uint32(i.A.V))
		binary.LittleEndian.PutUint32(data[6:], uint32(i.B.V))

		if i.Opcode == CAL {
			data[1] |= uint8(i.CalculateType) << 4
		}

	case 6: // JPC
		data[0] |= uint8(i.CompareType)
		data[1] |= uint8(i.A.AddressMode)
		binary.LittleEndian.PutUint32(data[2:], uint32(i.A.V))
	}
	return data, nil
}
func (i *Inst) UnmarshalBinary(data []byte) error {
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
		return ErrDataToShort.New("Data is not enough for inst need %v got %v", i.Opcode.Len(), len(data))
	}
	// clear
	i.A.AddressMode = 255
	i.B.AddressMode = 255

	switch i.Opcode.Len() {
	case 1: // 无操作数
	case 5: // 一个操作数
		i.A.AddressMode = AddressMode(data[0] & 0xF)
		i.A.V = int32(binary.LittleEndian.Uint32(data[1:]))
	case 10: // 两个操作数
		i.DataType = DataType(data[0] & 0xF)

		addrMode := data[1] & 0xF
		i.A.AddressMode = AddressMode(addrMode / 4)
		i.B.AddressMode = AddressMode(addrMode % 4)
		i.A.V = int32(binary.LittleEndian.Uint32(data[2:]))
		i.B.V = int32(binary.LittleEndian.Uint32(data[6:]))

		if i.Opcode == CAL {
			i.CalculateType = CalculateType(data[1] >> 4)
		}

	case 6: // JPC
		i.CompareType = CompareType(data[0] & 0xF)
		i.A.AddressMode = AddressMode(data[1] & 0xF)
		i.A.V = int32(binary.LittleEndian.Uint32(data[2:]))
	}
	return nil
}

func (i Inst) Assembly() (s string) {

	switch i.Opcode.Len() {
	case 1:
		s = fmt.Sprint(i.Opcode)
	case 5:
		s = fmt.Sprintf("%s %s", i.Opcode, i.A)
	case 10:
		switch i.Opcode {
		case CAL:
			s = fmt.Sprintf("%s %s %s %s, %s", i.Opcode, i.DataType, i.CalculateType, i.A, i.B)
		case LD:
			s = fmt.Sprintf("%s %s %s, %s", i.Opcode, i.DataType, i.A, i.B)
		case CMP:
			s = fmt.Sprintf("%s %s %s, %s", i.Opcode, i.CompareType, i.A, i.B)
		default:
			s = fmt.Sprintf("%s %s, %s", i.Opcode, i.A, i.B)
		}
	case 6:
		s = fmt.Sprintf("%s %s %s", i.Opcode, i.CompareType, i.A)
	default:
		panic(ErrWrongInst.New("Unknown opcode len"))
	}
	if i.Comment != "" {
		s += " ; " + i.Comment
	}
	return
}
