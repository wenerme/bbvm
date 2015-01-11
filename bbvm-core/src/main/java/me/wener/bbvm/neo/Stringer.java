package me.wener.bbvm.neo;

import static java.nio.charset.StandardCharsets.UTF_8;

import io.netty.buffer.ByteBuf;
import me.wener.bbvm.neo.define.Flags;
import me.wener.bbvm.neo.define.InstructionType;
import me.wener.bbvm.neo.inst.Inst;
import me.wener.bbvm.neo.inst.OneOperandInst;
import me.wener.bbvm.neo.inst.TowOperandInst;

/**
 * 将相应类型转换为字符串
 */
public class Stringer
{
    public static String string(byte[] bytes)
    {
        return new String(bytes, UTF_8);
    }

    public static String string(ByteBuf buf)
    {
        return buf == null ? null : buf.toString(UTF_8);
    }

    public static String string(Inst inst)
    {
        InstructionType type = Types.instructionType(inst);
        OneOperandInst one = null;
        TowOperandInst tow = null;
        if (inst instanceof OneOperandInst)
        {
            one = (OneOperandInst) inst;
        } else if (inst instanceof TowOperandInst)
        {
            tow = (TowOperandInst) inst;
        }

        switch (type.get())
        {
            case Flags.JPC:
            case Flags.POP:
            case Flags.PUSH:
            case Flags.CALL:
                return String.format("%s %s", type, one.a);
            case Flags.IN:
            case Flags.OUT:
            {
                String format = "%s %s, %s";
                return String.format(format, type, tow.a, tow.b);
            }
            case Flags.JMP:
            case Flags.CMP:
            case Flags.CAL:

            case Flags.RET:
            case Flags.NOP:
            case Flags.EXIT:
                return String.valueOf(type);
            case Flags.LD:
            {
                String format = "%s %s %s, %s";
                return String.format(format, type, Types.dataType(tow.dataType), tow.a, tow.b);
            }
            default:
                throw new AssertionError();
        }
    }
}
