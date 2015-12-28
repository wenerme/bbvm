package me.wener.bbvm.dev.swing;

import javax.swing.*;
import java.awt.*;
import java.awt.event.ComponentAdapter;
import java.awt.event.ComponentEvent;
import java.awt.image.BufferedImage;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.function.Supplier;

/**
 * @author wener
 * @since 15/12/28
 */
// Handle the image change on reset ---- by using supplier
// Create a custom component to contain image and refresh ----  by using image panel
// Handle image resize ---- by monitor image
public class MainFrame extends JFrame {
    private final Supplier<BufferedImage> image;
    private ImagePanel imagePanel;
    private ScheduledExecutorService refresher;

    public MainFrame(Supplier<BufferedImage> image) throws HeadlessException {
        super("BBVM");
        this.image = new MonitoredImage(image);
        initialize();
    }

    protected void initialize() {
        setFocusTraversalKeysEnabled(false);// Make VK_TAB works
        setDefaultCloseOperation(WindowConstants.EXIT_ON_CLOSE);
        setResizable(false);

        imagePanel = new ImagePanel(image);
        imagePanel.setLocation(0, 0);
        getContentPane().add(imagePanel);
        getImage();

        addComponentListener(new ComponentAdapter() {
            @Override
            public void componentShown(ComponentEvent e) {
                refresher = Executors.newSingleThreadScheduledExecutor();
                refresher.scheduleAtFixedRate(imagePanel::refresh, 0, 1000 / 16, TimeUnit.MILLISECONDS);
            }

            @Override
            public void componentHidden(ComponentEvent e) {
                refresher.shutdownNow();
            }
        });
    }

    private Image getImage() {
        return image.get();
    }

    public void setImageSize(int width, int height) {
        SwingUtilities.invokeLater(() -> {
            setSize(width, height);
            imagePanel.setPreferredSize(getSize());
            pack();
            setLocationRelativeTo(null);
        });
    }

    public ImagePanel getImagePanel() {
        return imagePanel;
    }

    class MonitoredImage implements Supplier<BufferedImage> {
        private final Supplier<BufferedImage> image;
        int w, h;

        MonitoredImage(Supplier<BufferedImage> image) {
            this.image = image;
        }

        @Override
        public BufferedImage get() {
            BufferedImage i = image.get();
            if (i.getWidth() != w || i.getHeight() != h) {
                setImageSize(w = i.getWidth(), h = i.getHeight());
            }
            return i;
        }
    }
}
