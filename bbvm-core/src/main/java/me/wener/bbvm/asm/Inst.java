package me.wener.bbvm.asm;

import me.wener.bbvm.vm.Instruction;

/**
 * @author wener
 * @since 15/12/11
 */
public class Inst extends AbstractAssembly implements Assembly {
    private Instruction instruction;

    public Inst(Instruction instruction) {
        this.instruction = instruction;
    }

    @Override
    public Type getType() {
        return Type.INST;
    }

    public Instruction getInstruction() {
        return instruction;
    }

    @Override
    public String toAssembly() {
        return instruction.toAssembly() + commentAssembly();
    }

    @Override
    public String toString() {
        return instruction.toString();
    }
}
