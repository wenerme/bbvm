package me.wener.bbvm.system.api;

/**
 * 虚拟机状态
 */
public interface VmStatus
{
    Register register(RegisterType type);

    Register rp();

    Register rb();

    Register rs();

    Register rf();

    Register r0();

    Register r1();

    Register r2();

    Register r3();
}
