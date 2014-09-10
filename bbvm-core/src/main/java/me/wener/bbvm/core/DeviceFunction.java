package me.wener.bbvm.core;

/**
 * BB 图形设备相关的函数,多用于图形图像的操作
 */
@SuppressWarnings("unused")
public interface DeviceFunction extends BBFunction
{
    /**
     * 在屏幕上显示数据<br>
     * <b>语法：</b>PRINT [表达式1],[表达式2]...
     *
     * @param v 需要打印的值,调用 {@link #toString}
     */
    void PRINT(Object... v);

    /**
     * 清空屏幕<br>
     * <b>语法：</b>CLS [0|1|2]<br>
     * CLS后边所的带的数字，目前BB不做区分，都只执行一个功能。清空整个屏幕，并将输出光标定在(0,0)
     */
    void CLS();

    /**
     * 移动屏幕输出光标到指定位置(象素坐标)<br>
     * 注意:这个函数的参数顺序是和LOCATE不同的，LOCATE是先行坐标即Y坐标，而这个函数是先X坐标。
     *
     * @param x X坐标
     * @param y Y坐标
     */
    void PIXLOCATE(int x, int y);

    /**
     * 设置当前的显示字体
     */
    void FONT(int type);

    int GETPENPOSX(int param);

    int GETPENPOSY(int param);

    void SETLCD(int WIDTH, int HEIGHT);

    /**
     * 从键盘获取一个字符，并返回该字符的字符串。
     * 本函数不暂停程序的执行。当键盘缓冲区中不存在字符时，本函数返回空字符串""；
     * 当键盘缓冲中存在输入字符时，本函数返回一个字符的字符串，比如返回"A"。返回字符串时，屏幕不会显示返回的字符串。
     *
     * @return 从键盘获取的字符
     */
    String INKEY$();

    /**
     * 从键盘获取一个字符，并返回该字符的ASCII。
     * 本函数除了返回值是字符的ASCII外，其他的和INKEY$完全相同。
     *
     * @return 字符的ASCII
     */
    int INKEY();

    /**
     * 从指定的rlb文件中载入图象资源。<br>
     * 资源库是由上面步骤生成的RLB文件，到具体的词典机型上时，可以使用由工具转换而成的lib文件。
     * 资源ID就是在工具界面中文件名前面的数字。
     * 例如这里的数字1就是bootpic.bmp的资源ID。
     *
     * @param FILE$ 资源库文件名
     * @param ID    资源ID
     * @return 资源句柄
     */
    int LOADRES(String FILE$, int ID);

    /**
     * @param FILE  文件句柄
     * @param ID    资源ID
     * @return 资源句柄
     */
    int LOADRES(int FILE, int ID);

    /**
     * 释放通过LOADRES函数载入的资源。
     *
     * @param PIC 资源句柄
     */
    void FREERES(int PIC);

    /**
     * 获取指定资源的象素宽度
     *
     * @param PIC 资源句柄
     * @return 资源的象素宽度
     */
    int GETPICWID(int PIC);

    /**
     * 获取指定资源的象素高度
     *
     * @param PIC 资源句柄
     * @return 资源的象素高度
     */
    int GETPICHGT(int PIC);

    /**
     * 创建一个和屏幕兼容的页面<br>
     * page为一个整型变量。用来存放返回的页面。
     * 最多创建10个页面。
     *
     * @return 页面句柄
     */
    int CREATEPAGE();

    /**
     * 删除指定的页面
     *
     * @param PAGE 页面句柄
     */
    void DELETEPAGE(int PAGE);

    /**
     * 显示指定的图象资源到指定的页面。<br>
     * 当page=-1时，图片将直接显示到前台屏幕上。
     *
     * @param PAGE 页面句柄
     * @param PIC  资源句柄
     * @param DX   把图片显示到的屏幕坐标 X
     * @param DY   把图片显示到的屏幕坐标 Y
     * @param W    显示的宽度
     * @param H    显示的高度
     * @param X    图片中开始显示的坐标 X
     * @param Y    图片中开始显示的坐标 Y
     * @param MODE 显示模式，目前只支持1，为透明模式。
     *             若图片使用了关键颜色，则关键颜色就呈透明显示。关键颜色为RGB=255,0,255的紫色。
     */
    void SHOWPIC(int PAGE, int PIC, int DX, int DY, int W, int H, int X, int Y, int MODE);

    /**
     * 将页面中的内容映射到屏幕上
     *
     * @param PAGE 页面句柄
     */
    void FLIPPAGE(int PAGE);

    /**
     * 页面对拷
     *
     * @param DEST 目标页面
     * @param SRC  源页面
     */
    void BITBLTPAGE(int DEST, int SRC);

    /**
     * 页面对拷 增强版
     *
     * @param X    SRC 到 DEST 的坐标 X
     * @param Y    SRC 到 DEST 的坐标 Y
     * @param DEST 目标页面
     * @param SRC  源页面
     */
    void STRETCHBLTPAGE(int X, int Y, int DEST, int SRC);

    void STRETCHBLTPAGEEX(int X, int Y, int WID, int HGT, int CX, int CY, int DEST, int SRC);

    /**
     * 用指定的颜色填充页面.
     * 当page=-1时，将直接填充前台屏幕。
     *
     * @param PAGE  页面句柄
     * @param X     坐标
     * @param Y     坐标
     * @param WID   宽
     * @param HGT   高
     * @param COLOR 颜色
     */
    void FILLPAGE(int PAGE, int X, int Y, int WID, int HGT, int COLOR);

