package parser

import (
	"bytes"
	"encoding"
	"fmt"
)

var _ = []encoding.BinaryMarshaler{
	&PseudoData{},
}

type Assembly interface {
	SetComment(string)
	GetComment() string

	// Len Byte code length of this asm
	Len() int
	// Assembly To assembly string
	Assembly() string
}

// Symbol represent a label of a location
type Symbol struct {
	Name      string
	Address   int
	Reference []func(int) // Referenced this symbol's address
}

func commentString(Comment string) string {
	if Comment == "" {
		return ""
	}
	return " ; " + Comment
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
	Len() int         // Byte length
	Assembly() string // Assembly represent
}

func (a *PseudoData) Len() int {
	n := 0
	for _, v := range a.Values {
		n += v.Len()
	}
	return n
}

func (a PseudoData) MarshalBinary() (data []byte, err error) {
	for _, v := range a.Values {
		b, err := v.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}
	return
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
