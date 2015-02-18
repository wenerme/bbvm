package me.wener.bbvm.system;

import lombok.Getter;
import lombok.experimental.Accessors;
import me.wener.bbvm.utils.val.IntegerHolder;

public class Register implements IntegerHolder
{
    @Getter
    @Accessors(chain = true)
    private final String name;
    private Integer value;

    public Register(String name) {this.name = name;}

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
