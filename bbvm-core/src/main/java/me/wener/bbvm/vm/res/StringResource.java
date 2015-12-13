package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/13
 */
public class StringResource implements Resource {
    @Override
    public int getHandler() {
        return 0;
    }

    @Override
    public StringManager getManager() {
        return null;
    }

    public String getValue() {
        return null;
    }
}
