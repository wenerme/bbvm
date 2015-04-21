package bbvm
import (
	"math"
	"strconv"
	"log"
)

/*
0 | 浮点数转换为整数 | 整数 | r3:浮点数 | IN(r3:float):int
1 | 整数转换为浮点数 | 浮点数 | r3:整数 | IN(r3:int):float
2 | 申请字符串句柄 | 申请到的句柄 |  |  IN():SHDL<br>从-1开始查询
3 | 字符串转换为整数 | 整数 | r3:字符串句柄,__地址__ | IN(r3:SHDL):int<br>若r3的值不是合法的字符串句柄则返回r3的值
4 | 整数转换为字符串 | 返回的值为r3:整数 | r2:目标字符串_句柄_<br>r3:整数 | IN(r2:SHDL,r3:int):int<br>r2所代表字符串的内容被修改
*/

func handleInConv(i Inst) {
	v, p, o := i.VM, i.B.Get(), i.A // port and out

	switch p {
		case 0:
		o.Set(int(math.Float32frombits(uint(i.VM.Register(REG_R3).Get()))))
		case 1:
		o.SetFloat(float32(i.VM.Register(REG_R3).Get()))
		case 3:
		if s, fine := v.GetStr(v.Register(REG_R3)); fine {
			if r, ok := strconv.Atoi(s); ok {
				o.Set(r)
			}else {
				log.Println("Convert atoi faield:"+s)
			}
		}else {
			log.Println("GetStr faield")
		}
		case 4:
		if s, fine := v.GetStr(v.Register(REG_R3)); fine {
			if r, ok := strconv.Atoi(s); ok {
				o.Set(r)
			}else {
				log.Println("Convert atoi faield:"+s)
			}
		}else {
			log.Println("GetStr faield")
		}

	}
}