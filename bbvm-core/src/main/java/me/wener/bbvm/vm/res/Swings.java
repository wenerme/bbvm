package me.wener.bbvm.vm.res;

import com.google.common.collect.Maps;
import com.google.common.eventbus.EventBus;
import com.google.common.eventbus.Subscribe;
import com.google.inject.AbstractModule;
import com.google.inject.Module;
import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.vm.event.ResetEvent;
import me.wener.bbvm.vm.event.VmTestEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.imageio.ImageIO;
import javax.inject.Inject;
import javax.inject.Singleton;
import java.awt.*;
import java.awt.image.BufferedImage;
import java.io.File;
import java.util.Map;
import java.util.NavigableSet;

/**
 * @author wener
 * @since 15/12/18
 */
public class Swings {
    public static Module graphModule() {
        return new SwingModule();
    }

    static class SwingModule extends AbstractModule {

        @Override
        protected void configure() {
            bind(PageManager.class).to(PageMgr.class);
            bind(ImageManager.class).to(ImageMgr.class);
        }

    }

    @Singleton
    static class PageMgr implements PageManager {
        private final static Logger log = LoggerFactory.getLogger(PageManager.class);
        private final Map<Integer, Page> resources = Maps.newConcurrentMap();
        private int handler = 0;
        private NavigableSet<Integer> handlers;
        private int width, height;

        @Override
        public PageManager reset() {
            resources.forEach((k, v) -> v.close());
            return this;
        }

        @Override
        public PageResource create() {
            Page page = new Page(handler++, this);
            resources.put(page.getHandler(), page);
            return page;
        }

        @Override
        public PageResource getScreen() {
            return getResource(0);
        }

        @Override
        public int getWidth() {
            return width;
        }

        @Override
        public int getHeight() {
            return height;
        }

        @Override
        public PageManager setSize(int w, int h) {
            log.info("{} set size to {},{}", getType(), w, h);
            // Clear all pages
            width = w;
            height = h;
            reset();
            Page page = new Page(-1, this);
            resources.put(-1, page);
            return this;
        }

        @Override
        public PageResource getResource(int handler) {
            Page page = resources.get(handler);
            if (page == null) {
                throw new ExecutionException(String.format("%s resource #%s not found", getType(), handler));
            }
            return page;
        }

        public void close(Page page) {
            resources.remove(page.getHandler());
        }

        @Inject
        public void init(EventBus eventBus) {
            eventBus.register(this);
        }

        @Subscribe
        public void onVmTest(VmTestEvent e) {
            try {
                log.debug("Dump screen to file");
                for (Map.Entry<Integer, Page> entry : resources.entrySet()) {
                    String fn = "page-" + (entry.getKey() == -1 ? "screen" : entry.getKey()) + ".png";
                    ImageIO.write(entry.getValue().image, "png", new File(fn));
                }
            } catch (Exception ex) {
                ex.printStackTrace();
            }
        }

        @Subscribe
        public void onReset(ResetEvent e) {
            log.debug("Reset {} resources", getType());
            reset();
        }
    }

    static class ImageMgr implements ImageManager {

        @Override
        public PageResource load(String file, int index) {
            return null;
        }

        @Override
        public ImageResource getResource(int handler) {
            return null;
        }

        public void close(Image image) {

        }
    }

    private static class Draw implements Drawable {
        protected final BufferedImage image;
        protected final Graphics2D g;

        private Draw(BufferedImage image) {
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


    static class Image extends Draw implements ImageResource {
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
            manager.close(this);
        }
    }

    static class Page extends Draw implements PageResource {
        private final int handler;
        private final PageMgr manager;
        private final Point pen = new Point();

        Page(int handler, PageMgr manager, BufferedImage image) {
            super(image);
            this.handler = handler;
            this.manager = manager;
        }

        private Page(int handler, PageMgr manager) {
            this(handler, manager, new BufferedImage(manager.getWidth(), manager.getHeight(), BufferedImage.TYPE_INT_RGB));
        }

        @Override
        public PageResource draw(Drawable o, int dx, int dy, int w, int h, int x, int y, int mode) {
            // TODO Ignore mode
            g.drawImage(o.unwrap(Draw.class).image, dx, dy, dx + w, dy + h, w, y, x + w, y + h, null);
            return this;
        }

        @Override
        public PageResource display() {
            manager.getScreen().draw(this);
            return this;
        }

        @Override
        public PageResource draw(PageResource resource) {
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
        public int getHandler() {
            return handler;
        }

        @Override
        public PageMgr getManager() {
            return manager;
        }

        @Override
        public void close() {
            manager.close(this);
        }
    }
}
