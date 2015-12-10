package me.wener.bbvm.impl;

import lombok.AccessLevel;
import lombok.Data;
import lombok.Getter;
import me.wener.bbvm.api.BBVm;
import me.wener.bbvm.util.Bins;
import me.wener.bbvm.util.val.IntEnums;
import me.wener.bbvm.util.val.Values;
import me.wener.bbvm.vm.DataType;
import me.wener.bbvm.vm.Opcode;

@Data
public class InstructionContext
{
    @Getter(AccessLevel.NONE)
    private final byte[] memory;
    @Getter(AccessLevel.NONE)
    private final BBVm vm;
    private Opcode instruction;
    private Operand op1;
    private Operand op2;
    private DataType dataType;
    private int specialByte;
    private int addressingType;
    private int firstByte;

    InstructionContext(BBVm vm)
    {
        this.memory = vm.getMemory();

        this.vm = vm;
    }

    void read(int pc)
    {
       /*
            指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数 + 第二个操作数
         0x 0       0         0           0        0000         0000
        */
        firstByte = Bins.uint16b(memory, pc);
        int opcode = firstByte >> 12;// 指令码
        instruction = IntEnums.fromInt(Opcode.class, opcode);
        Integer length = instruction.length();

        if (length == 1 || length == 5)
        {
            dataType = null;
            specialByte = 0;
            addressingType = (firstByte & 0x0F00) >> 8;
        } else if (length > 1)
        {
            specialByte = (firstByte & 0x00F0) >> 4;
            addressingType = firstByte & 0x000F;
            dataType = Values.fromValue(DataType.class, (firstByte & 0x0F00) >> 8);
        }

        op1 = op2 = Operand.invalid();
        int op1t;
        int op2t;
        if (length >= 10)
        {
            // 双操作数
            op1t = addressingType / 4;
            op2t = addressingType % 4;
        } else
        {
            // 单操作数
            op1t = addressingType % 4;
            op2t = 0;
        }

        switch (length)
        {
            case 5:
                op1 = operand(op1t, Bins.int32l(memory, pc + 1));
                break;
            case 6:
                op1 = operand(op1t, Bins.int32l(memory, pc + 2));
                break;
            case 10:
                op1 = operand(op1t, Bins.int32l(memory, pc + 2));
                op2 = operand(op2t, Bins.int32l(memory, pc + 6));
                break;
        }

    }

    /**
     * <pre>
     * rx	| 0x0 | 寄存器寻址
     * [rx]	| 0x1 | 寄存器间接寻址
     * n	| 0x2 | 立即数寻址
     * [n]	| 0x3 | 间接寻址
     * </pre>
     */
    private Operand operand(int type, int op)
    {
        switch (type)
        {
            case 0:
                return Operand.holder(vm.getRegister(op));
            case 1:
                return Operand.indirect(vm.getRegister(op), memory);
            case 2:
                return Operand.value(op);
            case 3:
                return Operand.address(op, memory);
            default:
                throw new UnsupportedOperationException(String.format("未知的寻址类型: %s 操作数为: %s", type, op));
        }
    }

}
