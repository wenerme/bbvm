package me.wener.bbvm.system;

import me.wener.bbvm.util.val.IntHolder;

public interface Register extends IntHolder
{
    /**
     * @return 寄存器名
     */
    String name();

    RegisterType type();
}
