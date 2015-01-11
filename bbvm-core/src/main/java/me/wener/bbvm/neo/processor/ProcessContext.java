package me.wener.bbvm.neo.processor;

import io.netty.util.AttributeMap;
import me.wener.bbvm.neo.BBVMContext;
import me.wener.bbvm.neo.define.InstructionType;
import me.wener.bbvm.neo.inst.Inst;

public interface ProcessContext extends AttributeMap
{
    ProcessContext instruction(Inst inst);

    BBVMContext vm();

    Inst instruction();

    InstructionType instructionType();
}
