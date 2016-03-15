package asm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLookup(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		t   interface{}
		v   string
		exp interface{}
	}{
		{T_INT, "Int", T_INT},
		{T_INT, "int", T_INT},
		{T_INT, "word", T_WORD},
		{REG_R0, "rP", REG_RP},
		{REG_R0, "rxz", nil},
	}
	for _, test := range tests {
		assert.Equal(test.exp, Lookup(test.t, test.v))
	}
}
