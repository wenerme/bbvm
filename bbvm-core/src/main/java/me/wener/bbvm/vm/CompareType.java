package me.wener.bbvm.vm;

import me.wener.bbvm.util.IsInt;

/**
 * Z    | 0x1 | 等于
 * B    | 0x2 | Below,小于
 * BE   | 0x3 | 小于等于
 * A    | 0x4 | Above,大于
 * AE   | 0x5 | 大于等于
 * NZ   | 0x6 | 不等于
 */
public enum CompareType implements IsInt {
    Z(0x1),
    B(0x2),
    BE(0x3),
    A(0x4),
    AE(0x5),
    NZ(0x6);
    private final int value;

    CompareType(int value) {
        this.value = value;
    }

    public int asInt() {
        return value;
    }

    public boolean isMatch(CompareType b) {
        // 比较后的结果只有 大于 等于 小于
        // 需要匹配的包括所有条件
        switch (this) {
            case A:
                if (b == CompareType.AE || b == CompareType.A || b == CompareType.NZ)
                    return true;
                break;
            case B:
                if (b == CompareType.BE || b == CompareType.B || b == CompareType.NZ)
                    return true;
                break;
            case Z:
                if (b == CompareType.Z || b == CompareType.AE || b == CompareType.BE)
                    return true;
                break;
        }
        return this == b;
    }
}
