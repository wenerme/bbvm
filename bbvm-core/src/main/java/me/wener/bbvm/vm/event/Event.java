package me.wener.bbvm.vm.event;

import me.wener.bbvm.vm.VM;

import java.util.EventObject;

/**
 * @author wener
 * @since 15/12/17
 */
public class Event extends EventObject {
    public Event(VM source) {
        super(source);
    }

    @Override
    public VM getSource() {
        return (VM) super.getSource();
    }
}
