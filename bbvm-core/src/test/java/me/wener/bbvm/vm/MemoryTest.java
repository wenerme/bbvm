package me.wener.bbvm.vm;

import org.junit.Test;

import static junit.framework.TestCase.assertEquals;

/**
 * @author wener
 * @since 15/12/10
 */
public class MemoryTest {

    @Test
    public void testOp() throws Exception {
        Memory mem = new Memory(32, 16);
        mem.rs = new Register(RegisterType.RS, null);
        mem.rb = new Register(RegisterType.RB, null);
        mem.reset();
        mem.push(1);
        assertEquals(1, mem.pop());
        for (int i = 0; i < 4; i++) {
            mem.push(i);
        }
        System.out.println(mem);
    }
}
