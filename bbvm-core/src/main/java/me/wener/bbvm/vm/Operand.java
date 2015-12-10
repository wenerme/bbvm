package me.wener.bbvm.vm;

import com.google.common.base.MoreObjects;
import me.wener.bbvm.util.val.IsInt;
import org.apache.commons.lang3.mutable.MutableInt;

import static me.wener.bbvm.util.val.IntEnums.fromInt;

/**
 * @author wener
 * @since 15/12/10
 */
public class Operand extends MutableInt {
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
        setValue(v.asInt());
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
                return vm.getRegister(fromInt(RegisterType.class, intValue())).intValue();
            case REGISTER_DEFERRED:
                return vm.getMemory().read(vm.getRegister(fromInt(RegisterType.class, intValue())).intValue());
            case IMMEDIATE:
                return intValue();
            case DIRECT:
                return vm.getMemory().read(intValue());
            default:
                throw new AssertionError();
        }
    }

    public Operand set(int v) {
        switch (addressingMode) {

            case REGISTER:
                vm.getRegister(fromInt(RegisterType.class, intValue())).setValue(v);
                break;
            case REGISTER_DEFERRED:
                vm.getMemory().write(vm.getRegister(fromInt(RegisterType.class, intValue())).intValue(), v);
                break;
            case IMMEDIATE:
                throw new AssertionError("Set a IMMEDIATE operand");
            case DIRECT:
                vm.getMemory().write(intValue(), v);
                break;
            default:
                throw new AssertionError();
        }
        return this;
    }

    public String toAssembly() {
        switch (addressingMode) {
            case REGISTER:
                return fromInt(RegisterType.class, intValue()).toString();
            case REGISTER_DEFERRED:
                return "[" + fromInt(RegisterType.class, intValue()) + "]";
            case IMMEDIATE:
                return String.valueOf(intValue());
            case DIRECT:
                return "[" + intValue() + "]";
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

    @Override
    public String toString() {
        return MoreObjects.toStringHelper(this)
                .add("value", intValue())
                .add("mode", addressingMode)
                .toString();
    }
}
