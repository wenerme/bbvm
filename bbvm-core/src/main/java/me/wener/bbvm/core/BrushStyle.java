package me.wener.bbvm.core;

public enum BrushStyle implements IsValue<Integer>
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
