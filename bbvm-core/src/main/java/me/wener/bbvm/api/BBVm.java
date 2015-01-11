package me.wener.bbvm.api;

import me.wener.bbvm.def.RegType;
import me.wener.bbvm.impl.Reg;

public interface BBVm
{
    byte[] getMemory();

    void reset();

    void start();

    void push(int v);

    int pop();

    void exit();

    Reg getRegister(int reg);

    Reg getRegister(RegType r);
}
