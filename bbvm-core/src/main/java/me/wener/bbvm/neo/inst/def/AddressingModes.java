package me.wener.bbvm.neo.inst.def;

public interface AddressingModes
{
    /*
    表示	| 字节码 | 说明
--------|-----|----
rx		| 0x0 | 寄存器寻址
[rx]	| 0x1 | 寄存器间接寻址
n		| 0x2 | 立即数寻址
[n]	| 0x3 | 直接寻址
     */
    /**
     * 寄存器寻址
     */
    public final static int REGISTER = 0x0;
    /**
     * 寄存器间接寻址
     */
    public final static int REGISTER_DEFERRED = 0x1;
    /**
     * 立即数
     */
    public final static int IMMEDIATE = 0x2;
    /**
     * 直接寻址
     */
    public final static int DIRECT = 0x3;
}
