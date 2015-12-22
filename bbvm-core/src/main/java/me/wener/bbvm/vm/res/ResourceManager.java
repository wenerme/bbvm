package me.wener.bbvm.vm.res;

import me.wener.bbvm.exception.ResourceMissingException;

/**
 * @author wener
 * @since 15/12/13
 */
public interface ResourceManager<M extends ResourceManager, R extends Resource> {
    R getResource(int handler) throws ResourceMissingException;

    /**
     * Close this resource, if the resource is not exists may not throw an exception
     */
//    R close(int handler);

    @SuppressWarnings("unchecked")
    default M reset() {
        // If this resources do not need to reset
        return (M) this;
    }

    default R create() {
        throw new UnsupportedOperationException();
    }

    String getType();

}
