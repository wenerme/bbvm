package me.wener.bbvm.vm.res;

import com.google.common.base.MoreObjects;
import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import com.google.common.collect.Sets;
import com.google.common.eventbus.EventBus;
import com.google.common.eventbus.Subscribe;
import com.google.inject.AbstractModule;
import com.google.inject.Module;
import me.wener.bbvm.dev.FontType;
import me.wener.bbvm.dev.Images;
import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.exception.ResourceMissingException;
import me.wener.bbvm.util.IntEnums;
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
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.Map;
import java.util.NavigableSet;

import static com.google.common.base.Preconditions.checkState;

/**
 * @author wener
 * @since 15/12/18
 */
public class Swings {
    static {
        IntEnums.cache(FontType.class);
    }

    public static Module graphModule() {
        return new SwingModule();
    }

    private static <T> T checkMissing(ResourceManager mgr, int handler, T v) {
        if (v == null) {
            throw new ResourceMissingException(String.format("%s #%s not exists", mgr.getType(), handler), handler);
        }
        return v;
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
        // TODO Need to reuse the page handler ?
//        private NavigableSet<Integer> handlers;
        private int width, height;

        @Override
        public PageManager reset() {
            setSize(320, 240);
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
            Page screen = resources.get(-1);
            if (screen == null) {
                throw new ExecutionException("Screen not found, may not initialize correctly.");
            }
            return screen;
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
            resources.forEach((k, v) -> v.close());
            checkState(resources.size() == 0, "%s resources should be cleared", getType());
            Page page = new Page(-1, this);
            resources.put(-1, page);
            return this;
        }

        @Override
        public PageResource getResource(int handler) {
            Page page = resources.get(handler);
            if (page == null) {
//                log.warn("{} #{} not found", getType(), handler);
                throw new ExecutionException(String.format("%s resource #%s not found", getType(), handler));
//                return getScreen();
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
                log.debug("Dump pages to file");
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

    @Singleton
    static class ImageMgr implements ImageManager {
        private final static Logger log = LoggerFactory.getLogger(ImageManager.class);
        private final Map<Integer, Image> resources = Maps.newConcurrentMap();
        private final List<String> directories = Lists.newArrayList(".");
        private final NavigableSet<Integer> handlers = Sets.newTreeSet();
        private int handle = 0;

        @Override
        public ImageResource load(String file, int index) {
            try {
                String fn = null;
                for (String directory : directories) {
                    Path path = Paths.get(directory, file);
                    if (Files.exists(path)) {
                        fn = path.toAbsolutePath().toString();
                    }
                }
                if (fn == null) {
                    throw new ExecutionException(String.format("Load %s resource not found #%s %s in %s", getType(), index, file, directories));
                }
                // Index start from 0
                Image image = new Image(nextHandler(), this, Images.read(fn, index));
                image.name = index + "@" + fn;
                log.debug("Load {} resource #{} {}@{}", getType(), handle, index, image);
                resources.put(image.handler, image);
                return image;
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
        }

        int nextHandler() {
            if (handlers.isEmpty()) {
                return handle++;
            }
            return handlers.pollFirst();
        }

        @Override
        public ImageManager reset() {
            resources.forEach((k, v) -> v.close());
            return this;
        }

        @Override
        public List<String> getResourceDirectory() {
            return directories;
        }

        @Override
        public ImageResource getResource(int handler) {
            return checkMissing(this, handler, resources.get(handler));
        }

        public void close(Image image) {
            int handler = image.getHandler();
            handlers.add(handler);
            resources.remove(handler);
        }

        @Inject
        public void init(EventBus eventBus) {
            eventBus.register(this);
        }

        @Subscribe
        public void onVmTest(VmTestEvent e) {
            log.debug("VmTest {} loaded {}", getType(), resources.size());
            resources.forEach((k, v) -> log.debug("Image #{} -> {}", k, v));
        }

        @Subscribe
        public void onReset(ResetEvent e) {
            log.debug("Reset {} resources", getType());
            reset();
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
        public String name;

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
        public void close() {
            manager.close(this);
        }

        @Override
        public String toString() {
            return MoreObjects.toStringHelper(this)
                    .add("handler", handler)
                    .add("width", getWidth())
                    .add("height", getHeight())
                    .add("name", name)
                    .toString();
        }
    }

    static class Page extends Draw implements PageResource {
        private final static Logger log = LoggerFactory.getLogger(PageResource.class);
        private final int handler;
        private final PageMgr manager;
        private final Point pen = new Point();
        private final StringDrawer stringDrawer;

        Page(int handler, PageMgr manager, BufferedImage image) {
            super(image);
            this.handler = handler;
            this.manager = manager;
            Font font = new Font("楷体", Font.PLAIN, 12);
            g.setFont(font);
            stringDrawer = new StringDrawer(g, getWidth(), getHeight());
        }

        private Page(int handler, PageMgr manager) {
            this(handler, manager, new BufferedImage(manager.getWidth(), manager.getHeight(), BufferedImage.TYPE_INT_RGB));
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
//            image.copyData(resource.unwrap(Draw.class).image.getRaster());
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
            stringDrawer.draw(text);
            return this;
        }

        @Override
        public PageResource font(int font) {
            log.info("Set font to {}", IntEnums.fromInt(FontType.class, font));
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

    static class StringDrawer {
        protected final Graphics2D g;
        private int x, y, w, h;
        private FontMetrics metrics;
        private int fontSize;
//        private Consumer<StringDrawer> nextPage = (d) -> d.g.fillRect(0, 0, w, h);

        StringDrawer(Graphics2D g, int w, int h) {
            this.g = g;
            this.w = w;
            this.h = h;
            getFontInfo();
            locate(1, 0);
        }

        public void setFont(Font f) {
            g.setFont(f);
        }

        private void getFontInfo() {
            metrics = g.getFontMetrics();
            Font font = g.getFont();
            fontSize = font.getSize();
        }

        public void locate(int row, int column) {
            int width = g.getFont().getSize();
            x = column * width / 2;
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
//            nextPage.accept(this);
        }
    }
}
