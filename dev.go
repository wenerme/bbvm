package bbvm
import (
	"image"
	"image/color"
)

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
26 | 释放图片句柄 | 0 | r3:资源句柄 |  FREERES(ID)
27 | 延迟一段时间 | 0 | r3:延迟时间 |  MSDELAY(MSEC)
32 | 用种子初始化随机数生成器 | 0 | r3:SEED |  RANDOMIZE(SEED)
33 | 获取范围内随机数 | 0 | r3:RANGE |  RND(RANGE)
34 | 判定某键是否按下 | 0;r3 | r3:KEY |  KEYPRESS(KEY)
35 | 清屏 | 0 |  |
36 | 按行列定位光标 | 0 | r2:行,r3:列 |  LOCATE(LINE,ROW)
37 | 设定文字颜色 | 0 | r3:参数地址 |  COLOR(FRONT,BACK,FRAME)
38 | 设定文字字体大小 | 0 | r3:FONT |  FONT(F)
39 | 等待按键 | r3:按键 | - |  WAITKEY()
40 | 获取图片宽度 | r3 | r3 |  GETPICWID(PIC)
41 | 获取图片高度 | r3 | r3 |  GETPICHGT(PIC)
42 | 按坐标定位光标 | - | r2:行,r3:列 |  PIXLOCATE(LINE,ROW)
43 | 复制部分画布 | - | r3:参数地址 |  STRETCHBLTPAGE(X,Y,DEST,SRC)
44 | 设定背景模式 | r3:MODE | - |  SETBKMODE(mode)
45 | 获取按键的字符串 | 0 | r3:字符串句柄,用于存储结果 |  InKey$
46 | 获取按键的ASCII码 | 0 | r3:KEYPRESS |  INKEY()
64 | 设置画笔 | 0 | r3:参数地址 |  SETPEN(PAGE,STYLE,WID,COLOR)
65 | 设置刷子 | 0 | r2:PAGE r3:STYLE |  SETBRUSH(PAGE,STYLE)
66 | 移动画笔 | 0 | r1,r2,r3:PAGE,X,Y |  MOVETO(PAGE,X,Y)
67 | 画线 | 0 | r1,r2,r3:PAGE,X,Y |  LINETO(PAGE,X,Y)
68 | 画矩形 | 0 | r3:参数地址 |  RECTANGLE(PAGE,LEFT,TOP,RIGHT,BOTTOM)
69 | 画圆 | 0 | r3:参数地址 |  CIRCLE(PAGE,CX,CY,CR)
80 | 复制部分画布扩展 | 0 | r3:参数地址 |  STRETCHBLTPAGEEX(X,Y,WID,HGT,CX,CY,DEST,SRC)
*/
type (
PenStyle int
Page interface {
	Graphic
	//	SetBrushStyle(BrushStyle)
	SetPen(PenStyle, int, color.Color)
	SetBrushStyle(BrushStyle)
}



Picture interface {
	image.Image
	// 图片名
	Name() string
}

GraphDev interface {
	LoadPicture(fn string, i int) (Picture, error)
	SetLcd(int, int)
	PagePool() ResPool
	PicPool() ResPool
	FlipPage(Page)
	SetFont(FontType)
	SetFontStyle(color.Color, color.Color)
	SetBackgroundMode(BackgroundMode)
	Print(string, ... interface{})
	Locate(int, int)
	LocatePixel(int, int)
	Screen() Page
}

KeyEventType int
KeyEvent struct {
	Type KeyEventType
	Code KeyCode
	Char rune
}
MouseEventType int
MouseEvent struct {
	Type MouseEventType
}
InputDev interface {
	IsPressed(KeyCode) bool
	//	InKey() KeyCode
	WaitKey() KeyEvent
	KeyEvent() chan KeyEvent
	MouseEvent() chan MouseEvent
}

Font interface {
	Render(string string, fg, bg color.Color) (image.Image, error)
	Height() int
	Width() int
}


)


const (
	KeyUp KeyEventType = iota
	KeyDown
)
const (
	MouseDown MouseEventType = iota
	MouseUp
	MouseMove
)
