package me.wener.bbvm.core;

/**
 * 指令集
 *
 * NOP	| 0x0
 * LD     | 0x1
 * PUSH   | 0x2
 * POP    | 0x3
 * IN     | 0x4
 * OUT    | 0x5
 * JMP    | 0x6
 * JPC    | 0x7
 * CALL   | 0x8
 * RET    | 0x9
 * CMP    | 0xA
 * CAL    | 0xB
 * EXIT   | 0xF
 */
public enum Instruction implements IsValue<Integer>
{
    NOP     (0x0),
    LD      (0x1),
    PUSH    (0x2),
    POP     (0x3),
    IN      (0x4),
    OUT     (0x5),
    JMP     (0x6),
    JPC     (0x7),
    CALL    (0x8),
    RET     (0x9),
    CMP     (0xA),
    CAL     (0xB),
    EXIT    (0xF);
    private final int value;

    Instruction(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }
}
