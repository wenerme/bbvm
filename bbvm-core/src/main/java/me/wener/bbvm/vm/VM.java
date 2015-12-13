package me.wener.bbvm.vm;

import io.netty.buffer.ByteBuf;
import me.wener.bbvm.util.val.IntEnums;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Iterator;

import static com.google.common.base.Preconditions.checkState;

/**
 * @author wener
 * @since 15/12/10
 */
public class VM {
    private final static Logger log = LoggerFactory.getLogger(VM.class);
    //    Injector injector;
    final Register r0 = new Register(RegisterType.R0, this);
    final Register r1 = new Register(RegisterType.R1, this);
    final Register r2 = new Register(RegisterType.R2, this);
    final Register r3 = new Register(RegisterType.R3, this);
    final Register rs = new Register(RegisterType.RS, this);
    final Register rf = new Register(RegisterType.RF, this);
    final Register rb = new Register(RegisterType.RB, this);
    final Register rp = new Register(RegisterType.RP, this);
    Memory memory;
    SymbolTable symbolTable;
    private boolean exit = false;

    public VM() {
    }

    private static Number cal(CalculateType calculateType, DataType dataType, Operand a, Operand b) {
        float va;
        float vb;
        if (dataType == DataType.FLOAT) {
            va = a.getFloat();
            vb = b.getFloat();
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

    public SymbolTable getSymbolTable() {
        return symbolTable;
    }

    public VM setSymbolTable(SymbolTable symbolTable) {
        this.symbolTable = symbolTable;
        return this;
    }

    boolean hasRemaining() {
        return rp.intValue() < memory.getMemorySize();
    }

    public void run() {
        Instruction instruction = new Instruction().setVm(this);
        ByteBuf buf = this.memory.getByteBuf();
        int last;
        while (hasRemaining()) {
            last = rp.intValue();
            instruction.reset().read(buf, last);
            run(instruction);
            if (exit) {
                return;
            }
            if (rp.intValue() == last) {
                rp.add(instruction.getOpcode().length());
            }
        }
    }

    public Iterable<Instruction> instructions(final Instruction instruction, final int position) {
        return new InstructionIterable(position, instruction);
    }

    public Memory getMemory() {
        return memory;
    }

    public VM setMemory(Memory memory) {
        this.memory = memory.setVm(this);
        return this;
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
        memory.reset();
        return this;
    }

    public void run(Instruction inst) {
        checkState(!exit, "Exited");
        log.debug("{}", inst);
        log.debug("{} ' A={} B={} {}",
                inst.toAssembly(),
                inst.hasA() ? inst.getA().get() : "NaN",
                inst.hasB() ? inst.getB().get() : "NaN",
                debugAsm());
        run(inst, inst.opcode, inst.a, inst.b);
    }

    String debugAsm() {
        return String.format("RP=%s RF=%s RS=%s RB=%s R0=%s R1=%s R2=%s R3=%s"
                , rp.intValue(), rf.intValue(), rs.intValue(), rb.intValue(), r0.intValue(), r1.intValue(), r2.intValue(), r3.intValue());
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
                exit = true;
                break;
        }
    }

    public boolean isExit() {
        return exit;
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

    public Symbol getSymbol(int address) {
        return symbolTable != null ? symbolTable.getSymbol(address) : null;
    }

    private class InstructionIterable implements Iterable<Instruction> {
        private final int position;
        private final Instruction instruction;

        public InstructionIterable(int position, Instruction instruction) {
            this.position = position;
            this.instruction = instruction;
        }

        @Override
        public Iterator<Instruction> iterator() {
            return new Iterator<Instruction>() {
                int pos = position;

                @Override
                public boolean hasNext() {
                    return pos < memory.getMemorySize();
                }

                @Override
                public Instruction next() {
                    instruction.reset().read(memory.getByteBuf(), pos);
                    pos += instruction.opcode.length();
                    return instruction;
                }
            };
        }
    }
}
