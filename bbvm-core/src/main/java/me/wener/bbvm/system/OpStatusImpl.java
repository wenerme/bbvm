package me.wener.bbvm.system;

import lombok.Data;
import lombok.experimental.Accessors;
import me.wener.bbvm.system.api.CalculateType;
import me.wener.bbvm.system.api.CompareType;
import me.wener.bbvm.system.api.DataType;
import me.wener.bbvm.system.api.OpStatus;
import me.wener.bbvm.system.api.Opcode;

@Data
@Accessors(chain = true, fluent = true)
public class OpStatusImpl implements OpStatus
{
    protected final OperandImpl a = new OperandImpl();
    protected final OperandImpl b = new OperandImpl();
    protected DataType dataType;
    protected CompareType compareType;
    protected CalculateType calculateType;
    protected Opcode opcode;

    public static String toString(OperandImpl operand)
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

    @Override
    public byte[] toBinary()
    {
        byte[] bytes = new byte[opcode.length()];
        switch (opcode)
        {
            case NOP:
            case RET:
            case EXIT:
                bytes[0] = opcode.get().byteValue();
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
}
