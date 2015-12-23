package me.wener.bbvm.vm.invoke;

import me.wener.bbvm.vm.Register;
import me.wener.bbvm.vm.SystemInvoke;

import javax.inject.Inject;
import javax.inject.Named;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.StringReader;

/**
 * @author wener
 * @since 15/12/17
 */
public class BufferedReaderInput {
    private BufferedReader reader;
    @Inject
    @Named("R3")
    private Register r3;
    @Inject
    @Named("R2")
    private Register r2;

    public BufferedReaderInput() {
    }

    public BufferedReaderInput(BufferedReader reader) {
        this.reader = reader;
    }

    public BufferedReaderInput setReader(BufferedReader reader) {
        this.reader = reader;
        return this;
    }

    public BufferedReaderInput setReader(String content) {
        return setReader(new BufferedReader(new StringReader(content)));
    }
/*
10 | 键入整数 | 0 |  | r3的值变为键入的整数
11 | 键入字符串 | 0 | r3:目标字符串句柄 | r3所指字符串的内容变为键入的字符串
12 | 键入浮点数 | 0 |  | r3的值变为键入的浮点数
 */

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 10, b = 0)
    public void inputInt() throws IOException {
        String line = reader.readLine();
        r3.set((int) Float.parseFloat(line));
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 11, b = 0)
    public void inputString() throws IOException {
        r3.set(reader.readLine());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 12, b = 0)
    public void inputFloat() throws IOException {
        String line = reader.readLine();
        r3.set(Float.parseFloat(line));
    }
}
