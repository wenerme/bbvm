package me.wener.bbvm.vm.invoke;

import me.wener.bbvm.dev.ImageManager;
import me.wener.bbvm.dev.PageManager;
import me.wener.bbvm.vm.Operand;
import me.wener.bbvm.vm.Register;
import me.wener.bbvm.vm.SystemInvoke;
import me.wener.bbvm.vm.VM;

import javax.inject.Inject;
import javax.inject.Named;

/**
 * @author wener
 * @since 15/12/18
 */
public class GraphInvoke {
    private final VM vm;
    private final Register r3;
    private final Register r2;
    private final Register r1;
    private final Register r0;
    private final PageManager pages;
    private final ImageManager images;

    @Inject
    public GraphInvoke(VM vm, @Named("R3") Register r3, @Named("R2") Register r2, @Named("R1") Register r1, @Named("R0") Register r0, PageManager pages, ImageManager images) {
        this.vm = vm;
        this.r3 = r3;
        this.r2 = r2;
        this.r1 = r1;
        this.r0 = r0;
        this.pages = pages;
        this.images = images;
    }

    /**
     * bgr2rgb &lt;-> rgb2bgr
     */
    static int color(int c) {
        return ((c & 0xff) << 16) | (c & 0xff00) | c >> 16 & 0xff;
    }

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
    */
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 16, b = 0)
    public void setSize() {
        pages.setSize(r2.get(), r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 17, b = 0)
    public void createPage() {
        r3.set(pages.create().getHandler());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 18, b = 0)
    public void deletePage() throws Exception {
        r3.get(pages).close();
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 19)
    public void loadRes(@Named("B") Operand o) {
        // Index start from 1
        r3.set(images.load(r3.getString(), r2.get() - 1).getHandler());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 20, b = 0)
    public void showPic() {
        // SHOWPIC(PAGE,PIC,DX,DY,W,H,X,Y,MODE)
        Params params = params(r3.get(), 9);
        pages.getResource(params.next()).draw(images.getResource(params.next()),
            params.next(), params.next(),
            params.next(), params.next(),
            params.next(), params.next(),
            params.next());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 21, b = 0)
    public void flipPage() {
        r3.get(pages).display();
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 22, b = 0)
    public void showPage() {
        r2.get(pages).draw(r3.get(pages));
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 23, b = 0)
    public void fillPage() {
        Params params = params(r3.get(), 6);
        pages.getResource(params.next()).fill(params.next(), params.next(), params.next(), params.next(), color(params.next()));
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 24, b = 0)
    public void pixel() {
        Params params = params(r3.get(), 4);
        pages.getResource(params.next()).pixel(params.next(), params.next(), color(params.next()));
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 25, b = 0)
    public void readPixel() {
        Params params = params(r3.get(), 3);
        int c = pages.getResource(params.next()).pixel(params.next(), params.next());
        r3.set(color(c));
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 26, b = 0)
    public void releaseImage() throws Exception {
        r3.get(images).close();
    }

    /*
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
     */
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 35, b = 0)
    public void pageClear() throws Exception {
        pages.getScreen().clear();
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 40)
    public void imageWidth() throws Exception {
        r3.set(r3.get(images).getWidth());
    }
    /*
64 | 设置画笔 | 0 | r3:参数地址 |  SETPEN(PAGE,STYLE,WID,COLOR)
65 | 设置刷子 | 0 | r2:PAGE r3:STYLE |  SETBRUSH(PAGE,STYLE)
66 | 移动画笔 | 0 | r1,r2,r3:PAGE,X,Y |  MOVETO(PAGE,X,Y)
67 | 画线 | 0 | r1,r2,r3:PAGE,X,Y |  LINETO(PAGE,X,Y)
68 | 画矩形 | 0 | r3:参数地址 |  RECTANGLE(PAGE,LEFT,TOP,RIGHT,BOTTOM)
69 | 画圆 | 0 | r3:参数地址 |  CIRCLE(PAGE,CX,CY,CR)
80 | 复制部分画布扩展 | 0 | r3:参数地址 |  STRETCHBLTPAGEEX(X,Y,WID,HGT,CX,CY,DEST,SRC)
     */

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 41)
    public void imageHeight() throws Exception {
        r3.set(r3.get(images).getHeight());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 43, b = 0)
    public void showPart() {
        Params params = params(r3.get(), 4);
        int x = params.next(), y = params.next(), dest = params.next(), src = params.next();
        pages.getResource(dest).draw(pages.getResource(src), x, y);
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 44, b = 0)
    public void setBgMode() {
        pages.getScreen().setBackgroundMode(r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 64)
    public void setPen() {
        Params params = params(r3.get(), 4);
        int page = params.next(), style = params.next(), wid = params.next(), color = params.next();
        pages.getResource(page).pen(wid, style, color(color));
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 66)
    public void moveTo() {
        r1.get(pages).move(r2.get(), r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 67)
    public void lineTo() {
        r1.get(pages).line(r2.get(), r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 68)
    public void drawRect() {
        Params params = params(r3.get(), 5);
        int page = params.next(), left = params.next(), top = params.next(), right = params.next(), bottom = params.next();
        pages.getResource(page).rectangle(left, top, right, bottom);
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 69)
    public void drawCircle() {
        Params params = params(r3.get(), 4);
        int page = params.next(), x = params.next(), y = params.next(), r = params.next();
        pages.getResource(page).circle(x, y, r);
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 80)
    public void showPartEx() {
        Params params = params(r3.get(), 8);
        int x = params.next(), y = params.next(), w = params.next(), h = params.next(), cx = params.next(), cy = params.next(), dest = params.next(), src = params.next();
        pages.getResource(dest).draw(pages.getResource(src), x, y, w, h, cx, cy);
    }

    Params params(int address, int n) {
        return new Params(vm, address, n);
    }

    /*
36 | 按行列定位光标 | 0 | r2:行,r3:列 |  LOCATE(LINE,ROW)
38 | 设定文字字体大小 | 0 | r3:FONT |  FONT(F)
42 | 按坐标定位光标 | - | r2:行,r3:列 |  PIXLOCATE(LINE,ROW)
     */
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 36, b = 0)
    public void locate() {
        pages.getScreen().locate(r2.get(), r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 38, b = 0)
    public void setFont() {
        pages.getScreen().font(r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 42, b = 0)
    public void cursor() {
        pages.getScreen().cursor(r2.get(), r3.get());
    }

    /*
37 | 设定文字颜色 | 0 | r3:参数地址 |  COLOR(FRONT,BACK,FRAME)
 */
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 37, b = 0)
    public void setFontColor() {
        // TODO Transparent background
        Params params = params(r3.get(), 3);
        pages.getScreen().font(color(params.next()), color(params.next()), params.next());
    }

    static class Params {
        private final VM vm;
        private final int limit;
        private int pos;

        Params(VM vm, int pos, int n) {
            this.vm = vm;
            this.pos = pos + n * 4;
            this.limit = pos;
        }

        public int next() {
            assert pos >= limit;
            return vm.getMemory().read(pos -= 4);
        }
    }
}
