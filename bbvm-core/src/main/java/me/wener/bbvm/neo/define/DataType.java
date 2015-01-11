package me.wener.bbvm.neo.define;

import java.util.EnumMap;
import java.util.Map;
import me.wener.bbvm.utils.val.IsInteger;

/**
 * 数据类型<pre>
 * dword    | 0x0
 * word   | 0x1
 * byte   | 0x2
 * float  | 0x3
 * int    | 0x4
 * </pre>
 */
public enum DataType implements IsInteger
{
    T_DWORD(DataTypes.T_DWORD),
    T_WORD(DataTypes.T_WORD),
    T_BYTE(DataTypes.T_BYTE),
    T_FLOAT(DataTypes.T_FLOAT),
    T_INT(DataTypes.T_INT);
    private final static Map<DataType, String> strings = new EnumMap<DataType, String>(DataType.class);

    static
    {
        strings.put(T_DWORD, "DWORD");
        strings.put(T_WORD, "WORD");
        strings.put(T_BYTE, "BYTE");
        strings.put(T_FLOAT, "FLOAT");
        strings.put(T_INT, "INT");
    }

    private final int value;

    DataType(int value)
    {
        this.value = value;
    }

    public Integer get()
    {
        return value;
    }

    @Override
    public String toString()
    {
        return strings.get(this);
    }
}
