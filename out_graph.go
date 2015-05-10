package bbvm
import (
	"image/color"
)

func (out)Graphic(v VM) {
	v.Attr()["graph-dev"] = NewGraphDev()
	v.SetOut(16, HANDLE_ALL, outGraphicFunc)
}
type args struct {
	addr int
	vm *vm
	step int
}
func (a *args)NextInt() int {
	i := a.vm.GetInt(a.addr)
	a.addr+=a.step
	return i
}
func (a *args)Next2Int() (int, int) {
	return a.NextInt(), a.NextInt()
}
func (a *args)Next3Int() (int, int, int) {
	return a.NextInt(), a.NextInt(), a.NextInt()
}
func (a *args)Next4Int() (int, int, int, int) {
	return a.NextInt(), a.NextInt(), a.NextInt(), a.NextInt()
}
func (a *args)Next5Int() (int, int, int, int, int) {
	return a.NextInt(), a.NextInt(), a.NextInt(), a.NextInt(), a.NextInt()
}
func (a *args)Next6Int() (int, int, int, int, int, int) {
	return a.NextInt(), a.NextInt(), a.NextInt(), a.NextInt(), a.NextInt(), a.NextInt()
}

func (a *args)Int() int {
	return a.vm.GetInt(a.addr)
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
		args := args{r3.Get()+36, v, -4}
		pageId, picId, dx, dy, w, h := args.Next6Int()
		x, y, mode := args.Next3Int()
		log.Debug("SHOWPIC(%d,%d,%d,%d,%d,%d,%d,%d,%d)", pageId, picId, dx, dy, w, h, x, y, mode)
	case 21:
		log.Debug("FLIPPAGE(%d)", r3.Get())
	case 22:
		log.Debug("BITBLTPAGE(%d,%d)", r2.Get(), r3.Get())
	case 23:
		args := args{r3.Get()+26, v, -4}
		pageId, x, y, wid, hgt, color := args.Next6Int()
		log.Debug("FILLPAGE(%d,%d,%d,%d,%d,%d)", pageId, x, y, wid, hgt, color)
	case 24:
		args := args{r3.Get()+16, v, -4}
		pageId, x, y, color := args.Next4Int()
		log.Debug("PIXEL(%d,%d,%d,%d)", pageId, x, y, color)
		if pg := getPage(pagePool, pageId); pg != nil {
			// TODO 边界检查
			pg.Set(x, y, rgbInt2Color(color))
		}
	case 25:
		args := args{r3.Get()+12, v, -4}
		pageId, x, y := args.Next3Int()
		log.Debug("READPIXEL(%d,%d,%d)", pageId, x, y)
		if pg := getPage(pagePool, pageId); pg != nil {
			// TODO 边界检查
			r3.Set(color2RGBInt(pg.At(x, y)))
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
func color2RGBInt(c color.Color) int {
	r, g, b, _ := c.RGBA()
	return int(r<<16|g<<8|b)
}
func rgbInt2Color(i int) color.Color {
	return color.RGBA{uint8(i >> 16&0xff), uint8(i>>8&0xff), uint8(i&0xff), 0xff}
}