package me.wener.bbvm.neo.processor;

import me.wener.bbvm.neo.Operand;
import me.wener.bbvm.neo.define.InstructionType;
import me.wener.bbvm.neo.inst.Inst;
import me.wener.bbvm.neo.inst.OneOperandInst;
import me.wener.bbvm.neo.inst.TowOperandInst;

public abstract class BasicProcessor extends VMContext implements Processor
{
    @Override
    public ProcessResult apply(ProcessContext ctx)
    {
        Inst inst = ctx.instruction();
        Operand a = null;
        Operand b = null;
        int dataType = -1;
        if (inst instanceof OneOperandInst)
        {
            a = ((OneOperandInst) inst).a;
        } else if (inst instanceof TowOperandInst)
        {
            a = ((TowOperandInst) inst).a;
            b = ((TowOperandInst) inst).b;
            dataType = ((TowOperandInst) inst).dataType;
        }
        return process(ctx, ctx.instructionType(), inst, a, b, dataType);
    }

    public abstract ProcessResult process(ProcessContext ctx,
                                          InstructionType instructionType,
                                          Inst instruction,
                                          Operand a,
                                          Operand b,
                                          int dataType);


}
