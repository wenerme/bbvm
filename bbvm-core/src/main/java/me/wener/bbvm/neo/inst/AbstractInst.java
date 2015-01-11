package me.wener.bbvm.neo.inst;

import me.wener.bbvm.neo.Stringer;

public class AbstractInst implements Inst
{
    @Override
    public String toString()
    {
        return Stringer.string(this);
    }
}
