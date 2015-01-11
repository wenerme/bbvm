package me.wener.bbvm.neo;

import me.wener.bbvm.neo.define.Flags;

public abstract class AbstractOperand implements Operand
{
    private int addressingMode;

    @Override
    public Integer get()
    {
        throw new UnsupportedOperationException();
    }

    /**
     * @throws UnsupportedOperationException 有可能该操作数不能被写
     */
    @Override
    public void set(Integer v)
    {
        throw new UnsupportedOperationException();
    }

    @Override
    public int addressingMode()
    {
        return addressingMode;
    }

    @Override
    public Operand addressingMode(int mode)
    {
        addressingMode = mode;
        return this;
    }

    @Override
    public String toString()
    {
        switch (addressingMode)
        {
            case Flags.IMMEDIATE:
            case Flags.REGISTER:
                return toString0();
            case Flags.REGISTER_DEFERRED:
            case Flags.DIRECT:
                return "[" + toString0() + "]";
        }
        return super.toString();
    }

    protected String toString0()
    {
        return String.valueOf(get());
    }
}
