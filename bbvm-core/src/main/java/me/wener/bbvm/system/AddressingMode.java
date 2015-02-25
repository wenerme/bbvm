package me.wener.bbvm.system;

import me.wener.bbvm.utils.val.IsInteger;
import me.wener.bbvm.utils.val.Values;

public enum AddressingMode implements IsInteger
{
    /*
    表示	| 字节码 | 说明
--------|-----|----
rx		| 0x0 | 寄存器寻址
[rx]	| 0x1 | 寄存器间接寻址
n		| 0x2 | 立即数寻址
[n]	| 0x3 | 直接寻址

op1/op2	|rx		| [rx]	| n	    | [n]
--------|-------|-------|-------|----
rx		| 0x0	| 0x1	| 0x2	| 0x3
[rx]	| 0x4	| 0x5	| 0x6	| 0x7
n		| 0x8	| 0x9	| 0xa	| 0xb
[n]	    | 0xc	| 0xd	| 0xe	| 0xf
     */
    /**
     * 寄存器寻址
     */
    REGISTER(0x0),
    /**
     * 寄存器间接寻址
     */
    REGISTER_DEFERRED(0x1),
    /**
     * 立即数
     */
    IMMEDIATE(0x2),
    /**
     * 直接寻址
     */
    DIRECT(0x3);
    private final int val;

    AddressingMode(int val)
    {
        this.val = val;
    }

    static
    {
        Values.cache(AddressingMode.class);
    }

    public static AddressingMode valueOf(int mode)
    {
        return Values.fromValue(AddressingMode.class, mode);
    }

    @Override
    public Integer get()
    {
        return val;
    }

}
