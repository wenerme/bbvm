package me.wener.bbvm.dev;

import com.google.common.collect.Maps;
import com.google.common.eventbus.EventBus;
import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.exception.ResourceMissingException;
import me.wener.bbvm.vm.event.ResetEvent;
import me.wener.bbvm.vm.event.VmTestEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import javax.inject.Singleton;
import java.util.Map;

/**
 * Standard string resource implementation
 *
 * @author wener
 * @since 15/12/13
 */
@Singleton
class DefaultStringManager implements StringManager {
    private final static Logger log = LoggerFactory.getLogger(DefaultStringManager.class);
    int handler = -1;
    // TODO Reuse handler ?
    Map<Integer, StringResource> resources = Maps.newHashMap();

    @Override
    public StringResource getResource(int handler) {
        StringResource resource = resources.get(handler);
        if (resource == null) {
            throw new ResourceMissingException(getType(), handler);
        }
        return resource;
    }

    @Override
    public DefaultStringManager reset() {
        handler = -1;
        resources.clear();
        return this;
    }

    void close(StringResource resource) {
        resources.remove(resource.getHandler());
    }

    @Override
    public StringResource create() {
        log.debug("Create string resource {}", handler);
        StringResource resource = new StringResourceImpl(this, handler--);
        resources.put(resource.getHandler(), resource);
        return resource;
    }

    @Inject
    void init(EventBus eventBus) {
        eventBus.register(this);
    }

    @Subscribe
    public void onReset(ResetEvent resetEvent) {
        log.debug("Reset all string resources");
        reset();
    }

    @Subscribe
    public void onVmTest(VmTestEvent e) {
        log.info("On vm test handler={}, resources={}", handler, resources.size());
        resources.forEach((k, v) -> log.info("#{} -> {}", k, v.getValue()));
    }

    /**
     * @author wener
     * @since 15/12/26
     */
    static class StringResourceImpl implements StringResource {
        final private DefaultStringManager manager;
        final private int handler;
        private String value;

        public StringResourceImpl(DefaultStringManager manager, int handler) {
            this.manager = manager;
            this.handler = handler;
        }

        @Override
        public int getHandler() {
            return handler;
        }

        @Override
        public StringManager getManager() {
            return manager;
        }

        @Override
        public void close() {
            manager.close(this);
        }

        public String getValue() {
            return value;
        }

        public StringResourceImpl setValue(String v) {
            value = v;
            return this;
        }

    }
}
