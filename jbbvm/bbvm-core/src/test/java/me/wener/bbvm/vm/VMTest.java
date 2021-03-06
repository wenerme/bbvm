package me.wener.bbvm.vm;

import me.wener.bbvm.vm.invoke.BasicInvoke;
import org.junit.Test;

/**
 * @author wener
 * @since 15/12/11
 */
public class VMTest {

    @Test
    public void testLoad() throws Exception {
        VM vm = VMConfig.newBuilder().invokeWithSystemOutput().invokeWith(BasicInvoke.class).create();
        vm.setMemory(Memory.load(Dumps.simpleInst().skipBytes(16))).run();
    }
}
