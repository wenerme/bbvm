package me.wener.bbvm.system;

import me.wener.bbvm.util.val.IsInt;

/**
 * 寄存器类型
 * <pre>
 * RP | 0x0 | 程序计数器
 * RF | 0x1 |
 * RS | 0x2 | 栈顶位置
 * RB | 0x3 | 栈底位置
 * R0 | 0x4 | #0 寄存器
 * R1 | 0x5 | #1 寄存器
 * R2 | 0x6 | #2 寄存器
 * R3 | 0x7 | #3 寄存器
 * </pre>
 */
public enum RegisterType implements IsInt
{
    /**
     * 程序计数器,指令寻址寄存器
     */
    rp(0x0),
    /**
     * 标志寄存器,存储比较操作结果
     */
    rf(0x1),
    /**
     * 栈寄存器	<br>空栈顶地址，指向的是下一个准备要压入数据的位置
     */
    rs(0x2),
    /**
     * 辅助栈寄存器<br>栈开始的地址（文件长度+2）
     */
    rb(0x3),
    /**
     * #0 寄存器
     */
    r0(0x4),
    /**
     * #1 寄存器
     */
    r1(0x5),
    /**
     * #2 寄存器
     */
    r2(0x6),
    /**
     * #3 寄存器
     */
    r3(0x7);

    private final int value;

    RegisterType(int value)
    {
        this.value = value;
    }

    public int asInt()
    {
        return value;
    }
}
