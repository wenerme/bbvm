package me.wener.bbvm.core;

import static me.wener.bbvm.core.Values.fromValue;

import me.wener.bbvm.core.constant.BackgroundMode;
import me.wener.bbvm.core.constant.DrawMode;
import me.wener.bbvm.core.constant.FontType;
import me.wener.bbvm.core.constant.PenStyle;

public class AbstractDeviceFunction<SCR extends Screen<PAGE>, PAGE extends Page, PIC extends Picture> implements DeviceFunction
{
    private final AbstractDevice<SCR, PAGE, PIC> device;
    private final SCR screen;
    private final ResourceHandlePool<PIC> picturePool;
    private final ResourceHandlePool<PAGE> pagePool;
    private final PAGE scrPage;

    public AbstractDeviceFunction(AbstractDevice<SCR, PAGE, PIC> device)
    {
        this.device = device;
        screen = device.getScreen();
        scrPage = screen.asPage();
        picturePool = device.getPicturePool();
        pagePool = device.getPagePool();
    }

    @Override
    public void PRINT(Object... v)
    {
        for (Object o : v)
        {
            if (o instanceof String)
                scrPage.print((String) o);
            else
                scrPage.print(o.toString());
        }
    }

    @Override
    public void CLS()
    {
        scrPage.clear();
    }

    @Override
    public void PIXLOCATE(int x, int y)
    {
        scrPage.cursor(x, y);
    }

    @Override
    public void FONT(int type)
    {
        scrPage.setFontType(fromValue(FontType.class, type));
    }

    @Override
    public int GETPENPOSX(int param)
    {
        return 0;
    }

    @Override
    public int GETPENPOSY(int param)
    {
        return 0;
    }

    @Override
    public void SETLCD(int WIDTH, int HEIGHT)
    {
        device.setScreenSize(WIDTH, HEIGHT);
    }

    @Override
    public String INKEY$()
    {
        return null;
    }

    @Override
    public int INKEY()
    {
        return 0;
    }

    @Override
    public int LOADRES(String FILE$, int ID)
    {
        return device.loadPicture(FILE$,ID);
    }

    @Override
    public int LOADRES(int FILE, int ID)
    {
        return device.loadPicture(FILE,ID);
    }

    @Override
    public void FREERES(int PIC)
    {
        picturePool.release(PIC);
    }

    @Override
    public int GETPICWID(int PIC)
    {
        return picturePool.getResource(PIC).getHeight();
    }

    @Override
    public int GETPICHGT(int PIC)
    {
        return picturePool.getResource(PIC).getHeight();
    }

    @Override
    public int CREATEPAGE()
    {
        return pagePool.acquire();
    }

    @Override
    public void DELETEPAGE(int PAGE)
    {
        pagePool.release(PAGE);
    }

    @Override
    public void SHOWPIC(int PAGE, int PIC, int DX, int DY, int W, int H, int X, int Y, int MODE)
    {
        pagePool.getResource(PAGE)
                .draw(picturePool.getResource(PIC), DX, DY, W, H, X, Y, fromValue(DrawMode.class, MODE));

    }

    @Override
    public void FLIPPAGE(int PAGE)
    {
        screen.showPage(pagePool.getResource(PAGE));
    }

    @Override
    public void BITBLTPAGE(int DEST, int SRC)
    {
        pagePool.getResource(DEST).draw(pagePool.getResource(SRC));
    }

    @Override
    public void STRETCHBLTPAGE(int X, int Y, int DEST, int SRC)
    {
        pagePool.getResource(DEST).draw(pagePool.getResource(SRC), X,Y);
    }

    @Override
    public void STRETCHBLTPAGEEX(int X, int Y, int WID, int HGT, int CX, int CY, int DEST, int SRC)
    {
        pagePool.getResource(DEST).draw(pagePool.getResource(SRC), X, Y, WID, HGT, CX, CY);
    }

    @Override
    public void FILLPAGE(int PAGE, int X, int Y, int WID, int HGT, int COLOR)
    {
        pagePool.getResource(PAGE).fill(X, Y, WID, HGT, Colour.fromARGB(COLOR));
    }

    @Override
    public void PIXEL(int PAGE, int X, int Y, int COLOR)
    {
        pagePool.getResource(PAGE).pixel(X,Y,Colour.fromARGB(COLOR));
    }

    @Override
    public int READPIXEL(int PAGE, int X, int Y)
    {
        return pagePool.getResource(PAGE).pixel(X,Y).getRGB();
    }

    @Override
    public void COLOR(int FRONT, int BACK, int FRAME)
    {
        scrPage.color(Colour.fromARGB(FRONT), Colour.fromARGB(BACK));
    }

    @Override
    public void SETBKMODE(int mode)
    {
        scrPage.setBgMode(fromValue(BackgroundMode.class, mode));
    }

    @Override
    public void SETPEN(int PAGE, int STYLE, int WID, int COLOR)
    {
        pagePool.getResource(PAGE)
                .pen(fromValue(PenStyle.class, STYLE), WID, Colour.fromARGB(COLOR));
    }

    @Override
    public void SETBRUSH(int PAGE, int STYLE)
    {
        pagePool.getResource(PAGE)
                .setBrushStyle(fromValue(BrushStyle.class, STYLE));
    }

    @Override
    public void MOVETO(int PAGE, int X, int Y)
    {
        pagePool.getResource(PAGE)
                .moveTo(X, Y);
    }

    @Override
    public void LINETO(int PAGE, int X, int Y)
    {
        pagePool.getResource(PAGE)
                .lineTo(X, Y);
    }

    @Override
    public void RECTANGLE(int PAGE, int LEFT, int TOP, int RIGHT, int BOTTOM)
    {
        pagePool.getResource(PAGE)
                .rectangle(LEFT, TOP, RIGHT, BOTTOM);
    }

    @Override
    public void CIRCLE(int PAGE, int CX, int CY, int CR)
    {
        pagePool.getResource(PAGE)
                .circle(CX, CY, CR);
    }

    @Override
    public String[] INPUT(String PROMOTE, int n)
    {
        return new String[0];
    }

    @Override
    public boolean KEYPRESS(int KEYCODE)
    {
        return device.isKeyPressed(KEYCODE);
    }

    @Override
    public int WAITKEY()
    {
        return device.waitkey();
    }

    @Override
    public void LOCATE(int row, int column)
    {
        scrPage.locate(row, column);
    }
}
