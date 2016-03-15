package asm

import (
	"fmt"
	"reflect"
)

var lookTable = make(map[reflect.Type]map[string]interface{})

func init() {
	for i := T_DWORD; i <= T_INT; i++ {
		putLookTable(i)
	}
	for i := CMP_Z; i <= CMP_NZ; i++ {
		putLookTable(i)
	}
	for i := CAL_ADD; i <= CAL_MOD; i++ {
		putLookTable(i)
	}
	for i := OP_NOP; i <= OP_CAL; i++ {
		putLookTable(i)
	}
	putLookTable(OP_EXIT)

	for i := AM_REGISTER; i <= AM_DIRECT; i++ {
		putLookTable(i)
	}
	for i := REG_RP; i <= REG_R3; i++ {
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

func Lookup(t interface{}, v string) (ret interface{}) {
	var lookType reflect.Type
	if ty, ok := t.(reflect.Type); ok {
		lookType = ty
	} else {
		lookType = reflect.TypeOf(t)
	}
	m := lookTable[lookType]
	if m != nil {
		ret = m[v]
	}
	return
}
