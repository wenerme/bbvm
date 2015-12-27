package me.wener.bbvm.dev.swing;

import me.wener.bbvm.dev.DeviceConstants;
import me.wener.bbvm.dev.DeviceConstants.FontType;
import me.wener.bbvm.dev.Drawable;
import me.wener.bbvm.dev.PageResource;
import me.wener.bbvm.util.IntEnums;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.awt.*;
import java.awt.image.BufferedImage;

/**
 * @author wener
 * @since 15/12/26
 */
class SwingPage extends Draw implements PageResource {
    private final static Logger log = LoggerFactory.getLogger(PageResource.class);
    private final int handler;
    private final SwingPageManager manager;
    private final Point pen = new Point();
    private final StringDrawer stringDrawer;

    // Use Graphics color as pen color
//        private Color penColor = Color.WHITE;

    SwingPage(int handler, SwingPageManager manager, BufferedImage image) {
        super(image);
        this.handler = handler;
        this.manager = manager;
        Font font = new Font("楷体", Font.PLAIN, 12);// Makes a 12 height
        g.setFont(deriveFont(font, 12));
        g.setColor(Color.BLACK);
        stringDrawer = new StringDrawer(g, getWidth(), getHeight());
        stringDrawer.setBackgroundVisible(true);// Default
    }

    SwingPage(int handler, SwingPageManager manager) {
        this(handler, manager, new BufferedImage(manager.getWidth(), manager.getHeight(), BufferedImage.TYPE_INT_RGB));
    }

    Font deriveFont(Font font, int lineHeight) {
        float size = (font.getSize() * 1f / g.getFontMetrics().getHeight()) * lineHeight;
        if (size == font.getSize()) {
            return font;
        }
        return font.deriveFont(size);
    }

    @Override
    public PageResource draw(Drawable o, int dx, int dy, int w, int h, int x, int y, int mode) {
        // TODO Ignore mode
        g.drawImage(o.unwrap(Draw.class).image, dx, dy, dx + w, dy + h, x, y, x + w, y + h, null);
        return this;
    }

    @Override
    public PageResource display() {
        manager.getScreen().draw(this);
        return this;
    }

    @Override
    public PageResource draw(PageResource resource) {
        g.drawImage(resource.unwrap(Draw.class).image, 0, 0, null);
        return this;
    }

    PageResource color(int color, Runnable run) {
        Color old = g.getColor();
        g.setColor(new Color(color));
        run.run();
        g.setColor(old);
        return this;
    }

    @Override
    public PageResource fill(int x, int y, int w, int h, int color) {
        color(color, () -> g.fillRect(x, y, w, h));
        return this;
    }

    @Override
    public PageResource pixel(int x, int y, int color) {
        image.setRGB(x, y, color);
        return this;
    }

    @Override
    public int pixel(int x, int y) {
        return image.getRGB(x, y);
    }

    @Override
    public PageResource fill() {
        g.fillRect(0, 0, getWidth(), getHeight());
        return this;
    }

    @Override
    public PageResource draw(PageResource resource, int x, int y) {
        g.drawImage(resource.unwrap(Draw.class).image, x, y, null);
        return this;
    }

    @Override
    public PageResource draw(PageResource resource, int x, int y, int w, int h, int cx, int cy) {
        g.drawImage(resource.unwrap(Draw.class).image, x, y, x + w, y + h, cx, cy, cx + w, cy + h, null);
        return this;
    }

    @Override
    public PageResource pen(int width, int style, int color) {
        // TODO Ignore width and style
        g.setColor(new Color(color));
        return this;
    }

    @Override
    public PageResource circle(int cx, int cy, int r) {
        int i = r * 2;
        g.drawOval(cx - r, cy - r, i, i);
        return this;
    }

    @Override
    public PageResource rectangle(int left, int top, int right, int bottom) {
        g.drawRect(left, top, right - left, bottom - top);
        return this;
    }

    @Override
    public PageResource line(int x, int y) {
        g.drawLine(pen.x, pen.y, x, y);
        return this;
    }

    @Override
    public PageResource move(int x, int y) {
        pen.setLocation(x, y);
        return this;
    }

    @Override
    public PageResource locate(int row, int column) {
        if (row == 0 || column == 0) {
            log.debug("Bad locate {},{}", row, column);
            return this;
        }
        stringDrawer.locate(row, column);
        return this;
    }

    @Override
    public PageResource cursor(int row, int column) {
        stringDrawer.cursor(row, column);
        return this;
    }

    @Override
    public PageResource draw(String text) {
        Color color = g.getColor();
        stringDrawer.draw(text);
        g.setColor(color);
        return this;
    }

    public PageResource draw(char text) {
        Color color = g.getColor();
        stringDrawer.draw(text);
        g.setColor(color);
        return this;
    }

    @Override
    public PageResource font(int font) {
        FontType type = IntEnums.fromInt(FontType.class, font);
        if (type == null) {
            log.warn("Font {} not found", font);
        } else {
            log.info("Set font to {}", type);
            int height = g.getFontMetrics().getAscent();
            g.setFont(deriveFont(g.getFont(), type.getSize()));
            stringDrawer.fontChanged();
            // Adjust for next draw
            stringDrawer.y += g.getFontMetrics().getAscent() - height;
        }
        return this;
    }

    @Override
    public PageResource font(int frontColor, int backColor, int frame) {
        // TODO Ignore frame
        stringDrawer.setFront(new Color(frontColor));
        stringDrawer.setBack(new Color(backColor));
        return this;
    }

    @Override
    public PageResource setBackgroundMode(int mode) {
        stringDrawer.setBackgroundVisible(mode == DeviceConstants.BackgroundMode.OPAQUE.asInt());
        return this;
    }

    @Override
    public int getHandler() {
        return handler;
    }

    @Override
    public SwingPageManager getManager() {
        return manager;
    }

    @Override
    public void close() {
        manager.close(this);
    }

    public SwingPage fill(int c) {
        fill(0, 0, getWidth(), getHeight(), c);
        return this;
    }
}
