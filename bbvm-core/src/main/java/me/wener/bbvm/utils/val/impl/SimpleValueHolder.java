package me.wener.bbvm.utils.val.impl;

import me.wener.bbvm.utils.val.ValueHolder;

public class SimpleValueHolder<T> implements ValueHolder<T>
{
    private T value;

    public SimpleValueHolder(T value)
    {
        this.value = value;
    }

    public SimpleValueHolder()
    {
    }

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

    @Override
    public String toString()
    {
        return "SimpleValueHolder{" +
                "value=" + value +
                '}';
    }
}
