package me.wener.bbvm.asm;

import io.netty.buffer.ByteBuf;
import me.wener.bbvm.vm.Instruction;

/**
 * @author wener
 * @since 15/12/11
 */
public class Inst extends AbstractAssembly implements Assembly {
    private Instruction instruction;
    private int line;

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

    @Override
    public void write(ByteBuf buf) {
        instruction.write(buf);
    }

    @Override
    public int getLine() {
        return line;
    }

    public Inst setLine(int line) {
        this.line = line;
        return this;
    }

    @Override
    public int length() {
        return instruction.getOpcode().length();
    }
}
