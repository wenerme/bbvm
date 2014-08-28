package me.wener.bbvm.core;

/**
 * 通用的句柄池
 *
 * @param <T> 句柄内部表示的资源
 */
public abstract class AbstractHandlePool<T> implements HandlePool
{
    private final int maxsize;

    protected AbstractHandlePool(int maxsize)
    {
        this.maxsize = maxsize;
    }

    @Override
    public int size()
    {
        return 0;
    }

    /**
     *
     * @return 如果返回 < 0, 则代表没有最大限制
     */
    @Override
    public int maxsize()
    {
        return maxsize;
    }

    public T getResource(int i)
    {
        return null;
    }

    public abstract T createResource();

    public abstract void freeResource(T resource);

    @Override
    public int acquire()
    {
        return 0;
    }

    @Override
    public void release(int handle)
    {

    }
}
