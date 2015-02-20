package me.wener.bbvm.system;

import static me.wener.bbvm.neo.inst.def.CompareTypes.*;

import java.io.IOException;
import java.io.Writer;
import java.nio.ByteBuffer;
import lombok.Getter;
import lombok.Setter;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.system.api.CPU;
import me.wener.bbvm.system.api.CompareType;
import me.wener.bbvm.system.api.DataType;
import me.wener.bbvm.system.api.OpStatus;
import me.wener.bbvm.system.api.Register;
import me.wener.bbvm.system.api.RegisterType;
import me.wener.bbvm.system.api.VmStatus;
import me.wener.bbvm.utils.Bins;
import me.wener.bbvm.utils.val.Values;

@Accessors(chain = true, fluent = true)
@Slf4j
public class VmCPU extends OpStatusImpl implements CPU, VmStatus
{
    @Getter
    private final RegisterImpl rp = new RegisterImpl("rp");
    @Getter
    private final RegisterImpl rb = new RegisterImpl("rb");
    @Getter
    private final RegisterImpl rs = new RegisterImpl("rs");
    @Getter
    private final RegisterImpl rf = new RegisterImpl("rf");
    @Getter
    private final RegisterImpl r0 = new RegisterImpl("r0");
    @Getter
    private final RegisterImpl r1 = new RegisterImpl("r1");
    @Getter
    private final RegisterImpl r2 = new RegisterImpl("r2");
    @Getter
    private final RegisterImpl r3 = new RegisterImpl("r3");

    @Getter
    private VmMemory memory;
    @Getter
    @Setter
    private Writer asmWriter;
    private ByteBuffer stack = ByteBuffer.allocate(1024);

    void reset() {}

    public void step()
    {
        if (asmWriter != null)
        {
            try
            {
                asmWriter.write(disasm());
            } catch (IOException e)
            {
                log.info("Write asm failed.", e);
            }
        }
        process();
    }

    private String disasm()
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
            case LD:
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
                break;
            case CMP:
            case CAL:
                break;
        }
        return sb.toString();
    }

    private void process()
    {
        switch (opcode)
        {
            case NOP:
                NOP();
                break;
            case LD:
                LD();
                break;
            case PUSH:
                PUSH();
                break;
            case POP:
                POP();
                break;
            case IN:
                IN();
                break;
            case OUT:
                OUT();
                break;
            case JMP:
                JMP();
                break;
            case JPC:
                JPC();
                break;
            case CALL:
                CALL();
                break;
            case RET:
                RET();
                break;
            case CMP:
                CMP();
                break;
            case CAL:
                CAL();
                break;
            case EXIT:
                EXIT();
                break;
        }
    }

    void NOP()
    {
    }

    void LD()
    {
        switch (dataType)
        {
            case T_DWORD:
            case T_FLOAT:
            case T_INT:
                a.set(b.get());
                break;
            case T_BYTE:
                a.set((a.get() & 0xffffff00) | (b.get() & 0xff));
                break;
            case T_WORD:
                a.set((a.get() & 0xffff0000) | (b.get() & 0xffff));
                break;
            default:
                throw new AssertionError("未知的数据类型:" + dataType);
        }
    }

    void PUSH()
    {
        push(a.get());
    }

    void POP()
    {
        a.set(pop());
    }

    void IN()
    {
    }

    void OUT()
    {
    }

    void JMP()
    {
        rp.set(a.get());
    }

    void JPC()
    {
        // JPC 的数据类型为比较操作
        CompareType org = compareType;
        CompareType flag = Values.fromValue(CompareType.class, rf.get());


        if (CompareType.isMatch(org, flag))
            rp.set(a.get());
    }

    void CALL()
    {
        push(rp.get() + opcode.length());
        rp.set(a.get());
    }

    void RET()
    {
        rp.set(pop());
    }

    void CMP()
    {
        float oa;
        float ob;

        if (dataType == DataType.T_FLOAT)
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

    void CAL()
    {
        // 返回结果为 r0
        double oa;
        double ob;
        double oc;

        if (dataType == DataType.T_FLOAT)
        {
            oa = Bins.float32(a.get());
            ob = Bins.float32(b.get());
        } else
        {
            oa = a.get();
            ob = b.get();
        }
        switch (calculateType)
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

    void EXIT()
    {
    }

    public void push(int v)
    {
        stack.putInt(v);
    }

    public int pop()
    {
        return stack.getInt();
    }

    @Override
    public OpStatus opstatus()
    {
        return this;
    }

    @Override
    public VmStatus vmstatus()
    {
        return this;
    }

    @Override
    public Register register(RegisterType type)
    {
        switch (type)
        {
            case rp:
                return rp;
            case rf:
                return rf;
            case rs:
                return rs;
            case rb:
                return rb;
            case r0:
                return r0;
            case r1:
                return r1;
            case r2:
                return r2;
            case r3:
                return r3;
        }
        throw new UnsupportedOperationException();
    }
}
