package me.wener.bbvm.utils.val.impl;

public class ReadonlyHolder<T> extends SimpleValueHolder<T>
{
    public ReadonlyHolder(T value)
    {
        super(value);
    }

    public ReadonlyHolder()
    {
    }

    @Override
    public void set(T v)
    {
        throw new UnsupportedOperationException("readonly");
    }
}
