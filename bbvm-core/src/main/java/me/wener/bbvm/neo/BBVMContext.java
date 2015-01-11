package me.wener.bbvm.neo;

import io.netty.buffer.ByteBuf;
import me.wener.bbvm.neo.define.CalculateType;
import me.wener.bbvm.neo.define.CompareType;
import me.wener.bbvm.neo.define.DataType;
import me.wener.bbvm.neo.define.Flags;
import me.wener.bbvm.neo.define.InstructionType;
import me.wener.bbvm.neo.define.RegisterType;
import me.wener.bbvm.utils.val.IntegerHolder;
import me.wener.bbvm.utils.val.Values;

public class BBVMContext
{
    private final Register rp = Registers.create("rp", new MemoryReaderIndexHolder());
    private final Register rb = Registers.create("rb");
    private final Register rs = Registers.create("rs");
    private final Register rf = Registers.create("rf");
    private final Register r0 = Registers.create("r0");
    private final Register r1 = Registers.create("r1");
    private final Register r2 = Registers.create("r2");
    private final Register r3 = Registers.create("r3");

    static
    {
        Values.cache(InstructionType.class,
                RegisterType.class,
                CompareType.class,
                CalculateType.class,
                DataType.class);
    }

    private ByteBuf memory;


    public BBVMContext(ByteBuf memory)
    {
        this.memory = memory;
    }

    public ByteBuf memory()
    {
        return memory;
    }

    public BBVMContext push(int value)
    {
        return this;
    }

    public int pop()
    {
        return 0;
    }

    public Register register(RegisterType type)
    {
        return register(type.get());
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

    private class MemoryReaderIndexHolder implements IntegerHolder
    {
        @Override
        public Integer get()
        {
            return memory.readerIndex();
        }

        @Override
        public void set(Integer v)
        {
            memory.readerIndex(v);
        }
    }
}
