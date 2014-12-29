package me.wener.bbvm.def;

public enum PenStyle implements IsIntValue
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
