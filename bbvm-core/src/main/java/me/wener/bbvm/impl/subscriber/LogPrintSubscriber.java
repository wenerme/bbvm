package me.wener.bbvm.impl.subscriber;

import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.def.CmpOP;
import me.wener.bbvm.def.DataType;
import me.wener.bbvm.def.Instruction;
import me.wener.bbvm.event.InstEvent;
import me.wener.bbvm.impl.InstructionContext;
import me.wener.bbvm.impl.Operand;
import me.wener.bbvm.utils.Bins;
import me.wener.bbvm.utils.val.Values;

public class LogPrintSubscriber extends AbstractVMSubscriber
{

    private boolean logInst;

    protected void log(Object... objects)
    {
        System.out.println(logString(true, objects));
    }

    protected String logString(boolean debug, Object... objects)
    {
        StringBuilder builder = new StringBuilder();
        boolean lastIsOperand = false;
        for (Object object : objects)
        {
            if (lastIsOperand)
            {
                builder.append(", ");
            }
            builder.append(object).append(" ");
            lastIsOperand = object instanceof Operand;
        }
        if (debug)
        {
            builder.append("\n;")
                   .append(String.format("r0= %s, r1= %s, r2= %s, r3= %s, rs= %s, rb= %s, rp= %s, rf= %s",
                           r0.get(), r1.get(), r2.get(), r3.get(), rs.get(), rb.get(), rp.get(), rf.get()));
        }
        return builder.toString();
    }

    @Subscribe
    protected void subscribeInst(InstEvent e)
    {
        doInstruction(e.getContext());
    }

    protected void doInstruction(InstructionContext ctx)
    {
        final Operand op1 = ctx.getOp1();
        final Operand op2 = ctx.getOp2();
        final DataType dataType = ctx.getDataType();
        final Instruction instruction = ctx.getInstruction();

        switch (instruction)
        {
            case NOP:
                if (logInst)
                    log(instruction);
                break;
            case LD:
                if (logInst)
                    log(instruction, dataType, op1, op2);
                break;
            case PUSH:
                if (logInst)
                    log(instruction, op1);
                break;
            case POP:
                if (logInst)
                    log(instruction, op1);
                break;
            case IN:
                if (logInst)
                    log(instruction, op1, op2);
                break;
            case OUT:
                if (logInst)
                    log(instruction, op1, op2);
                break;
            case JMP:
                if (logInst)
                    log(instruction, op1);

                rp.set(op1.get());
                break;
            case JPC:
            {
                // JPC 的数据类型为比较操作
                CmpOP org = Values.fromValue(CmpOP.class, (int) Bins.int4(ctx.getFirstByte(), 2));
                CmpOP flag = Values.fromValue(CmpOP.class, rf.get());

                if (logInst)
                    log(instruction, org, op1);
            }
            break;
            case CALL:
                if (logInst)
                    log(instruction, op1);
                break;
            case RET:
                if (logInst)
                    log(instruction);

                break;
            case CMP:
                if (logInst)
                    log(instruction, dataType, op1, op2);
                break;
            case CAL:
                if (logInst)
                    log(instruction, op1);
                break;
            case EXIT:
                if (logInst)
                    log(instruction);
                break;
            default:
                throw unsupport("未知指令: %s", instruction);
        }
    }
}
