package me.wener.bbvm.core.asm;

import com.google.common.collect.Maps;
import java.util.EnumMap;
import me.wener.bbvm.core.IsValue;

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
    private final static EnumMap<Instruction, Integer> length;

    static
    {
        length = Maps.newEnumMap(Instruction.class);
        length.put(NOP, 1);
        length.put(LD, 10);
        length.put(PUSH, 5);
        length.put(POP, 5);
        length.put(IN, 10);
        length.put(OUT, 10);
        length.put(JMP, 6);
        length.put(JPC, 6);
        length.put(CALL, 5);
        length.put(RET, 1);
        length.put(CMP, 10);
        length.put(CAL, 10);
        length.put(EXIT, 1);
    }


    Instruction(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }

    /**
     * 获取对应指令的长度
     */
    public static Integer length(Instruction instruction)
    {
        return length.get(instruction);
    }
}
