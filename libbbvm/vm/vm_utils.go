package vm

import (
	"fmt"
	"github.com/wenerme/bbvm/libbbvm/asm"
	"math"
	"strconv"
)

func calculate(a int, b int, t asm.CalculateType, d asm.DataType) int {
	var oa, ob, oc float64
	if d == asm.T_FLOAT {
		oa = float64(math.Float32frombits(uint32(a)))
		ob = float64(math.Float32frombits(uint32(b)))
	} else {
		oa = float64(a)
		ob = float64(b)
	}
	switch t {
	case asm.CAL_ADD:
		oc = oa + ob
	case asm.CAL_SUB:
		oc = oa - ob
	case asm.CAL_MUL:
		oc = oa * ob
	case asm.CAL_DIV:
		oc = oa / ob
	case asm.CAL_MOD:
		oc = float64(int(oa) % int(ob))
	default:
		log.Error("Unhandled calculate")
	}
	var ret int
	switch d {
	case asm.T_DWORD:
		ret = int(oc)
	case asm.T_WORD:
		ret = int(int16(oc))
	case asm.T_BYTE:
		ret = int(byte(oc))
	case asm.T_FLOAT:
		ret = int(math.Float32bits(float32(oc)))
	case asm.T_INT:
		ret = int(oc)
	}
	return ret
}

func multiPortHandler(a, b int, m map[string]InstHandler) InstHandler {
	h, exists := m[multiPortKey(a, b)]
	if !exists {
		h, exists = m[multiPortKey(HANDLE_ALL, b)]
	}
	if !exists {
		h, exists = m[multiPortKey(a, HANDLE_ALL)]
	}
	return h
}
func multiPortKey(a, b int) string {
	var sa, sb string
	if a == HANDLE_ALL {
		sa = "*"
	} else {
		sa = strconv.Itoa(a)
	}
	if b == HANDLE_ALL {
		sb = "*"
	} else {
		sb = strconv.Itoa(b)
	}
	return sa + "," + sb
}

func newStrPool() ResPool {
	return &resPool{
		pool:    make(map[int]Res),
		current: -1,
		start:   -1,
		step:    -1,
		limit:   math.MaxInt32,
		creator: func(ResPool) interface{} { return "" },
	}
}

func float32ToStr(f float32) string {
	return fmt.Sprintf("%.6f", f)
}
