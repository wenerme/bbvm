package me.wener.bbvm.dev;

import me.wener.bbvm.util.IsInt;

/**
 * 环境类型
 */
public enum EnvType implements IsInt {
    ENV_SIM(0),
    ENV_9288(9288),
    ENV_9188(9188),
    ENV_9288T(9287),
    ENV_9288S(9286),
    ENV_9388(9388);
    private final int value;

    EnvType(int value) {
        this.value = value;
    }

    public int asInt() {
        return value;
    }
}
