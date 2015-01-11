package me.wener.bbvm.neo.define;

import me.wener.bbvm.utils.val.IsInteger;

/**
 * 算数操作符
 */
public enum CalculateType implements IsInteger
{
    ADD(0x0),
    SUB(0x1),
    MUL(0x2),
    DIV(0x3),
    MOD(0x4);

    private final int value;

    CalculateType(int value)
    {
        this.value = value;
    }

    public Integer get()
    {
        return value;
    }
}
