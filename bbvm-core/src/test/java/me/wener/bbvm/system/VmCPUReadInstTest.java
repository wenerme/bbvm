package me.wener.bbvm.system;

import io.netty.buffer.ByteBuf;
import java.io.OutputStreamWriter;
import me.wener.bbvm.neo.TestUtil;
import me.wener.bbvm.system.internal.VmCPU;
import me.wener.bbvm.utils.Dumper;
import org.junit.Test;

public class VmCPUReadInstTest extends TestUtil
{

    @Test
    public void testWithJMPCAL()
    {
        ByteBuf mem = Dumps.jmpAndCal();
        mem.skipBytes(16);
        System.out.println(Dumper.hexDumpReadable(mem));
        VmCPU cpu = new VmCPU().ignoreProcess(true);
        cpu.load(readableToBytes(mem));
        OutputStreamWriter writer = new OutputStreamWriter(System.out);
        while (cpu.step())
        {
//            System.out.println(cpu.toAssembly());
        }
    }

    @Test
    public void testSimpleInst()
    {
        ByteBuf mem = Dumps.simpleInst();
        mem.skipBytes(16);
        System.out.println(Dumper.hexDumpReadable(mem));
        VmCPU cpu = new VmCPU();
        cpu.load(readableToBytes(mem));
        OutputStreamWriter writer = new OutputStreamWriter(System.out);
        System.out.println(Dumps.simpleInstAsm());
        while (cpu.step())
        {
//            System.out.println(cpu.toAssembly());
        }
    }
}
