package me.wener.bbvm.neo.inst.def;

import me.wener.bbvm.utils.val.IsInteger;


public enum RegisterType implements IsInteger
{
    rp(RegisterTypes.rp),
    rf(RegisterTypes.rf),
    rs(RegisterTypes.rs),
    rb(RegisterTypes.rb),
    r0(RegisterTypes.r0),
    r1(RegisterTypes.r1),
    r2(RegisterTypes.r2),
    r3(RegisterTypes.r3);

    private final int value;

    RegisterType(int value)
    {
        this.value = value;
    }

    public Integer get()
    {
        return value;
    }
}
