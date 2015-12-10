package me.wener.bbvm.system.internal;

import lombok.Data;
import lombok.experimental.Accessors;
import me.wener.bbvm.system.*;
import me.wener.bbvm.util.Bins;
import me.wener.bbvm.util.val.IntEnums;

import java.nio.ByteBuffer;

public class OpStates
{
    public static <T extends WritableOpState> T readBinary(T s, ByteBuffer buf)
    {
        /*
   指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数 + 第二个操作数
0x 0       0         0           0        00000000     00000000

无操作数 1byte
   指令码 + 无用
0x 0       0
一个操作数 5byte
   指令码 + 寻址方式 + 第一个操作数
0x 0       0        00000000
两个操作数 10byte
   指令码 + 数据类型 + 保留字节 + 寻址方式 + 第一个操作数 + 第二个操作数
0x 0       0         0        0        00000000     00000000
JPC指令 6byte
   指令码 + 比较操作 + 保留字节 + 寻址方式 + 第一个操作数
0x 0       0         0        0        00000000
        */

        short first = Bins.unsigned(buf.get());
        Opcode opcode = IntEnums.fromInt(Opcode.class, first >> 4);
        me.wener.bbvm.system.Operand a = s.a();
        me.wener.bbvm.system.Operand b = s.b();
        s.opcode(opcode);
        switch (opcode)
        {
            case RET:
            case NOP:
            case EXIT:
            {
                // 无操作数
            }
            break;
            case POP:
            case PUSH:
            case CALL:
            case JMP:
            {
                a.addressingMode(IntEnums.fromInt(AddressingMode.class, first & 0xf));
                // 一个操作数
                a.value(buf.getInt());
            }
            break;
            case LD:
            case IN:
            case OUT:
            case CAL:
            case CMP:
            {
                // 两个操作数
                s.dataType(IntEnums.fromInt(DataType.class, first & 0xf));
                short second = Bins.unsigned(buf.get());
                int special = second >> 4;
                int addressingMode = second & 0xf;

                a.addressingMode(IntEnums.fromInt(AddressingMode.class, addressingMode / 4));
                b.addressingMode(IntEnums.fromInt(AddressingMode.class, addressingMode % 4));
                a.value(buf.getInt());
                b.value(buf.getInt());

                if (opcode == Opcode.CAL)
                {
                    s.calculateType(IntEnums.fromInt(CalculateType.class, special));
                }
            }
            break;

            case JPC:
            {
                short second = Bins.unsigned(buf.get());
                int addressingMode = second & 0xf;
                // JPC A R1
                // 数据类型为比较操作
                s.compareType(IntEnums.fromInt(CompareType.class, first & 0xf));
                a.addressingMode(IntEnums.fromInt(AddressingMode.class, addressingMode));
                a.value(buf.getInt());
            }
            break;

            default:
                throw new UnsupportedOperationException();
        }
        return s;
    }

    public static String toAssembly(OpState s)
    {
        me.wener.bbvm.system.Operand a = s.a();
        me.wener.bbvm.system.Operand b = s.b();
        CompareType compareType = s.compareType();
        DataType dataType = s.dataType();
        CalculateType calculateType = s.calculateType();
        Opcode opcode = s.opcode();

        StringBuilder sb = new StringBuilder();
        sb.append(opcode);
        switch (opcode)
        {
            // 没有操作数
            case NOP:
            case RET:
            case EXIT:
                break;

            case PUSH:
            case POP:
            case JMP:
            case CALL:
                // 一个操作数
                sb.append(' ').append(a.toAssembly());
                break;
            case IN:
            case OUT:
                // 标准的两个操作数
                sb.append(' ').append(a.toAssembly())
                  .append(", ").append(b.toAssembly());
                break;

            case JPC:
                sb.append(' ').append(compareType)
                  .append(' ').append(a.toAssembly());
                break;
            case CMP:
                sb.append(' ').append(compareType)
                  .append(' ').append(a.toAssembly())
                  .append(", ").append(b.toAssembly())
                ;
                break;
            case LD:
                sb.append(' ').append(dataType)
                  .append(' ').append(a.toAssembly())
                  .append(", ").append(b.toAssembly())
                ;
                break;
            case CAL:
                sb.append(' ').append(dataType)
                  .append(' ').append(calculateType)
                  .append(' ').append(a.toAssembly())
                  .append(", ").append(b.toAssembly())
                ;
                break;
        }
        return sb.toString();
    }

    public static byte[] toBinary(OpState s)
    {
        me.wener.bbvm.system.Operand a = s.a();
        me.wener.bbvm.system.Operand b = s.b();
        CompareType compareType = s.compareType();
        DataType dataType = s.dataType();
        CalculateType calculateType = s.calculateType();
        Opcode opcode = s.opcode();


        byte[] bytes = new byte[opcode.length()];
        switch (opcode)
        {
            case NOP:
            case RET:
            case EXIT:
//                bytes[0] = opcode.asInt().byteValue();
                break;
            case LD:
                break;
            case PUSH:
                break;
            case POP:
                break;
            case IN:
                break;
            case OUT:
                break;
            case JMP:
                break;
            case JPC:
                break;
            case CALL:
                break;
            case CMP:
                break;
            case CAL:
                break;
        }
        return bytes;
    }

    private void readOperand(me.wener.bbvm.system.Operand o, ByteBuffer buffer)
    {
        int v = buffer.getInt();
        o.value(v);
    }

    public interface WritableOpState extends OpState
    {

        WritableOpState dataType(DataType dataType);

        WritableOpState compareType(CompareType compareType);

        WritableOpState calculateType(CalculateType calculateType);

        WritableOpState opcode(Opcode opcode);

        WritableOpState a(me.wener.bbvm.system.Operand a);

        WritableOpState b(me.wener.bbvm.system.Operand b);
    }

    @Data
    @Accessors(chain = true, fluent = true)
    public static class DefaultOpState implements WritableOpState
    {
        protected me.wener.bbvm.system.Operand a;
        protected me.wener.bbvm.system.Operand b;
        protected DataType dataType;
        protected CompareType compareType;
        protected CalculateType calculateType;
        protected Opcode opcode;

        @Override
        public String toAssembly()
        {
            return OpStates.toAssembly(this);
        }

        @Override
        public byte[] toBinary()
        {
            return OpStates.toBinary(this);
        }
    }
}
