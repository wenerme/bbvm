package me.wener.bbvm.neo.define;

public interface CompareTypes
{
    /**
     * 等于
     */
    public static final int Z = 0x1;
    /**
     * Below,小于
     */
    public static final int B = 0x2;
    /**
     * 小于等于
     */
    public static final int BE = 0x3;
    /**
     * Above,大于
     */
    public static final int A = 0x4;
    /**
     * 大于等于
     */
    public static final int AE = 0x5;
    /**
     * 不等于
     */
    public static final int NZ = 0x6;
}
