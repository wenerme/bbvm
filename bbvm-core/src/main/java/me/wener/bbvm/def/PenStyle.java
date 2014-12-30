package me.wener.bbvm.def;

import me.wener.bbvm.utils.val.IsInteger;

public enum PenStyle implements IsInteger
{
    PEN_SOLID(0), PEN_DASH(1);
    private final int value;

    PenStyle(int value)
    {
        this.value = value;
    }

    public Integer get()
    {
        return value;
    }
}
