package me.wener.bbvm.impl;

import me.wener.bbvm.api.IntHolder;

public class Reg extends SimpleValueHolder<Integer> implements IntHolder
{
    private final String name;
    public Reg(String name) {this.name = name;}
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
