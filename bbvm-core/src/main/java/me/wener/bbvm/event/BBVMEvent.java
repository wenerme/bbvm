package me.wener.bbvm.event;

import java.util.EventObject;
import me.wener.bbvm.api.BBVm;

public class BBVmEvent extends EventObject
{
    public BBVmEvent(BBVm source)
    {
        super(source);
    }

    public BBVmEvent()
    {
        super("");
        setSource(null);
    }

    @Override
    public BBVm getSource()
    {
        return (BBVm) super.getSource();
    }

    public BBVmEvent setSource(BBVm source)
    {
        this.source = source;
        return this;
    }

}
