package me.wener.bbvm.system;

import com.google.common.collect.Maps;
import me.wener.bbvm.util.val.IsInt;

import java.util.EnumMap;

/**
 * 操作码
 * <pre>
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
 * </pre>
 */
public enum Opcode implements IsInt
{
    NOP(0x0),
    LD(0x1),
    PUSH(0x2),
    POP(0x3),
    IN(0x4),
    OUT(0x5),
    JMP(0x6),
    JPC(0x7),
    CALL(0x8),
    RET(0x9),
    CMP(0xA),
    CAL(0xB),
    EXIT(0xF);
    private final static EnumMap<Opcode, Integer> length;

    static
    {
        length = Maps.newEnumMap(Opcode.class);
        length.put(NOP, 1);
        length.put(RET, 1);
        length.put(EXIT, 1);

        length.put(PUSH, 5);
        length.put(POP, 5);
        length.put(JMP, 5);
        length.put(CALL, 5);
        length.put(JPC, 6);
        length.put(LD, 10);
        length.put(IN, 10);
        length.put(OUT, 10);
        length.put(CMP, 10);
        length.put(CAL, 10);
    }

    static {
//        Values.cache(Opcode.class);
    }

    private final int value;

    Opcode(int value)
    {
        this.value = value;
    }

    /**
     * 获取对应指令的长度
     */
    public Integer length()
    {
        return length.get(this);
    }

    public int asInt()
    {
        return value;
    }
}
