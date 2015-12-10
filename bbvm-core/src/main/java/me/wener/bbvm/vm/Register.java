package me.wener.bbvm.vm;

import me.wener.bbvm.util.val.IsInt;
import org.apache.commons.lang3.mutable.MutableInt;

/**
 * @author wener
 * @since 15/12/10
 */
public class Register extends MutableInt {
    RegisterType type;
    VM vm;

    public Register(RegisterType type, VM vm) {
        this.type = type;
        this.vm = vm;
    }

    public Register() {
    }

    public VM getVm() {
        return vm;
    }

    public Register setVm(VM vm) {
        this.vm = vm;
        return this;
    }


    public RegisterType getType() {
        return type;
    }

    public Register setType(RegisterType type) {
        this.type = type;
        return this;
    }

    public void setValue(IsInt v) {
        setValue(v.asInt());
    }
}
