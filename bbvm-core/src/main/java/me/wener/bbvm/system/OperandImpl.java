package me.wener.bbvm.system;

import java.io.Serializable;
import lombok.Data;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.system.api.AddressingMode;
import me.wener.bbvm.system.api.Operand;
import me.wener.bbvm.system.api.Register;
import me.wener.bbvm.utils.Bins;
import me.wener.bbvm.utils.val.IntegerHolder;
import me.wener.bbvm.utils.val.IsInteger;

@Data
@Accessors(chain = true, fluent = true)
@Slf4j
public class OperandImpl implements Operand, Serializable
{
    private VmCPU cpu;
    private Integer value;
    private IsInteger indirect;
    private AddressingMode addressingMode;

    public OperandImpl(VmCPU cpu)
    {
        this.cpu = cpu;
    }

    public OperandImpl()
    {
    }

    @Override
    public String asString()
    {
        return cpu.memory().readString(get());
    }

    public Integer get()
    {
        switch (addressingMode)
        {
            case REGISTER:
                return indirect.get();
            case REGISTER_DEFERRED:
                return cpu.memory().readInt(indirect.get());
            case IMMEDIATE:
                return value;
            case DIRECT:
                return cpu.memory().readInt(value);
        }
        throw new UnsupportedOperationException();
    }

    @Override
    public float asFloat()
    {
        return Bins.float32(get());
    }

    @Override
    public OperandImpl asFloat(float v)
    {
        set(Bins.int32(v));
        return this;
    }

    public void set(Integer value)
    {
        switch (addressingMode)
        {
            case REGISTER:
                if (indirect instanceof IntegerHolder)
                    ((IntegerHolder) indirect).set(value);
                else
                    throw new UnsupportedOperationException();
                break;
            case REGISTER_DEFERRED:
            case DIRECT:
                cpu.memory().writeInt(get(), value);
                break;
            default:
            case IMMEDIATE:
                throw new UnsupportedOperationException();
        }
    }

    public String toAssembly()
    {
        switch (addressingMode)
        {
            case REGISTER:
                return ((Register) indirect()).name();
            case REGISTER_DEFERRED:
                return "[ " + ((Register) indirect()).name() + " ]";
            case IMMEDIATE:
                return value.toString();
            case DIRECT:
                return "[ " + value + " ]";
        }
        throw new UnsupportedOperationException();
    }
}
