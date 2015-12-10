package me.wener.bbvm.api;

import me.wener.bbvm.impl.Reg;
import me.wener.bbvm.vm.RegisterType;

public interface BBVm
{
    byte[] getMemory();

    void reset();

    void start();

    void push(int v);

    int pop();

    void exit();

    Reg getRegister(int reg);

    Reg getRegister(RegisterType r);
}
