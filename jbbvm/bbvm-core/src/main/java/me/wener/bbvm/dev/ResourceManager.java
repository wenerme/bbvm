package me.wener.bbvm.dev;

import com.google.common.base.Throwables;
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
    default R close(int handler) {
        try {
            R resource = getResource(handler);
            resource.close();
        } catch (ResourceMissingException e) {
            // ignored
        } catch (Exception e) {
            Throwables.propagate(e);
        }

        return null;
    }

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
