package me.wener.bbvm.system;

public interface CPU
{
    OpStatus opstatus();

    VmStatus vmstatus();

    void step();
}
