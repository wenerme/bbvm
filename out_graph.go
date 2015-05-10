package bbvm
import (
	"image/color"
	"image"
)

func (out)Graphic(v VM) {
	v.Attr()["graph-dev"] = NewGraphDev(240, 320)
	v.SetOut(16, HANDLE_ALL, outGraphicFunc)
	v.SetOut(17, HANDLE_ALL, outGraphicFunc)
	v.SetOut(18, HANDLE_ALL, outGraphicFunc)
	v.SetOut(19, HANDLE_ALL, outGraphicFunc)
	v.SetOut(20, HANDLE_ALL, outGraphicFunc)
	v.SetOut(21, HANDLE_ALL, outGraphicFunc)
	v.SetOut(22, HANDLE_ALL, outGraphicFunc)
	v.SetOut(23, HANDLE_ALL, outGraphicFunc)
	v.SetOut(24, HANDLE_ALL, outGraphicFunc)
	v.SetOut(25, HANDLE_ALL, outGraphicFunc)
	v.SetOut(64, HANDLE_ALL, outGraphicFunc)
	v.SetOut(68, HANDLE_ALL, outGraphicFunc)
	v.SetOut(69, HANDLE_ALL, outGraphicFunc)
}
func outGraphicFunc(i *Inst) {
	v, p, _ := i.VM, i.A.Get(), i.B // port and param
	r2, r3 := &v.r2, &v.r3
	gdev := v.Attr()["graph-dev"].(GraphDev)
	pagePool := gdev.PagePool()
	picPool := gdev.PicPool()
	_ = picPool
	_ = gdev
	/*
16 | 设定模拟器屏幕 | 0 | r2:宽, r3:高 |  SETLCD(WIDTH,HEIGHT)
17 | 申请画布句柄 | 0 ,r3:PAGE句柄 | - | CREATEPAGE()
18 | 释放画布句柄 | 0 | r3:PAGE句柄 |  DELETEPAGE(PAGE)
19 | 申请图片句柄并从文件载入像素资源 | r3:资源句柄 | r3:文件名, r2:资源索引 |  LOADRES(FILE$,ID)
20 | 复制图片到画布上 | 0 | r3:地址,其他参数在该地址后 |  SHOWPIC(PAGE,PIC,DX,DY,W,H,X,Y,MODE)
21 | 显示画布 | 0 | r3:PAGE句柄 |  FLIPPAGE(PAGE)
22 | 复制画布 | 0 | r2:目标PAGE句柄,r3:源PAGE句柄 |  BITBLTPAGE(DEST,SRC)
23 | 填充画布 | 0 | r3:参数地址 |  FILLPAGE(PAGE,X,Y,WID,HGT,COLOR)
24 | 写入画布某点颜色 | 0 | r3:参数地址 |  PIXEL(PAGE,X,Y,COLOR)
25 | 读取画布某点颜色 | 0 | r3:参数地址 |  READPIXEL(PAGE,X,Y)
	 */
	switch p{
	case 16:
		gdev.SetLcd(r2.Get(), r3.Get())
	case 17:
		r, err := pagePool.Acquire()
		if err != nil {log.Warning("CREATEPAGE() faield: %s", err.Error()); break}
		log.Debug("CREATEPAGE() -> %d", r.Id())
		r3.Set(r.Id())
	case 18:
		r := pagePool.Get(r3.Get())
		if r == nil {log.Warning("DELETEPAGE(%s) faield: Page not exists", r3.Get()); break}
		pagePool.Release(r)
		log.Debug("DELETEPAGE(%s)", r3.Get())
	case 19:
		fn, idx := r3.Str(), r2.Get()
		r3.Set(0)
		log.Debug("LOADRES(%s,%d) -> %d", fn, idx, r3.Get())
	case 20:
		args := newArgs(r3.Get(), v, 9)
		pi, picId, dx, dy, w, h, x, y, mode := args.Next9Int()
		log.Debug("SHOWPIC(%d,%d,%d,%d,%d,%d,%d,%d,%d)", pi, picId, dx, dy, w, h, x, y, mode)
	case 21:
		log.Debug("FLIPPAGE(%d)", r3.Get())
	case 22:
		log.Debug("BITBLTPAGE(%d,%d)", r2.Get(), r3.Get())
	case 23:
		args := newArgs(r3.Get(), v, 6)
		pi, x, y, wid, hgt, c := args.Next6Int()
		log.Debug("FILLPAGE(%d,%d,%d,%d,%d,%d)", pi, x, y, wid, hgt, c)
	case 24:
		args := newArgs(r3.Get(), v, 4)
		pi, x, y, c := args.Next4Int()
		log.Debug("PIXEL(%d,%d,%d,%d)", pi, x, y, c)
		if pg := getPage(pagePool, pi); pg != nil {
			// TODO 边界检查
			pg.Set(x, y, rgbInt2Color(c))
		}
	case 25:
		args := newArgs(r3.Get(), v, 4)
		pageId, x, y := args.Next3Int()
		log.Debug("READPIXEL(%d,%d,%d)", pageId, x, y)
		if pg := getPage(pagePool, pageId); pg != nil {
			// TODO 边界检查
			r3.Set(color2BGRInt(pg.At(x, y)))
		}
	/*
26 | 释放图片句柄 | 0 | r3:资源句柄 |  FREERES(ID)
35 | 清屏 | 0 |  |
36 | 按行列定位光标 | 0 | r2:行,r3:列 |  LOCATE(LINE,ROW)
37 | 设定文字颜色 | 0 | r3:参数地址 |  COLOR(FRONT,BACK,FRAME)
38 | 设定文字字体大小 | 0 | r3:FONT |  FONT(F)
40 | 获取图片宽度 | r3 | r3 |  GETPICWID(PIC)
41 | 获取图片高度 | r3 | r3 |  GETPICHGT(PIC)
42 | 按坐标定位光标 | - | r2:行,r3:列 |  PIXLOCATE(LINE,ROW)
43 | 复制部分画布 | - | r3:参数地址 |  STRETCHBLTPAGE(X,Y,DEST,SRC)
44 | 设定背景模式 | r3:MODE | - |  SETBKMODE(mode)
	 */
	case 40:
		if pic := getPic(picPool, r3.Get()); pic != nil {
			log.Debug("GETPICWID(%d) -> %d", r3.Get(), pic.Bounds().Dx())
			r3.Set(pic.Bounds().Dx())
		}else {
			log.Warning("GETPICWID(%d) faield: Picture not exists", r3.Get())
			r3.Set(0)
		}
	case 41:
		if pic := getPic(picPool, r3.Get()); pic != nil {
			log.Debug("GETPICHGT(%d) -> %d", r3.Get(), pic.Bounds().Dy())
			r3.Set(pic.Bounds().Dy())
		}else {
			log.Warning("GETPICHGT(%d) faield: Picture not exists", r3.Get())
			r3.Set(0)
		}
	case 64:
		args := newArgs(r3.Get(), v, 4)
		pi, style, w, c := args.Next4Int()
		log.Debug("SETPEN(%d,%d,%d,%d)", pi, style, w, c)
		if pg := getPage(pagePool, pi); pg != nil {
			pg.SetPen(PenStyle(style), w, bgrIntColor(c))
		}
	case 68:
		// RECTANGLE(PAGE,LEFT,TOP,RIGHT,BOTTOM)
		args := newArgs(r3.Get(), v, 5)
		pi, left, top, right, bottom := args.Next5Int()
		log.Debug("RECTANGLE(%d,%d,%d,%d,%d)", pi, left, top, right, bottom)
		if pg := getPage(pagePool, pi); pg != nil {
			pg.Rect(image.Rect(left, top, right, bottom))
		}
	case 69:
		args := newArgs(r3.Get(), v, 4)
		pi, cx, cy, cr := args.Next4Int()
		log.Debug("CIRCLE(%d,%d,%d,%d)", pi, cx, cy, cr)
		if pg := getPage(pagePool, pi); pg != nil {
			pg.Circle(cx, cy, cr)
		}
	}
}

