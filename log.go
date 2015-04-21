package bbvm
import (
	"fmt"
	"strconv"
	"io"
)

func LogOut(v *vm, o io.Writer) {
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