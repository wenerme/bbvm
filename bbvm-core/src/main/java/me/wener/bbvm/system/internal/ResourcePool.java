package me.wener.bbvm.system.internal;

import com.google.common.collect.Maps;
import java.io.IOException;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;
import lombok.extern.slf4j.Slf4j;

/**
 * 资源池,需要维护句柄,主要用于
 */
@Slf4j
public class ResourcePool implements me.wener.bbvm.system.ResourcePool
{
    protected final Map<Integer, me.wener.bbvm.system.Resource> resources = Maps.newConcurrentMap();
    protected AtomicInteger handler = new AtomicInteger();

    @Override
    public Resource request()
    {
        me.wener.bbvm.system.internal.Resource res = new me.wener.bbvm.system.internal.Resource();
        res.handler(next());
        resources.put(res.handler(), res);
        return res;
    }

    protected int next()
    {
        return handler.getAndIncrement();
    }

    @Override
    public Map<Integer, me.wener.bbvm.system.Resource> resources()
    {
        return resources;
    }

    /**
     * 资源被回收
     */
    public void recycle(Resource resource)
    {

    }

    @Override
    public void close() throws IOException
    {
        for (me.wener.bbvm.system.Resource resource : resources.values())
        {
            try
            {
                resource.close();
            } catch (IOException e)
            {
                log.warn("Close resource failed.", e);
            }
        }
    }

    @Override
    public me.wener.bbvm.system.Resource get(int handler)
    {
        return resources.get(handler);
    }
}
