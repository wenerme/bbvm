package me.wener.bbvm.impl;

import me.wener.bbvm.api.BBVm;
import me.wener.bbvm.util.Bins;
import me.wener.bbvm.util.val.IntEnums;
import me.wener.bbvm.vm.RegisterType;

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
        rp = vm.getRegister(RegisterType.RP);
        rb = vm.getRegister(RegisterType.RB);
        rs = vm.getRegister(RegisterType.RS);
        rf = vm.getRegister(RegisterType.RF);
        r0 = vm.getRegister(RegisterType.R0);
        r1 = vm.getRegister(RegisterType.R1);
        r2 = vm.getRegister(RegisterType.R2);
        r3 = vm.getRegister(RegisterType.R3);
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
        Bins.int32l(memory, rs.asInt(), v);
        rs.set(rs.asInt() - 4);
    }

    public int pop()
    {
        rs.set(rs.asInt() + 4);
        return Bins.int32l(memory, rs.asInt());
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
        return getRegister(IntEnums.fromInt(RegisterType.class, reg));
    }

    public Reg getRegister(RegisterType r)
    {
        switch (r)
        {
            case RP:
                return rp;
            case RF:
                return rf;
            case RS:
                return rs;
            case RB:
                return rb;
            case R0:
                return r0;
            case R1:
                return r1;
            case R2:
                return r2;
            case R3:
                return r3;
            default:
                throw new IllegalArgumentException("未知的寄存器 :" + r);
        }
    }

}
