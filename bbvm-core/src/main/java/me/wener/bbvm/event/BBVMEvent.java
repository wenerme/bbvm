package me.wener.bbvm.event;

import java.util.EventObject;
import me.wener.bbvm.api.BBVm;

public class BBVMEvent extends EventObject
{
    public BBVMEvent(BBVm source)
    {
        super(source);
    }

    public BBVMEvent()
    {
        super("");
        setSource(null);
    }

    @Override
    public BBVm getSource()
    {
        return (BBVm) super.getSource();
    }

    public BBVMEvent setSource(BBVm source)
    {
        this.source = source;
        return this;
    }

}
