package me.wener.bbvm.vm.res.swing;

import java.awt.*;

/**
 * @author wener
 * @since 15/12/26
 */
class StringDrawer {
    protected final Graphics2D g;
    private int x, y, w, h;
    private FontMetrics metrics;
    private int fontSize;
    private Color front = Color.WHITE;
    private Color back = Color.BLACK;
    private boolean backgroundVisible;
//        private Consumer<StringDrawer> nextPage = (d) -> d.g.fillRect(0, 0, w, h);

    StringDrawer(Graphics2D g, int w, int h) {
        this.g = g;
        this.w = w;
        this.h = h;
        getFontInfo();
        locate(1, 1);
    }

    public void fontChanged() {
        getFontInfo();
    }

    private void getFontInfo() {
        metrics = g.getFontMetrics();
        Font font = g.getFont();
        fontSize = font.getSize();
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
        int width = metrics.charWidth(ch);
        if (x + width > w) {
            nextCursorLine();
        }

        if (backgroundVisible) {
            g.setColor(back);
            g.fillRect(x, y - fontSize, width, fontSize);
        }

        // Draw char
        g.setColor(front);
        g.drawString(String.valueOf((char) ch), x, y);
        advanceCursor(width);
    }

    void nextCursorLine() {
        x = 0;
        y += fontSize;
        if (y + fontSize > h) {// Not enough
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
