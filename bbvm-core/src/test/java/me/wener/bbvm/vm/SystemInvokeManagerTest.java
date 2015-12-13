package me.wener.bbvm.vm;

import me.wener.bbvm.vm.invoke.PrintStreamOutput;
import org.junit.Test;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

import static org.junit.Assert.assertEquals;

/**
 * @author wener
 * @since 15/12/13
 */
public class SystemInvokeManagerTest {

    @Test
    public void testRegister() throws Exception {
        SystemInvokeManager manager = new SystemInvokeManagerImpl();
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        manager.register(new PrintStreamOutput(new PrintStream(out)));
        Instruction inst = new Instruction();
        inst.a.setAddressingMode(AddressingMode.IMMEDIATE).setInternal(0);
        inst.b.setAddressingMode(AddressingMode.IMMEDIATE).setInternal(10);
        inst.opcode = Opcode.OUT;
        manager.invoke(inst);
        inst.a.setInternal(5);
        inst.b.setInternal(Float.floatToRawIntBits(0.333f));
        manager.invoke(inst);

        assertEquals("10\n0.333000", out.toString());
    }

    @Test
    public void testInvoke() throws Exception {

    }
}
