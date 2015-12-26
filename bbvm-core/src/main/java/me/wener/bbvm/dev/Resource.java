package me.wener.bbvm.dev;

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
    default void close() throws Exception {
    }
}
