package me.wener.bbvm.vm.res.swing;

import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import com.google.common.collect.Sets;
import com.google.common.eventbus.EventBus;
import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.dev.Images;
import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.vm.event.ResetEvent;
import me.wener.bbvm.vm.event.VmTestEvent;
import me.wener.bbvm.vm.res.ImageManager;
import me.wener.bbvm.vm.res.ImageResource;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import javax.inject.Singleton;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.Map;
import java.util.NavigableSet;

/**
 * @author wener
 * @since 15/12/26
 */
@Singleton
class SwingImageManager implements ImageManager {
    private final static Logger log = LoggerFactory.getLogger(ImageManager.class);
    private final Map<Integer, SwingImage> resources = Maps.newConcurrentMap();
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
            SwingImage image = new SwingImage(nextHandler(), this, Images.read(fn, index));
            image.name = index + "@" + fn;
            log.debug("Load {} resource #{} {}@{}", getType(), handle, index, image);
            resources.put(image.getHandler(), image);
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
        return Swings.checkMissing(this, handler, resources.get(handler));
    }

    public void close(SwingImage image) {
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
