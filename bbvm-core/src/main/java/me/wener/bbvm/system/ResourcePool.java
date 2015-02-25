package me.wener.bbvm.system;

import com.google.common.collect.ImmutableMap;
import com.google.common.collect.Maps;
import java.io.Closeable;
import java.io.IOException;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;
import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.system.api.Resource;

/**
 * 资源池,需要维护句柄,主要用于
 */
@Slf4j
public class ResourcePool implements Closeable
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

    @Override
    public void close() throws IOException
    {
        for (Resource resource : resources.values())
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

    public Resource get(int handler)
    {
        return resources.get(handler);
    }
}
