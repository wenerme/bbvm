package me.wener.bbvm.vm;

/**
 * @author wener
 * @since 15/12/10
 */
public class Register extends AbstractValue<Register> {
    final RegisterType type;
    final VM vm;
    int value;

    public Register(RegisterType type, VM vm) {
        this.type = type;
        this.vm = vm;
    }

    public VM getVm() {
        return vm;
    }

    public int get() {
        return value;
    }

    public Register set(int value) {
        this.value = value;
        return this;
    }

    public RegisterType getType() {
        return type;
    }
}
