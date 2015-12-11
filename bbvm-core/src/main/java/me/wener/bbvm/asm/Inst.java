package me.wener.bbvm.asm;

import me.wener.bbvm.vm.Instruction;

/**
 * @author wener
 * @since 15/12/11
 */
public class Inst implements Assembly {
    private me.wener.bbvm.vm.Instruction instruction;

    public Inst(Instruction instruction) {
        this.instruction = instruction;
    }

    @Override
    public Type getType() {
        return Type.INST;
    }

    public me.wener.bbvm.vm.Instruction getInstruction() {
        return instruction;
    }

    @Override
    public String toAssembly() {
        return instruction.toAssembly();
    }

    @Override
    public String toString() {
        return instruction.toString();
    }
}
