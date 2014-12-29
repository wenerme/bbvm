package me.wener.bbvm.impl;

public class AdapterResourceHandlePool<T> extends ResourceHandlePool<T>
{
    private final Adapter<T> adapter;
    public AdapterResourceHandlePool(int maxsize, Adapter<T> adapter)
    {
        super(maxsize);
        this.adapter = adapter;
    }

    @Override
    public T createResource()
    {
        return adapter.createResource();
    }

    @Override
    public void freeResource(T resource)
    {
        adapter.freeResource(resource);
    }

    public interface Adapter<T>
    {
        T createResource();
        void freeResource(T resource);
    }
}
