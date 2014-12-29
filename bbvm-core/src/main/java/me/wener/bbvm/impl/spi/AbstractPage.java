package me.wener.bbvm.impl.spi;

import me.wener.bbvm.api.Page;
import me.wener.bbvm.def.BackgroundMode;
import me.wener.bbvm.def.BrushStyle;
import me.wener.bbvm.def.FontType;
import me.wener.bbvm.def.PenStyle;
import me.wener.bbvm.impl.plaf.Colour;

public abstract class AbstractPage implements Page
{
    protected FontType fontType = FontType.FONT_12SONG;
    protected BackgroundMode bgMode = BackgroundMode.OPAQUE;
    protected Colour background = Colour.black;
    protected Colour foreground = Colour.white;
    protected int cursorX;
    protected int cursorY;
    protected PenStyle penStyle = PenStyle.PEN_SOLID;
    protected int penSize;
    protected Colour penColor = Colour.white;
    protected BrushStyle brushStyle = BrushStyle.BRUSH_SOLID;
    protected int penX;
    protected int penY;
    protected int width;
    protected int height;

    protected AbstractPage(int height, int width)
    {
        this.height = height;
        this.width = width;
    }

    @Override
    public int getHeight()
    {
        return height;
    }

    @Override
    public int getWidth()
    {
        return width;
    }

    @Override
    public void pen(PenStyle penStyle, int thickness, Colour color)
    {
        this.penStyle = penStyle;
        this.penSize = thickness;
        this.penColor = color;
    }

    @Override
    public void lineTo(int x, int y)
    {
        drawLine(penX, penY, x, y);
        penX = x;
        penY = y;
    }

    @Override
    public void moveTo(int x, int y)
    {
        penX = x;
        penY = y;
    }

    public int getCharSize(char c)
    {
        return c < 127 ? getFontSize() / 2 : getFontSize();
    }

    /**
     * 字体大小为一个全角字符的大小
     */
    public abstract int getFontSize();

    @Override
    public void print(String... strings)
    {
        for (String string : strings)
        {
            print(string);
        }
    }

    @Override
    public void print(String v)
    {
        for (char c : v.toCharArray())
        {
            int size = -1;
            switch (c)
            {
                case '\n':
                    // 换行
                    nextCursorLine();
                    continue;
                case '\t':
                    // 一个制表符算四个
                    size = getFontSize() * 4/2;
                    break;
            }

            if(size < 0)
                size = getCharSize(c);
            if (size + cursorX > width)
                nextCursorLine();

            drawChar(c, cursorX, cursorY);
        }
    }

    private void testCursorPosition()
    {
        if (cursorX > width)
            nextCursorLine();
    }
    private void nextCursorLine()
    {
        cursorX = 0;
        cursorY += getFontSize();
    }

    @Override
    public void drawChar(char c, int x, int y)
    {
        drawString(String.valueOf(c), x, y);
    }

    @Override
    public void clear()
    {
        fill(0, 0, getWidth(), getHeight(), background);
    }

    @Override
    public void locate(int row, int column)
    {
        cursorX = column * getFontSize();
        cursorY = row * getFontSize();
    }

    @Override
    public void color(Colour front, Colour back)
    {
        this.foreground = front;
        this.background = back;
    }

    @Override
    public void cursor(int x, int y)
    {
        cursorX = x;
        cursorY = y;
    }

    public BrushStyle getBrushStyle()
    {
        return brushStyle;
    }

    @Override
    public void setBrushStyle(BrushStyle brushStyle)
    {
        this.brushStyle = brushStyle;
    }

    public PenStyle getPenStyle()
    {
        return penStyle;
    }

    public BackgroundMode getBgMode()
    {
        return bgMode;
    }

    @Override
    public void setBgMode(BackgroundMode backgroundMode)
    {
        this.bgMode = backgroundMode;
    }

    public FontType getFontType()
    {
        return fontType;
    }

    @Override
    public void setFontType(FontType fontType)
    {
        this.fontType = fontType;
    }

    public Colour getBackground()
    {
        return background;
    }

    public Colour getForeground()
    {
        return foreground;
    }
}
