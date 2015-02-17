package me.wener.bbvm.neo.processor;

import com.google.common.collect.HashBasedTable;
import com.google.common.collect.Table;
import me.wener.bbvm.neo.BBVMContext;
import me.wener.bbvm.neo.Operand;
import me.wener.bbvm.neo.inst.IN;
import me.wener.bbvm.neo.inst.Inst;
import me.wener.bbvm.neo.inst.OUT;
import me.wener.bbvm.neo.inst.def.InstructionType;

public abstract class ComposedProcessor extends BasicProcessor
{
    private final Table<Integer, Integer, Processor> processor = HashBasedTable.create();

    public ComposedProcessor register(Processor p)
    {
        return register(null, null, p);
    }

    public ComposedProcessor register(Integer first, Integer second, Processor p)
    {
        processor.put(first, second, p);
        tryInitialize(vm, p);
        return this;
    }

    public ComposedProcessor register(Integer first, Processor p)
    {
        return register(first, null, p);
    }

    public Processor get(Integer first, Integer second)
    {
        return processor.get(first, second);
    }

    @Override
    public void initialize(BBVMContext vm)
    {
        super.initialize(vm);
        for (Processor p : processor.values())
        {
            tryInitialize(vm, p);
        }
    }

    private void tryInitialize(BBVMContext vm, Processor p)
    {
        if (vm != null && p instanceof VMContextAware)
        {
            ((VMContextAware) p).initialize(vm);
        }
    }

    public static class InProcessor extends ComposedProcessor
    {

        @Override
        public ProcessResult process(ProcessContext ctx, InstructionType instructionType, Inst instruction, Operand a, Operand b, int dataType)
        {
            if (instructionType == InstructionType.IN)
            {
                IN in = (IN) instruction;
                Processor p = get(in.a.get(), in.b.get());
                if (p != null)
                {
                    return p.apply(ctx);
                }
            }
            return Results.keepGoing();
        }
    }

    public static class OutProcessor extends ComposedProcessor
    {

        @Override
        public ProcessResult process(ProcessContext ctx, InstructionType instructionType, Inst instruction, Operand a, Operand b, int dataType)
        {
            if (instructionType == InstructionType.OUT)
            {
                OUT in = (OUT) instruction;
                Processor p = get(in.a.get(), in.b.get());
                if (p != null)
                {
                    return p.apply(ctx);
                }
            }
            return Results.keepGoing();
        }
    }
}
