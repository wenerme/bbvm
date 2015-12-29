package me.wener.bbvm.dev.swing;

import com.google.common.collect.Maps;
import com.google.common.eventbus.EventBus;
import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.dev.PageManager;
import me.wener.bbvm.dev.PageResource;
import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.exception.ResourceMissingException;
import me.wener.bbvm.vm.event.ResetEvent;
import me.wener.bbvm.vm.event.VmTestEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.imageio.ImageIO;
import javax.inject.Inject;
import javax.inject.Singleton;
import java.io.File;
import java.util.Map;

import static com.google.common.base.Preconditions.checkState;

/**
 * @author wener
 * @since 15/12/26
 */
@Singleton
class SwingPageManager implements PageManager {
    private final static Logger log = LoggerFactory.getLogger(PageManager.class);
    private final Map<Integer, SwingPage> resources = Maps.newConcurrentMap();
    private int handler = 0;
    // TODO Need to reuse the page handler ?
//        private NavigableSet<Integer> handlers;
    private int width = 240, height = 320;

    public SwingPageManager() {
        // Create a default screen
        SwingPage page = new SwingPage(-1, this);
        resources.put(-1, page);
    }

    @Override
    public PageManager reset() {
        log.info("Reset {} resource", getType());
        handler = 0;
        resources.forEach((k, v) -> v.close());
        checkState(resources.size() == 0, "%s resources should be cleared", getType());
        SwingPage page = new SwingPage(-1, this);
        resources.put(-1, page);
        return this;
    }

    @Override
    public PageResource create() {
        SwingPage page = new SwingPage(handler++, this);
        resources.put(page.getHandler(), page);
        return page;
    }

    @Override
    public SwingPage getScreen() {
        SwingPage screen = resources.get(-1);
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
        if (w < 0 || h < 0 || w > 640 || h > 480) {
            throw new ExecutionException(String.format("Bad page size %s,%s", w, h));
        }
        setSize0(w, h);
        return this;
    }

    private boolean setSize0(int w, int h) {
        if (width == w && height == h) {
            log.debug("{} manager size already in {},{}", getType(), w, h);
            return false;
        } else {
            log.info("{} set size to {},{}", getType(), w, h);
            // Clear all pages
            width = w;
            height = h;
            reset();
            return true;
        }
    }

    @Override
    public PageResource getResource(int handler) {
        SwingPage page = resources.get(handler);
        if (page == null) {
//                log.warn("{} #{} not found", getType(), handler);
            throw new ResourceMissingException(getType(), handler);
//                return getScreen();
        }
        return page;
    }

    public void close(SwingPage page) {
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
            for (Map.Entry<Integer, SwingPage> entry : resources.entrySet()) {
                String fn = "page-" + (entry.getKey() == -1 ? "screen" : entry.getKey()) + ".png";
                ImageIO.write(entry.getValue().image, "png", new File(fn));
            }
        } catch (Exception ex) {
            ex.printStackTrace();
        }
    }

    @Subscribe
    public void onReset(ResetEvent e) {
        if (!setSize0(240, 320)) {
            reset();
        }
    }
}
