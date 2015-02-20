package me.wener.bbvm.system.api;

import me.wener.bbvm.utils.val.IntegerHolder;
import me.wener.bbvm.utils.val.IsInteger;

public interface Operand extends IntegerHolder
{
    String asString();

    Integer value();

    Operand value(Integer v);

    IsInteger indirect();

    Operand indirect(IsInteger v);

    AddressingMode addressingMode();

    Operand addressingMode(AddressingMode mode);
}
