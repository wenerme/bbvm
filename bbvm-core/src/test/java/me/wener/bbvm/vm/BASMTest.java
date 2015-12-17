package me.wener.bbvm.vm;

import com.google.inject.Guice;
import com.google.inject.Injector;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import me.wener.bbvm.asm.BBAsmParser;
import me.wener.bbvm.asm.ParseException;
import me.wener.bbvm.util.Dumper;
import me.wener.bbvm.vm.invoke.BufferedReaderInput;
import me.wener.bbvm.vm.invoke.PrintStreamOutput;
import me.wener.bbvm.vm.res.Resources;
import org.junit.Test;

import javax.inject.Inject;
import javax.inject.Named;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.IOException;
import java.io.StringReader;
import java.nio.ByteOrder;
import java.nio.charset.Charset;
import java.nio.file.Files;

import static org.junit.Assert.assertNull;

/**
 * @author wener
 * @since 15/12/13
 */
public class BASMTest {
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
    @Inject
    @Named("R3")
    Register r3;
    @Inject
    private VM vm;
    private String basmContent;
    private BBAsmParser parser;
    private TestSpec io = new TestSpec();

    public BASMTest() {
        VMConfig.Builder builder = new VMConfig.Builder()
                .withModule(Resources.fileModule());
        Injector injector = Guice.createInjector(new VirtualMachineModule(builder.build()));
        injector.injectMembers(this);
        out = new ByteArrayOutputStream();
        in = new BufferedReaderInput();
        systemInvokeManager.register(new PrintStreamOutput(out), in);
    }

    public static void main(String[] args) throws IOException, ParseException {
        BASMTest test = new BASMTest();
        test.init(new File("bbvm-test/case/in/13.basm")).run();
    }

    @Test
    public void test() throws IOException, ParseException {
        BASMTest test = new BASMTest();
//        System.out.println(new File(".").getAbsoluteFile());
        test.init(new File("../bbvm-test/case/read-restore.basm")).run();
    }

    public BASMTest init(File basm) throws IOException, ParseException {
        basmFile = basm;
        basmContent = new String(Files.readAllBytes(basm.toPath()), charset);
        parser = new BBAsmParser(new StringReader(basmContent));
        io.accept(basmContent);
        return this;
    }

    public void run() throws ParseException {
        out.reset();
        in.setReader(io.output().toString());
        parser.Parse();
        int length = parser.estimateAddress();
        System.out.printf("Estimate length is %s\n", length);
        System.out.printf("Expected output \n%s\nWith input\n%s\n", io.output(), io.input());
        parser.checkLabel();
        ByteBuf buf = Unpooled.buffer(length).order(ByteOrder.LITTLE_ENDIAN);
        parser.write(buf);
        System.out.println(basmContent);
        System.out.println(Dumper.hexDumpReadable(buf));
        try {
            vm
                    .reset()
                    .setAddressTable(parser.getAddressTable())
                    .setSymbolTable(parser.createSymbolTable())
                    .setMemory(Memory.load(buf)).run();
            assertNull(vm.getLastError());
            io.assertMatch(out.toString());
            System.out.printf("Output\n%s\n", out.toString());
        } catch (Throwable e) {
            throw e;
        } finally {
            System.err.flush();
        }
    }
}
