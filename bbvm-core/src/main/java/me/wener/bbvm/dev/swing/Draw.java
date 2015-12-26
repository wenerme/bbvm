package me.wener.bbvm.dev.swing;

import me.wener.bbvm.dev.Drawable;

import java.awt.*;
import java.awt.image.BufferedImage;

/**
 * @author wener
 * @since 15/12/26
 */
class Draw implements Drawable {
    protected final BufferedImage image;
    protected final Graphics2D g;

    Draw(BufferedImage image) {
        this.image = image;
        g = image.createGraphics();
//            g.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);
    }

    public int getWidth() {
        return image.getWidth();
    }

    public int getHeight() {
        return image.getHeight();
    }
}
