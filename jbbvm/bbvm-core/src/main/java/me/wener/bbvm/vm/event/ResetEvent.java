package me.wener.bbvm.vm.event;

import me.wener.bbvm.vm.VM;

import java.util.EventObject;

/**
 * @author wener
 * @since 15/12/18
 */
public class ResetEvent extends EventObject {
    public ResetEvent(VM source) {
        super(source);
    }
}
