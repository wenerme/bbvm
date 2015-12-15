package me.wener.bbvm.vm;

import com.google.common.base.MoreObjects;
import com.google.common.base.Objects;
import me.wener.bbvm.util.val.IsInt;

import static me.wener.bbvm.util.val.IntEnums.fromInt;

/**
 * @author wener
 * @since 15/12/10
 */
public class Operand extends AbstractValue<Operand> {
    protected int value;
    AddressingMode addressingMode;
    transient VM vm;
    transient Symbol symbol;

    public Symbol getSymbol() {
        return symbol;
    }

    public Operand setSymbol(Symbol symbol) {
        this.symbol = symbol;
        return this;
    }

    /**
     * Set this operand like {@code other}
     */
    Operand set(Operand other) {
        if (other == null) {
            return this;
        }
        this.value = other.value;
        this.addressingMode = other.addressingMode;
        this.vm = other.vm;
        this.symbol = other.symbol;
        return this;
    }

    /**
     * @return The interval value of this operand
     */
    public int getInterval() {
        return value;
    }

    public Operand setInternal(IsInt v) {
        setInternal(v.asInt());
        return this;
    }

    public Operand setInternal(int value) {
        this.value = value;
        return this;
    }

    public AddressingMode getAddressingMode() {
        return addressingMode;
    }

    public Operand setAddressingMode(AddressingMode addressingMode) {
        this.addressingMode = addressingMode;
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
                return vm.getRegister(fromInt(RegisterType.class, getInterval())).get();
            case REGISTER_DEFERRED:
                return vm.getMemory().read(vm.getRegister(fromInt(RegisterType.class, getInterval())).get());
            case IMMEDIATE:
                return getInterval();
            case DIRECT:
                return vm.getMemory().read(getInterval());
            default:
                throw new AssertionError();
        }
    }

    public Operand set(int v) {
        switch (addressingMode) {

            case REGISTER:
                vm.getRegister(fromInt(RegisterType.class, getInterval())).set(v);
                break;
            case REGISTER_DEFERRED:
                vm.getMemory().write(vm.getRegister(fromInt(RegisterType.class, getInterval())).get(), v);
                break;
            case IMMEDIATE:
                throw new AssertionError("Set a IMMEDIATE operand");
            case DIRECT:
                vm.getMemory().write(getInterval(), v);
                break;
            default:
                throw new AssertionError();
        }
        return this;
    }

    public String toAssembly() {
        switch (addressingMode) {
            case REGISTER:
                return fromInt(RegisterType.class, getInterval()).toString();
            case REGISTER_DEFERRED:
                return "[" + fromInt(RegisterType.class, getInterval()) + "]";
            case IMMEDIATE:
                if ((symbol == null || symbol.getValue() != getInterval()) && vm != null)
                    symbol = vm.getSymbol(getInterval());
                return String.valueOf(symbol != null ? symbol.getName() : getInterval());
            case DIRECT:
                if ((symbol == null || symbol.getValue() != getInterval()) && vm != null)
                    symbol = vm.getSymbol(getInterval());
                return "[" + (symbol != null ? symbol.getName() : getInterval()) + "]";
            default:
                throw new AssertionError();
        }
    }

    public void reset() {
        symbol = null;
        setInternal(0);
        addressingMode = null;
    }

    @Override
    public String toString() {
        return MoreObjects.toStringHelper(this)
                .add("value", getInterval())
                .add("mode", addressingMode)
                .toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof Operand)) return false;
        if (!super.equals(o)) return false;
        Operand operand = (Operand) o;
        return addressingMode == operand.addressingMode;
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(super.hashCode(), addressingMode);
    }
}
