package me.wener.bbvm.vm;

import me.wener.bbvm.asm.ParseException;
import org.junit.Test;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;

/**
 * @author wener
 * @since 15/12/18
 */
public class BasmSpecTest {
    @Test
    public void in() throws IOException, ParseException {
        BasmTester test = new BasmTester();
        Files.walk(Paths.get("../bbvm-test/case/in"))
                .filter(p -> !p.toFile().isDirectory())
                .filter(p -> p.toFile().getName().endsWith(".basm"))
                .forEach(p -> test.init(p.toFile()).run());
    }
}
