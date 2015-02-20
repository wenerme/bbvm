package me.wener.bbvm.system.api;

import me.wener.bbvm.utils.val.IntegerHolder;

public interface Register extends IntegerHolder
{
    /**
     * @return 寄存器名
     */
    String name();

    RegisterType type();
}
