package me.wener.bbvm.core;

/**
 * 该句柄池实现提供的句柄均是  < 0 的
 *
 * {@inheritDoc}
 */
public abstract class ResourceNegativeHandlePool<T> extends ResourceHandlePool<T>
{
    protected ResourceNegativeHandlePool(int maxsize)
    {
        super(maxsize);
    }
    /**
     * index to handle
     */
    protected int i2h(int i)
    {
        return -(i + 1);
    }
    /**
     * handle to index
     */
    protected int h2i(int h)
    {
        return -h-1;
    }
}
