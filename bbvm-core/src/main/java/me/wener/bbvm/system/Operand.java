package me.wener.bbvm.system;

import lombok.Data;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.utils.val.IntegerHolder;
import me.wener.bbvm.utils.val.IsInteger;

@Data
@Accessors(chain = true, fluent = true)
@Slf4j
public class Operand
{
    private CPU cpu;
    private Integer value;
    private IsInteger indirect;
    private AddressingMode addressingMode;

    public Operand(CPU cpu)
    {
        this.cpu = cpu;
    }

    public Operand()
    {
    }

    public Operand asInteger(int value)
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
                cpu.memory().writeInt(asInteger(), value);
                break;
            default:
            case IMMEDIATE:
                throw new UnsupportedOperationException();
        }
        return this;
    }

    public int asInteger()
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
        return cpu.memory().readString(asInteger());
    }
}
