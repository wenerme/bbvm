package me.wener.bbvm.core;

public class SimpleValueHolder<T> implements ValueHolder<T>
{
    private T value;
    @Override
    public T get()
    {
        return value;
    }

    @Override
    public void set(T v)
    {
        value = v;
    }
}
