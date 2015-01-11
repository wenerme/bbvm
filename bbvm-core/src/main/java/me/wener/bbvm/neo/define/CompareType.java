package me.wener.bbvm.neo.define;

import me.wener.bbvm.utils.val.IsInteger;

/**
 * Z    | 0x1 | 等于
 * B    | 0x2 | Below,小于
 * BE   | 0x3 | 小于等于
 * A    | 0x4 | Above,大于
 * AE   | 0x5 | 大于等于
 * NZ   | 0x6 | 不等于
 */
public enum CompareType implements IsInteger
{
    Z(0x1),
    B(0x2),
    BE(0x3),
    A(0x4),
    AE(0x5),
    NZ(0x6);
    private final int value;

    CompareType(int value)
    {
        this.value = value;
    }

    public Integer get()
    {
        return value;
    }
}