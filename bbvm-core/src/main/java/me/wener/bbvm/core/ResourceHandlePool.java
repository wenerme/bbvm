package me.wener.bbvm.core;

import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import java.util.Queue;

/**
 * 通用的资源句柄池<br>
 * 该句柄池实现会对释放的句柄重复使用.<br>
 *
 * @param <T> 句柄内部表示的资源
 */
public abstract class ResourceHandlePool<T> implements HandlePool
{
    private final int maxsize;
    /**
     * 资源池
     */
    private final List<T> pool = new ArrayList<>();
    /**
     * 空闲句柄
     */
    private final Queue<Integer> idle = new LinkedList<>();

    protected ResourceHandlePool(int maxsize)
    {
        this.maxsize = maxsize;
    }

    /**
     * 预创建所有资源
     */
    public void prepare()
    {
        if (maxsize < 0)
            throw new IllegalStateException("没有最大句柄数量限制,无法预创建句柄池");
        for (int i = 0; i < maxsize - size(); i++)
        {
            acquire();
        }
    }

    @Override
    public int size()
    {
        return pool.size();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public int maxsize()
    {
        return maxsize;
    }

    /**
     * 根据句柄获取资源
     */
    public T getResource(int handle)
    {
        return pool.get(h2i(handle));
    }

    /**
     * 创建资源
     */
    public abstract T createResource();

    /**
     * 往资源池中添加用户自己创建的资源
     *
     * @return 代表该资源的句柄
     */
    public int addResource(T res)
    {
        Integer index = idle.poll();
        if (index == null)
        {
            index = pool.size();
            pool.add(res);
        } else
            pool.set(index, res);

        return i2h(index);
    }

    /**
     * 在资源池中移除该资源,返回句柄
     */
    public int removeResource(T res)
    {
        int index = pool.indexOf(res);
        if (index >= 0)
        {
            pool.set(index, null);
            idle.add(index);
        }
        return i2h(index);
    }

    /**
     * 释放资源
     */
    public abstract void freeResource(T resource);

    /**
     * {@inheritDoc}
     */
    @Override
    public int acquire()
    {
        return addResource(createResource());
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void release(int handle)
    {
        T res = getResource(handle);
        removeResource(res);
        freeResource(res);
    }

    /**
     * index to handle
     */
    protected int i2h(int i)
    {
        return i;
    }
    /**
     * handle to index
     */
    protected int h2i(int h)
    {
        return h;
    }

}
