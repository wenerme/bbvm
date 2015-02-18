package me.wener.bbvm.system;

import java.io.IOException;
import java.io.Writer;
import lombok.Getter;
import lombok.Setter;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;

@Accessors(chain = true, fluent = true)
@Slf4j
public class CPU
{
    private final Register rp = new Register("rp");
    private final Register rb = new Register("rb");
    private final Register rs = new Register("rs");
    private final Register rf = new Register("rf");
    private final Register r0 = new Register("r0");
    private final Register r1 = new Register("r1");
    private final Register r2 = new Register("r2");
    private final Register r3 = new Register("r3");
    @Getter
    private final Operand a = new Operand(this);
    @Getter
    private final Operand b = new Operand(this);
    @Getter
    private Memory memory;
    @Getter
    private Opcode opcode;
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
}
