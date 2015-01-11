package me.wener.bbvm.neo;

import me.wener.bbvm.utils.val.IntegerHolder;

public interface Operand extends IntegerHolder
{
    int addressingMode();

    Operand addressingMode(int mode);
}
