package me.wener.bbvm.system.api;


/**
 * 操作状态信息,可获取汇编内容,也可以直接获取编译后内容,代表一句汇编指令<br>
 * 因此在编译的时候可以直接通过解析出该状态来生成编译
 */
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
