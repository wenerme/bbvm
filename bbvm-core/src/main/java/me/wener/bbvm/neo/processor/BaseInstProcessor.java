package me.wener.bbvm.neo.processor;

import static me.wener.bbvm.neo.inst.def.Flags.*;

import me.wener.bbvm.neo.Operand;
import me.wener.bbvm.neo.inst.Inst;
import me.wener.bbvm.neo.inst.def.InstructionType;
import me.wener.bbvm.utils.Bins;

public class BaseInstProcessor extends BasicProcessor
{
    @Override
    public ProcessResult process(ProcessContext ctx, InstructionType instructionType, Inst instruction, Operand a, Operand b, int dataType)
    {
        switch (instructionType)
        {
            case NOP:
                break;
            case LD:
            {
                switch (dataType)
                {
                    case T_DWORD:
                    case T_FLOAT:
                    case T_INT:
                        a.set(b.get());
                        break;
                    case T_BYTE:
                        a.set(a.get() & 0xffffff00 | (b.get() & 0xff));
                        break;
                    case T_WORD:
                        a.set(a.get() & 0xffff0000 | (b.get() & 0xffff));
                        break;
                    default:
                        throw new AssertionError("未知的数据类型:" + dataType);
                }
            }
            break;
            case PUSH:
                push(a.get());
                break;
            case POP:
                a.set(pop());
            case JMP:
                rp.set(a.get());
                break;
            case JPC:
            {
                me.wener.bbvm.neo.inst.JPC jpc = (me.wener.bbvm.neo.inst.JPC) instruction;
                // JPC 的数据类型为比较操作
                int org = jpc.compare;
                int flag = rf.get();
                boolean valid = false;
                // 判断是否兼容
                switch (flag)
                {
                    case A:
                        if (org == AE || org == A || org == NZ)
                            valid = true;
                        break;
                    case B:
                        if (org == BE || org == B || org == NZ)
                            valid = true;
                        break;
                    case Z:
                        if (org == Z || org == AE || org == BE)
                            valid = true;
                        break;
                    default:
                        if (org == flag)
                            valid = true;
                }

                if (valid)
                    rp.set(a.get());
            }
            break;
            case CALL:
                // 设置返回位置为下一句的开始
                push(rp.get() + instructionType.length());
                rp.set(a.get());
                break;
            case RET:
                rp.set(pop());
                break;

            case CMP:
            {
                float oa;
                float ob;

                if (dataType == T_FLOAT)
                {
                    oa = Bins.float32(a.get());
                    ob = Bins.float32(b.get());
                } else
                {
                    oa = a.get();
                    ob = b.get();
                }
                float oc = oa - ob;
                if (oc > 0)
                    rf.set(A);
                else if (oc < 0)
                    rf.set(B);
                else
                    rf.set(Z);
            }
            break;
            case CAL:
            {
                me.wener.bbvm.neo.inst.CAL cal = (me.wener.bbvm.neo.inst.CAL) instruction;
                // 返回结果为 r0
                double oa;
                double ob;
                double oc;

                if (dataType == T_FLOAT)
                {
                    oa = Bins.float32(a.get());
                    ob = Bins.float32(b.get());
                } else
                {
                    oa = a.get();
                    ob = b.get();
                }
                switch (cal.operator)
                {
                    case ADD:
                        oc = oa + ob;
                        break;
                    case DIV:
                        oc = oa / ob;
                        break;
                    case MOD:
                        oc = oa % ob;
                        break;
                    case MUL:
                        oc = oa * ob;
                        break;
                    case SUB:
                        oc = oa - ob;
                        break;
                    default:
                        throw new AssertionError("未知计算操作: " + cal.operator);
                }
                int ret = (int) oc;
                // 值返回归约
                switch (dataType)
                {
                    case T_FLOAT:
                        ret = Bins.int32((float) oc);
                        break;
                    case T_BYTE:
                        ret &= 0xff;
                        break;
                    case T_WORD:
                        ret &= 0xffff;
                        break;
                }
                a.set(ret);
            }
            break;
            case EXIT:
                break;
            default:
                return Results.keepGoing();
        }
        return Results.processed();
    }
}
