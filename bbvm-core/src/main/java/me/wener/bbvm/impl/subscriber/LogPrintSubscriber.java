package me.wener.bbvm.impl.subscriber;

import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.event.InstEvent;
import me.wener.bbvm.impl.InstructionContext;
import me.wener.bbvm.impl.Operand;
import me.wener.bbvm.impl.VMContext;
import me.wener.bbvm.util.Bins;
import me.wener.bbvm.util.val.Values;
import me.wener.bbvm.vm.CompareType;
import me.wener.bbvm.vm.DataType;
import me.wener.bbvm.vm.Opcode;

public class LogPrintSubscriber extends VMContext
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
                    .append(String.format("R0= %s, R1= %s, R2= %s, R3= %s, RS= %s, RB= %s, RP= %s, RF= %s",
                            r0.asInt(), r1.asInt(), r2.asInt(), r3.asInt(), rs.asInt(), rb.asInt(), rp.asInt(), rf.asInt()));
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
        final Opcode instruction = ctx.getInstruction();

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

                rp.set(op1.asInt());
                break;
            case JPC:
            {
                // JPC 的数据类型为比较操作
                CompareType org = Values.fromValue(CompareType.class, (int) Bins.int4(ctx.getFirstByte(), 2));
                CompareType flag = Values.fromValue(CompareType.class, rf.asInt());

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
