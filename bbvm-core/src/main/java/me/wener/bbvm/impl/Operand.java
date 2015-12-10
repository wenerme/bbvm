package me.wener.bbvm.impl;

import me.wener.bbvm.util.Bins;
import me.wener.bbvm.util.val.IntHolder;

public abstract class Operand implements IntHolder
{
    public final static Operand INVALID = new InvalidOperand();

    /**
     * 获取地址相关的操作数
     */
    public static Operand address(int v, byte[] memory)
    {
        return new AddressedOperand(v, memory);
    }

    /**
     * 直接数,只能取不能读
     */
    public static Operand value(int v)
    {
        return new ValueOperand(v);
    }

    /**
     * 将一个 Holder 包装为一个操作数类
     */
    public static Operand holder(IntHolder v)
    {
        return new HolderOperand(v);
    }

    public static Operand invalid()
    {
        return INVALID;
    }

    /**
     * 包装一个间接寻址
     */
    public static Operand indirect(IntHolder v, byte[] memory)
    {
        return new IndirectOperand(v, memory);
    }

    @Override
    public Integer get() {
        return asInt();
    }

    @Override
    public int asInt()
    {
        throw new UnsupportedOperationException();
    }

    /**
     * @throws java.lang.UnsupportedOperationException 有可能该操作数不能被写
     */
    @Override
    public void set(Integer v)
    {
        throw new UnsupportedOperationException();
    }

    private static class InvalidOperand extends Operand
    {
        @Override
        public int asInt()
        {
            return 0;
        }

        @Override
        public String toString()
        {
            return "invalid";
        }
    }

    private static class HolderOperand extends Operand
    {
        private final IntHolder value;

        HolderOperand(IntHolder value)
        {
            this.value = value;
        }

        @Override
        public int asInt()
        {
            return value.asInt();
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

    private static class ValueOperand extends Operand
    {
        private final int value;

        ValueOperand(int value)
        {
            this.value = value;
        }

        @Override
        public int asInt()
        {
            return value;
        }

        @Override
        public String toString()
        {
            return String.valueOf(value);
        }
    }

    private static class AddressedOperand extends Operand
    {
        private final int address;
        private final byte[] memory;

        AddressedOperand(int address, byte[] memory)
        {
            this.address = address;
            this.memory = memory;
        }

        @Override
        public int asInt()
        {
            return Bins.int32l(memory, address);
        }

        @Override
        public void set(Integer v)
        {
            Bins.int32l(memory, address, v);
        }

        @Override
        public String toString()
        {
            return String.valueOf(address);
        }
    }

    private static class IndirectOperand extends AddressedOperand
    {

        private final IntHolder origin;

        IndirectOperand(IntHolder v, byte[] memory)
        {
            super(v.asInt(), memory);
            this.origin = v;
        }

        @Override
        public String toString()
        {
            return "[ " + origin.toString() + " ]";
        }
    }
}
