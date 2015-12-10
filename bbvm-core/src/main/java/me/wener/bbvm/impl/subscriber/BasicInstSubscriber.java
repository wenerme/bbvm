package me.wener.bbvm.impl.subscriber;

import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.event.InstEvent;
import me.wener.bbvm.impl.InstructionContext;
import me.wener.bbvm.impl.Operand;
import me.wener.bbvm.impl.VMContext;
import me.wener.bbvm.util.Bins;
import me.wener.bbvm.util.val.Values;
import me.wener.bbvm.vm.CalculateType;
import me.wener.bbvm.vm.CompareType;
import me.wener.bbvm.vm.DataType;
import me.wener.bbvm.vm.Opcode;

public class BasicInstSubscriber extends VMContext
{

    @Subscribe
    protected void subscribeBasicInst(InstEvent e)
    {
        doInstruction(e.getContext());
    }

    protected void doInstruction(InstructionContext ctx)
    {
        final Operand op1 = ctx.getOp1();
        final Operand op2 = ctx.getOp2();
        final Integer opv1 = op1.asInt();
        final Integer opv2 = op2.asInt();
        final DataType dataType = ctx.getDataType();
        final Opcode instruction = ctx.getInstruction();

        switch (instruction)
        {
            case NOP:
                break;
            case LD:
                switch (dataType)
                {
                    case DWORD:
                    case FLOAT:
                    case INT:
                        op1.set(op2.asInt());
                        break;
                    case BYTE:
                        op1.set(op1.asInt() & 0xffffff00 | (op2.asInt() & 0xff));
                        break;
                    case WORD:
                        op1.set(op1.asInt() & 0xffff0000 | (op2.asInt() & 0xffff));
                        break;
                    default:
                        throw unsupport("未知的数据类型: %s", dataType);
                }
                break;
            case PUSH:
                vm.push(op1.asInt());
                break;
            case POP:
                op1.set(vm.pop());
                break;
            case JMP:
                rp.set(op1.asInt());
                break;
            case JPC:
            {
                // JPC 的数据类型为比较操作
                CompareType org = Values.fromValue(CompareType.class, (int) Bins.int4(ctx.getFirstByte(), 2));
                CompareType flag = Values.fromValue(CompareType.class, rf.asInt());
                boolean valid = false;

                switch (flag)
                {
                    case A:
                        if (org == CompareType.AE || org == CompareType.A || org == CompareType.NZ)
                            valid = true;
                        break;
                    case B:
                        if (org == CompareType.BE || org == CompareType.B || org == CompareType.NZ)
                            valid = true;
                        break;
                    case Z:
                        if (org == CompareType.Z || org == CompareType.AE || org == CompareType.BE)
                            valid = true;
                        break;
                    default:
                        if (org.equals(flag))
                            valid = true;
                }

                if (valid)
                    rp.set(opv1);
            }
            break;
            case CALL:
                // 设置返回位置为下一句的开始
                vm.push(rp.asInt() + Opcode.length(instruction));
                rp.set(opv1);
                break;
            case RET:
                rp.set(vm.pop());
                break;
            case CMP:
            {
                float a = opv1;
                float b = opv2;
                if (dataType == DataType.FLOAT)
                {
                    a = Bins.float32(opv1);
                    b = Bins.float32(opv2);
                }
                float c = a - b;
                if (c > 0)
                    rf.set(CompareType.A.asInt());
                else if (c < 0)
                    rf.set(CompareType.B.asInt());
                else
                    rf.set(CompareType.Z.asInt());
            }
            break;
            case CAL:
            {

                CalculateType op = Values.fromValue(CalculateType.class, ctx.getSpecialByte());
                // 返回结果为 R0
                double a = opv1;
                double b = opv2;
                if (dataType.equals(DataType.FLOAT))
                {
                    a = Bins.float32(opv1);
                    b = Bins.float32(opv2);
                }
                double c;
                switch (op)
                {
                    case ADD:
                        c = a + b;
                        break;
                    case DIV:
                        c = a / b;
                        break;
                    case MOD:
                        c = a % b;
                        break;
                    case MUL:
                        c = a * b;
                        break;
                    case SUB:
                        c = a - b;
                        break;
                    default:
                        throw unsupport("未知计算操作: %s", op);
                }
                int ret = (int) c;
                // 值返回归约
                switch (dataType)
                {
                    case FLOAT:
                        ret = Bins.int32((float) c);
                        break;
                    case BYTE:
                        ret &= 0xff;
                        break;
                    case WORD:
                        ret &= 0xffff;
                        break;
                }
                op1.set(ret);
            }
            break;
            case EXIT:
                break;
            default:
                throw unsupport("未知指令: %s", instruction);
        }
    }

}
