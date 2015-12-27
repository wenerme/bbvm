package me.wener.bbvm;

import com.google.common.base.Splitter;
import com.typesafe.config.Config;
import com.typesafe.config.ConfigFactory;
import me.wener.bbvm.asm.ParseException;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

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
    private final static Logger log = LoggerFactory.getLogger(BasmSpecTest.class);

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
        Splitter.MapSplitter separator = Splitter.on('&').omitEmptyStrings().withKeyValueSeparator('=');
        Files.walk(Paths.get(first))
                .filter(p -> !p.toFile().isDirectory())
                .filter(p -> p.toFile().getName().endsWith(".basm"))
                .forEach(p -> {
                    BasmTester tester = test;
                    String fn = com.google.common.io.Files.getNameWithoutExtension(p.toString());
                    if (fn.contains("=")) {
                        Config config = ConfigFactory.parseMap(separator.split(fn));
                        tester = new BasmTester(config);
                        log.info("Create BasmTester for {} -> {}", config, p);
                    }

                    // When test failed, we need the output
                    try {
                        out.reset();
                        tester.init(p.toFile()).run();
                    } catch (Throwable e) {
                        System.out.println(out.toString());
                        System.out.println("Test failed for " + p);
                        e.printStackTrace();
                        throw e;
                    }
                });
    }
}
