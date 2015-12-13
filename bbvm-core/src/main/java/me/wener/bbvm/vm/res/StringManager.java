package me.wener.bbvm.vm.res;

import com.google.common.collect.Maps;

import java.util.Map;

/**
 * Standard string resource implementation
 *
 * @author wener
 * @since 15/12/13
 */
public class StringManager implements ResourceManager<StringManager, StringResource> {
    int handler = -1;
    // TODO Reuse handler ?
    Map<Integer, StringResource> resources = Maps.newHashMap();

    @Override
    public StringResource getResource(int handler) {
        return resources.get(handler);
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
        return new StringResource(this, handler++);
    }

    @Override
    public String getType() {
        return "String";
    }

    private class Call {

    }
}
