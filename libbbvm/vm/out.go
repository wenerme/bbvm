package vm
import (
	"fmt"
	"strconv"
	"io"
	"bufio"
	"math"
	"math/rand"
	"time"
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
		case 5: msg = float32ToStr(i.B.Float32())
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

/*
27 | 延迟一段时间 | 0 | r3:延迟时间 |  MSDELAY(MSEC)
32 | 用种子初始化随机数生成器 | 0 | r3:SEED |  RANDOMIZE(SEED)
33 | 获取范围内随机数 | 0 | r3:RANGE |  RND(RANGE)
255 | 虚拟机测试 | 0 | 0 |  VmTest
 */
func (out)Misc(v VM) {
	v.Attr()["rand"]=rand.New(rand.NewSource(0))
	v.SetOut(27, 0, outMiscFunc)
	v.SetOut(32, 0, outMiscFunc)
	v.SetOut(33, 0, outMiscFunc)
	v.SetOut(255, 0, outMiscFunc)
}
func outMiscFunc(i *Inst) {
	v, p, _ := i.VM, i.A.Get(), i.B // port and param
	r3 := &v.r3
	rand := v.Attr()["rand"].(*rand.Rand)
	switch p{
	case 27:
		log.Info("MSDELAY(%d)", r3.Get())
		time.Sleep(time.Duration(r3.Get()) * time.Millisecond)
	case 32:
		log.Info("RANDOMIZE(%d)", r3.Get())
		rand.Seed(int64(r3.Get()))
	case 33:
		r := rand.Int31n(int32(r3.Get()))
		log.Info("RND(%d) -> %d", r3.Get(), r)
		r3.Set(int(r))
	case 255:
		log.Info("VMTEST()")
	}
}
