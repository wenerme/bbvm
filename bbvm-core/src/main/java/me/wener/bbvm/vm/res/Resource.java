package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/13
 */
public interface Resource extends AutoCloseable, Wrapper {
    int getHandler();

    ResourceManager getManager();

    /**
     * Destroy this resource
     */
    @Override
    void close() throws Exception;
}
