package bbvm

import (
	"context"
	"github.com/wenerme/bbvm/bbasm"
	"reflect"
	"time"
)

type StdBuilder func(rt bbasm.Runtime, std *Std) *Std

func StdFrom(ctx context.Context) *Std {
	v, _ := ctx.Value("STDLIB").(*Std)
	return v
}
func WithStd(ctx context.Context, std *Std) context.Context {
	return context.WithValue(ctx, "STDLIB", std)
}

type Handler interface {
	Handler() int
}
type StringHandler = Handler
type PageHandler = Handler
type ResourceHandler = Handler
type FileHandler = Handler

type Std struct {
	FloatToInt func(ctx context.Context, v float32) int
	IntToFloat func(ctx context.Context, v int) float32

	AllocString        func(ctx context.Context) StringHandler
	StringToInt        func(ctx context.Context, hdr StringHandler) (int, error)
	IntToString        func(ctx context.Context, dst StringHandler, v int)
	StringCopy         func(ctx context.Context, dst StringHandler, src StringHandler)
	StringConcat       func(ctx context.Context, a StringHandler, b StringHandler)
	StringLength       func(ctx context.Context, hdr StringHandler) int
	FreeString         func(ctx context.Context, hdr StringHandler)
	StringCompare      func(ctx context.Context, a StringHandler, b StringHandler) int
	IntToFloatToString func(ctx context.Context, dst StringHandler, v int)
	StringToFloat      func(ctx context.Context, hdr StringHandler) (float32, error)
	StringGetAscii     func(ctx context.Context, hdr StringHandler, idx int) int
	StringSetAscii     func(ctx context.Context, hdr StringHandler, idx, v int)
	// extra
	StringGet func(ctx context.Context, hdr StringHandler) string
	StringSet func(ctx context.Context, hdr StringHandler, v string)
	StringOf  func(ctx context.Context, hdr int) StringHandler

	// in 14,unknown

	Tick func(ctx context.Context) int

	Sin      func(ctx context.Context, a float32) float32
	Cos      func(ctx context.Context, a float32) float32
	Tan      func(ctx context.Context, a float32) float32
	Sqrt     func(ctx context.Context, a float32) float32
	IntAbs   func(ctx context.Context, a int) int
	FloatAbs func(ctx context.Context, a float32) float32

	DataPtrSet func(ctx context.Context, v int)

	Read  func(ctx context.Context, addr int) int
	Write func(ctx context.Context, addr int, v int)

	GetEnv func(ctx context.Context) int

	StringLeft       func(ctx context.Context, dst StringHandler, hdr StringHandler, len int)
	StringRight      func(ctx context.Context, dst StringHandler, hdr StringHandler, len int)
	StringMid        func(ctx context.Context, dst StringHandler, hdr StringHandler, idx int, len int)
	StringFirstAscii func(ctx context.Context, hdr StringHandler) int
	StringFind       func(ctx context.Context, hdr StringHandler, sub StringHandler, offset int) int

	// out

	PrintLnInt    func(ctx context.Context, v int)
	PrintLnString func(ctx context.Context, v StringHandler)
	PrintString   func(ctx context.Context, v StringHandler)
	PrintInt      func(ctx context.Context, v int)
	PrintChar     func(ctx context.Context, v int)
	PrintFloat    func(ctx context.Context, v float32)

	InputInt    func(ctx context.Context) int
	InputString func(ctx context.Context, dst StringHandler)
	InputFloat  func(ctx context.Context) float32

	DataReadInt    func(ctx context.Context) int
	DataReadString func(ctx context.Context, hdr StringHandler)
	DataReadFloat  func(ctx context.Context) float32

	SetLcd        func(ctx context.Context, w int, h int)
	AllocPage     func(ctx context.Context) PageHandler
	FreePage      func(ctx context.Context, hdr PageHandler)
	LoadImage     func(ctx context.Context, fn StringHandler, idx int) ResourceHandler
	ShowPic       func(ctx context.Context, page PageHandler, res ResourceHandler, dx, dy, w, h, x, y int, mode int)
	FlipPage      func(ctx context.Context, hdr PageHandler)
	PageCopy      func(ctx context.Context, dst PageHandler, src PageHandler)
	PageFill      func(ctx context.Context, hdr PageHandler, x, y, w, h int, color int)
	PagePixel     func(ctx context.Context, hdr PageHandler, x, y, color int)
	PageReadPixel func(ctx context.Context, hdr PageHandler, x, y int) int
	FreeRes       func(ctx context.Context, hdr ResourceHandler)

	Delay    func(ctx context.Context, msec int)  `out:"27,0"`
	RandSeed func(ctx context.Context, seed int)  `out:"32,0"`
	Rand     func(ctx context.Context, n int) int `out:"33,0"`

	IsKeyPressed func(ctx context.Context, k int) int
	Clear        func(ctx context.Context)
	LocateCursor func(ctx context.Context, line, row int)
	SetColor     func(ctx context.Context, font, back, frame int)
	SetFont      func(ctx context.Context, font int)
	WaitKey      func(ctx context.Context) int

	GetImageWidth  func(ctx context.Context, hdr ResourceHandler) int
	GetImageHeight func(ctx context.Context, hdr ResourceHandler) int

	PixelLocateCursor func(ctx context.Context, x, y int)
	PageCopyExt       func(ctx context.Context, dst, src ResourceHandler, x, y int)
	SetBackgroundMode func(ctx context.Context, mod int)

	InputKeyCode func(ctx context.Context, dst StringHandler)

	OpenFile        func(ctx context.Context, fd int, fn StringHandler, mode int)
	CloseFile       func(ctx context.Context, fd int)
	FileReadInt     func(ctx context.Context, fd int, offset int) int
	FileReadFloat   func(ctx context.Context, fd int, offset int) float32
	FileReadString  func(ctx context.Context, fd int, offset int, dst StringHandler)
	FileWriteInt    func(ctx context.Context, fd int, offset int, v int)
	FileWriteFloat  func(ctx context.Context, fd int, offset int, v float32)
	FileWriteString func(ctx context.Context, fd int, offset int, v StringHandler)
	FileEof         func(ctx context.Context, fd int) int
	FileLof         func(ctx context.Context, fd int) int
	FileLoc         func(ctx context.Context, fd int) int
	FileSeek        func(ctx context.Context, fd int, loc int)

	SetPen        func(ctx context.Context, page PageHandler, style, wid, color int)
	SetBrush      func(ctx context.Context, page PageHandler, style int)
	MoveTo        func(ctx context.Context, page PageHandler, x, y int)
	LineTo        func(ctx context.Context, page PageHandler, x, y int)
	DrawRectangle func(ctx context.Context, page PageHandler, left, top, right, bottom int)
	DrawCircle    func(ctx context.Context, page PageHandler, cx, cy, cr int)

	PageCopyExt2 func(ctx context.Context, dst, src PageHandler, x, y, w, h, cx, cy int)
	PageOf       func(ctx context.Context, hdr int) PageHandler
	ResOf        func(ctx context.Context, hdr int) ResourceHandler

	VmTest func(ctx context.Context)

	// extra
	BytesToString func(b []byte) (string, error)
	StringToBytes func(s string) ([]byte, error)
}

