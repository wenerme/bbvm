package me.wener.bbvm.vm.res;

import java.awt.*;
import java.awt.image.BufferedImage;

/**
 * @author wener
 * @since 15/12/18
 */
public class Swings {
    private static class PageMgr implements PageManager {
        int handler = 1;

        @Override
        public PageResource screen() {
            return getResource(0);
        }

        @Override
        public int getWidth() {
            return 0;
        }

        @Override
        public int getHeight() {
            return 0;
        }

        @Override
        public PageResource getResource(int handler) {
            return null;
        }

        public void close(Page page) {

        }
    }

    private static class ImageMgr implements ImageManager {

        @Override
        public PageResource load(String file, int index) {
            return null;
        }

        @Override
        public ImageResource getResource(int handler) {
            return null;
        }
    }

    private static class Draw implements Drawable {
        protected final BufferedImage image;
        protected final Graphics2D g;

        private Draw(BufferedImage image) {
            this.image = image;
            g = image.createGraphics();
        }

        public int getWidth() {
            return image.getWidth();
        }

        public int getHeight() {
            return image.getHeight();
        }
    }


    private static class Image extends Draw implements ImageResource {
        private final int handler;
        private final ImageMgr manager;

        private Image(int handler, ImageMgr manager, BufferedImage image) {
            super(image);
            this.handler = handler;
            this.manager = manager;
        }

        @Override
        public int getHandler() {
            return handler;
        }

        @Override
        public ImageManager getManager() {
            return manager;
        }

        @Override
        public void close() throws Exception {
        }
    }

    private static class Page extends Draw implements PageResource {
        private final int handler;
        private final PageMgr manager;

        private Page(int handler, PageMgr manager, BufferedImage image) {
            super(image);
            this.handler = handler;
            this.manager = manager;
        }

        private Page(int handler, PageMgr manager) {
            this(handler, manager, new BufferedImage(manager.getWidth(), manager.getHeight(), BufferedImage.TYPE_INT_RGB));
        }

        @Override
        public PageResource show(Drawable o, int dx, int dy, int w, int h, int x, int y, int mode) {
            // TODO Ignore mode
            g.drawImage(o.unwrap(Draw.class).image, dx, dy, dx + w, dy + h, w, y, x + w, y + h, null);
            return this;
        }

        @Override
        public PageResource display() {
            manager.screen().show(this);
            return this;
        }

        @Override
        public PageResource show(PageResource resource) {
            image.copyData(resource.unwrap(Draw.class).image.getRaster());
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
            // TODO Color transform
            color(color, () -> {
                g.fillRect(x, y, w, h);
            });
            return this;
        }

        @Override
        public PageResource pixel(int x, int y, int color) {
            // TODO Color transform
            image.setRGB(x, y, color);
            return this;
        }

        @Override
        public int pixel(int x, int y) {
            // TODO Color transform
            return image.getRGB(x, y);
        }

        @Override
        public PageResource clear() {
            java.awt.Color old = g.getColor();
            g.setColor(java.awt.Color.BLACK);
            g.fillRect(0, 0, getWidth(), getHeight());
            g.setColor(old);
            return this;
        }

        @Override
        public PageResource show(PageResource resource, int x, int y) {
            g.drawImage(resource.unwrap(Draw.class).image, x, y, null);
            return this;
        }

        @Override
        public PageResource show(PageResource resource, int x, int y, int w, int h, int cx, int cy) {
            g.drawImage(resource.unwrap(Draw.class).image, x, y, x + w, y + h, cx, cy, cx + w, cy + h, null);
            return this;
        }

        @Override
        public int getHandler() {
            return handler;
        }

        @Override
        public PageMgr getManager() {
            return manager;
        }

        @Override
        public void close() throws Exception {
            manager.close(this);
        }
    }
}
