package me.wener.bbvm.vm.res;

import com.google.common.collect.Maps;
import me.wener.bbvm.exception.ResourceMissingException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

/**
 * Standard string resource implementation
 *
 * @author wener
 * @since 15/12/13
 */
public class StringManager implements ResourceManager<StringManager, StringResource> {
    private final static Logger log = LoggerFactory.getLogger(StringManager.class);
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
    public StringManager reset() {
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
        StringResource resource = new StringResource(this, handler--);
        resources.put(resource.getHandler(), resource);
        return resource;
    }

    @Override
    public String getType() {
        return "String";
    }

    private class Call {

    }
}
