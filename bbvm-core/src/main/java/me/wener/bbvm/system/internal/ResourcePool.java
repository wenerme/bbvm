package me.wener.bbvm.system.internal;

import com.google.common.base.Preconditions;
import com.google.common.collect.Maps;
import java.io.IOException;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;
import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.system.Resource;

/**
 * 资源池,需要维护句柄,主要用于
 */
@Slf4j
public class ResourcePool implements me.wener.bbvm.system.ResourcePool
{
    protected final Map<Integer, Resource> resources = Maps.newConcurrentMap();
    protected AtomicInteger handler = new AtomicInteger();

    @Override
    public Resource request()
    {
        Resource res = createResource();
        setHandler(res);
        Preconditions.checkNotNull(res.handler());
        resources.put(res.handler(), res);
        return res;
    }

    /**
     * @return 创建的资源
     */
    protected Resource createResource()
    {
        return new DefaultResource();
    }

    /**
     * 设置句柄值
     */
    protected void setHandler(Resource resource)
    {
        resource.handler(handler.getAndIncrement());
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

    @Override
    public me.wener.bbvm.system.Resource get(int handler)
    {
        return resources.get(handler);
    }
}
