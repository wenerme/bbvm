package me.wener.bbvm.core;

import me.wener.bbvm.utils.Bins;

public class Reg implements IntHolder
{
    private int value;
    private final String name;

    public Reg(String name) {this.name = name;}

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
