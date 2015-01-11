package me.wener.bbvm.impl;

import lombok.Data;
import lombok.experimental.Accessors;
import me.wener.bbvm.def.DataType;
import me.wener.bbvm.def.InstructionType;

/**
 * <pre>
 *    指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数 + 第二个操作数
 * 0x 0       0         0           0        0000         0000
 * </pre>
 */
@Data
@Accessors(chain = true)
public class Instruct
{
    private InstructionType instruction;
    private Operand op1;
    private Operand op2;
    private DataType dataType;
    private int specialByte;
    private int addressingType;
    private int firstByte;

}
