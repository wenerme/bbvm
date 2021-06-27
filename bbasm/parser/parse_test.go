package parser_test

import (
	"encoding"
	"encoding/hex"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/wenerme/bbvm/bbasm/parser"
	"log"
	"testing"
)

func TestParseLine(t *testing.T) {
	for _, test := range []struct {
		L string
		S string
		B []byte
	}{
		{
			L: `DATA STR CHAR "Hello, BBvm",0`,
			B: []byte("Hello, BBvm\x00\x00\x00\x00"),
			S: `DATA STR "Hello, BBvm", 0`,
		},
	} {
		assembly, err := parser.ParseLine(test.L + "\n")
		assert.NoError(t, err)
		if !assert.NotNil(t, assembly) {
			continue
		}
		bm, isBin := assembly.(encoding.BinaryMarshaler)
		var bin []byte
		if isBin {
			bin, err = bm.MarshalBinary()
			assert.NoError(t, err)
		}
		asm := assembly.Assembly()

		log.Printf("PARSE %v -> %v", test.L, asm)
		log.Printf("\t  %v", hex.Dump(bin))

		if test.S != "" {
			assert.Equal(t, test.S, asm)
		}
		if test.B != nil {
			if !assert.Equal(t, test.B, bin) {
				spew.Dump(assembly)
			}
		}
	}
}

