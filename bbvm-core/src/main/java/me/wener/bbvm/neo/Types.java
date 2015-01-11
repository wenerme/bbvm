package me.wener.bbvm.neo;

import com.google.common.collect.Maps;
import java.util.Map;
import me.wener.bbvm.neo.define.CalculateType;
import me.wener.bbvm.neo.define.CompareType;
import me.wener.bbvm.neo.define.DataType;
import me.wener.bbvm.neo.define.Flags;
import me.wener.bbvm.neo.define.InstructionType;
import me.wener.bbvm.neo.define.RegisterType;
import me.wener.bbvm.neo.inst.*;
import me.wener.bbvm.utils.val.Values;

public class Types
{
    public static final Map<Class, Integer> CLASS_TO_TYPE = Maps.newHashMap();

    static
    {
        CLASS_TO_TYPE.put(NOP.class, Flags.NOP);
        CLASS_TO_TYPE.put(LD.class, Flags.LD);
        CLASS_TO_TYPE.put(PUSH.class, Flags.PUSH);
        CLASS_TO_TYPE.put(POP.class, Flags.POP);
        CLASS_TO_TYPE.put(IN.class, Flags.IN);
        CLASS_TO_TYPE.put(OUT.class, Flags.OUT);
        CLASS_TO_TYPE.put(JMP.class, Flags.JMP);
        CLASS_TO_TYPE.put(JPC.class, Flags.JPC);
        CLASS_TO_TYPE.put(CALL.class, Flags.CALL);
        CLASS_TO_TYPE.put(RET.class, Flags.RET);
        CLASS_TO_TYPE.put(CMP.class, Flags.CMP);
        CLASS_TO_TYPE.put(CAL.class, Flags.CAL);
        CLASS_TO_TYPE.put(EXIT.class, Flags.EXIT);
    }

    public static int instructionTypeOfClass(Class<? extends Inst> inst)
    {
        return CLASS_TO_TYPE.get(inst);
    }

    public static int instructionTypeOfObject(Inst inst)
    {
        return CLASS_TO_TYPE.get(inst.getClass());
    }

    public static InstructionType instructionType(int v)
    {
        return Values.fromValue(InstructionType.class, v);
    }

    public static InstructionType instructionType(Inst inst)
    {
        return instructionType(instructionTypeOfObject(inst));
    }

    public static DataType dataType(int v)
    {
        return Values.fromValue(DataType.class, v);
    }

    public static CompareType compareType(int v)
    {
        return Values.fromValue(CompareType.class, v);
    }

    public static CalculateType calculateType(int v)
    {
        return Values.fromValue(CalculateType.class, v);
    }

    public static RegisterType registerType(int v)
    {
        return Values.fromValue(RegisterType.class, v);
    }

    public static boolean hasDataType(int inst)
    {
        switch (inst)
        {
            case Flags.JPC:
            case Flags.POP:
            case Flags.PUSH:
            case Flags.RET:
            case Flags.CALL:
            case Flags.IN:
            case Flags.OUT:
            case Flags.JMP:
            case Flags.CMP:
            case Flags.CAL:
            case Flags.NOP:
            case Flags.EXIT:
                return false;
            case Flags.LD:
                return true;
            default:
                throw new AssertionError();
        }
    }

}
