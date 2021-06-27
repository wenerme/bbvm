package parser

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolve(t *testing.T) {
	asm := `
JMP CODE
DATA STR CHAR "Hello, BBvm",0
CODE:

OUT 1, STR
EXIT
`
	lines, err := Parse(asm)
	assert.NoError(t, err)
	a := &Assembler{Lines: lines}
	_, err = a.Assemble()
	assert.NoError(t, err)
	spew.Dump(a.Labels)
	spew.Dump(a.Symbols)
}
