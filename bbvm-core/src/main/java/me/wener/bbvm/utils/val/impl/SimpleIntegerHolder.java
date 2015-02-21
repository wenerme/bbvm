package me.wener.bbvm.utils.val.impl;

import java.util.concurrent.atomic.AtomicInteger;
import me.wener.bbvm.utils.val.IntegerHolder;

/**
 * 通过 {@link AtomicInteger} 实现的 {@link IntegerHolder} 实现
 */
public class SimpleIntegerHolder
        implements IntegerHolder
{
    protected final AtomicInteger value = new AtomicInteger();

    @Override
    public Integer get()
    {
        return value.get();
    }

    @Override
    public void set(Integer v)
    {
        value.set(v);
    }
}
