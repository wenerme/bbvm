package parser

import (
	"fmt"
	"github.com/wenerme/bbvm/bbasm"
	"reflect"
	"strings"
)

var lookTable = make(map[reflect.Type]map[string]interface{})

func init() {
	for i := bbasm.DWORD; i <= bbasm.INT; i++ {
		putLookTable(i)
	}
	for i := bbasm.Z; i <= bbasm.NZ; i++ {
		putLookTable(i)
	}
	for i := bbasm.ADD; i <= bbasm.MOD; i++ {
		putLookTable(i)
	}
	for i := bbasm.NOP; i <= bbasm.CAL; i++ {
		putLookTable(i)
	}
	putLookTable(bbasm.EXIT)

	for i := bbasm.AddressRegister; i <= bbasm.AddressDirect; i++ {
		putLookTable(i)
	}
	for i := bbasm.RP; i <= bbasm.R3; i++ {
		putLookTable(i)
	}
}

func putLookTable(v interface{}) {
	m := lookTable[reflect.TypeOf(v)]
	if m == nil {
		m = make(map[string]interface{})
		lookTable[reflect.TypeOf(v)] = m
	}
	m[v.(fmt.Stringer).String()] = v
}

// Lookup type t by string v, will upper case v.
//
// Return nil if not found
func Lookup(t interface{}, v string) (ret interface{}) {
	var lookType reflect.Type
	if ty, ok := t.(reflect.Type); ok {
		lookType = ty
	} else {
		lookType = reflect.TypeOf(t)
	}
	m := lookTable[lookType]
	if m != nil {
		ret = m[strings.ToUpper(v)]
	}
	return
}
