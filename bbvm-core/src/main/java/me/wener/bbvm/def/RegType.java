package me.wener.bbvm.def;

import me.wener.bbvm.utils.val.IsInteger;

/**
 * 寄存器类型
 * <pre>
 * rp | 0x0 | 程序计数器
 * rf | 0x1 |
 * rs | 0x2 | 栈顶位置
 * rb | 0x3 | 栈底位置
 * r0 | 0x4 | #0 寄存器
 * r1 | 0x5 | #1 寄存器
 * r2 | 0x6 | #2 寄存器
 * r3 | 0x7 | #3 寄存器
 * </pre>
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

    public Integer get()
    {
        return value;
    }
}
