package me.wener.bbvm.system.internal;

import me.wener.bbvm.system.Resource;

public class NegativeHandlerResourcePool extends ResourcePool
{
    @Override
    protected void setHandler(Resource resource)
    {
        resource.handler(handler.decrementAndGet());
    }
}
