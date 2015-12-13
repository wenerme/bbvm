package me.wener.bbvm.exception;

import com.google.common.base.MoreObjects;

/**
 * @author wener
 * @since 15/12/13
 */
public class ResourceMissingException extends ExecutionException {
    private String type;
    private int handler;

    public ResourceMissingException(String message, String type, int handler) {
        super(message);
        this.type = type;
        this.handler = handler;
    }

    public ResourceMissingException(String type, int handler) {
        this.type = type;
        this.handler = handler;
    }

    @Override
    public String toString() {
        return MoreObjects.toStringHelper(this)
                .add("type", type)
                .add("handler", handler)
                .add("message", getMessage())
                .toString();
    }
}
