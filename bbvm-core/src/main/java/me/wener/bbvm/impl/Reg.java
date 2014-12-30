package me.wener.bbvm.impl;

import lombok.Data;
import lombok.EqualsAndHashCode;
import me.wener.bbvm.utils.val.IntegerHolder;
import me.wener.bbvm.utils.val.SimpleValueHolder;

@EqualsAndHashCode(callSuper = true)
@Data
public class Reg extends SimpleValueHolder<Integer> implements IntegerHolder
{
    private final String name;

    public Reg(String name) {this.name = name;}
}
