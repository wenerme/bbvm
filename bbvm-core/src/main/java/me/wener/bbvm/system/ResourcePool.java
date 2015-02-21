package me.wener.bbvm.system;

import com.google.common.collect.ImmutableMap;
import com.google.common.collect.Maps;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;
import me.wener.bbvm.system.api.Resource;

/**
 * 资源池,需要维护句柄,主要用于
 */
public class ResourcePool
{
    protected final Map<Integer, Resource> resources = Maps.newConcurrentMap();
    protected AtomicInteger handler = new AtomicInteger();

    public Resource request()
    {
        me.wener.bbvm.system.Resource res = new me.wener.bbvm.system.Resource();
        res.handler(next());
        resources.put(res.handler(), res);
        return res;
    }

    protected int next()
    {
        return handler.getAndIncrement();
    }

    public Map<Integer, Resource> resources()
    {
        return ImmutableMap.copyOf(resources);
    }

    /**
     * 资源被回收
     */
    public void recycle(Resource resource)
    {

    }
}
