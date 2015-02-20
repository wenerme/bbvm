package me.wener.bbvm.system.api;

public interface CPU
{
    OpStatus opstatus();

    VmStatus vmstatus();

    void step();
}
