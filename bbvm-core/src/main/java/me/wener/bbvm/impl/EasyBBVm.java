package me.wener.bbvm.impl;

import com.google.common.base.Preconditions;
import java.nio.charset.Charset;
import java.util.Scanner;
import me.wener.bbvm.api.BBVm;
import me.wener.bbvm.api.Device;
import me.wener.bbvm.api.DeviceFunction;
import me.wener.bbvm.def.CalOP;
import me.wener.bbvm.def.CmpOP;
import me.wener.bbvm.def.DataType;
import me.wener.bbvm.def.EnvType;
import me.wener.bbvm.def.InstructionType;
import me.wener.bbvm.def.RegType;
import me.wener.bbvm.utils.Bins;
import me.wener.bbvm.utils.val.Values;

@SuppressWarnings("ConstantConditions")
public class EasyBBVm extends VMContext implements BBVm
{
    public static final Charset DEFAULT_CHARSET = Charset.forName("GBK");
    private final Device device;
    @SuppressWarnings("FieldCanBeLocal")
    private final DeviceFunction deviceFunction;
    private final StringHandlePool stringPool = new StringHandlePool();
    private final Reg rp = new Reg("rp");
    private final Reg rb = new Reg("rb");
    private final Reg rs = new Reg("rs");
    private final Reg rf = new Reg("rf");
    private final Reg r0 = new Reg("r0");
    private final Reg r1 = new Reg("r1");
    private final Reg r2 = new Reg("r2");
    private final Reg r3 = new Reg("r3");
    private long startTick;
    private byte[] memory;
    private InstructionContext context;
    private int stackSize = 1000;

    private EnvType envType = EnvType.ENV_SIM;

    static
    {
        Values.cache(CalOP.class);
        Values.cache(CmpOP.class);
        Values.cache(DataType.class);
        Values.cache(InstructionType.class);
        Values.cache(RegType.class);
    }

    private int romSize;
    /**
     * 数据指针位置
     */
    private int dataPtr = 0;
    private boolean running = false;
    private boolean useConsoleIO = false;

    public EasyBBVm(Device device)
    {

        this.device = device;
        deviceFunction = device.getFunction();
    }

    public byte[] getMemory()
    {
        return memory;
    }

    public void load(byte[] bytes)
    {
        memory = new byte[bytes.length + stackSize];
        this.romSize = bytes.length;
        System.arraycopy(bytes, 0, memory, 0, bytes.length);
    }

    @Override
    public void reset()
    {
        context = new InstructionContext(this);
        rp.set(0);
        rb.set(memory.length - stackSize - 4);
        rs.set(memory.length - 4);
        rf.set(0);
        r0.set(0);
        r1.set(0);
        r2.set(0);
        r3.set(0);
    }

    @Override
    public void start()
    {
        startTick = System.currentTimeMillis();
        running = true;
        while (running)
            loop();
    }

    boolean loop()
    {
        final int pc = rp.get();// 记录位置

        if (pc >= romSize)
            return false;

        context.read(pc);

        try
        {
            doInstruction(context);
        } catch (Exception e)
        {
            e.printStackTrace();
        }

        // 如果 rp 没变,则自增
        if (pc == rp.get())
            rp.set(pc + context.getInstruction().length());

        return true;
    }

    protected void doInstruction(InstructionContext ctx)
    {
        final Operand op1 = ctx.getOp1();
        final Operand op2 = ctx.getOp2();
        final Integer opv1 = op1.get();
        final Integer opv2 = op2.get();
        final DataType dataType = ctx.getDataType();
        final InstructionType instruction = ctx.getInstruction();

        switch (instruction)
        {
            case NOP:
            case LD:

            case PUSH:
            case POP:

            case IN:

            case OUT:

            case JMP:

            case JPC:

            case CALL:

            case RET:

            case CMP:

            case CAL:

            case EXIT:

            default:
        }
    }

