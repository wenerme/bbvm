package me.wener.bbvm;

import com.google.common.base.Throwables;
import com.google.inject.Guice;
import com.google.inject.Injector;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import me.wener.bbvm.asm.BBAsmParser;
import me.wener.bbvm.asm.ParseException;
import me.wener.bbvm.util.Dumper;
import me.wener.bbvm.vm.*;
import me.wener.bbvm.vm.invoke.BufferedReaderInput;
import me.wener.bbvm.vm.invoke.GraphInvoke;
import me.wener.bbvm.vm.invoke.PrintStreamOutput;
import me.wener.bbvm.vm.res.ImageManager;
import me.wener.bbvm.vm.res.Resources;
import me.wener.bbvm.vm.res.Swings;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import java.io.*;
import java.nio.ByteOrder;
import java.nio.charset.Charset;
import java.nio.file.Files;

import static org.junit.Assert.assertNull;

/**
 * @author wener
 * @since 15/12/13
 */
public class BasmTester {
    private final static Logger log = LoggerFactory.getLogger(BasmTester.class);
    private final ByteArrayOutputStream out;
    private final BufferedReaderInput in;
    // Parse basm
    // Compare with bin
    // Extract io from basm
    // Run
    // Compare io
    File basmFile;
    Charset charset = Charset.forName("UTF-8");
    @Inject
    SystemInvokeManager systemInvokeManager;
    private PrintStream printStream = System.out;
    @Inject
    private VM vm;
    private String basmContent;
    private BBAsmParser parser;
    private TestSpec io = new TestSpec();

    public BasmTester() {
        VMConfig.Builder builder = new VMConfig.Builder()
                .withModule(Resources.fileModule())
                .withModule(Swings.graphModule())
                .invokeWith(GraphInvoke.class);
        Injector injector = Guice.createInjector(new VirtualMachineModule(builder.build()));
        injector.injectMembers(this);
        injector.getInstance(ImageManager.class).getResourceDirectory().add("../bbvm-test/image");
        out = new ByteArrayOutputStream();
        in = new BufferedReaderInput();
        systemInvokeManager.register(new PrintStreamOutput(out), in);
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
        io.clear().accept(basmContent);
        return this;
    }

    public void run() {
        out.reset();
        in.setReader(io.output().toString());
        try {
            parser.Parse();
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
                    .reset()
                    .setAddressTable(parser.getAddressTable())
                    .setSymbolTable(parser.createSymbolTable())
                    .setMemory(Memory.load(buf)).run();
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
