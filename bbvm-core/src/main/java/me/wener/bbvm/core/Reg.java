package me.wener.bbvm.core;

import me.wener.bbvm.utils.Bins;

public class Reg
{
    private int value;
    private byte[] internal = new byte[8];

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
}
