package me.wener.bbvm.neo;

import lombok.Data;
import lombok.EqualsAndHashCode;
import me.wener.bbvm.utils.val.IntegerHolder;
import me.wener.bbvm.utils.val.impl.SimpleValueHolder;

@EqualsAndHashCode(callSuper = true)
@Data
public class Register extends SimpleValueHolder<Integer> implements IntegerHolder
{
    private final String name;

    public Register(String name) {this.name = name;}

    @Override
    public String toString()
    {
        return name;
    }
}
