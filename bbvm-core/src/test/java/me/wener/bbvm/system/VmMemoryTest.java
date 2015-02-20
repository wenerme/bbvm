package me.wener.bbvm.system;

import static org.junit.Assert.assertEquals;

import org.junit.Before;
import org.junit.Test;

public class VmMemoryTest
{
    VmMemory mem = new VmMemory();

    @Before
    public void before()
    {
        mem.reset();
    }

    @Test
    public void basicOP()
    {
        mem.writeInt(2, 123);
        assertEquals(123, mem.readInt(2));
    }


}
