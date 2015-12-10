package me.wener.bbvm.util.val.impl;

import me.wener.bbvm.util.val.IsValue;

public class SimpleValue<T> implements IsValue<T>
{
    private T value;

    public SimpleValue(T value)
    {
        this.value = value;
    }

    @Override
    public T get()
    {
        return value;
    }
}
