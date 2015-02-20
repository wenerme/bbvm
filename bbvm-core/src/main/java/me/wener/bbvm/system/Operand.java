package me.wener.bbvm.system;

import lombok.Data;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.utils.val.IntegerHolder;
import me.wener.bbvm.utils.val.IsInteger;

@Data
@Accessors(chain = true, fluent = true)
@Slf4j
public class Operand implements IntegerHolder
{
    private VmCPU cpu;
    private Integer value;
    private IsInteger indirect;
    private AddressingMode addressingMode;

    public Operand(VmCPU cpu)
    {
        this.cpu = cpu;
    }

    public Operand()
    {
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

    public String asString()
    {
        return cpu.memory().readString(get());
    }
}
