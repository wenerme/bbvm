package me.wener.bbvm.dev.swing;

import java.awt.*;
import java.awt.geom.Rectangle2D;

/**
 * @author wener
 * @since 15/12/26
 */
class StringDrawer {
    protected final Graphics2D g;
    int x, y, w, h;
    private FontMetrics metrics;
    private int fontHeight;
    private Color front = Color.WHITE;
    private Color back = Color.BLACK;
    private boolean backgroundVisible;
    private FontMetrics fm;
    //        private Consumer<StringDrawer> nextPage = (d) -> d.g.fillRect(0, 0, w, h);

    StringDrawer(Graphics2D g, int w, int h) {
        this.g = g;
        this.w = w;
        this.h = h;
        fontChanged();
        locate(1, 1);
    }

    /**
     * Notify font changed
     */
    public void fontChanged() {
        metrics = g.getFontMetrics();
//        Font font = g.getFont();
        fm = g.getFontMetrics();
        fontHeight = fm.getHeight();
    }

    /**
     * Start from 1
     */
    public void locate(int row, int column) {
        int width = g.getFont().getSize();
        x = (column - 1) * width / 2;
        y = row * width;
    }

    public void cursor(int row, int column) {
        x = column;
        y = row;
    }

    public void draw(String text) {
        text.chars().forEach(this::draw);
    }

    public void draw(int ch) {
        switch (ch) {
            case '\n':
                nextCursorLine();
                break;
            case '\t':
                for (int i = 0; i < 8; i++) {
                    drawNormal(' ');
                }
                break;
            default:
                drawNormal(ch);
        }
    }

    private void advanceCursor(int width) {
        x += width;
        if (x > w) {
            nextCursorLine();
        }
    }

    void drawNormal(int ch) {
        String s = String.valueOf((char) ch);
        int width = metrics.charWidth(ch);
        if (x + width > w) {
            nextCursorLine();
        }
        if (backgroundVisible) {
            Rectangle2D rect = fm.getStringBounds(s, g);
            g.setColor(back);
            g.fillRect(x,
                    y - fm.getAscent(),
                    (int) rect.getWidth(),
                    fontHeight);
        }
        // Draw char
        g.setColor(front);
        g.drawString(s, x, y);
        advanceCursor(width);
    }

    void nextCursorLine() {
        x = 0;
        y += fontHeight;
        if (y + fontHeight > h) {// Not enough
            nextPage();
        }
    }

    private void nextPage() {
        locate(1, 0);
    }

    public StringDrawer setFront(Color front) {
        this.front = front;
        return this;
    }

    public StringDrawer setBack(Color back) {
        this.back = back;
        return this;
    }

    public void setBackgroundVisible(boolean b) {
        this.backgroundVisible = b;
    }
}
