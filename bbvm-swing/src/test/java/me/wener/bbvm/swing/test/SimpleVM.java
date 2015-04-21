package me.wener.bbvm.swing.test;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Arrays;
import me.wener.bbvm.impl.BBVmImpl;
import me.wener.bbvm.impl.spi.DeviceProvider;

public class SimpleVM
{
    public static void main(String[] args) throws IOException, InterruptedException
    {
        System.out.println("当前目录: "+System.getProperty("user.dir"));
        byte[] bytes = Files.readAllBytes(Paths.get("doc/testsuit/BB/Sim/BBasic/test.bin"));
        BBVmImpl vm = new BBVmImpl(DeviceProvider.getProvider().createDevice(240, 320));
        byte[] mem = Arrays.copyOfRange(bytes, 16, bytes.length);
        vm.load(mem);
        vm.reset();

        vm.start();
    }
}
