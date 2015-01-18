package me.wener.bbvm.neo.inst.def;

import com.google.common.collect.Maps;
import java.util.EnumMap;
import me.wener.bbvm.utils.val.IsInteger;

/**
 * 指令集
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
public enum InstructionType implements IsInteger
{
    NOP(InstructionTypes.NOP),
    LD(InstructionTypes.LD),
    PUSH(InstructionTypes.PUSH),
    POP(InstructionTypes.POP),
    IN(InstructionTypes.IN),
    OUT(InstructionTypes.OUT),
    JMP(InstructionTypes.JMP),
    JPC(InstructionTypes.JPC),
    CALL(InstructionTypes.CALL),
    RET(InstructionTypes.RET),
    CMP(InstructionTypes.CMP),
    CAL(InstructionTypes.CAL),
    EXIT(InstructionTypes.EXIT);
    private final static EnumMap<InstructionType, Integer> length;

    static
    {
        length = Maps.newEnumMap(InstructionType.class);
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

    private final int value;


    InstructionType(int value)
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

    public Integer get()
    {
        return value;
    }
}
