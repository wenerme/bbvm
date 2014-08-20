package me.wener.bbvm.core;

import me.wener.bbvm.utils.Bins;

public class Reg
{
    private int value;

    public int getInt()
    {
        return value;
    }
    public Reg setInt(int v)
    {
        value = v;
        return this;
    }
    public Reg setFloat(float v)
    {
        value = Bins.int32(v);
        return this;
    }
    public float getFloat()
    {
        return Bins.float32(value);
    }
}