func (std *Std) in(ctx context.Context, rt bbasm.Runtime, inst *bbasm.Inst) {
	a := inst.A
	b := inst.B
	r0 := rt.Register(bbasm.R0)
	r1 := rt.Register(bbasm.R1)
	r2 := rt.Register(bbasm.R2)
	r3 := rt.Register(bbasm.R3)
	// a is result, b is port
	bv := b.Get()
	switch bv {
	case 0:
		a.Set(std.FloatToInt(ctx, r3.Float()))
	case 1:
		a.SetFloat(float32(r3.Get()))
	case 2: // 申请字符串句柄 | 申请到的句柄 |  |  IN():SHDL<br>从-1开始查询
		a.Set(std.AllocString(ctx).Handler())
	case 3: // 字符串转换为整数 | 整数 | r3:字符串句柄,__地址__ | float(r3.str);若r3的值不是合法的字符串句柄则返回r3的值
		hdr := std.StringOf(ctx, r3.Get())
		if hdr == nil {
			a.Set(r3.Get())
			break
		}
		v, err := std.StringToInt(ctx, hdr)
		if err != nil {
			a.Set(v)
		} else {
			// todo warn
		}
	case 4: // 整数转换为字符串 | 返回的值为r3:整数 | r2:目标字符串_句柄_<br>r3:整数 | r2.str=str(r3.int);return r3.int;r2所代表字符串的内容被修改
		std.IntToString(ctx, std.StringOf(ctx, r2.Get()), r3.Get())
		a.Set(r3.Get())
	case 5: // 复制字符串 | r3的值 | r2:源字符串句柄<br>r3:目标字符串句柄 | IN(r2:SHDL,r3:SHDL):SHDL<br>r3所代表字符串的内容被修改
		std.StringCopy(ctx, std.StringOf(ctx, r3.Get()), std.StringOf(ctx, r2.Get()))
	case 6: // 连接字符串 | r3的值 | r2:源字符串<br>r3:目标字符串 | r3.str=r3.str+r2.str
		std.StringConcat(ctx, std.StringOf(ctx, r3.Get()), std.StringOf(ctx, r2.Get()))
	case 7: // 获取字符串长度 | 字符串长度 | r3:字符串 | strlen(r3.str)
		a.Set(std.StringLength(ctx, std.StringOf(ctx, r3.Get())))
	case 8: // 释放字符串句柄 | r3的值 | r3:字符串句柄 | IN(r3:SHDL):SHDL
		std.FreeString(ctx, std.StringOf(ctx, r3.Get()))
		a.Set(r3.Get())
	case 9: // 比较字符串 | 两字符串的差值 相同为0，大于为1,小于为-1 | r2:基准字符串<br>r3:比较字符串 | IN(r2:SHDL,r3:SHDL):int
		a.Set(std.StringCompare(ctx, std.StringOf(ctx, r2.Get()), std.StringOf(ctx, r3.Get())))
	case 10: // 整数转换为浮点数再转换为字符串 | r3的值 | r2:目标字符串<br>r3:整数 | r2所代表字符串的内容被修改
		std.IntToFloatToString(ctx, std.StringOf(ctx, r2.Get()), r3.Get())
		a.Set(r3.Get())
	case 11:
		v, err := std.StringToFloat(ctx, std.StringOf(ctx, r3.Get()))
		if err == nil {
			a.SetFloat(v)
		} else {
			a.Set(0)
		}
	case 12: // 获取字符的ASCII码 | ASCII码 | r2:字符位置<br>r3:字符串 | GBK编码,返回的结果范围为有符号的 8bit值,因此对中文操作时返回负数
		a.Set(std.StringGetAscii(ctx, std.StringOf(ctx, r3.Get()), r2.Get()))
	case 13: // 将给定字符串中指定索引的字符替换为给定的ASCII代表的字符 | r3的值 | r1:ASCII码<br>r2:字符位置<br>r3:目标字符串 | r3所代表字符串的内容被修改, 要求r3是句柄才能修改r3的值,给出的ASCII会进行模256的处理
		std.StringSetAscii(ctx, std.StringOf(ctx, r3.Get()), r2.Get(), r1.Get())
		a.Set(r3.Get())
	case 14: // （功用不明） | 65535
		a.Set(65535)
	case 15: // 获取嘀嗒计数 | 嘀嗒计数 |  | 这里不知道他是怎么算的这个数字,但是会随着时间增长就是了
		a.Set(int(time.Now().Unix()))
	case 16: // 求正弦值 | X!的正弦值 | r3:X!
		a.SetFloat(std.Sin(ctx, r3.Float()))
	case 17: // 求余弦值 | X!的余弦值 | r3:X!
		a.SetFloat(std.Cos(ctx, r3.Float()))
	case 18: // 求正切值 | X!的正切值 | r3:X!
		a.SetFloat(std.Tan(ctx, r3.Float()))
	case 19: // 求平方根值 | X!的平方根值 | r3:X!
		a.SetFloat(std.Sqrt(ctx, r3.Float()))
	case 20: // 求绝对值 | X%的绝对值 | r3:X%
		a.Set(std.IntAbs(ctx, r3.Get()))
	case 21: // 求绝对值 | X!的绝对值 | r3:X!
		a.SetFloat(std.FloatAbs(ctx, r3.Float()))
	case 23: // 读内存数据 | 地址内容 | r3:地址
		a.Set(std.Read(ctx, r3.Get()))
	case 24: // 写内存数据 | r3的值 | r2:待写入数据<br>r3:待写入地址
		std.Write(ctx, r3.Get(), r2.Get())
		a.Set(r3.Get())
	case 25: // 获取环境值 | 环境值
		a.Set(std.GetEnv(ctx))
	case 32: // 整数转换为字符串 | r3的值 | r1:整数<br>r3:目标字符串 | r3所代表字符串的内容被修改
		std.IntToString(ctx, std.StringOf(ctx, r3.Get()), r1.Get())
		a.Set(r3.Get())
	case 33: // 字符串转换为整数 | 整数 | r3:字符串 |
		v, err := std.StringToInt(ctx, std.StringOf(ctx, r3.Get()))
		if err == nil {
			a.Set(v)
		} else {
			a.Set(0)
		}
	case 34: // 获取字符的ASCII码 | ASCII码 | r3:字符串
		a.Set(std.StringGetAscii(ctx, std.StringOf(ctx, r3.Get()), 0))
	case 35: // 左取字符串 | r3的值 | r1:截取长度<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改 （此端口似乎不正常）
		std.StringLeft(ctx, std.StringOf(ctx, r3.Get()), std.StringOf(ctx, r2.Get()), r1.Get())
		a.Set(r3.Get())
	case 36: // 右取字符串 | r3的值 | r1:截取长度<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改
		std.StringRight(ctx, std.StringOf(ctx, r3.Get()), std.StringOf(ctx, r2.Get()), r1.Get())
		a.Set(r3.Get())
	case 37: // 中间取字符串 | r0截取长度 | r0:截取长度<br>r1:截取位置<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改
		std.StringMid(ctx, std.StringOf(ctx, r3.Get()), std.StringOf(ctx, r2.Get()), r1.Get(), r0.Get())
		a.Set(r0.Get())
	case 38: // 查找字符串 | 位置 | r1:起始位置<br>r2:子字符串<br>r3:父字符串 |
		a.Set(std.StringFind(ctx, std.StringOf(ctx, r3.Get()), std.StringOf(ctx, r2.Get()), r1.Get()))
	case 39: // 获取字符串长度 | 字符串长度 | r3:字符串
		a.Set(std.StringLength(ctx, std.StringOf(ctx, r3.Get())))
	}
}
func (std *Std) out(ctx context.Context, rt bbasm.Runtime, inst *bbasm.Inst) {
	a := inst.A
	b := inst.B
	r0 := rt.Register(bbasm.R0)
	r1 := rt.Register(bbasm.R1)
	r2 := rt.Register(bbasm.R2)
	r3 := rt.Register(bbasm.R3)
	av := a.Get()
	switch av {
	case 0:
		std.PrintLnInt(ctx, b.Get())
	case 1:
		std.PrintLnString(ctx, std.StringOf(ctx, b.Get()))
	case 2:
		std.PrintString(ctx, std.StringOf(ctx, b.Get()))
	case 3:
		std.PrintInt(ctx, b.Get())
	case 4:
		std.PrintChar(ctx, b.Get())
	case 5:
		std.PrintFloat(ctx, b.Float())
	case 10: // 键入整数 | 0 |  | r3的值变为键入的整数
		r3.Set(std.InputInt(ctx))
	case 11: // 键入字符串 | 0 | r3:目标字符串句柄 | r3所指字符串的内容变为键入的字符串
		std.InputString(ctx, std.StringOf(ctx, r3.Get()))
	case 12: // 键入浮点数 | 0 |  | r3的值变为键入的浮点数
		r3.SetFloat(std.InputFloat(ctx))
	case 27: // 延迟一段时间 | 0 | r3:延迟时间 |  MSDELAY(MSEC)
		std.Delay(ctx, r3.Get())
	case 32: // 用种子初始化随机数生成器 | 0 | r3:SEED |  RANDOMIZE(SEED)
		std.RandSeed(ctx, r3.Get())
	case 33: // 获取范围内随机数 | 0 | r3:RANGE |  RND(RANGE)
		r3.Set(std.Rand(ctx, r3.Get()))
	case 255: // 虚拟机测试 | 0 | 0 |  VmTest
		std.VmTest(ctx)

	/* ========== File ========== */
	case 48: // 打开文件 | 0 | r0:打开方式<br>r1:文件号<br>r3:文件名字符串 | 打开方式目前只能为1
		std.OpenFile(ctx, r1.Get(), std.StringOf(ctx, r3.Get()), r0.Get())
	case 49: // 关闭文件 | 文件号
		std.CloseFile(ctx, b.Get())
	case 50: // 从文件读取数据
		bv := b.Get()
		switch bv {
		case 16: // 读取整数 | r1:文件号<br>r2:位置偏移量 | r3的值变为读取的整数
			r3.Set(std.FileReadInt(ctx, r1.Get(), r2.Get()))
		case 17: // 读取浮点数 | r1:文件号<br>r2:位置偏移量 | r3的值变为读取的浮点数
			r3.SetFloat(std.FileReadFloat(ctx, r1.Get(), r2.Get()))
		case 18: // 读取字符串 | r1:文件号<br>r2:位置偏移量<br>r3:目标字符串句柄 | r3所指字符串的内容变为读取的字符串
			std.FileReadString(ctx, r1.Get(), r2.Get(), std.StringOf(ctx, r3.Get()))
		default:
			panic("invalid out 50")
		}
	case 51: // 向文件写入数据
		bv := b.Get()
		switch bv {
		case 16: // 写入整数 | r1:文件号<br>r2:位置偏移量<br>r3:整数 |
			std.FileWriteInt(ctx, r1.Get(), r2.Get(), r3.Get())
		case 17: // 写入浮点数 | r1:文件号<br>r2:位置偏移量<br>r3:浮点数 |
			std.FileWriteFloat(ctx, r1.Get(), r2.Get(), r3.Float())
		case 18: // 写入字符串 | r1:文件号<br>r2:位置偏移量<br>r3:字符串 |
			std.FileWriteString(ctx, r1.Get(), r2.Get(), std.StringOf(ctx, r3.Get()))
		default:
			panic("invalid out 51")
		}
	case 52: // 判断文件位置指针是否指向文件尾 | 0 | r3:文件号 |  Eof
		r3.Set(std.FileEof(ctx, r3.Get()))
	case 53: // 获取文件长度 | 0 | r3:文件号 |  Lof
		r3.Set(std.FileLof(ctx, r3.Get()))
	case 54: // 获取文件位置指针的位置 | 0 | r3:文件号 |  Loc
		r3.Set(std.FileLoc(ctx, r3.Get()))
	case 55: // 定位文件位置指针 | 0 | r2:文件号<br>r3:目标位置 |
		std.FileSeek(ctx, r2.Get(), r3.Get())

	/* ========== GUI ========== */
	case 16: // 设定模拟器屏幕 | 0 | r2:宽, r3:高 |  SETLCD(WIDTH,HEIGHT)
		std.SetLcd(ctx, r2.Get(), r3.Get())
	case 17: // 申请画布句柄 | 0 ,r3:PAGE句柄 | - | CREATEPAGE()
		r3.Set(std.AllocPage(ctx).Handler())
	case 18: // 释放画布句柄 | 0 | r3:PAGE句柄 |  DELETEPAGE(PAGE)
		std.FreePage(ctx, std.PageOf(ctx, r3.Get()))
	case 21: // 显示画布 | 0 | r3:PAGE句柄 |  FLIPPAGE(PAGE)
		std.FlipPage(ctx, std.PageOf(ctx, r3.Get()))
	case 22: // 复制画布 | 0 | r2:目标PAGE句柄,r3:源PAGE句柄 |  BITBLTPAGE(DEST,SRC)
		std.PageCopy(ctx, std.PageOf(ctx, r2.Get()), std.PageOf(ctx, r3.Get()))
	case 23: // 填充画布 | 0 | r3:参数地址 |  FILLPAGE(PAGE,X,Y,WID,HGT,COLOR)
		p, x, y, w, h, color := ArgOf(rt, r3.Get()).Next6()
		std.PageFill(ctx, std.PageOf(ctx, p), x, y, w, h, color)
	case 24: // 写入画布某点颜色 | 0 | r3:参数地址 |  PIXEL(PAGE,X,Y,COLOR)
		p, x, y, color := ArgOf(rt, r3.Get()).Next4()
		std.PagePixel(ctx, std.PageOf(ctx, p), x, y, color)
	case 25: // 读取画布某点颜色 | 0 | r3:参数地址 |  READPIXEL(PAGE,X,Y)
		p, x, y := ArgOf(rt, r3.Get()).Next3()
		r3.Set(std.PageReadPixel(ctx, std.PageOf(ctx, p), x, y))
	case 26:
		std.FreeRes(ctx, std.ResOf(ctx, r3.Get()))
	case 35: // 清屏
		std.Clear(ctx)
	case 36: // 按行列定位光标 | 0 | r2:行,r3:列 |  LOCATE(LINE,ROW)
		std.LocateCursor(ctx, r2.Get(), r3.Get())
	case 37: // 设定文字颜色 | 0 | r3:参数地址 |  COLOR(FRONT,BACK,FRAME)
		front, back, frame := ArgOf(rt, r3.Get()).Next3()
		std.SetColor(ctx, front, back, frame)
	case 38: // 设定文字字体大小 | 0 | r3:FONT |  FONT(F)
		std.SetFont(ctx, r3.Get())
	case 40: // 获取图片宽度 | r3 | r3 |  GETPICWID(PIC)
		r3.Set(std.GetImageWidth(ctx, std.ResOf(ctx, r3.Get())))
	case 41: // 获取图片高度 | r3 | r3 |  GETPICHGT(PIC)
		r3.Set(std.GetImageHeight(ctx, std.ResOf(ctx, r3.Get())))
	case 42: // 按坐标定位光标 | - | r2:行,r3:列 |  PIXLOCATE(LINE,ROW)
		std.PixelLocateCursor(ctx, r2.Get(), r3.Get())
	case 43: // 复制部分画布 | - | r3:参数地址 |  STRETCHBLTPAGE(X,Y,DEST,SRC)
		x, y, dst, src := ArgOf(rt, r3.Get()).Next4()
		std.PageCopyExt(ctx, std.PageOf(ctx, dst), std.PageOf(ctx, src), x, y)
	case 44: // 设定背景模式 | r3:MODE | - |  SETBKMODE(mode)
		std.SetBackgroundMode(ctx, r3.Get())
	case 64: // 设置画笔 | 0 | r3:参数地址 |  SETPEN(PAGE,STYLE,WID,COLOR)
		p, style, w, color := ArgOf(rt, r3.Get()).Next4()
		std.SetPen(ctx, std.PageOf(ctx, p), style, w, color)
	case 65: // 设置刷子 | 0 | r2:PAGE r3:STYLE |  SETBRUSH(PAGE,STYLE)
		std.SetBrush(ctx, std.PageOf(ctx, r2.Get()), r3.Get())
	case 66: // 移动画笔 | 0 | r1,r2,r3:PAGE,X,Y |  MOVETO(PAGE,X,Y)
		std.MoveTo(ctx, std.PageOf(ctx, r1.Get()), r2.Get(), r3.Get())
	case 67: // 画线 | 0 | r1,r2,r3:PAGE,X,Y |  LINETO(PAGE,X,Y)
		std.LineTo(ctx, std.PageOf(ctx, r1.Get()), r2.Get(), r3.Get())
	case 68: // 画矩形 | 0 | r3:参数地址 |  RECTANGLE(PAGE,LEFT,TOP,RIGHT,BOTTOM)
		p, l, t, r, b := ArgOf(rt, r3.Get()).Next5()
		std.DrawRectangle(ctx, std.PageOf(ctx, p), l, t, r, b)
	case 69: // 画圆 | 0 | r3:参数地址 |  CIRCLE(PAGE,CX,CY,CR)
		p, cx, cy, cr := ArgOf(rt, r3.Get()).Next4()
		std.DrawCircle(ctx, std.PageOf(ctx, p), cx, cy, cr)
	}
}
func (std *Std) Execute(ctx context.Context, rt bbasm.Runtime, inst *bbasm.Inst) {
	switch inst.Opcode {
	case bbasm.IN:
		std.in(ctx, rt, inst)
	case bbasm.OUT:
		std.out(ctx, rt, inst)
	}
}

func (std *Std) Use(neo *Std) {
	a := reflect.ValueOf(std).Elem()
	b := reflect.ValueOf(neo).Elem()
	n := a.NumField()

	for i := 0; i < n; i++ {
		af := a.Field(i)
		bf := b.Field(i)
		if !bf.IsNil() {
			af.Set(bf)
		}
	}
}

type arg struct {
	mem  bbasm.Memory
	addr int
}

func ArgOf(mem bbasm.Memory, addr int) *arg {
	return &arg{mem: mem, addr: addr}
}
func (a *arg) Next() int {
	v := a.mem.GetInt(a.addr)
	a.addr += 4
	return v
}
func (a *arg) Next3() (int, int, int) {
	return a.Next(), a.Next(), a.Next()
}
func (a *arg) Next4() (int, int, int, int) {
	return a.Next(), a.Next(), a.Next(), a.Next()
}
func (a *arg) Next5() (int, int, int, int, int) {
	return a.Next(), a.Next(), a.Next(), a.Next(), a.Next()
}
func (a *arg) Next6() (int, int, int, int, int, int) {
	return a.Next(), a.Next(), a.Next(), a.Next(), a.Next(), a.Next()
}
