package me.wener.bbvm.neo.inst.def;

public interface InstructionTypes
{
    /**
     * 无操作
     */
    public static final int NOP = 0x0;
    /**
     * 数据读写
     */
    public static final int LD = 0x1;
    /**
     * 入栈
     */
    public static final int PUSH = 0x2;
    /**
     * 出栈
     */
    public static final int POP = 0x3;
    /**
     * IN端口调用
     */
    public static final int IN = 0x4;
    /**
     * OUT端口调用
     */
    public static final int OUT = 0x5;
    /**
     * 跳转
     */
    public static final int JMP = 0x6;
    /**
     * 条件跳转
     */
    public static final int JPC = 0x7;
    /**
     * 调用跳转
     */
    public static final int CALL = 0x8;
    /**
     * 返回
     */
    public static final int RET = 0x9;
    /**
     * 比较
     */
    public static final int CMP = 0xA;
    /**
     * 算术运算
     */
    public static final int CAL = 0xB;
    /**
     * 退出
     */
    public static final int EXIT = 0xF;
}
