package me.wener.bbvm.neo.define;


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
public interface RegisterTypes
{
    /**
     * 程序计数器,指令寻址寄存器
     */
    public static final int rp = 0x0;
    /**
     * 标志寄存器,存储比较操作结果
     */
    public static final int rf = 0x1;
    /**
     * 栈寄存器	<br>空栈顶地址，指向的是下一个准备要压入数据的位置
     */
    public static final int rs = 0x2;
    /**
     * 辅助栈寄存器<br>栈开始的地址（文件长度+2）
     */
    public static final int rb = 0x3;
    /**
     * #0 寄存器
     */
    public static final int r0 = 0x4;
    /**
     * #1 寄存器
     */
    public static final int r1 = 0x5;
    /**
     * #2 寄存器
     */
    public static final int r2 = 0x6;
    /**
     * #3 寄存器
     */
    public static final int r3 = 0x7;
}
