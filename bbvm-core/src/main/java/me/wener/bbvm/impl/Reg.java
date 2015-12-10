package me.wener.bbvm.impl;

import lombok.Data;
import lombok.EqualsAndHashCode;
import me.wener.bbvm.util.val.IntHolder;
import me.wener.bbvm.util.val.impl.SimpleValueHolder;

@EqualsAndHashCode(callSuper = true)
@Data
public class Reg extends SimpleValueHolder<Integer> implements IntHolder
{
    private final String name;

    public Reg(String name) {this.name = name;}

    @Override
    public String toString()
    {
        return name;
    }
}
