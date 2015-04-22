package bbvm
import (
	"fmt"
	"strconv"
	"io"
	"bufio"
	"math"
)

type out struct { }
var OUT out

//10 | 键入整数 | 0 |  | r3的值变为键入的整数
//11 | 键入字符串 | 0 | r3:目标字符串句柄 | r3所指字符串的内容变为键入的字符串
//12 | 键入浮点数 | 0 |  | r3的值变为键入的浮点数
func (out)InputByReader(v VM, input io.Reader) {
	r := bufio.NewReader(input)
	hdl := func(i *Inst) {
		v, p, _ := i.VM, i.A.Get(), i.B // port and param

		s, err := r.ReadString('\n')
		if err != nil {
			log.Error("Input failed, got '%s': %s", s, err)
		}
		s = s[:len(s)-1]
		switch p{
			//			case 10:
			//			i, err := strconv.Atoi(s)
			//			if err != nil {
			//				log.Error("Input int failed, got '%s': %s", s, err)
			//			}
			//			v.r3.Set(i)
			case 11:
			v.StrPool().Get(v.r3.Get()).Set(s)
			case 10:
			fallthrough
			case 12:
			i, err := strconv.ParseFloat(s, 32)
			if err != nil {
				log.Error("Input %d failed, got '%s': %s", p, s, err)
			}
			if p == 10 {
				v.r3.Set(int(i))
			}else {
				v.r3.Set(int(math.Float32bits(float32(i))))
			}
		}
	}
	v.SetOut(10, 0, hdl)
	v.SetOut(11, 0, hdl)
	v.SetOut(12, 0, hdl)
}

func (out)OutputToWriter(v VM, o io.Writer) {
	hdl := func(i *Inst) {
		var msg string
		switch i.A.Get(){
			case 0: msg = strconv.Itoa(i.B.Get()) +"\n"
			case 1: msg = i.B.Str() + "\n"
			case 2: msg = i.B.Str()
			case 3: msg = strconv.Itoa(i.B.Get())
			case 4: msg = fmt.Sprintf("%c", i.B.Get())
			case 5: msg = fmt.Sprintf("%.6f", i.B.Float())
		}
		fmt.Fprint(o, msg)
	}
	v.SetOut(0, HANDLE_ALL, hdl)
	v.SetOut(1, HANDLE_ALL, hdl)
	v.SetOut(2, HANDLE_ALL, hdl)
	v.SetOut(3, HANDLE_ALL, hdl)
	v.SetOut(4, HANDLE_ALL, hdl)
	v.SetOut(5, HANDLE_ALL, hdl)
}