package me.wener.bbvm.vm.event;

import me.wener.bbvm.vm.VM;

/**
 * @author wener
 * @since 15/12/18
 */
public class ResetEvent extends Event {
    public ResetEvent(VM source) {
        super(source);
    }
}
