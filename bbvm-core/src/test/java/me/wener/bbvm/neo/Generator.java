package me.wener.bbvm.neo;

import me.wener.bbvm.def.InstructionType;
import org.junit.Test;

public class Generator
{
    @Test
    public void genInstCreate()
    {
        for (InstructionType type : InstructionType.values())
        {
//            System.out.printf("case Flags.%s:return new %1$s();", type);
            System.out.printf("map.put(%s.class, Flags.%1$s);", type);
        }
    }
}
