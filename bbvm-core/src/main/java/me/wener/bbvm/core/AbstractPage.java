package me.wener.bbvm.core;

public abstract class AbstractPage<P extends Page> implements Page
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

    @Override
    public void pen(PenStyle penStyle, int width, Colour color)
    {
        this.penStyle = penStyle;
        this.penSize = width;
        this.penColor = color;
    }

    @Override
    public void setBrushStyle(BrushStyle brushStyle)
    {
        this.brushStyle = brushStyle;
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


    public int getFontSize()
    {
        return 0;
    }

    @Override
    public void print(String... v)
    {

    }

    @Override
    public void print(String v)
    {

    }

    @Override
    public void clear()
    {
        fill(0, 0, getWidth(), getHeight(), background);
    }

    @Override
    public void setBgMode(BackgroundMode backgroundMode)
    {
        this.bgMode = backgroundMode;
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
    public void setFontType(FontType fontType)
    {
        this.fontType = fontType;
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

    public PenStyle getPenStyle()
    {
        return penStyle;
    }

    public BackgroundMode getBgMode()
    {
        return bgMode;
    }

    public FontType getFontType()
    {
        return fontType;
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
