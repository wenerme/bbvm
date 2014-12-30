package me.wener.bbvm.def;

import me.wener.bbvm.utils.val.IsInteger;

public enum DrawMode implements IsInteger
{
    KEY_COLOR(1);
    private final int value;

    DrawMode(int value)
    {
        this.value = value;
    }

    public Integer get()
    {
        return value;
    }
}
