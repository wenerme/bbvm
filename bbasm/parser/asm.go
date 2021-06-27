package parser

import (
	"encoding"
	"fmt"
	"github.com/wenerme/bbvm/bbasm"
)

type Assembler struct {
	Lines   []Assembly
	Symbols map[string]*Symbol
	Labels  map[string]int
}

func (asm *Assembler) Assemble() ([]byte, error) {
	var o []byte
	locs := map[string]int{}
	syms := map[string]*Symbol{}
	// resolve
	track := func(a string, cb func(int)) {
		if a == "" {
			return
		}
		if syms[a] == nil {
			syms[a] = &Symbol{
				Name: a,
			}
		}
		if v, ok := locs[a]; ok {
			cb(v)
			return
		}
		s := syms[a]
		s.Reference = append(s.Reference, cb)
	}

	loc := 0
	for _, v := range asm.Lines {
		switch a := v.(type) {
		case *Label:
			locs[a.Name] = loc
		case *PseudoData:
			locs[a.Label] = loc
		case *bbasm.Inst:
			track(a.A.Symbol, func(addr int) {
				a.A.V = int32(addr)
			})
			track(a.B.Symbol, func(addr int) {
				a.B.V = int32(addr)
			})
		}
		loc += v.Len()
	}

	asm.Symbols = syms
	asm.Labels = locs

	for _, v := range syms {
		if addr, ok := locs[v.Name]; ok {
			v.Address = addr
			for _, ref := range v.Reference {
				ref(addr)
			}
		} else {
			return nil, fmt.Errorf("symbol not reolve: %q", v.Name)
		}
	}

	for _, v := range asm.Lines {
		if b, ok := v.(encoding.BinaryMarshaler); ok {
			data, err := b.MarshalBinary()
			if err != nil {
				return nil, err
			}
			o = append(o, data...)
		}
	}
	return o, nil
}
