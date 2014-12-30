package me.wener.bbvm.impl.subscriber;

import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.api.BBVm;
import me.wener.bbvm.event.VMStateEvent;
import me.wener.bbvm.impl.Reg;

public class AbstractVMSubscriber
{
    protected BBVm vm;
    protected Reg rp;
    protected Reg rb;
    protected Reg rs;
    protected Reg rf;
    protected Reg r0;
    protected Reg r1;
    protected Reg r2;
    protected Reg r3;

    @Subscribe
    protected void initState(VMStateEvent e)
    {
        BBVm vm = e.getSource();
    }


    protected UnsupportedOperationException unsupport(String format, Object... args)
    {
        return unsupport(String.format(format, args));
    }

    protected UnsupportedOperationException unsupport(String str)
    {
        return new UnsupportedOperationException(str);
    }

}
