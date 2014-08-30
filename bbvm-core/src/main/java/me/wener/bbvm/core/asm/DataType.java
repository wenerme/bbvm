package me.wener.bbvm.core.asm;

import me.wener.bbvm.core.IsValue;

/**
 * 数据类型
 * dword	| 0x0
 * word   | 0x1
 * byte   | 0x2
 * float  | 0x3
 * int    | 0x4
 */
public enum DataType implements IsValue<Integer>
{
    T_DWORD (0x0),
    T_WORD  (0x1),
    T_BYTE  (0x2),
    T_FLOAT (0x3),
    T_INT   (0x4);
    private final int value;

    DataType(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }
}
