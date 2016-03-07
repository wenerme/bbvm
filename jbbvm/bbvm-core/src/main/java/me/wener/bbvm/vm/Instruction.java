package me.wener.bbvm.vm;

import com.google.common.base.MoreObjects;
import io.netty.buffer.ByteBuf;
import me.wener.bbvm.util.IntEnums;

import static me.wener.bbvm.util.IntEnums.fromInt;


/**
 * A uniformed instruction target
 *
 * @author wener
 * @since 15/12/10
 */
@SuppressWarnings("unused")
public class Instruction {
    static {
        IntEnums.cache(AddressingMode.class, CalculateType.class, CompareType.class, DataType.class, Opcode.class, RegisterType.class);
    }

    Opcode opcode;
    Operand a = new Operand();
    Operand b = new Operand();
    CalculateType calculateType;
    CompareType compareType;
    DataType dataType;
    transient VM vm;
    /**
     * The address of this instruction in memory
     */
    int address = -1;

    public VM getVm() {
        return vm;
    }

    public Instruction setVm(VM vm) {
        this.vm = vm;
        if (a != null) {
            a.setVm(vm);
        }
        if (b != null) {
            b.setVm(vm);
        }
        return this;
    }

    public int getAddress() {
        return address;
    }

    public Instruction setAddress(int address) {
        this.address = address;
        return this;
    }

    public Opcode getOpcode() {
        return opcode;
    }

    public Instruction setOpcode(Opcode opcode) {
        this.opcode = opcode;
        return this;
    }

    public Operand getA() {
        return a;
    }

    public Instruction setA(Operand a) {
        this.a = a;
        return this;
    }

    public boolean hasA() {
        switch (opcode) {
            case RET:
            case NOP:
            case EXIT:
                return false;
        }
        return true;
    }

    public boolean hasB() {
        switch (opcode) {
            case LD:
            case IN:
            case OUT:
            case CAL:
            case CMP:
                return true;
        }
        return false;
    }

    public Operand getB() {
        return b;
    }

    public Instruction setB(Operand b) {
        this.b = b;
        return this;
    }

    public CalculateType getCalculateType() {
        return calculateType;
    }

    public Instruction setCalculateType(CalculateType calculateType) {
        this.calculateType = calculateType;
        return this;
    }

    public CompareType getCompareType() {
        return compareType;
    }

    public Instruction setCompareType(CompareType compareType) {
        this.compareType = compareType;
        return this;
    }

    public DataType getDataType() {
        return dataType;
    }

    public Instruction setDataType(DataType dataType) {
        this.dataType = dataType;
        return this;
    }

    @Override
    public String toString() {
        return MoreObjects.toStringHelper(this)
            .add("opcode", opcode)
            .add("a", a)
            .add("b", b)
            .add("calculateType", calculateType)
            .add("compareType", compareType)
            .add("dataType", dataType)
            .toString();
    }

