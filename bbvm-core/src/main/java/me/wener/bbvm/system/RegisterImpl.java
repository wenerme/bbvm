package me.wener.bbvm.system;

import java.io.Serializable;
import lombok.Getter;
import lombok.experimental.Accessors;
import me.wener.bbvm.system.api.Register;
import me.wener.bbvm.system.api.RegisterType;

@Accessors(chain = true, fluent = true)
public class RegisterImpl implements Register, Serializable
{
    @Getter
    private String name;
    @Getter
    private RegisterType type;
    private Integer value;

    public RegisterImpl()
    {
    }

    public RegisterImpl(String name)
    {
        this.name = name;
        type = RegisterType.valueOf(name);
    }

    public RegisterImpl(RegisterType type)
    {
        this.type = type;
        name = type.toString();
    }

    @Override
    public Integer get()
    {
        return value;
    }

    @Override
    public void set(Integer v)
    {
        value = v;
    }
}
