package me.wener.bbvm.impl.subscriber;

import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.api.DeviceFunction;
import me.wener.bbvm.event.InstEvent;
import me.wener.bbvm.impl.InstructionContext;
import me.wener.bbvm.impl.VMContext;

public class DeviceFunctionSubscriber extends VMContext
{
    private final DeviceFunction deviceFunction;

    public DeviceFunctionSubscriber(DeviceFunction deviceFunction) {this.deviceFunction = deviceFunction;}

    @Subscribe
    protected void subscribeInst(InstEvent e)
    {
        doInstruction(e.getContext());
    }

    protected void doInstruction(InstructionContext ctx)
    {

    }


    protected boolean out(InstructionContext ctx)
    {
        Integer input = ctx.getOp2().get();

        return true;
    }
}
