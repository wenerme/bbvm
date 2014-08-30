package me.wener.bbvm.core.constant;

import me.wener.bbvm.core.IsValue;

public enum BBEnv implements IsValue<Integer>
{
    ENV_SIM(0),
    ENV_9288(9288),
    ENV_9188(9188),
    ENV_9288T(9287),
    ENV_9288S(9286),
    ENV_9388(9388);
    private final int value;

    BBEnv(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }
}
