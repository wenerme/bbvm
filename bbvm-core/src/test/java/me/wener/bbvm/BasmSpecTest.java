package me.wener.bbvm;

import me.wener.bbvm.asm.ParseException;
import org.junit.Test;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.PrintStream;
import java.nio.file.Files;
import java.nio.file.Paths;

/**
 * @author wener
 * @since 15/12/18
 */
public class BasmSpecTest {
    @Test
    public void in() throws IOException, ParseException {
        doTest("../bbvm-test/case/in");
    }

    @Test
    public void out() throws IOException, ParseException {
        doTest("../bbvm-test/case/out");
    }

    @Test
    public void file() throws IOException, ParseException {
        doTest("../bbvm-test/case/file");
    }

    @Test
    public void graph() throws IOException, ParseException {
        doTest("../bbvm-test/case/graph");
    }

    @Test
    public void basic() throws IOException, ParseException {
        doTest("../bbvm-test/case/basic");
    }

    private void doTest(String first) throws IOException {
        BasmTester test = new BasmTester();
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        test.setPrintStream(new PrintStream(out));
        Files.walk(Paths.get(first))
                .filter(p -> !p.toFile().isDirectory())
                .filter(p -> p.toFile().getName().endsWith(".basm"))
                .forEach(p -> {
                    // When test failed, we need the output
                    try {
                        out.reset();
                        test.init(p.toFile()).run();
                    } catch (Throwable e) {
                        System.out.println(out.toString());
                        throw e;
                    }
                });
    }
}
