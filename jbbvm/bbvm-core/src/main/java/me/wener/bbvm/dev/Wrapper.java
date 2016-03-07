package me.wener.bbvm.dev;

/**
 * @author wener
 * @since 15/12/18
 */
public interface Wrapper {
    @SuppressWarnings("unchecked")
    default <T> T unwrap(Class<T> iface) {
        return (T) this;
    }

    default boolean isWrapperFor(java.lang.Class<?> iface) {
        return iface.isAssignableFrom(this.getClass());
    }
}
