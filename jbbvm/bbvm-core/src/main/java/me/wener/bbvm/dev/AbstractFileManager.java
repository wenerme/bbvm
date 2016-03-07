package me.wener.bbvm.dev;

import com.google.common.eventbus.EventBus;
import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.vm.event.ResetEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import java.nio.charset.Charset;

/**
 * @author wener
 * @since 15/12/26
 */
public abstract class AbstractFileManager implements FileManager {
    protected final static Logger log = LoggerFactory.getLogger(FileManager.class);
    static final int HANDLER_NUMBER = 10;
    protected final FileResource[] resources = new FileResource[HANDLER_NUMBER];
    @Inject
    protected Charset charset;

    protected abstract FileResource createResource(int i);

    @Override
    public FileResource getResource(int handler) {
        if (handler < 0 || handler > HANDLER_NUMBER) {
            throw new ExecutionException("No file resource for handler " + handler);
        }
        return resources[handler];
    }

    @Override
    public AbstractFileManager reset() {
        for (FileResource resource : resources) {
            resource.close();
        }
        return this;
    }

    @Override
    public FileResource create() {
        throw new ExecutionException("Resource for file is fixed");
    }

    @Inject
    void init(EventBus eventBus) {
        eventBus.register(this);

        for (int i = 0; i < HANDLER_NUMBER; i++) {
            resources[i] = createResource(i);
        }

    }

    @Subscribe
    public void onReset(ResetEvent resetEvent) {
        log.debug("Reset all file resources");
        reset();
    }
}
