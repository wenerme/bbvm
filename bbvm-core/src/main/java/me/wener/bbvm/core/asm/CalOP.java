package me.wener.bbvm.core.asm;

import me.wener.bbvm.core.IsValue;

/**
 * 算数操作符
 */
public enum CalOP implements IsValue<Integer>
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
