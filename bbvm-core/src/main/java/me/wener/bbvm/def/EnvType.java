package me.wener.bbvm.def;

import me.wener.bbvm.utils.val.IsInteger;

/**
 * 环境类型
 */
public enum EnvType implements IsInteger
{
    ENV_SIM(0),
    ENV_9288(9288),
    ENV_9188(9188),
    ENV_9288T(9287),
    ENV_9288S(9286),
    ENV_9388(9388);
    private final int value;

    EnvType(int value)
    {
        this.value = value;
    }

    public Integer get()
    {
        return value;
    }
}
