package me.wener.bbvm.neo.processor;

import io.netty.util.DefaultAttributeMap;
import lombok.Getter;
import lombok.experimental.Accessors;
import me.wener.bbvm.neo.BBVMContext;
import me.wener.bbvm.neo.inst.Inst;
import me.wener.bbvm.neo.inst.def.InstructionType;

@Getter
@Accessors(chain = true)
public class DefaultProcessContext extends DefaultAttributeMap implements ProcessContext
{
    private Inst inst;

    public DefaultProcessContext()
    {

    }

    @Override
    public ProcessContext instruction(Inst inst)
    {
        this.inst = inst;
        return this;
    }

    @Override
    public BBVMContext vm()
    {
        return null;
    }

    @Override
    public Inst instruction()
    {
        return inst;
    }

    @Override
    public InstructionType instructionType()
    {
        return null;
    }

}
