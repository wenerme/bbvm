package me.wener.bbvm.system;

import me.wener.bbvm.utils.val.IntegerHolder;

public interface Operand extends IntegerHolder
{
    Integer value();

    Operand value(Integer v);

    RegisterType asRegisterType();

    Operand value(RegisterType v);

    AddressingMode addressingMode();

    Operand addressingMode(AddressingMode mode);

    float asFloat();

    Operand asFloat(float v);

    /**
     * 将该操作数的值作为资源句柄
     *
     * @param res 资源名
     * @return 资源对象
     */
    Resource asResource(String res);

    Operand asString(String v);

    /**
     * 将该操作数的值作为字符串地址或字符串句柄句柄
     */
    String asString();

    /**
     * 将该操作数以字符串的形式呈现
     */
    String toAssembly();

    /**
     * 获取该操作数的实际值,与 value 不同,该值会取实际值
     */
    @Override
    Integer get();

    /**
     * 设置该操作数所指向的值,例如:寄存器值或相应的地址值
     */
    @Override
    void set(Integer v);
}
