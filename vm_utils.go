package bbvm
import (
	"math"
	"strconv"
	"github.com/op/go-logging"
	"os"
	"fmt"
)


var log = logging.MustGetLogger("bbvm")

// 初始化 Log
func init() {
	format := logging.MustStringFormatter("%{color}%{time:15:04:05} %{level:.4s} %{shortfunc} %{color:reset} %{message}", )
	//	format := logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{longfile} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}", )
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1Formatter)
}

func calculate(a int, b int, t CalculateType, d DataType) int {
	var oa, ob, oc float64
	if d == T_FLOAT {
		oa = float64(math.Float32frombits(uint32(a)))
		ob = float64(math.Float32frombits(uint32(b)))
	}else {
		oa = float64(a)
		ob = float64(b)
	}
	switch t{
	case CAL_ADD: oc = oa + ob
	case CAL_SUB: oc = oa - ob
	case CAL_MUL: oc = oa * ob
	case CAL_DIV: oc = oa / ob
	case CAL_MOD: oc = float64(int(oa) % int(ob))
	default:
		log.Error("Unhandled calculate")
	}
	var ret int
	switch d{
	case T_DWORD: ret = int(oc)
	case T_WORD: ret = int(int16(oc))
	case T_BYTE: ret = int(byte(oc))
	case T_FLOAT: ret = int(math.Float32bits(float32(oc)))
	case T_INT: ret = int(oc)
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
	}else {
		sa = strconv.Itoa(a)
	}
	if b == HANDLE_ALL {
		sb = "*"
	}else {
		sb = strconv.Itoa(b)
	}
	return sa+","+sb
}


func newStrPool() ResPool {
	return &resPool{
		pool: make(map[int]Res),
		current:-1,
		start:-1,
		step: -1,
		limit: math.MaxInt32,
		creator:func(_ ResPool) interface{} {return ""},
	}
}

func float32ToStr(f float32) string {
	return fmt.Sprintf("%.6f", f)
}