    /**
     * 像指定的页面画点
     * 当page=-1时，将直接画到前台屏幕上
     *
     * @param PAGE  页面句柄
     * @param X     坐标
     * @param Y     坐标
     * @param COLOR 颜色
     */
    void PIXEL(int PAGE, int X, int Y, int COLOR);

    /**
     * 读取指定页面上一点的颜色值
     * 当page=-1时，将直接从前台屏幕上读取
     *
     * @param PAGE 页面句柄
     * @param X    坐标
     * @param Y    坐标
     * @return 颜色
     */
    int READPIXEL(int PAGE, int X, int Y);

    /**
     * 设置屏幕上的字体颜色，字体背景颜色和边框颜色
     * 由于这个函数是和print语句搭配使用的，而print只能将内容显示到前台屏幕上。所以这里没有page参数
     *
     * @param FRONT 字体颜色
     * @param BACK  背景颜色
     * @param FRAME 边框颜色,由于目前没提供创建子窗口的window指令，所以这个参数并无实际意义。
     */
    void COLOR(int FRONT, int BACK, int FRAME);

    /**
     * 设置字体是否透明显示
     * 可以为：TRANSPARENT或OPAQUE，这两个值在stdlib.lib文件中一定预定义好了。
     * TRANSPARENT为透明显示，即字体的背景颜色无效。
     * OPAQUE为不透明显示，即字体的背景颜色有效
     *
     * @param mode TRANSPARENT或OPAQUE
     */
    void SETBKMODE(int mode);

    /**
     * 设置画笔<br>
     * 设置的画笔，是和PAGE关联的。就是说，对一个PAGE设置了一种画笔后，不会影响另一个PAGE的画笔。画笔被设置后，虽有绘图函数都将使用画笔的属性（颜色，线宽和样式）绘制边框线(直线是特殊的边框线)。目前画笔只支持PEN_SOLID(实线模式)，线宽只支持1。
     *
     * @param PAGE  页面句柄
     * @param STYLE 画笔央视
     * @param WID   线宽
     * @param COLOR 颜色
     */
    void SETPEN(int PAGE, int STYLE, int WID, int COLOR);

    /**
     * 设置画刷
     * 画刷和画笔的作用机制类似。画刷影响的是对图形的填充样式。目前支持的样式只有BRUSH_SOLID(实心填充模式)
     *
     * @param PAGE  页面句柄
     * @param STYLE 样式
     */
    void SETBRUSH(int PAGE, int STYLE);

    /**
     * 设置画线的起始点<br>
     * 设置画线的起始点。当执行过画线函数后，该点将被自动更新为画线函数的终点。
     *
     * @param PAGE 页面
     * @param X    X 坐标
     * @param Y    Y 坐标
     */
    void MOVETO(int PAGE, int X, int Y);

    /**
     * 从起始点画线到目标点<br>
     * 画线的起始点是MOVETO函数设置的。执行本函数后，下次画线的起始点，将被更新成本函数中的坐标
     *
     * @param PAGE 页面
     * @param X    X 坐标
     * @param Y    Y 坐标
     */
    void LINETO(int PAGE, int X, int Y);

    /**
     * 画矩形
     *
     * @param PAGE   页面
     * @param LEFT   左
     * @param TOP    上
     * @param RIGHT  右
     * @param BOTTOM 下
     */
    void RECTANGLE(int PAGE, int LEFT, int TOP, int RIGHT, int BOTTOM);

    /**
     * 画圆
     *
     * @param PAGE 页面
     * @param CX   圆心 X
     * @param CY   圆心 Y
     * @param CR   半径
     */
    void CIRCLE(int PAGE, int CX, int CY, int CR);

    /**
     * 程序执行时从键盘输入数据
     * <b>语法：</b>INPUT "输入提示字语";变量列表
     * BB中INPUT语句没做相应的数据类型检查。所以输入数据的类型和个数请在输入提示语中明确指出。
     *
     * @param PROMOTE 输入提示字语
     * @return 输入的数值列表
     */
    String[] INPUT(String PROMOTE, int n);

    /**
     * 检测指定键的状态<br>
     * 该函数检测到指定键的状态后立即返回，程序不会在函数内部停留等待。
     *
     * @param KEYCODE 指定的KEYCODE是ASCII码。
     * @return 若指定键按下返回1，否则返回0
     */
    boolean KEYPRESS(int KEYCODE);

    /**
     * 等待并返回一个按键的键值<br>
     * 该程序进入函数后，若没有按键按下，程序会一直在这个函数中等待，
     * 直到有按键按下才将被按下的键值返回。
     *
     * @return 按键的键值
     */
    int WAITKEY();

    /**
     * 移动屏幕输出光标到指定位置。<br>
     * 行列的计算方法是根据当前被设置的字体(BB中可通过FONT函数设置)的字符宽度和高度来计算的。
     * <pre>
     * 例如：当前被设置的字体是12×12的宋体字，则对应行列的象素坐标是：
     *     行象素坐标 ＝ (行 - 1) × 12
     *     列象素坐标 ＝ (列 - 1) × 6
     * </pre>
     *
     * @param row    行
     * @param column 列
     */
    void LOCATE(int row, int column);
}
