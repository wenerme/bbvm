package me.wener.bbvm.core;

import me.wener.bbvm.core.constant.DrawMode;
import me.wener.bbvm.core.constant.FontType;
import me.wener.bbvm.core.constant.PenStyle;

public interface Page
{
    void pen(PenStyle penStyle, int width, Colour color);

    void setBrushStyle(BrushStyle brushStyle);

    void lineTo(int x, int y);

    void moveTo(int x, int y);

    void rectangle(int left, int top, int right, int bottom);

    void circle(int cx, int cy, int cr);

    Colour pixel(int x, int y);

    void pixel(int x, int y, Colour color);

    void fill(int x, int y, int width, int height, Colour color);

    void draw(Page src, int x, int y, int width, int height, int srcX, int srcY);

    void draw(Page src, int x, int y);

    void drawLine(int x, int y, int dx, int dy);

    void draw(Page src);

    void drawString(String content, int x, int y);

    void draw(Picture picture, int destX, int destY, int width, int height, int x, int y, DrawMode drawMode);

    int getWidth();
    int getHeight();

    void print(String... v);

    void print(String v);

    void clear();

    void setBgMode(BackgroundMode backgroundMode);

    void locate(int row, int column);

    void color(Colour front, Colour back);

    void setFontType(FontType fontType);

    void cursor(int x, int y);
}
