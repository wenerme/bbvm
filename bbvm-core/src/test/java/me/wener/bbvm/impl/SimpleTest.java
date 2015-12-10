package me.wener.bbvm.impl;

import me.wener.bbvm.impl.spi.DeviceProvider;
import me.wener.bbvm.util.Bins;
import org.junit.Test;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Arrays;

public class SimpleTest
{
    public static void main(String[] args) throws IOException
    {
        System.out.println("当前目录: " + System.getProperty("user.dir"));
        byte[] bytes = Files.readAllBytes(Paths.get("doc/testsuit/BB/Sim/BBasic/test.bin"));
        BBVmImpl vm = new BBVmImpl(DeviceProvider.getProvider().createDevice(240, 320));
        byte[] mem = Arrays.copyOfRange(bytes, 16, bytes.length);
        vm.load(mem);
        vm.reset();

        vm.start();
    }

    @Test
    public void test()
    {
        System.out.println(Bins.int32(1.23456f));
        System.out.println(Bins.float32(Bins.int32(1.23456f)));
        System.out.println(Bins.int32(4.123f));
    }
    @Test
    public void str()
    {
        String s = "123456";
        System.out.println(s.substring(s.length()-3, s.length()));
    }
    void testFile()
    {
        File file = new File("");
    }
}
