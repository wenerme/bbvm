package me.wener.bbvm.system;

import java.io.IOException;
import java.io.Writer;
import lombok.Getter;
import lombok.Setter;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;

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



    /*
     * Logic
     */

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
    }

    void POP()
    {
    }

    void IN()
    {
    }

    void OUT()
    {
    }

    void JMP()
    {
    }

    void JPC()
    {
    }

    void CALL()
    {
    }

    void RET()
    {
    }

    void CMP()
    {
    }

    void CAL()
    {
    }

    void EXIT()
    {
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
