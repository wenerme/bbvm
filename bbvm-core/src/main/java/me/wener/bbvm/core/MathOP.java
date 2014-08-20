package me.wener.bbvm.core;

/**
 * 算数操作符
 */
public enum MathOP implements IsValue<Integer>
{
    ADD(0x0),
    SUB(0x1),
    MUL(0x2),
    DIV(0x3),
    MOD(0x4);

    private final int value;

    MathOP(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }
}
