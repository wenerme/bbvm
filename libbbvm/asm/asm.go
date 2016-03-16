package asm

import (
	"bytes"
	"encoding"
	"fmt"
)

type Assembly interface {
	//fmt.Stringer
	//encoding.BinaryMarshaler
	//encoding.BinaryUnmarshaler
	SetComment(string)
	GetComment() string

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
	return " ; " + Comment
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
	return fmt.Sprintf(":%s%s", a.Name, commentString(a.Comment))
}

type PseudoData struct {
	Label  string
	Values []PseudoDataValue

	Line    int
	Comment string
}

type PseudoDataValue interface {
	encoding.BinaryMarshaler
	// Byte length
	Len() int
	// Assembly represent
	Assembly() string
}

func (a *PseudoData) Len() int {
	// TODO
	return 0
}

func (a PseudoData) Assembly() string {
	buf := bytes.NewBufferString("DATA ")
	buf.WriteString(a.Label)
	buf.WriteRune(' ')
	switch len(a.Values) {
	case 0:
	// FIXME No data, what should I do
	case 1:
		buf.WriteString(a.Values[0].Assembly())
	default:
		buf.WriteString(a.Values[0].Assembly())
		for _, v := range a.Values[1:] {
			buf.WriteString(", ")
			buf.WriteString(v.Assembly())
		}
	}

	buf.WriteString(commentString(a.Comment))
	return buf.String()
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
	return fmt.Sprintf(".BLOCK %v %v%s", a.Size, a.Byte, commentString(a.Comment))
}
func (a *PseudoBlock) MarshalBinary() (data []byte, err error) {
	data = make([]byte, a.Size)
	for i := range data {
		data[i] = a.Byte
	}
	return
}

/*
	Common method
*/

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
func (a *PseudoData) GetComment() string {
	return a.Comment
}
func (a *PseudoData) SetComment(v string) {
	a.Comment = v
}
func (a *Comment) GetComment() string {
	return a.Content
}
func (a *Comment) SetComment(v string) {
	a.Content = v
}