    protected boolean out(InstructionContext ctx)
    {
        Integer input = ctx.getOp2().get();
        switch (ctx.getOp1().get())
        {
            // 0 | 显示整数 | 整数 |  | 会换行
            // 1 | 显示字符串 | 字符串 |  | 会换行
            // 2 | 显示字符串 | 字符串 |  |
            // 3 | 显示整数 | 整数 |  |
            // 4 | 显示字符 | 字符ASCII码 |  |
            // 5 | 显示浮点数 | 浮点数 |  |
            case 0:
                deviceFunction.PRINT(input, '\n');
                break;
            case 1:
                deviceFunction.PRINT(string(input), '\n');
                break;
            case 2:
                deviceFunction.PRINT(string(input));
                break;
            case 3:
                deviceFunction.PRINT(input);
                break;
            case 4:
                deviceFunction.PRINT(Character.toChars(input)[0]);
                break;
            case 5:
                deviceFunction.PRINT(String.format("%.6f", Bins.float32(input)));
                break;
            //  10 | 键入整数 | 0 |  | r3的值变为键入的整数
            case 10:
            {
                String in = deviceFunction.INPUT();
                try
                {
                    r3.set((int) Float.parseFloat(in));
                } catch (NumberFormatException ignored)
                {
                    r3.set(0);
                }
            }
            break;
            //  11 | 键入字符串 | 0 | r3:目标字符串句柄 | r3所指字符串的内容变为键入的字符串
            case 11:
            {
                String in = deviceFunction.INPUT();
                stringHandle(r3.get()).set(deviceFunction.INPUT());
            }
            break;
            //  12 | 键入浮点数 | 0 |  | r3的值变为键入的浮点数
            case 12:
            {
                String in = deviceFunction.INPUT();
                try
                {
                    r3.set(Bins.int32(Float.parseFloat(in)));
                } catch (NumberFormatException ignored)
                {
                    r3.set(0);
                }
            }
            break;
            // 13 | 从数据区读取整数 | 0 |  | r3的值变为读取的整数
            case 13:
            {
                Preconditions.checkState(input == 0, "输入的值为 %s, 要求为 0", input);
                r3.set(Bins.int32l(memory, dataPtr));
                dataPtr += 4;
            }
            break;
            // 14 | 从数据区读取字符串 | 0 | r3:目标字符串句柄 | r3所指字符串的内容变为读取的字符串
            case 14:
            {
                Preconditions.checkState(input == 0, "输入的值为 %s, 要求为 0", input);
                byte[] bytes = Bins.zByte(memory, dataPtr);
                dataPtr += bytes.length + 1;
                stringHandle(r3.get()).set(new String(bytes, DEFAULT_CHARSET));
            }
            break;
            // 15 | 从数据区读取浮点数 | 0 |  | r3的值变为读取的浮点数
            // 注意: 读取浮点数和读取整数的内部表示是一样的
            case 15:
            {
                Preconditions.checkState(input == 0, "输入的值为 %s, 要求为 0", input);
                r3.set(Bins.int32l(memory, dataPtr));
                dataPtr += 4;
            }
            break;
            // 16 | 设定模拟器屏幕 | 0 | r2:宽, r3:高 |  SetLcd
            case 16:
            {
                deviceFunction.SETLCD(r2.get(), r3.get());
            }
            break;
            // 17 | 申请画布句柄 | r3:PAGE句柄 | - | CreatPage
            case 17:
            {
                r3.set(deviceFunction.CREATEPAGE());
            }
            break;
            // 18 | 释放画布句柄 | 0 | r3:PAGE句柄 |  DeletePage
            case 18:
            {
                deviceFunction.DELETEPAGE(r3.get());
            }
            break;
            // 19 | 申请图片句柄并从文件载入像素资源 | r3:资源句柄 | r3:文件名, r2:资源索引 |  LoadRes
            case 19:
            {
                int handle = deviceFunction.LOADRES(string(r3.get()), r2.get());
                r3.set(handle);
            }
            break;
            // ; 20 | 复制图片到画布上 | 0 | r3:地址,其他参数在该地址后<br>(PAGE,PIC,DX,DY,W,H,X,Y,MODE) |  ShowPic
            case 20:
            {
                Integer[] args = readParameters(9, r3.get());
                deviceFunction.SHOWPIC(args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8]);
            }
            break;
            // 21 | 显示画布 | 0 | r3:PAGE句柄 |  FlipPage
            case 21:
            {
                deviceFunction.FLIPPAGE(r3.get());
            }
            break;

            default:
                return false;
        }
        return true;
    }

    protected boolean beforeOut(InstructionContext ctx)
    {
        if (!useConsoleIO)
            return false;

        Integer input = ctx.getOp2().get();
        switch (ctx.getOp1().get())
        {
            // 0 | 显示整数 | 整数 |  | 会换行
            // 1 | 显示字符串 | 字符串 |  | 会换行
            // 2 | 显示字符串 | 字符串 |  |
            // 3 | 显示整数 | 整数 |  |
            // 4 | 显示字符 | 字符ASCII码 |  |
            // 5 | 显示浮点数 | 浮点数 |  |
            case 0:
                System.out.println(input);
                break;
            case 1:
                System.out.println(string(input));
                break;
            case 2:
                System.out.printf("%s", string(input));
                break;
            case 3:
                System.out.printf("%s", input);
                break;
            case 4:
                System.out.printf("%c", Character.toChars(input)[0]);
                break;
            case 5:
                System.out.printf("%.6f", Bins.float32(input));
                break;
            //  10 | 键入整数 | 0 |  | r3的值变为键入的整数
            case 10:
            {
                String line = new Scanner(System.in).nextLine();
                try
                {
                    r3.set((int) Float.parseFloat(line));
                } catch (NumberFormatException ignored)
                {
                    r3.set(0);
                }
            }
            break;
            //  11 | 键入字符串 | 0 | r3:目标字符串句柄 | r3所指字符串的内容变为键入的字符串
            case 11:
            {
                String line = new Scanner(System.in).nextLine();
                stringHandle(r3.get()).set(line);
            }
            break;
            //  12 | 键入浮点数 | 0 |  | r3的值变为键入的浮点数
            case 12:
            {
                String line = new Scanner(System.in).nextLine();
                try
                {
                    r3.set(Bins.int32(Float.parseFloat(line)));
                } catch (NumberFormatException ignored)
                {
                    r3.set(0);
                }
            }
            break;
            default:
                return false;
        }

        return true;
    }


    /**
     * 处理 in 端口操作
     *
     * @return 如果被处理了, 返回 true 否则 false
     */
    protected boolean in(InstructionContext ctx)
    {
        // o for out
        Operand o = ctx.getOp1();
        Operand op2 = ctx.getOp2();
        Integer opv2 = op2.get();
        switch (opv2)
        {
            //0 | 浮点数转换为整数 | 整数 | r3:浮点数
            case 0:
                o.set((int) Bins.float32(r3.get()));
                break;
            //1 | 整数转换为浮点数 | 浮点数 | r3:整数
            case 1:
                o.set(Bins.int32((float) r3.get()));
                break;
            //2 | 申请字符串句柄 | 申请到的句柄 |  |  IN():SHDL
            case 2:
                o.set(stringPool.acquire());
                break;
            //3 | 字符串转换为整数 | 整数 | r3:字符串句柄,__地址__ | IN(r3:SHDL):int<br>若r3的值不是合法的字符串句柄则返回r3的值
            case 3:
            {
                String s = string(r3.get());
                try
                {
                    o.set((int) Float.parseFloat(s));
                } catch (NumberFormatException e)
                {
                    o.set(r3.get());
                }
            }
            break;
            //4 | 整数转换为字符串 | 返回的值为r3:整数 | r2:目标字符串_句柄_<br>r3:整数 | IN(r2:SHDL,r3:int):int<br>r2所代表字符串的内容被修改
            case 4:
                stringPool.getResource(r2.get())// r2 所代表的字符串
                        .set(r3.get().toString());// 被修改
                o.set(r3.get());// 返回 r3
                break;
            // 5 | 复制字符串 | r3的值 | r2:源字符串句柄<br>r3:目标字符串句柄 | IN(r2:SHDL,r3:SHDL):SHDL<br>r3所代表字符串的内容被修改
            case 5:
            {
                stringHandle(r3.get()).set(string(r2.get()));
                o.set(r3.get());// 返回 r3
            }
            break;
            // 6 | 连接字符串 | r3的值 | r2:源字符串<br>r3:目标字符串 | IN(r2:SHDL,r3:SHDL):SHDL<br>r3所代表字符串的内容被修改 r3+r2
            case 6:
            {
                stringHandle(r3.get()).concat(string(r2.get()));
                o.set(r3.get());// 返回 r3
            }
            break;
            // 7 | 获取字符串长度 | 字符串长度 | r3:字符串 | IN(r3:SHDL):int
            case 7:
            {
                r3.set(string(r3.get()).length());
                o.set(r3.get());// 返回 r3
            }
            break;
            // 8 | 释放字符串句柄 | r3的值 | r3:字符串句柄 | IN(r3:SHDL):SHDL
            case 8:
            {
                stringPool.release(r3.get());
                o.set(r3.get());// 返回 r3
            }
            break;
            // 9 | 比较字符串 | 两字符串的差值 相同为0，大于为1,小于为-1 | r2:基准字符串<br>r3:比较字符串 | IN(r2:SHDL,r3:SHDL):int
            case 9:
            {
                int v = string(r2.get()).compareTo(string(r3.get()));
                // 对结果进行约束
                if (v > 0)
                    v = 1;
                else if (v < 0)
                    v = -1;
                o.set(v);
            }
            break;
            // 10 | 整数转换为浮点数再转换为字符串 | r3的值 | r2:目标字符串<br>r3:整数 | r2所代表字符串的内容被修改
            case 10:
            {
                stringHandle(r2.get()).set(String.format("%.6f", (float) r3.get()));

                o.set(r3.get());
            }
            break;
            // 11 | 字符串转换为浮点数 | 浮点数 | r3:字符串 |
            case 11:
            {
                o.set(Bins.int32(Float.parseFloat(string(r3.get()))));
            }
            break;
            // 12 | 获取字符的ASCII码 | ASCII码 | r2:字符位置<br>r3:字符串 |
            // 备注: 返回的结果范围为有符号的 8bit值,因此对中文操作时返回负数
            case 12:
            {
                byte b = string(r3.get()).getBytes(DEFAULT_CHARSET)[r2.get()];

                o.set((int) b);
            }
            break;
            // 13 | 将给定字符串中指定索引的字符替换为给定的ASCII代表的字符 |
            // r3的值 | r1:ASCII码<br>r2:字符位置<br>r3:目标字符串 |
            // r3所代表字符串的内容被修改, 要求r3是句柄才能修改r3的值,给出的ASCII会进行模256的处理
            case 13:
            {
                StringHandle handle = stringHandle(r3.get());
                char[] chars = handle.get().toCharArray();
                int ascii = r1.get() % 256;
                char[] c = Character.toChars(ascii);
                chars[r2.get()] = c[0];

                handle.set(new String(chars));

                o.set(r3.get());
            }
            break;
            // 14 | （功用不明） | 65536 |  |
            case 14:
            {
                o.set(65535);
            }
            break;
            // 15 | 获取嘀嗒计数 | 嘀嗒计数 |  | 这里不知道他是怎么算的这个数字,但是会随着时间增长就是了
            case 15:
            {
                o.set(getTick());
            }
            break;
            // 16 | 求正弦值 | X!的正弦值 | r3:X! |
            case 16:
            {
                o.set(Bins.int32((float) Math.sin(Bins.float32(r3.get()))));
            }
            break;
            // 17 | 求余弦值 | X!的余弦值 | r3:X! |
            case 17:
            {
                o.set(Bins.int32((float) Math.cos(Bins.float32(r3.get()))));
            }
            break;
            // 18 | 求正切值 | X!的正切值 | r3:X! |
            case 18:
            {
                o.set(Bins.int32((float) Math.tan(Bins.float32(r3.get()))));
            }
            break;
            // 19 | 求平方根值 | X!的平方根值 | r3:X! |
            case 19:
            {
                o.set(Bins.int32((float) Math.sqrt(Bins.float32(r3.get()))));
            }
            break;
            // 20 | 求绝对值 | X%的绝对值 | r3:X% |
            case 20:
            {
                o.set(Math.abs(r3.get()));
            }
            break;
            // 21 | 求绝对值 | X!的绝对值 | r3:X! |
            case 21:
            {
                o.set(Bins.int32(Math.abs(Bins.float32(r3.get()))));
            }
            break;
            // 22 | 重定位数据指针 | r3的值 | r2:数据位置 | r3中为任意值
            case 22:
            {
                dataPtr = r2.get();
                o.set(r3.get());
            }
            break;
            // 23 | 读内存数据 | 地址内容 | r3:地址 |
            case 23:
            {
                // 由于端口 23 和 24 在虚拟机内不能使用,所以无法测试,
                o.set(Bins.int32l(memory, r3.get()));
            }
            break;
            // 24 | 写内存数据 | r3的值 | r2:待写入数据<br>r3:待写入地址 |
            case 24:
            {
                Bins.int32l(memory, r3.get(), r2.get());

                o.set(r3.get());
            }
            break;
            // 25 | 获取环境值 | 环境值 |  |
            case 25:
            {
                o.set(envType.get());
            }
            break;
            // 32 | 整数转换为字符串 | r3的值 | r1:整数<br>r3:目标字符串 | r3所代表字符串的内容被修改
            case 32:
            {
                stringHandle(r3.get()).set(r1.get().toString());
                o.set(r3.get());
            }
            break;
            // 33 | 字符串转换为整数 | 整数 | r3:字符串 |
            case 33:
            {
                o.set((int) Float.parseFloat(string(r3.get())));
            }
            break;
            // 34 | 获取字符的ASCII码 | ASCII码 | r3:字符串 |
            case 34:
            {
                o.set(string(r3.get()).codePointAt(0));
            }
            break;
            // 35 | 左取字符串 | r3的值 | r1:截取长度<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改 （此端口似乎不正常）
            case 35:
            {
                stringHandle(r3.get()).set(string(r2.get()).substring(0, r1.get()));

                o.set(r3.get());
            }
            break;
            // 36 | 右取字符串 | r3的值 | r1:截取长度<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改
            case 36:
            {
                Integer len = r1.get();
                String str = string(r2.get());
                int start = str.length() - len;
                int end = str.length();
                stringHandle(r3.get()).set(str.substring(start, end));

                o.set(r3.get());
            }
            break;
            // 37 | 中间取字符串 | r0:截取长度 | r0:截取长度<br>r1:截取位置<br>r2:源字符串<br>r3:目标字符串 | r3所代表字符串的内容被修改
            case 37:
            {
                Integer len = r0.get();
                String str = string(r2.get());
                int start = r1.get();
                int end = start + len;

                stringHandle(r3.get()).set(str.substring(start, end));

                o.set(r0.get());
            }
            break;
            // 38 | 查找字符串 | 位置 | r1:起始位置<br>r2:子字符串<br>r3:父字符串 |
            case 38:
            {
                int i = string(r3.get()).indexOf(string(r2.get()), r1.get());
                // FIXME PC 虚拟机中有这个BUG,不知道小机中有这个BUG没
                if (i < 0)
                    i = 0;
                o.set(i);
            }
            break;
            // 39 | 获取字符串长度 | 字符串长度 | r3:字符串 |
            case 39:
            {
                o.set(string(r3.get()).length());
            }
            break;
            default:
                return false;
        }
        return true;
    }

    public int getTick()
    {
        return (int) (System.currentTimeMillis() - startTick);
    }

    /**
     * 获取内存中的字符串
     *
     * @param o 字符串偏移量或字符串句柄
     */
    protected String string(Integer o)
    {
        if (o < 0)
            return stringPool.getResource(o).get();
        return Bins.zString(memory, o, DEFAULT_CHARSET);
    }

    protected StringHandle stringHandle(Integer o)
    {
        if (o < 0)
            return stringPool.getResource(o);
        return StringHandle.valueOf(Bins.zString(memory, o, DEFAULT_CHARSET));
    }

    @Override
    public void exit()
    {
        running = false;
    }
}