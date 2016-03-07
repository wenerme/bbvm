package me.wener.bbvm.exception;

import me.wener.bbvm.vm.VM;

/**
 * @author wener
 * @since 15/12/13
 */
public class ExecutionException extends RuntimeException {
    private VM vm;

    public ExecutionException() {
    }

    public ExecutionException(String message) {
        super(message);
    }

    public ExecutionException(String message, Throwable cause) {
        super(message, cause);
    }

    public ExecutionException(Throwable cause) {
        super(cause);
    }

    public VM getVm() {
        return vm;
    }

    public ExecutionException setVm(VM vm) {
        this.vm = vm;
        return this;
    }
}
