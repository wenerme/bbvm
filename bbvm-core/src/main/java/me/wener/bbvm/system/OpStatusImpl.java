package me.wener.bbvm.system;

import lombok.Data;
import lombok.experimental.Accessors;

@Data
@Accessors(chain = true, fluent = true)
public class OpStatusImpl implements OpStatus
{
    private final Operand a = new Operand();
    private final Operand b = new Operand();
    private DataType dataType;
    private CompareType compareType;
    private CalculateType calculateType;
    private Opcode opcode;

    public static String toString(Operand operand)
    {
        switch (operand.addressingMode())
        {
            case REGISTER:
                return operand.indirect().toString();
            case REGISTER_DEFERRED:
                return "[ " + operand.indirect() + " ]";
            case IMMEDIATE:
                return operand.value().toString();
            case DIRECT:
                return "[ " + operand.value() + " ]";
        }
        throw new UnsupportedOperationException();
    }

    @Override
    public String toAssembly()
    {
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
                sb.append(' ').append(a);
                break;
            case IN:
            case OUT:
                // 标准的两个操作数
                sb.append(' ').append(a)
                  .append(", ").append(b);
                break;

            case JPC:
                sb.append(' ').append(compareType)
                  .append(' ').append(a);
                break;
            case CMP:
                sb.append(' ').append(compareType)
                  .append(' ').append(a)
                  .append(", ").append(b)
                ;
                break;
            case LD:
                sb.append(' ').append(dataType)
                  .append(' ').append(a)
                  .append(", ").append(b)
                ;
                break;
            case CAL:
                sb.append(' ').append(dataType)
                  .append(' ').append(calculateType)
                  .append(' ').append(a)
                  .append(", ").append(b)
                ;
                break;
        }
        return sb.toString();
    }

    @Override
    public byte[] toBinary()
    {
        return new byte[0];
    }
}
