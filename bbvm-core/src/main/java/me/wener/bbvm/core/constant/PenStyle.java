package me.wener.bbvm.core.constant;

import me.wener.bbvm.core.IsValue;

public enum PenStyle implements IsValue<Integer>
{
    PEN_SOLID(0), PEN_DASH(1);
    private final int value;

    PenStyle(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }
}
