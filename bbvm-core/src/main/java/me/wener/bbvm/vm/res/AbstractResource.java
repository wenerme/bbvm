package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/13
 */
public abstract class AbstractResource implements Resource {
    protected final int handler;

    protected AbstractResource(int handler) {
        this.handler = handler;
    }

    @Override
    public int getHandler() {
        return handler;
    }

    @Override
    public void close() throws Exception {

    }
}
