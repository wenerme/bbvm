package me.wener.bbvm.system;

import java.io.Serializable;
import lombok.AccessLevel;
import lombok.Data;
import lombok.Setter;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.system.api.AddressingMode;
import me.wener.bbvm.system.api.Defines;
import me.wener.bbvm.system.api.Operand;
import me.wener.bbvm.system.api.Register;
import me.wener.bbvm.system.api.RegisterType;
import me.wener.bbvm.utils.Bins;
import me.wener.bbvm.utils.val.Values;

@Data
@Accessors(chain = true, fluent = true)
@Slf4j
public class OperandImpl implements Operand, Serializable
{
    private VmCPU cpu;
    @Setter(AccessLevel.NONE)//lombok 生成的方法导致冲突
    private Integer value;
    private AddressingMode addressingMode;

    public OperandImpl(VmCPU cpu)
    {
        this.cpu = cpu;
    }

    public OperandImpl()
    {
    }

    @Override
    public me.wener.bbvm.system.api.Resource asResource(String res)
    {
        return cpu.resources(res).get(get());
    }

    @Override
    public OperandImpl asString(String v)
    {
        if (addressingMode == AddressingMode.IMMEDIATE)
        {
            throw new UnsupportedOperationException();
        }
        asResource(Defines.RES_STRING).set(v);
        return this;
    }

    @Override
    public String asString()
    {
        Integer v = get();
        if (v < 0)
        {
            return cpu.resources(Defines.RES_STRING).get(v).as();
        } else
        {
            return cpu.memory().readString(v);
        }
    }

    public Integer get()
    {
        switch (addressingMode)
        {
            case REGISTER:
                return asRegister().get();
            case REGISTER_DEFERRED:
                return cpu.memory().readInt(asRegister().get());
            case IMMEDIATE:
                return value;
            case DIRECT:
                return cpu.memory().readInt(value);
        }
        throw new UnsupportedOperationException();
    }

    private Register asRegister()
    {
        return cpu.register(Values.fromValue(RegisterType.class, value));
    }

    @Override
    public Operand value(RegisterType v)
    {
        return value(v.get());
    }

    public Operand value(Integer v)
    {
        value = v;
        return this;
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
                asRegister().set(value);
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

    @Override
    public String toAssembly()
    {
        switch (addressingMode)
        {
            case REGISTER:
                return Values.fromValue(RegisterType.class, value).toString();
            case REGISTER_DEFERRED:
                return "[ " + Values.fromValue(RegisterType.class, value) + " ]";
            case IMMEDIATE:
                return value.toString();
            case DIRECT:
                return "[ " + value + " ]";
        }
        throw new UnsupportedOperationException();
    }
}
