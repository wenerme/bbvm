package me.wener.bbvm.core;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Arrays;
import me.wener.bbvm.utils.Bins;
import org.junit.Test;

public class SimpleTest
{
    public static void main(String[] args) throws IOException
    {
        byte[] bytes = Files.readAllBytes(Paths.get("D:\\dev\\projects\\bbvm\\doc\\ignored\\test\\BB\\Sim\\BBasic\\test.bin"));
        BBVm vm = new BBVm(null);
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
