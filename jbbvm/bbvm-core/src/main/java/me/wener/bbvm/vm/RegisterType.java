package me.wener.bbvm.vm;

import me.wener.bbvm.util.IsInt;

/**
 * 寄存器类型
 * <pre>
 * RP | 0x0 | 程序计数器
 * RF | 0x1 | 比较标示符
 * RS | 0x2 | 栈顶位置
 * RB | 0x3 | 栈底位置
 * R0 | 0x4 | #0 寄存器
 * R1 | 0x5 | #1 寄存器
 * R2 | 0x6 | #2 寄存器
 * R3 | 0x7 | #3 寄存器
 * </pre>
 */
public enum RegisterType implements IsInt {
    RP(0x0),
    RF(0x1),
    RS(0x2),
    RB(0x3),
    R0(0x4),
    R1(0x5),
    R2(0x6),
    R3(0x7);

    private final int value;

    RegisterType(int value) {
        this.value = value;
    }

    public int asInt() {
        return value;
    }
}
