package me.wener.bbvm.vm;

import me.wener.bbvm.util.Bins;
import me.wener.bbvm.util.val.IntEnums;

import static com.google.common.base.Preconditions.checkState;

/**
 * @author wener
 * @since 15/12/10
 */
public class VM {
    static {
        IntEnums.cache(AddressingMode.class, CalculateType.class, CompareType.class, DataType.class, Opcode.class, RegisterType.class);
    }

    final Register r0 = new Register(RegisterType.R0, this);
    final Register r1 = new Register(RegisterType.R1, this);
    final Register r2 = new Register(RegisterType.R2, this);
    final Register r3 = new Register(RegisterType.R3, this);
    final Register rs = new Register(RegisterType.RS, this);
    final Register rf = new Register(RegisterType.RF, this);
    final Register rb = new Register(RegisterType.RB, this);
    final Register rp = new Register(RegisterType.RP, this);
    Memory memory = new Memory().setVm(this);
    private boolean exit = false;
    private boolean jumped;

    public VM() {
    }

    private static Number cal(CalculateType calculateType, DataType dataType, Operand a, Operand b) {
        float va;
        float vb;
        if (dataType == DataType.FLOAT) {
            va = Bins.float32(a.get());
            vb = Bins.float32(b.get());
        } else {
            va = a.get();
            vb = b.get();
        }
        Number vc;
        switch (calculateType) {
            case ADD:
                vc = va + vb;
                break;
            case SUB:
                vc = va - vb;
                break;
            case MUL:
                vc = va * vb;
                break;
            case DIV:
                vc = va / vb;
                break;
            case MOD:
                vc = va % vb;
                break;
            default:
                throw new UnsupportedOperationException();
        }
        switch (dataType) {
            case DWORD:
            case INT:
                vc = vc.intValue();
                break;
            case WORD:
                vc = vc.shortValue();
                break;
            case BYTE:
                vc = vc.byteValue();
                break;
        }
        return vc;
    }

    public Memory getMemory() {
        return memory;
    }

    public VM reset() {
        r0.setValue(0);
        r1.setValue(0);
        r2.setValue(0);
        r3.setValue(0);
        rs.setValue(0);
        rf.setValue(0);
        rb.setValue(0);
        rp.setValue(0);
        exit = false;
        return this;
    }

    public void run(Instruction inst) {
        checkState(!exit, "Exited");
        run(inst, inst.opcode, inst.a, inst.b);
    }

    private void run(Instruction inst, Opcode opcode, Operand a, Operand b) {
        switch (opcode) {
            case NOP:
                break;
            case LD:
                // TODO Data type overflow check
                a.set(b.get());
                break;
            case PUSH:
                push(a.get());
                break;
            case POP:
                a.set(pop());
                break;
            case IN:
                in(a.get(), b.get());
                break;
            case OUT:
                out(a.get(), b.get());
                break;
            case JMP:
                jmp(a.get());
                break;
            case JPC:
                if (IntEnums.fromInt(CompareType.class, rf.intValue()).isMatch(inst.compareType)) {
                    jmp(a.get());
                }
                break;
            case CALL:
                push(rp.intValue() + inst.getOpcode().length());
                jmp(a.get());
                break;
            case RET:
                ret();
                break;
            case CMP: {
                float vc = cal(CalculateType.SUB, inst.getDataType(), a, b).intValue();
                if (vc > 0)
                    rf.setValue(CompareType.A);
                else if (vc < 0)
                    rf.setValue(CompareType.B);
                else
                    rf.setValue(CompareType.Z);
            }
            break;
            case CAL: {
                Number vc = cal(CalculateType.SUB, inst.getDataType(), a, b).intValue();
                if (inst.getDataType() == DataType.FLOAT) {
                    a.set(vc.floatValue());
                }
            }
            break;
            case EXIT:
                exit = false;
                break;
        }
    }

    public void jmp(int i) {
        rp.setValue(i);
    }

    public void ret() {
        rp.setValue(pop());
    }

    public VM push(int v) {
        memory.push(v);
        return this;
    }

    public int pop() {
        return memory.pop();
    }

    public void in(int a, int b) {

    }

    public void out(int a, int b) {

    }

    public Register getRegister(RegisterType type) {
        switch (type) {
            case RP:
                return rp;
            case RF:
                return rf;
            case RS:
                return rs;
            case RB:
                return rb;
            case R0:
                return r0;
            case R1:
                return r1;
            case R2:
                return r2;
            case R3:
                return r3;
        }
        throw new UnsupportedOperationException();
    }
}
