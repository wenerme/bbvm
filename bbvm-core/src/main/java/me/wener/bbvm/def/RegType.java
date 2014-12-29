package me.wener.bbvm.def;

import me.wener.bbvm.utils.val.IsInteger;

/**
 * 寄存器类型
 */
public enum RegType implements IsInteger
{
    rp(0x0),
    rf(0x1),
    rs(0x2),
    rb(0x3),
    r0(0x4),
    r1(0x5),
    r2(0x6),
    r3(0x7);

    private final int value;

    RegType(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }
}
