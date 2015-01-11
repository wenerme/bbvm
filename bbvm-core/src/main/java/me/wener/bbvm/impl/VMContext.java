package me.wener.bbvm.impl;

import me.wener.bbvm.api.BBVm;
import me.wener.bbvm.def.RegType;
import me.wener.bbvm.utils.Bins;
import me.wener.bbvm.utils.val.Values;

public class VMContext
{
    protected BBVm vm;
    protected Reg rp;
    protected Reg rb;
    protected Reg rs;
    protected Reg rf;
    protected Reg r0;
    protected Reg r1;
    protected Reg r2;
    protected Reg r3;
    protected byte[] memory;
    protected StringHandlePool stringPool;

    public void initialize(BBVm vm)
    {
        this.vm = vm;
        memory = vm.getMemory();
        rp = vm.getRegister(RegType.rp);
        rb = vm.getRegister(RegType.rb);
        rs = vm.getRegister(RegType.rs);
        rf = vm.getRegister(RegType.rf);
        r0 = vm.getRegister(RegType.r0);
        r1 = vm.getRegister(RegType.r1);
        r2 = vm.getRegister(RegType.r2);
        r3 = vm.getRegister(RegType.r3);
    }


    protected UnsupportedOperationException unsupport(String format, Object... args)
    {
        return unsupport(String.format(format, args));
    }

    protected UnsupportedOperationException unsupport(String str)
    {
        return new UnsupportedOperationException(str);
    }

    public void push(int v)
    {
        Bins.int32l(memory, rs.get(), v);
        rs.set(rs.get() - 4);
    }

    public int pop()
    {
        rs.set(rs.get() + 4);
        return Bins.int32l(memory, rs.get());
    }

    protected Integer[] readParameters(int n, int offset)
    {
        Integer[] parameters = new Integer[n];

        for (int i = 0; i < n; i++)
        {
            parameters[n - i - 1] = Bins.int32l(memory, offset);
            offset += 4;
        }

        return parameters;
    }

    public Reg getRegister(int reg)
    {
        return getRegister(Values.fromValue(RegType.class, reg));
    }

    public Reg getRegister(RegType r)
    {
        switch (r)
        {
            case rp:
                return rp;
            case rf:
                return rf;
            case rs:
                return rs;
            case rb:
                return rb;
            case r0:
                return r0;
            case r1:
                return r1;
            case r2:
                return r2;
            case r3:
                return r3;
            default:
                throw new IllegalArgumentException("未知的寄存器 :" + r);
        }
    }

}
