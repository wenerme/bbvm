package me.wener.bbvm.def;

import me.wener.bbvm.utils.val.IsInteger;

public enum BrushStyle implements IsInteger
{
    BRUSH_SOLID(0);
    private final int value;

    BrushStyle(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }
}
