package me.wener.bbvm.util.val.impl;

import me.wener.bbvm.util.val.ValueHolder;

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

    public int asInt()
    {
        return 0;
    }

    @Override
    public T get() {
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
