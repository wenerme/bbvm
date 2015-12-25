package me.wener.bbvm;

import com.google.common.base.Throwables;
import com.google.inject.Guice;
import com.google.inject.Injector;
import com.typesafe.config.Config;
import com.typesafe.config.ConfigFactory;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import me.wener.bbvm.asm.BBAsmParser;
import me.wener.bbvm.asm.ParseException;
import me.wener.bbvm.util.Dumper;
import me.wener.bbvm.vm.*;
import me.wener.bbvm.vm.invoke.GraphInvoke;
import me.wener.bbvm.vm.invoke.InputInvoke;
import me.wener.bbvm.vm.invoke.OutputInvoke;
import me.wener.bbvm.vm.res.ImageManager;
import me.wener.bbvm.vm.res.PageManager;
import me.wener.bbvm.vm.res.Resources;
import me.wener.bbvm.vm.res.swing.Swings;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import java.io.*;
import java.nio.ByteOrder;
import java.nio.charset.Charset;
import java.nio.file.Files;
import java.util.Scanner;

import static org.junit.Assert.assertNull;

/**
 * @author wener
 * @since 15/12/13
 */
public class BasmTester {
    private final static Logger log = LoggerFactory.getLogger(BasmTester.class);
    private static final Config DEFAULT_CONFIG = ConfigFactory.parseString("charset=UTF-8");
    private final ByteArrayOutputStream out;
    private final InputInvoke in;
    private final Charset charset;
    // Parse basm
    // Compare with bin
    // Extract io from basm
    // Run
    // Compare io
    File basmFile;
    @Inject
    SystemInvokeManager systemInvokeManager;
    private Config c;
    private PrintStream printStream = System.out;
    @Inject
    private VM vm;
    private String basmContent;
    private BBAsmParser parser;
    private TestSpec io = new TestSpec();

    public BasmTester() {
        this(DEFAULT_CONFIG);
    }

    public BasmTester(Config config) {
        this.c = config;
        if (c == DEFAULT_CONFIG) {
            c = c.withFallback(DEFAULT_CONFIG);
        }
        this.charset = Charset.forName(c.getString("charset"));
        VMConfig.Builder builder = new VMConfig.Builder()
                .withModule(Resources.fileModule())
                .withModule(Swings.graphModule())
                .charset(charset)
                .invokeWith(GraphInvoke.class);
        Injector injector = Guice.createInjector(new VirtualMachineModule(builder.build()));
        injector.injectMembers(this);
        injector.getInstance(ImageManager.class).getResourceDirectory().add("../bbvm-test/image");
        out = new ByteArrayOutputStream();
        in = new InputInvoke();
        // TODO Need a way to makeup the input and output
        PageManager manager = injector.getInstance(PageManager.class);
        systemInvokeManager.register(new OutputInvoke((s) -> {
            try {
                out.write(s.getBytes());
            } catch (IOException e) {
                Throwables.propagate(e);
            }
            manager.getScreen().draw(s);
        }), in);
    }

    public BasmTester setPrintStream(PrintStream printStream) {
        this.printStream = printStream;
        return this;
    }

    public BasmTester init(File basm) {
        log.info("Init basm tester {}", basm);
        basmFile = basm;
        try {
            basmContent = new String(Files.readAllBytes(basm.toPath()), charset);
        } catch (IOException e) {
            throw Throwables.propagate(e);
        }
        parser = new BBAsmParser(new StringReader(basmContent));
        parser.setCharset(charset);
        io.clear().accept(basmContent);
        return this;
    }

    public void run() {
        out.reset();
        Scanner scanner = new Scanner(io.output().toString());
        in.setSupplier(scanner::nextLine);
        try {
            parser.Parse();
            parser.getAssemblies().stream().forEach(s -> printStream.printf("%02d %s\n", s.getLine(), s.toAssembly()));
        } catch (ParseException e) {
            Throwables.propagate(e);
        }
        int length = parser.estimateAddress();
        printStream.printf("Estimate length is %s\n", length);
        printStream.printf("Expected output \n%s\nWith input\n%s\n", io.output(), io.input());
        parser.checkLabel();
        ByteBuf buf = Unpooled.buffer(length).order(ByteOrder.LITTLE_ENDIAN);
        parser.write(buf);
        printStream.println(basmContent);
        printStream.println(Dumper.hexDumpReadable(buf));
        try {
            vm
                    .setAddressTable(parser.getAddressTable())
                    .setSymbolTable(parser.createSymbolTable())
                    .setMemory(Memory.load(buf))
                    .reset()
                    .run();
            assertNull(vm.getLastError());
            io.assertMatch(out.toString());
            printStream.printf("Output\n%s\n", out.toString());
        } catch (Throwable e) {
            throw e;
        } finally {
            System.err.flush();
        }
    }
}
