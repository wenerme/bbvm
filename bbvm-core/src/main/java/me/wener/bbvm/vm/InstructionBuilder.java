package me.wener.bbvm.vm;

import com.google.common.collect.Lists;

import java.util.LinkedList;

/**
 * @author wener
 * @since 15/12/10
 */
public class InstructionBuilder {
    VM vm;
    Instruction last;
    LinkedList<Instruction> instructions = Lists.newLinkedList();

    public VM vm() {
        return vm;
    }

    public InstructionBuilder vm(VM vm) {
        this.vm = vm;
        for (Instruction i : instructions) {
            i.setVm(vm);
        }
        return this;
    }

    public Instruction last() {
        return last;
    }

    public LinkedList<Instruction> instructions() {
        return instructions;
    }

    public InstructionBuilder instructions(LinkedList<Instruction> instructions) {
        this.instructions = instructions;
        return this;
    }

    public InstructionBuilder jpc(Operand operand, CompareType compareType) {
        create(Opcode.JPC, operand).setCompareType(compareType);
        return this;
    }

    public InstructionBuilder cal(CalculateType calculateType, DataType dataType, Operand a, Operand b) {
        create(Opcode.CAL, a, b).setCalculateType(calculateType).setDataType(dataType);
        return this;
    }

    public InstructionBuilder jmp(Operand operand) {
        create(Opcode.JMP, operand);
        return this;
    }

    public InstructionBuilder push(Operand operand) {
        create(Opcode.PUSH, operand);
        return this;
    }

    public InstructionBuilder in(Operand a, Operand b) {
        create(Opcode.IN, a, b);
        return this;
    }

    public InstructionBuilder out(Operand a, Operand b) {
        create(Opcode.OUT, a, b);
        return this;
    }

    public InstructionBuilder ld(DataType dataType, Operand a, Operand b) {
        create(Opcode.LD, a, b).setDataType(dataType);
        return this;
    }

    private Instruction create(Opcode opcode, Operand a, Operand b) {
        return last = new Instruction().setOpcode(opcode).setA(a).setB(b).setVm(vm);
    }

    private Instruction create(Opcode opcode, Operand a) {
        return last = new Instruction().setOpcode(opcode).setA(a).setVm(vm);
    }

    private Instruction create(Opcode opcode) {
        last = new Instruction().setVm(vm).setOpcode(opcode);
        instructions.add(last);
        return last;
    }

    public InstructionBuilder pop(Operand operand) {
        create(Opcode.POP, operand);
        return this;
    }

    public InstructionBuilder nop() {
        create(Opcode.NOP);
        return this;
    }

    public InstructionBuilder exit() {
        create(Opcode.EXIT);
        return this;
    }

    public InstructionBuilder ret() {
        create(Opcode.RET);
        return this;
    }
}
