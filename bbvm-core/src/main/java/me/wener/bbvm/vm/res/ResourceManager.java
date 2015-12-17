package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/13
 */
public interface ResourceManager<M extends ResourceManager, R extends Resource> {
    R getResource(int handler);

    M reset();

    default R create() {
        throw new UnsupportedOperationException();
    }

    default R create(int handler) {
        throw new UnsupportedOperationException();
    }

    String getType();
}
