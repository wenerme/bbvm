package me.wener.bbvm.impl;

import com.google.common.eventbus.EventBus;
import me.wener.bbvm.api.Device;

@SuppressWarnings("ConstantConditions")
public class BBVmEventBasedImpl extends BBVmImpl
{
    EventBus eventBus;

    public BBVmEventBasedImpl(Device device)
    {
        super(device);
    }

    @Override
    protected void doInstruction(InstructionContext ctx)
    {
        eventBus.post(ctx);
    }
}
