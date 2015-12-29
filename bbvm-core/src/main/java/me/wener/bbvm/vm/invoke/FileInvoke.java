package me.wener.bbvm.vm.invoke;

import me.wener.bbvm.dev.FileManager;
import me.wener.bbvm.vm.Operand;
import me.wener.bbvm.vm.Register;
import me.wener.bbvm.vm.SystemInvoke;

import javax.inject.Inject;
import javax.inject.Named;
import java.io.IOException;

/**
 * @author wener
 * @since 15/12/17
 */
public class FileInvoke {
    private final Register r3;
    private final Register r2;
    private final Register r1;
    private final FileManager manager;
    private final int CURRENT_ADDRESS = 0x7FFFFFFF;

    @Inject
    public FileInvoke(@Named("R3") Register r3, @Named("R2") Register r2, @Named("R1") Register r1, FileManager manager) {
        this.r3 = r3;
        this.r2 = r2;
        this.r1 = r1;
        this.manager = manager;
    }

    /*
48 | 打开文件 | 0 | r0:打开方式<br>r1:文件号<br>r3:文件名字符串 | 打开方式目前只能为1
49 | 关闭文件 | 文件号 |  |
50	| 从文件读取数据 | 16:读取整数 | r1:文件号<br>r2:位置偏移量 | r3的值变为读取的整数
-	|  | 17:读取浮点数 | r1:文件号<br>r2:位置偏移量 | r3的值变为读取的浮点数
-	|  | 18:读取字符串 | r1:文件号<br>r2:位置偏移量<br>r3:目标字符串句柄 | r3所指字符串的内容变为读取的字符串
51	| 向文件写入数据 | 16:写入整数 | r1:文件号<br>r2:位置偏移量<br>r3:整数 |
-	|  | 17:写入浮点数 | r1:文件号<br>r2:位置偏移量<br>r3:浮点数 |
-	|  | 18:写入字符串 | r1:文件号<br>r2:位置偏移量<br>r3:字符串 |
     */
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 48, b = 0)
    public void open() throws IOException {
        r1.get(manager).open(r3.getString());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 49)
    public void close(@Named("B") Operand b) {
        b.get(manager).close();
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 50, b = 16)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 50, b = 17)
    public void readInt() throws IOException {
        if (r2.get() == CURRENT_ADDRESS) {
            r3.set(r1.get(manager).readInt());
        } else {
            r3.set(r1.get(manager).readInt(r2.get()));
        }
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 50, b = 18)
    public void readString() throws IOException {
        if (r2.get() == CURRENT_ADDRESS) {
            r3.set(r1.get(manager).readString());
        } else {
            r3.set(r1.get(manager).readString(r2.get()));
        }
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 51, b = 16)
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 51, b = 17)
    public void writeInt() throws IOException {
        if (r2.get() == CURRENT_ADDRESS) {
            r1.get(manager).writeInt(r3.get());
        } else {
            r1.get(manager).writeInt(r2.get(), r3.get());
        }
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 51, b = 18)
    public void writeString() throws IOException {
        if (r2.get() == CURRENT_ADDRESS) {
            r1.get(manager).writeString(r3.getString());
        } else {
            r1.get(manager).writeString(r2.get(), r3.getString());
        }
    }

    /*
52 | 判断文件位置指针是否指向文件尾 | 0;r3为0或1 | r3:文件号 |  Eof
53 | 获取文件长度 | 0 | r3:文件号,返回在 r3 |  Lof
54 | 获取文件位置指针的位置 | 0;返回值在r3 | r3:文件号 |  LOC(FILE)
55 | 定位文件位置指针 | 16 | r2:文件号<br>r3:目标位置 |
     */
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 52)
    public void isEof(@Named("B") Operand b) throws IOException {
        r3.set(r3.get(manager).isEof() ? 1 : 0);
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 53, b = 0)
    public void fileLength(@Named("B") Operand b) throws IOException {
        r3.set(r3.get(manager).length());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 54, b = 0)
    public void tell() throws IOException {
        r3.set(r3.get(manager).tell());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 55, b = 16)
    public void seek() throws IOException {
        r2.get(manager).seek(r3.get());
    }
}
