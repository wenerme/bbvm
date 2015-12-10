package me.wener.bbvm.vm;

import me.wener.bbvm.util.val.IsInt;
import org.apache.commons.lang3.mutable.MutableInt;

import static me.wener.bbvm.util.val.IntEnums.fromInt;

/**
 * @author wener
 * @since 15/12/10
 */
public class Operand extends MutableInt {
    int value;
    AddressingMode addressingMode;
    VM vm;

    public AddressingMode getAddressingMode() {
        return addressingMode;
    }

    public Operand setAddressingMode(AddressingMode addressingMode) {
        this.addressingMode = addressingMode;
        return this;
    }

    public Operand setValue(IsInt v) {
        value = v.asInt();
        return this;
    }

    public VM getVm() {
        return vm;
    }

    public Operand setVm(VM vm) {
        this.vm = vm;
        return this;
    }

    public int get() {
        switch (addressingMode) {
            case REGISTER:
                return vm.getRegister(fromInt(RegisterType.class, value)).intValue();
            case REGISTER_DEFERRED:
                return vm.getMemory().read(vm.getRegister(fromInt(RegisterType.class, value)).intValue());
            case IMMEDIATE:
                return value;
            case DIRECT:
                return vm.getMemory().read(value);
            default:
                throw new AssertionError();
        }
    }

    public Operand set(int v) {
        switch (addressingMode) {

            case REGISTER:
                vm.getRegister(fromInt(RegisterType.class, value)).setValue(v);
                break;
            case REGISTER_DEFERRED:
                vm.getMemory().write(vm.getRegister(fromInt(RegisterType.class, value)).intValue(), v);
                break;
            case IMMEDIATE:
                throw new AssertionError("Set a IMMEDIATE operand");
            case DIRECT:
                vm.getMemory().write(value, v);
                break;
            default:
                throw new AssertionError();
        }
        return this;
    }

    public String toAssembly() {
        switch (addressingMode) {
            case REGISTER:
                return fromInt(RegisterType.class, value).toString();
            case REGISTER_DEFERRED:
                return "[" + fromInt(RegisterType.class, value) + "]";
            case IMMEDIATE:
                return String.valueOf(value);
            case DIRECT:
                return "[" + value + "]";
            default:
                throw new AssertionError();
        }
    }

    public Operand set(float v) {
        set(Float.floatToIntBits(v));
        return this;
    }

    public float getFloat() {
        return Float.intBitsToFloat(get());
    }
}
