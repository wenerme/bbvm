package me.wener.bbvm.system;

public interface OpStatus
{
    DataType dataType();

    CompareType compareType();

    CalculateType calculateType();

    Opcode opcode();

    Operand a();

    Operand b();

    String toAssembly();

    byte[] toBinary();
}
