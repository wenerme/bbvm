package me.wener.bbvm.core;

import me.wener.bbvm.utils.Bins;

public class Reg implements IntHolder
{
    private int value;
    private final String name;
    private byte[] internal = new byte[8];

    public Reg(String name) {this.name = name;}

    public int getInt()
    {
        return Bins.int32(internal, 0);
    }

    public Reg setInt(int v)
    {
        Bins.int32(internal, 0, v);
        return this;
    }

    public float getFloat()
    {
        return Bins.float32(internal, 0);
    }

    public Reg setFloat(float v)
    {
        Bins.float32(internal, 0, v);
        return this;
    }

    @Override
    public Integer get()
    {
        return value;
    }

    @Override
    public void set(Integer v)
    {
        value = v;
    }

    public String getName()
    {
        return name;
    }

    @Override
    public String toString()
    {
        return name;
    }
}
