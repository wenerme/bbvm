package me.wener.bbvm.system;

import me.wener.bbvm.util.val.IsInt;

/**
 * 算数操作符
 */
public enum CalculateType implements IsInt
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

    public int asInt()
    {
        return value;
    }
}
