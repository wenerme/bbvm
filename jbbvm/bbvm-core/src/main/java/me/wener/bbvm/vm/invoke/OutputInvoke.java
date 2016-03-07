package me.wener.bbvm.vm.invoke;

import me.wener.bbvm.vm.Operand;
import me.wener.bbvm.vm.SystemInvoke;

import java.util.function.Consumer;

/**
 * @author wener
 * @since 15/12/22
 */
public class OutputInvoke {
    private final Consumer<String> consumer;

    public OutputInvoke(Consumer<String> consumer) {
        this.consumer = consumer;
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 0)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 1)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 2)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 3)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 4)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 5)
    public void out(Operand a, Operand b) {
        switch (a.get()) {
            case 0:
                println(String.valueOf(b.get()));
                break;
            case 1:
                println(b.getString());
                break;
            case 2:
                print(b.getString());
                break;
            case 4:
                print(String.valueOf(Character.toChars(b.get())[0]));
                break;
            case 3:
                print(String.valueOf(b.get()));
                break;
            case 5:
                float v = b.getFloat();
                print(String.format("%.6f", v));
                break;
        }
    }

    protected void print(String s) {
        consumer.accept(s);
    }

    protected void println(String s) {
        consumer.accept(s);
        consumer.accept("\n");
    }
}