    /**
     * Read instruction at address, will not change readerIndex of bu
     */
    public Instruction read(ByteBuf buf, int address) {
        /*
   指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数 + 第二个操作数
0x 0       0         0           0        00000000     00000000

无操作数 1byte
   指令码 + 无用
0x 0       0
一个操作数 5byte
   指令码 + 寻址方式 + 第一个操作数
0x 0       0        00000000
两个操作数 10byte
   指令码 + 数据类型 + 保留字节 + 寻址方式 + 第一个操作数 + 第二个操作数
0x 0       0         0        0        00000000     00000000
JPC指令 6byte
   指令码 + 比较操作 + 保留字节 + 寻址方式 + 第一个操作数
0x 0       0         0        0        00000000
        */
        int offset = address;
        this.address = address;
        short first = buf.getUnsignedByte(offset++);
        opcode = fromInt(Opcode.class, first >> 4);
        switch (opcode) {
            case RET:
            case NOP:
            case EXIT: {
                // 无操作数
            }
            break;
            case POP:
            case PUSH:
            case CALL:
            case JMP: {
                a.setAddressingMode(fromInt(AddressingMode.class, first & 0xf));
                // 一个操作数
                a.setInternal(buf.getInt(offset));
            }
            break;
            case LD:
            case IN:
            case OUT:
            case CAL:
            case CMP: {
                // 两个操作数
                dataType = fromInt(DataType.class, first & 0xf);
                short second = buf.getUnsignedByte(offset++);
                int special = second >> 4;
                int addressingMode = second & 0xf;

                a.addressingMode = fromInt(AddressingMode.class, addressingMode / 4);
                b.addressingMode = fromInt(AddressingMode.class, addressingMode % 4);
                a.setInternal(buf.getInt(offset));

                offset += 4;
                b.setInternal(buf.getInt(offset));

                if (opcode == Opcode.CAL) {
                    calculateType = fromInt(CalculateType.class, special);
                }
            }
            break;

            case JPC: {
                short second = buf.getUnsignedByte(offset++);
                int addressingMode = second & 0xf;
                // JPC A R1
                // 数据类型为比较操作
                compareType = fromInt(CompareType.class, first & 0xf);
                a.addressingMode = fromInt(AddressingMode.class, addressingMode);
                a.setInternal(buf.getInt(offset));
            }
            break;
            default:
                throw new UnsupportedOperationException();
        }
        return this;
    }

    public Instruction write(ByteBuf buf) {
        switch (opcode) {
            case RET:
            case NOP:
            case EXIT: {
                // 无操作数
                buf.writeByte(opcode.asInt() << 4);
            }
            break;
            case POP:
            case PUSH:
            case CALL:
            case JMP: {
                // 一个操作数
                buf.writeByte(opcode.asInt() << 4 | a.addressingMode.asInt());
                buf.writeInt(a.getInterval());
            }
            break;
            case LD:
            case IN:
            case OUT:
            case CAL:
            case CMP: {
                // 两个操作数
                if (dataType == null) {
                    buf.writeByte(opcode.asInt() << 4);
                } else {
                    buf.writeByte(opcode.asInt() << 4 | dataType.asInt());
                }

                if (opcode == Opcode.CAL) {
                    buf.writeByte(calculateType.asInt() << 4 | (a.addressingMode.asInt() << 2 | b.addressingMode.asInt()));
                } else {
                    buf.writeByte(a.addressingMode.asInt() << 2 | b.addressingMode.asInt());
                }
                buf.writeInt(a.getInterval()).writeInt(b.getInterval());
            }
            break;

            case JPC: {
                buf.writeByte(opcode.asInt() << 4 | compareType.asInt());
                buf.writeByte(a.addressingMode.asInt());
                buf.writeInt(a.getInterval());
            }
            break;
            default:
                throw new UnsupportedOperationException();
        }
        return this;
    }

    public String toAssembly() {
        StringBuilder sb = new StringBuilder();
        sb.append(opcode);
        switch (opcode) {
            // 没有操作数
            case NOP:
            case RET:
            case EXIT:
                break;

            case PUSH:
            case POP:
            case JMP:
            case CALL:
                // 一个操作数
                sb.append(' ').append(a.toAssembly());
                break;
            case IN:
            case OUT:
                // 标准的两个操作数
                sb.append(' ').append(a.toAssembly()).append(", ").append(b.toAssembly());
                break;
            case JPC:
                sb.append(' ').append(compareType).append(' ').append(a.toAssembly());
                break;
            case CMP:
                sb.append(' ').append(compareType)
                    .append(' ').append(a.toAssembly())
                    .append(", ").append(b.toAssembly())
                ;
                break;
            case LD:
                sb.append(' ').append(dataType)
                    .append(' ').append(a.toAssembly())
                    .append(", ").append(b.toAssembly())
                ;
                break;
            case CAL:
                sb.append(' ').append(dataType)
                    .append(' ').append(calculateType)
                    .append(' ').append(a.toAssembly())
                    .append(", ").append(b.toAssembly())
                ;
                break;
        }
        return sb.toString();
    }


    public Instruction reset() {
        opcode = null;
        a.reset();
        b.reset();
        calculateType = null;
        compareType = null;
        dataType = null;
        return this;
    }
}
