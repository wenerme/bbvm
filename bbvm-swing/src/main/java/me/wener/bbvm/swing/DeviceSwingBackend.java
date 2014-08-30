package me.wener.bbvm.swing;

import java.awt.Image;
import me.wener.bbvm.core.constant.Device;
import me.wener.bbvm.core.DeviceFunction;

public class DeviceSwingBackend implements Device
{
    private final SwingDeviceFunction function = new SwingDeviceFunction();
    @Override
    public DeviceFunction getFunction()
    {
        return function;
    }

    public Image getScreen()
    {
        return null;
    }

    class SwingDeviceFunction implements DeviceFunction
    {
        @Override
        public void PRINT(Object... v)
        {

        }

        @Override
        public void CLS()
        {

        }

        @Override
        public void PIXLOCATE(int x, int y)
        {

        }

        @Override
        public void FONT(int type)
        {

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

        }

        @Override
        public int INSTR(int N, String S$, String B$)
        {
            return 0;
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
            return 0;
        }

        @Override
        public void FREERES(int PIC)
        {

        }

        @Override
        public int GETPICWID(int PIC)
        {
            return 0;
        }

        @Override
        public int GETPICHGT(int PIC)
        {
            return 0;
        }

        @Override
        public int CREATEPAGE()
        {
            return 0;
        }

        @Override
        public void DELETEPAGE(int PAGE)
        {

        }

        @Override
        public void SHOWPIC(int PAGE, int PIC, int DX, int DY, int W, int H, int X, int Y, int MODE)
        {

        }

        @Override
        public void FLIPPAGE(int PAGE)
        {

        }

        @Override
        public void BITBLTPAGE(int DEST, int SRC)
        {

        }

        @Override
        public void STRETCHBLTPAGE(int X, int Y, int DEST, int SRC)
        {

        }

        @Override
        public void STRETCHBLTPAGEEX(int X, int Y, int WID, int HGT, int CX, int CY, int DEST, int SRC)
        {

        }

        @Override
        public void FILLPAGE(int PAGE, int X, int Y, int WID, int HGT, int COLOR)
        {

        }

        @Override
        public void PIXEL(int PAGE, int X, int Y, int COLOR)
        {

        }

        @Override
        public int READPIXEL(int PAGE, int X, int Y)
        {
            return 0;
        }

        @Override
        public void COLOR(int FRONT, int BACK, int FRAME)
        {

        }

        @Override
        public void SETBKMODE(int mode)
        {

        }

        @Override
        public void SETPEN(int PAGE, int STYLE, int WID, int COLOR)
        {

        }

        @Override
        public void SETBRUSH(int PAGE, int STYLE)
        {

        }

        @Override
        public void MOVETO(int PAGE, int X, int Y)
        {

        }

        @Override
        public void LINETO(int PAGE, int X, int Y)
        {

        }

        @Override
        public void RECTANGLE(int PAGE, int LEFT, int TOP, int RIGHT, int BOTTOM)
        {

        }

        @Override
        public void CIRCLE(int PAGE, int CX, int CY, int CR)
        {

        }

        @Override
        public String[] INPUT(String PROMOTE, int n)
        {
            return new String[0];
        }

        @Override
        public boolean KEYPRESS(int KEYCODE)
        {
            return false;
        }

        @Override
        public int WAITKEY()
        {
            return 0;
        }

        @Override
        public void LOCATE(int row, int column)
        {

        }
    }
}
