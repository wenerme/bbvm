package me.wener.bbvm.vm;

/**
 * 数据类型<pre>
 * DWORD  | 0x0
 * WORD   | 0x1
 * BYTE   | 0x2
 * FLOAT  | 0x3
 * INT    | 0x4
 * </pre>
 */
public enum DataType implements IsMachineCode {
    DWORD(0x0),
    WORD(0x1),
    BYTE(0x2),
    FLOAT(0x3),
    INT(0x4);

    private final int value;

    DataType(int value) {
        this.value = value;
    }

    public int asInt() {
        return value;
    }

}
