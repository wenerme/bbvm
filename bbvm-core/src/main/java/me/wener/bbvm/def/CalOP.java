package me.wener.bbvm.def;

import me.wener.bbvm.utils.val.IsInteger;

/**
 * 算数操作符
 */
public enum CalOP implements IsInteger
{
    ADD(0x0),
    SUB(0x1),
    MUL(0x2),
    DIV(0x3),
    MOD(0x4);

    private final int value;

    CalOP(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }
}
