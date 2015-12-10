package me.wener.bbvm.system;

import me.wener.bbvm.util.val.IsInt;

/**
 * Z    | 0x1 | 等于
 * B    | 0x2 | Below,小于
 * BE   | 0x3 | 小于等于
 * A    | 0x4 | Above,大于
 * AE   | 0x5 | 大于等于
 * NZ   | 0x6 | 不等于
 */
public enum CompareType implements IsInt
{
    Z(0x1),
    B(0x2),
    BE(0x3),
    A(0x4),
    AE(0x5),
    NZ(0x6);
    private final int value;

    CompareType(int value)
    {
        this.value = value;
    }

    public static boolean isMatch(CompareType a, CompareType b)
    {
        boolean valid = false;
        // 判断是否兼容
        switch (a)
        {
            case A:
                if (b == AE || b == A || b == NZ)
                    valid = true;
                break;
            case B:
                if (b == BE || b == B || b == NZ)
                    valid = true;
                break;
            case Z:
                if (b == Z || b == AE || b == BE)
                    valid = true;
                break;
            default:
                if (a == b)
                    valid = true;
        }
        return valid;
    }

    public int asInt()
    {
        return value;
    }
}
