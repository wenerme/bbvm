package me.wener.bbvm.vm;


/**
 * @author wener
 * @since 15/12/15
 */
public class Reg implements Register {
    final RegisterType type;
    final VM vm;
    int value;

    public Reg(RegisterType type, VM vm) {
        this.type = type;
        this.vm = vm;
    }

    public VM getVm() {
        return vm;
    }

    @Override
    public int get() {
        return value;
    }

    public Reg set(int value) {
        this.value = value;
        return this;
    }

    @Override
    public RegisterType getType() {
        return type;
    }
}
