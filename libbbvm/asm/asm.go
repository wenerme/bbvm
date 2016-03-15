package asm

import (
	"fmt"
)

type Assembly interface {
	//fmt.Stringer
	//encoding.BinaryMarshaler
	//encoding.BinaryUnmarshaler

	// Byte code length of this asm
	Len() int
	// To assembly string
	Assembly() string
}
type Symbol struct {
	Name    string
	Address int
	// Referenced this symbol's address
	Reference []func(int)
}

func commentString(Comment string) string {
	if Comment == "" {
		return ""
	}
	return "; " + Comment
}

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

func (a *Inst) Len() int {
	return a.Opcode.Len()
}

func (a *Inst) SetComment(v string) {
	a.Comment = v
}
func (a *Inst) GetComment() string {
	return a.Comment
}
func (a *Label) GetComment() string {
	return a.Comment
}
func (a *Label) SetComment(v string) {
	a.Comment = v
}
func (a *PseudoBlock) GetComment() string {
	return a.Comment
}
func (a *PseudoBlock) SetComment(v string) {
	a.Comment = v
}

type Comment struct {
	Content string
	Line    int
}

func (*Comment) Len() int {
	return 0
}
func (a *Comment) Assembly() string {
	return "; " + a.Content
}

type Label struct {
	Name string

	Line    int
	Comment string
}

func (a *Label) Len() int {
	return 0
}
func (a *Label) Assembly() string {
	return fmt.Sprintf(":%s %s", a.Name, commentString(a.Comment))
}

type PseudoData struct {
}

type PseudoBlock struct {
	Size int
	Byte byte

	Line    int
	Comment string
}

func (a *PseudoBlock) Len() int {
	return a.Size
}
func (a PseudoBlock) Assembly() string {
	return fmt.Sprintf(".BLOCK %v %v %v", a.Size, a.Byte, commentString(a.Comment))
}
func (a *PseudoBlock) MarshalBinary() (data []byte, err error) {
	data = make([]byte, a.Size)
	for i := range data {
		data[i] = a.Byte
	}
	return
}
