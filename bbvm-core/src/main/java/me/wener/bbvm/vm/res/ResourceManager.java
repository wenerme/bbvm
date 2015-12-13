package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/13
 */
public interface ResourceManager<T extends Resource> {
    T getResource(int handler);
}
