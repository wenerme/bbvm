package me.wener.bbvm.system;

import java.util.Map;
import me.wener.bbvm.system.internal.ResourcePool;

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

    ResourcePool resources(String resourceName);

    Map<String, ResourcePool> resources();
}
