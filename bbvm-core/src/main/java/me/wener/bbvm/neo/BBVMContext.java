package me.wener.bbvm.neo;

import io.netty.buffer.ByteBuf;
import me.wener.bbvm.neo.define.CalculateType;
import me.wener.bbvm.neo.define.CompareType;
import me.wener.bbvm.neo.define.DataType;
import me.wener.bbvm.neo.define.Flags;
import me.wener.bbvm.neo.define.InstructionType;
import me.wener.bbvm.neo.define.RegisterType;
import me.wener.bbvm.utils.val.Values;

public class BBVMContext
{
    private final Register rp = new Register("rp");
    private final Register rb = new Register("rb");
    private final Register rs = new Register("rs");
    private final Register rf = new Register("rf");
    private final Register r0 = new Register("r0");
    private final Register r1 = new Register("r1");
    private final Register r2 = new Register("r2");
    private final Register r3 = new Register("r3");
    private ByteBuf memory;

    static
    {
        Values.cache(InstructionType.class,
                RegisterType.class,
                CompareType.class,
                CalculateType.class,
                DataType.class);
    }

    public BBVMContext(ByteBuf memory)
    {
        this.memory = memory;
    }

    public ByteBuf memory()
    {
        return memory;
    }

    public Register register(int type)
    {
        switch (type)
        {
            case Flags.rp:
                return rp;
            case Flags.rf:
                return rf;
            case Flags.rs:
                return rs;
            case Flags.rb:
                return rb;
            case Flags.r0:
                return r0;
            case Flags.r1:
                return r1;
            case Flags.r2:
                return r2;
            case Flags.r3:
                return r3;
            default:
                throw new IllegalArgumentException("未知的寄存器 :" + type);
        }
    }
}