func getPage(p ResPool, id int) Page {
	r := p.Get(id)
	if r == nil || r.Get() == nil {return nil}
	return r.Get().(Page)
}
func getPic(p ResPool, id int) Picture {
	r := p.Get(id)
	if r == nil || r.Get() == nil {return nil}
	return r.Get().(Picture)
}
func color2BGRInt(c color.Color) int {
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	return int(b<<16|g<<8|r)
}
func rgbInt2Color(i int) color.Color {
	return color.RGBA{uint8(i >> 16&0xff), uint8(i>>8&0xff), uint8(i&0xff), 0xff}
}

func newArgs(addr int, vm VM, n int) args {
	a := args(make([]int, n))
	for i := 0; i < n; i ++ {
		a[i] = vm.GetInt(addr)
		addr += 4
	}
	return a
}
type args []int
func (a args)Next2Int() (int, int) {
	return a[1], a[0]
}
func (a args)Next3Int() (int, int, int) {
	return a[2], a[1], a[0]
}
func (a args)Next4Int() (int, int, int, int) {
	return a[3], a[2], a[1], a[0]
}
func (a args)Next5Int() (int, int, int, int, int) {
	return a[4], a[3], a[2], a[1], a[0]
}
func (a args)Next6Int() (int, int, int, int, int, int) {
	return a[5], a[4], a[3], a[2], a[1], a[0]
}

func (a args)Next9Int() (int, int, int, int, int, int, int, int, int) {
	return a[8], a[7], a[6], a[5], a[4], a[3], a[2], a[1], a[0]
}
