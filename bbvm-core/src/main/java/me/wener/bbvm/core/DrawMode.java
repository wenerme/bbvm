package me.wener.bbvm.core;

public enum DrawMode implements IsValue<Integer>
{
    KEY_COLOR(1);
    private final int value;

    DrawMode(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }
}
