package me.wener.bbvm.vm.invoke;

import me.wener.bbvm.vm.Operand;
import me.wener.bbvm.vm.SystemInvoke;

import java.io.OutputStream;
import java.io.PrintStream;

/**
 * @author wener
 * @since 15/12/13
 */
public class PrintStreamOutput extends PrintStream {

    public PrintStreamOutput(OutputStream out) {
        super(out);
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 0, b = SystemInvoke.ANY)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 1, b = SystemInvoke.ANY)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 2, b = SystemInvoke.ANY)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 3, b = SystemInvoke.ANY)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 4, b = SystemInvoke.ANY)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 5, b = SystemInvoke.ANY)
    public void out(Operand a, Operand b) {
        switch (a.get()) {
            case 0:
                println(b.get());
                break;
            case 1:
                println(b.getString());
                break;
            case 2:
                print(b.getString());
                break;
            case 4:
                print(Character.toChars(b.get())[0]);
                break;
            case 3:
                print(b.get());
                break;
            case 5:
                float v = b.getFloat();
                if (v == 0) {
                    print("000.000000");// Match the same result
                } else {
                    print(String.format("%.6f", v));
                }
                break;
        }
    }

}
