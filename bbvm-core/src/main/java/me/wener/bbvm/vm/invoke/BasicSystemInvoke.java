package me.wener.bbvm.vm.invoke;

import me.wener.bbvm.vm.Operand;
import me.wener.bbvm.vm.Register;
import me.wener.bbvm.vm.SystemInvoke;
import me.wener.bbvm.vm.VM;
import me.wener.bbvm.vm.res.StringManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import javax.inject.Named;

/**
 * @author wener
 * @since 15/12/13
 */
public class BasicSystemInvoke {
    private final static Logger log = LoggerFactory.getLogger(BasicSystemInvoke.class);
    private final VM vm;
    private final Register r3;
    private final Register r2;
    private final Register r1;
    private final Register r0;
    private int pointer;

    @Inject
    public BasicSystemInvoke(VM vm, @Named("R3") Register r3, @Named("R2") Register r2, @Named("R1") Register r1, @Named("R0") Register r0) {
        this.vm = vm;
        this.r3 = r3;
        this.r2 = r2;
        this.r1 = r1;
        this.r0 = r0;
    }

    /*
;0 | 浮点数转换为整数 | 整数 | r3:浮点数 | int(r3.float)
;1 | 整数转换为浮点数 | 浮点数 | r3:整数 | float(r3.int)
;2 | 申请字符串句柄 | 申请到的句柄 |  |  strPool.acquire
;3 | 字符串转换为整数 | 整数 | r3:字符串句柄,__地址__ | float(r3.str);若r3的值不是合法的字符串句柄则返回r3的值
;4 | 整数转换为字符串 | 返回的值为r3:整数 | r2:目标字符串_句柄_<br>r3:整数 | r2.str=str(r3.int);return r3.int;r2所代表字符串的内容被修改
     */
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 0)
    public void float2int(@Named("A") Operand o) {
        o.set((int) r3.getFloat());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 1)
    public void int2float(@Named("A") Operand o) {
        o.set((float) r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 2)
    public void acquireStringResource(@Named("A") Operand o, StringManager stringManager) {
        o.set(stringManager.create().getHandler());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 3)
    public void string2int(StringManager stringManager, @Named("A") Operand o) {
        o.set(Integer.parseInt(r3.getString()));
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 4)
    public void int2string(StringManager stringManager, @Named("A") Operand o) {
        r2.set(String.valueOf(r3.get()));
        o.set(r3.get());
    }

    /*
;5 | 复制字符串 | r3的值 | r2:源字符串句柄<br>r3:目标字符串句柄 | r3.str=r2.str;return r3
;6 | 连接字符串 | r3的值 | r2:源字符串<br>r3:目标字符串 | r3.str=r3.str+r2.str
;7 | 获取字符串长度 | 字符串长度 | r3:字符串 | strlen(r3.str)
;8 | 释放字符串句柄 | r3的值 | r3:字符串句柄 | strPool.release(r3);return r3
     */
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 5)
    public void stringCopy(@Named("A") Operand o) {
        r3.set(r2.getString());
        o.set(r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 6)
    public void stringConcat(@Named("A") Operand o) {
        r3.set(r3.getString() + r2.getString());
        o.set(r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 7)
    public void stringLength(Operand o) {
        o.set(r3.getString().length());//TODO Char length or bytes length ?
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 8)
    public void releaseString(StringManager stringManager, Operand o) {
        stringManager.getResource(r3.get()).close();
        o.set(r3.get());
    }

    /*
;9 | 比较字符串
; 返回: 两字符串的差值 相同为0，大于为1,小于为-1
; 参数: r2:基准字符串, r3:比较字符串
     */
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 9)
    public void stringCompare(@Named("A") Operand o) {
        int cmp = r2.getString().compareTo(r3.getString());
        if (cmp > 0) {
            o.set(1);
        } else if (cmp < 0) {
            o.set(-1);
        } else {
            o.set(0);
        }
    }

    /*
; 10 | 整数转换为浮点数再转换为字符串
; 返回: r3的值
; 参数: r2:目标字符串<br>r3:整数
; 备注: r2所代表字符串的内容被修改
     */
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 10)
    public void int2floatString(@Named("A") Operand o) {
        r2.set(String.format("%.6f", (float) r3.get()));
        o.set(r3.get());
    }

    /*
; 11 | 字符串转换为浮点数
; 返回: 浮点数
; 参数: r3:字符串
     */
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 11)
    public void string2float(@Named("A") Operand o) {
        try {
            o.set(Float.parseFloat(r3.getString()));
        } catch (NumberFormatException e) {
            log.warn("{}: {}", e.getClass().getSimpleName(), e.getMessage());
            o.set(0);
        }
    }

    /*
13 | 将给定字符串中指定索引的字符替换为给定的ASCII代表的字符 | r3的值 | r1:ASCII码<br>r2:字符位置<br>r3:目标字符串 | r3所代表字符串的内容被修改, 要求r3是句柄才能修改r3的值,给出的ASCII会进行模256的处理
14 | （功用不明） | 65535 |  |
15 | 获取嘀嗒计数 | 嘀嗒计数 |  | 这里不知道他是怎么算的这个数字,但是会随着时间增长就是了
16 | 求正弦值 | X!的正弦值 | r3:X! |
17 | 求余弦值 | X!的余弦值 | r3:X! |
18 | 求正切值 | X!的正切值 | r3:X! |
19 | 求平方根值 | X!的平方根值 | r3:X! |
20 | 求绝对值 | X%的绝对值 | r3:X% |
21 | 求绝对值 | X!的绝对值 | r3:X! |
     */
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 13)
    public void stringReplace(Operand o) {
        // TODO Charset
        byte[] bytes = r3.getString().getBytes();
        bytes[r2.get()] = (byte) (r1.get() & 0xff);
        r3.set(new String(bytes));
        o.set(r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 14)
    public void unknownIn14(Operand o) {
        o.set(65535);
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 15)
    public void getTick(Operand o, VM vm) {
        o.set(vm.getTick());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 16)
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 17)
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 18)
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 19)
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 20)
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 21)
    public void math(Operand a, Operand b) {
        switch (b.get()) {
            case 16:
                a.set((float) Math.sin(r3.getFloat()));
                break;
            case 17:
                a.set((float) Math.cos(r3.getFloat()));
                break;
            case 18:
                a.set((float) Math.tan(r3.getFloat()));
                break;
            case 19:
                a.set((float) Math.sqrt(r3.getFloat()));
                break;
            case 20:
                a.set(Math.abs(r3.get()));
                break;
            case 21:
                a.set(Math.abs(r3.getFloat()));
                break;
        }
    }

    /*
22 | 重定位数据指针 | r3的值 | r2:数据位置 | r3中为任意值
23 | 读内存数据 | 地址内容 | r3:地址 |
24 | 写内存数据 | r3的值 | r2:待写入数据<br>r3:待写入地址 |
     */
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 22)
    public void pointerReset(Operand o) {
        pointer = r2.get();
        o.set(r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 23)
    public void pointerRead(Operand o) {
        o.set(vm.getMemory().read(r3.get()));
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 24)
    public void pointerWrite(Operand o) {
        vm.getMemory().write(r3.get(), r2.get());
        o.set(r3.get());
    }

    /*
32 | 整数转换为字符串 | r3的值 | r1:整数<br>r3:目标字符串 | r3所代表字符串的内容被修改
33 | 字符串转换为整数 | 整数 | r3:字符串 |
34 | 获取字符第一个字符的ASCII码 | ASCII码 | r3:字符串 |
35 | 左取字符串 | r3的值 | r1:截取长度<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改 （此端口似乎不正常）
36 | 右取字符串 | r3的值 | r1:截取长度<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改
37 | 中间取字符串 | r0截取长度 | r0:截取长度<br>r1:截取位置<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改
38 | 查找字符串 | 位置 | r1:起始位置<br>r2:子字符串<br>r3:父字符串 |
39 | 获取字符串长度 | 字符串长度 | r3:字符串 |
     */
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 32)
    public void int2string(@Named("A") Operand o) {
        r3.set(Integer.toString(r1.get()));
        o.set(r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 33)
    public void string2int(@Named("A") Operand o) {
        try {
            o.set((int) Float.parseFloat(r3.getString()));
        } catch (NumberFormatException e) {
            // TODO should I do this ?
            o.set(0);
        }
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 34)
    public void firstCharCode(@Named("A") Operand o) {
        //TODO How to handler unicode ?
        o.set(r3.getString().codePointAt(0));
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 35)
    public void stringLeft(Operand o) {
        //TODO Exceptions
        r3.set(r2.getString().substring(0, r1.get()));

        o.set(r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 36)
    public void stringRight(Operand o) {
        String s = r2.getString();
        int start = s.length() - r1.get();
        if (start < 0) {
            log.warn("\"{}\".right(%s)", s, r1.get());
            r3.set(s.substring(0, s.length()));
        } else {
            r3.set(s.substring(start, s.length()));
        }

        o.set(r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 37)
    public void stringMid(Operand o) {
        int start = r1.get();
        int end = start + r0.get();
        r3.set(r2.getString().substring(start, end));

        o.set(r0.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 38)
    public void indexOfString(Operand o) {
        int i = r3.getString().indexOf(r2.getString(), r1.get());
        // FIXME PC 虚拟机中有这个BUG,不知道小机中有这个BUG没
        if (i < 0)
            i = 0;
        o.set(i);
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 39)
    public void stringLength2(Operand o) {
        // TODO Same as 7 ?
        stringLength(o);
    }
}
