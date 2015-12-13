package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/13
 */
public class StringResource implements Resource {
    final private StringManager manager;
    final private int handler;
    private String value;

    public StringResource(StringManager manager, int handler) {
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

    public StringResource setValue(String v) {
        value = v;
        return this;
    }

}
