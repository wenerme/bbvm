package me.wener.bbvm.dev.swing;

import javax.swing.*;
import java.awt.*;
import java.util.function.Supplier;

/**
 * The container for the image
 *
 * @author wener
 * @since 15/12/28
 */
class ImagePanel extends JPanel {
    private final Supplier<? extends Image> image;

    public ImagePanel(Supplier<? extends Image> image) {
        super(true);
        this.image = image;
    }

    @Override
    protected void paintComponent(Graphics g) {
        super.paintComponent(g);

        g.drawImage(image.get(), 0, 0, null);
    }

    public void refresh() {
        repaint();
    }
}
