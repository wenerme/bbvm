package me.wener.bbvm.core;

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
