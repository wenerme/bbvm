package me.wener.bbvm.system.internal;

import java.io.Serializable;
import lombok.AccessLevel;
import lombok.Data;
import lombok.Setter;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.system.AddressingMode;
import me.wener.bbvm.system.Defines;
import me.wener.bbvm.system.RegisterType;
import me.wener.bbvm.utils.Bins;
import me.wener.bbvm.utils.val.Values;

@Data
@Accessors(chain = true, fluent = true)
@Slf4j
class Operand implements me.wener.bbvm.system.Operand, Serializable
{
    private VmCPU cpu;
    @Setter(AccessLevel.NONE)//lombok 生成的方法导致冲突
    private Integer value;
    private AddressingMode addressingMode;

    public Operand(VmCPU cpu)
    {
        this.cpu = cpu;
    }

    public Operand()
    {
    }

    @Override
    public me.wener.bbvm.system.Resource asResource(String res)
    {
        return cpu.resources(res).get(get());
    }

    @Override
    public Operand asString(String v)
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

    private me.wener.bbvm.system.Register asRegister()
    {
        return cpu.register(asRegisterType());
    }

    @Override
    public RegisterType asRegisterType()
    {
        return Values.fromValue(RegisterType.class, value);
    }

    @Override
    public me.wener.bbvm.system.Operand value(RegisterType v)
    {
        return value(v.get());
    }

    public me.wener.bbvm.system.Operand value(Integer v)
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
    public Operand asFloat(float v)
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
                return asRegisterType().toString();
            case REGISTER_DEFERRED:
                return "[ " + asRegisterType() + " ]";
            case IMMEDIATE:
                return value.toString();
            case DIRECT:
                return "[ " + value + " ]";
        }
        throw new UnsupportedOperationException();
    }
}
