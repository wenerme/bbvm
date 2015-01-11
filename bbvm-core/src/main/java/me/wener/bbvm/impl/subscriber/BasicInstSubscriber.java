package me.wener.bbvm.impl.subscriber;

import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.def.CalOP;
import me.wener.bbvm.def.CmpOP;
import me.wener.bbvm.def.DataType;
import me.wener.bbvm.def.InstructionType;
import me.wener.bbvm.event.InstEvent;
import me.wener.bbvm.impl.InstructionContext;
import me.wener.bbvm.impl.Operand;
import me.wener.bbvm.impl.VMContext;
import me.wener.bbvm.utils.Bins;
import me.wener.bbvm.utils.val.Values;

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
        final Integer opv1 = op1.get();
        final Integer opv2 = op2.get();
        final DataType dataType = ctx.getDataType();
        final InstructionType instruction = ctx.getInstruction();

        switch (instruction)
        {
            case NOP:
                break;
            case LD:
                switch (dataType)
                {
                    case T_DWORD:
                    case T_FLOAT:
                    case T_INT:
                        op1.set(op2.get());
                        break;
                    case T_BYTE:
                        op1.set(op1.get() & 0xffffff00 | (op2.get() & 0xff));
                        break;
                    case T_WORD:
                        op1.set(op1.get() & 0xffff0000 | (op2.get() & 0xffff));
                        break;
                    default:
                        throw unsupport("未知的数据类型: %s", dataType);
                }
                break;
            case PUSH:
                vm.push(op1.get());
                break;
            case POP:
                op1.set(vm.pop());
                break;
            case JMP:
                rp.set(op1.get());
                break;
            case JPC:
            {
                // JPC 的数据类型为比较操作
                CmpOP org = Values.fromValue(CmpOP.class, (int) Bins.int4(ctx.getFirstByte(), 2));
                CmpOP flag = Values.fromValue(CmpOP.class, rf.get());
                boolean valid = false;

                switch (flag)
                {
                    case A:
                        if (org == CmpOP.AE || org == CmpOP.A || org == CmpOP.NZ)
                            valid = true;
                        break;
                    case B:
                        if (org == CmpOP.BE || org == CmpOP.B || org == CmpOP.NZ)
                            valid = true;
                        break;
                    case Z:
                        if (org == CmpOP.Z || org == CmpOP.AE || org == CmpOP.BE)
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
                vm.push(rp.get() + InstructionType.length(instruction));
                rp.set(opv1);
                break;
            case RET:
                rp.set(vm.pop());
                break;
            case CMP:
            {
                float a = opv1;
                float b = opv2;
                if (dataType == DataType.T_FLOAT)
                {
                    a = Bins.float32(opv1);
                    b = Bins.float32(opv2);
                }
                float c = a - b;
                if (c > 0)
                    rf.set(CmpOP.A.get());
                else if (c < 0)
                    rf.set(CmpOP.B.get());
                else
                    rf.set(CmpOP.Z.get());
            }
            break;
            case CAL:
            {

                CalOP op = Values.fromValue(CalOP.class, ctx.getSpecialByte());
                // 返回结果为 r0
                double a = opv1;
                double b = opv2;
                if (dataType.equals(DataType.T_FLOAT))
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
                    case T_FLOAT:
                        ret = Bins.int32((float) c);
                        break;
                    case T_BYTE:
                        ret &= 0xff;
                        break;
                    case T_WORD:
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
