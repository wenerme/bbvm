package me.wener.bbvm.neo;

import io.netty.buffer.ByteBuf;
import me.wener.bbvm.utils.Bins;
import me.wener.bbvm.utils.val.IntegerHolder;

public class Operands
{
    public final static Operand INVALID = new InvalidOperand();

    /**
     * 获取地址相关的操作数
     */
    public static Operand direct(int v, ByteBuf memory)
    {
        return new ByteBufAddressedOperand(v, memory);
    }

    /**
     * 直接数,只能取不能读
     */
    public static Operand immediate(int v)
    {
        return new ValueOperand(v);
    }

    /**
     * 将一个 Holder 包装为一个操作数类
     */
    public static Operand wrap(IntegerHolder v)
    {
        return new HolderOperand(v);
    }

    public static Operand invalid()
    {
        return INVALID;
    }

    /**
     * 包装一个直接寻址
     */
    public static Operand direct(IntegerHolder v, ByteBuf memory)
    {
        return new DirectOperand(v, memory);
    }

    private static class InvalidOperand extends AbstractOperand
    {
        @Override
        public Integer get()
        {
            return 0;
        }

        @Override
        public String toString()
        {
            return "INVALID";
        }
    }

    private static class HolderOperand extends AbstractOperand
    {
        private final IntegerHolder value;

        HolderOperand(IntegerHolder value)
        {
            this.value = value;
        }

        @Override
        public Integer get()
        {
            return value.get();
        }

        @Override
        public void set(Integer v)
        {
            value.set(v);
        }

        @Override
        public String toString()
        {
            return value.toString();
        }
    }

    private static class ValueOperand extends AbstractOperand
    {
        private final int value;

        ValueOperand(int value)
        {
            this.value = value;
        }

        @Override
        public Integer get()
        {
            return value;
        }

        @Override
        public String toString()
        {
            return String.valueOf(value);
        }
    }

    private static class AddressedOperand extends AbstractOperand
    {
        private final int address;
        private final byte[] memory;

        AddressedOperand(int address, byte[] memory)
        {
            this.address = address;
            this.memory = memory;
        }

        @Override
        public Integer get()
        {
            return Bins.int32l(memory, address);
        }

        @Override
        public void set(Integer v)
        {
            Bins.int32l(memory, address, v);
        }

        @Override
        public String toString0()
        {
            return String.valueOf(address);
        }
    }

    private static class ByteBufAddressedOperand extends AbstractOperand
    {
        private final int address;
        private final ByteBuf memory;

        ByteBufAddressedOperand(int address, ByteBuf memory)
        {
            this.address = address;
            this.memory = memory;
        }

        @Override
        public Integer get()
        {
            return memory.getInt(address);
        }

        @Override
        public void set(Integer v)
        {
            memory.setInt(address, v);
        }

        @Override
        public String toString0()
        {
            return String.valueOf(address);
        }
    }

    private static class DirectOperand extends ByteBufAddressedOperand
    {

        private final IntegerHolder origin;

        DirectOperand(IntegerHolder v, ByteBuf memory)
        {
            super(v.get(), memory);
            this.origin = v;
        }

        @Override
        public String toString0()
        {
            return origin.toString();
        }
    }
}
