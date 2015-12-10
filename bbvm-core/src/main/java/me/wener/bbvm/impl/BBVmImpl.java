package me.wener.bbvm.impl;

import com.google.common.base.Preconditions;
import me.wener.bbvm.api.BBVm;
import me.wener.bbvm.api.Device;
import me.wener.bbvm.api.DeviceFunction;
import me.wener.bbvm.dev.EnvType;
import me.wener.bbvm.util.Bins;
import me.wener.bbvm.util.val.Values;
import me.wener.bbvm.vm.*;

import java.nio.charset.Charset;
import java.util.Scanner;

/*
//
//                       _oo0oo_
//                      o8888888o
//                      88" . "88
//                      (| -_- |)
//                      0\  =  /0
//                    ___/`---'\___
//                  .' \\|     |// '.
//                 / \\|||  :  |||// \
//                / _||||| -:- |||||- \
//               |   | \\\  -  /// |   |
//               | \_|  ''\---/''  |_/ |
//               \  .-\__  '-'  ___/-. /
//             ___'. .'  /--.--\  `. .'___
//          ."" '<  `.___\_<|>_/___.' >' "".
//         | | :  `- \`.;`\ _ /`;.`/ - ` : | |
//         \  \ `_.   \_ __\ /__ _/   .-` /  /
//     =====`-.____`.___ \_____/___.-`___.-'=====
//                       `=---='
//
//
//     ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//
//               佛祖保佑         永无BUG
//
//
//   █████▒█    ██  ▄████▄   ██ ▄█▀       ██████╗ ██╗   ██╗ ██████╗
// ▓██   ▒ ██  ▓██▒▒██▀ ▀█   ██▄█▒        ██╔══██╗██║   ██║██╔════╝
// ▒████ ░▓██  ▒██░▒▓█    ▄ ▓███▄░        ██████╔╝██║   ██║██║  ███╗
// ░▓█▒  ░▓▓█  ░██░▒▓▓▄ ▄██▒▓██ █▄        ██╔══██╗██║   ██║██║   ██║
// ░▒█░   ▒▒█████▓ ▒ ▓███▀ ░▒██▒ █▄       ██████╔╝╚██████╔╝╚██████╔╝
//  ▒ ░   ░▒▓▒ ▒ ▒ ░ ░▒ ▒  ░▒ ▒▒ ▓▒       ╚═════╝  ╚═════╝  ╚═════╝
//  ░     ░░▒░ ░ ░   ░  ▒   ░ ░▒ ▒░
//  ░ ░    ░░░ ░ ░ ░        ░ ░░ ░
//           ░     ░ ░      ░  ░
//                 ░
//
// WRITTEN BY
//  __  _  __ ____   ____   ___________
//  \ \/ \/ // __ \ /    \_/ __ \_  __ \
//   \     /\  ___/|   |  \  ___/|  | \/
//    \/\_/  \___  >___|  /\___  >__|
//               \/     \/     \/
*/
@SuppressWarnings("ConstantConditions")
public class BBVmImpl implements BBVm
{
    public static final Charset DEFAULT_CHARSET = Charset.forName("GBK");
    private final Device device;
    @SuppressWarnings("FieldCanBeLocal")
    private final boolean logInst = true;//log.getLevel() == Level.INFO;
    private final DeviceFunction deviceFunction;
    private final StringHandlePool stringPool = new StringHandlePool();
    private final byte[] stack = new byte[1024];
    private final Reg rp = new Reg("RP");
    private final Reg rb = new Reg("RB");
    private final Reg rs = new Reg("RS");
    private final Reg rf = new Reg("RF");
    private final Reg r0 = new Reg("R0");
    private final Reg r1 = new Reg("R1");
    private final Reg r2 = new Reg("R2");
    private final Reg r3 = new Reg("R3");
    private long startTick;
    private byte[] memory;
    private InstructionContext context;
    private int stackSize = 1000;

    private EnvType envType = EnvType.ENV_SIM;

    private int romSize;
    /**
     * 数据指针位置
     */
    private int dataPtr = 0;
    private boolean running = false;
    private boolean useConsoleIO = false;

    public BBVmImpl(Device device)
    {

        this.device = device;
        deviceFunction = device.getFunction();
        //deviceFunction = device.getFunction();
    }

    @Override
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
    @SuppressWarnings("StatementWithEmptyBody")
    public void start()
    {
        startTick = System.currentTimeMillis();
        running = true;
        while (running)
            loop();
    }

    boolean loop()
    {
        final int pc = rp.asInt();// 记录位置

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

        // 如果 RP 没变,则自增
        if (pc == rp.asInt())
            rp.set(pc + Opcode.length(context.getInstruction()));

        return true;
    }

    protected void doInstruction(InstructionContext ctx)
    {
        final Operand op1 = ctx.getOp1();
        final Operand op2 = ctx.getOp2();
        final Integer opv1 = op1.asInt();
        final Integer opv2 = op2.asInt();
        final DataType dataType = ctx.getDataType();
        final Opcode instruction = ctx.getInstruction();

        switch (instruction)
        {
            case NOP:
                if (logInst)
                    log(instruction);
                break;
            case LD:
                if (logInst)
                    log(instruction, dataType, op1, op2);
                switch (dataType)
                {
                    case DWORD:
                    case FLOAT:
                    case INT:
                        op1.set(op2.asInt());
                        break;
                    case BYTE:
                        op1.set(op1.asInt() & 0xffffff00 | (op2.asInt() & 0xff));
                        break;
                    case WORD:
                        op1.set(op1.asInt() & 0xffff0000 | (op2.asInt() & 0xffff));
                        break;
                    default:
                        throw unsupport("未知的数据类型: %s", dataType);
                }
                break;
            case PUSH:
                if (logInst)
                    log(instruction, op1);

                push(op1.asInt());
                break;
            case POP:
                if (logInst)
                    log(instruction, op1);

                op1.set(pop());
                break;
            case IN:
                if (logInst)
                    log(instruction, op1, op2);

                in(ctx);
                break;
            case OUT:
                if (logInst)
                    log(instruction, op1, op2);

                out(ctx);
                break;
            case JMP:
                if (logInst)
                    log(instruction, op1);

                rp.set(op1.asInt());
                break;
            case JPC:
            {
                // JPC 的数据类型为比较操作
                CompareType org = Values.fromValue(CompareType.class, (int) Bins.int4(ctx.getFirstByte(), 2));
                CompareType flag = Values.fromValue(CompareType.class, rf.asInt());
                boolean valid = false;

                if (logInst)
                    log(instruction, org, op1);


                switch (flag)
                {
                    case A:
                        if (org == CompareType.AE || org == CompareType.A || org == CompareType.NZ)
                            valid = true;
                        break;
                    case B:
                        if (org == CompareType.BE || org == CompareType.B || org == CompareType.NZ)
                            valid = true;
                        break;
                    case Z:
                        if (org == CompareType.Z || org == CompareType.AE || org == CompareType.BE)
                            valid = true;
                        break;
                    default:
                        if (org.equals(flag))
                            valid = true;
                }

                if (valid)
                    rp.set(opv1);
            }
            break;
            case CALL:
                if (logInst)
                    log(instruction, op1);

                // 设置返回位置为下一句的开始
                push(rp.asInt() + Opcode.length(instruction));
                rp.set(opv1);
                break;
            case RET:
                if (logInst)
                    log(instruction);

                rp.set(pop());
                break;
            case CMP:
            {
                if (logInst)
                    log(instruction, dataType, op1, op2);

                float a = opv1;
                float b = opv2;
                if (dataType == DataType.FLOAT)
                {
                    a = Bins.float32(opv1);
                    b = Bins.float32(opv2);
                }
                float c = a - b;
                if (c > 0)
                    rf.set(CompareType.A.asInt());
                else if (c < 0)
                    rf.set(CompareType.B.asInt());
                else
                    rf.set(CompareType.Z.asInt());
            }
            break;
            case CAL:
            {
                if (logInst)
                    log(instruction, op1);

                CalculateType op = Values.fromValue(CalculateType.class, ctx.getSpecialByte());
                // 返回结果为 R0
                double a = opv1;
                double b = opv2;
                if (dataType.equals(DataType.FLOAT))
                {
                    a = Bins.float32(opv1);
                    b = Bins.float32(opv2);
                }
                double c;
                switch (op)
                {
                    case ADD:
                        c = a + b;
                        break;
                    case DIV:
                        c = a / b;
                        break;
                    case MOD:
                        c = a % b;
                        break;
                    case MUL:
                        c = a * b;
                        break;
                    case SUB:
                        c = a - b;
                        break;
                    default:
                        throw unsupport("未知计算操作: %s", op);
                }
                int ret = (int) c;
                // 值返回归约
                switch (dataType)
                {
                    case FLOAT:
                        ret = Bins.int32((float) c);
                        break;
                    case BYTE:
                        ret &= 0xff;
                        break;
                    case WORD:
                        ret &= 0xffff;
                        break;
                }
                op1.set(ret);
            }
            break;
            case EXIT:
                if (logInst)
                    log(instruction);

                exit();
                break;
            default:
                throw unsupport("未知指令: %s", instruction);
        }
    }

    protected UnsupportedOperationException unsupport(String format, Object... args)
    {
        return unsupport(String.format(format, args));
    }

    protected UnsupportedOperationException unsupport(String str)
    {
        return new UnsupportedOperationException(str);
    }

    @Override
    public void push(int v)
    {
        Bins.int32l(memory, rs.asInt(), v);
        rs.set(rs.asInt() - 4);
    }

    @Override
    public int pop()
    {
        rs.set(rs.asInt() + 4);
        return Bins.int32l(memory, rs.asInt());
    }

    protected Integer[] readParameters(int n, int offset)
    {
        Integer[] parameters = new Integer[n];

        for (int i = 0; i < n; i++)
        {
            parameters[n - i - 1] = Bins.int32l(memory, offset);
            offset += 4;
        }

        return parameters;
    }

    /**
     * 处理 out 端口操作
     *
     * @return 如果被处理了, 返回 true 否则 false
     */
    protected boolean out(InstructionContext ctx)
    {
        return beforeOut(ctx) || out0(ctx) || afterOut(ctx);
    }

    protected boolean out0(InstructionContext ctx)
    {
        Integer input = ctx.getOp2().asInt();
        switch (ctx.getOp1().asInt())
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
            //  11 | 键入字符串 | 0 | R3:目标字符串句柄 | r3所指字符串的内容变为键入的字符串
            case 11:
            {
                String in = deviceFunction.INPUT();
                stringHandle(r3.asInt()).set(deviceFunction.INPUT());
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
            // 14 | 从数据区读取字符串 | 0 | R3:目标字符串句柄 | r3所指字符串的内容变为读取的字符串
            case 14:
            {
                Preconditions.checkState(input == 0, "输入的值为 %s, 要求为 0", input);
                byte[] bytes = Bins.zByte(memory, dataPtr);
                dataPtr += bytes.length + 1;
                stringHandle(r3.asInt()).set(new String(bytes, DEFAULT_CHARSET));
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
            // 16 | 设定模拟器屏幕 | 0 | R2:宽, R3:高 |  SetLcd
            case 16:
            {
                deviceFunction.SETLCD(r2.asInt(), r3.asInt());
            }
            break;
            // 17 | 申请画布句柄 | R3:PAGE句柄 | - | CreatPage
            case 17:
            {
                r3.set(deviceFunction.CREATEPAGE());
            }
            break;
            // 18 | 释放画布句柄 | 0 | R3:PAGE句柄 |  DeletePage
            case 18:
            {
                deviceFunction.DELETEPAGE(r3.asInt());
            }
            break;
            // 19 | 申请图片句柄并从文件载入像素资源 | R3:资源句柄 | R3:文件名, R2:资源索引 |  LoadRes
            case 19:
            {
                int handle = deviceFunction.LOADRES(string(r3.asInt()), r2.asInt());
                r3.set(handle);
            }
            break;
            // ; 20 | 复制图片到画布上 | 0 | R3:地址,其他参数在该地址后<br>(PAGE,PIC,DX,DY,W,H,X,Y,MODE) |  ShowPic
            case 20:
            {
                Integer[] args = readParameters(9, r3.asInt());
                deviceFunction.SHOWPIC(args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8]);
            }
            break;
            // 21 | 显示画布 | 0 | R3:PAGE句柄 |  FlipPage
            case 21:
            {
                deviceFunction.FLIPPAGE(r3.asInt());
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

        Integer input = ctx.getOp2().asInt();
        switch (ctx.getOp1().asInt())
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
            //  11 | 键入字符串 | 0 | R3:目标字符串句柄 | r3所指字符串的内容变为键入的字符串
            case 11:
            {
                String line = new Scanner(System.in).nextLine();
                stringHandle(r3.asInt()).set(line);
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

    protected boolean afterOut(InstructionContext ctx)
    {
        return false;
    }

    /**
     * 处理 in 端口操作
     *
     * @return 如果被处理了, 返回 true 否则 false
     */
    protected boolean in(InstructionContext ctx)
    {
        return beforeIn(ctx) || in0(ctx) || afterIn(ctx);
    }

    protected boolean in0(InstructionContext ctx)
    {
        // o for out
        Operand o = ctx.getOp1();
        Operand op2 = ctx.getOp2();
        Integer opv2 = op2.asInt();
        switch (opv2)
        {
            //0 | 浮点数转换为整数 | 整数 | R3:浮点数
            case 0:
                o.set((int) Bins.float32(r3.asInt()));
                break;
            //1 | 整数转换为浮点数 | 浮点数 | R3:整数
            case 1:
                o.set(Bins.int32((float) r3.asInt()));
                break;
            //2 | 申请字符串句柄 | 申请到的句柄 |  |  IN():SHDL
            case 2:
                o.set(stringPool.acquire());
                break;
            //3 | 字符串转换为整数 | 整数 | R3:字符串句柄,__地址__ | IN(R3:SHDL):int<br>若r3的值不是合法的字符串句柄则返回r3的值
            case 3:
            {
                String s = string(r3.asInt());
                try
                {
                    o.set((int) Float.parseFloat(s));
                } catch (NumberFormatException e)
                {
                    o.set(r3.asInt());
                }
            }
            break;
            //4 | 整数转换为字符串 | 返回的值为r3:整数 | R2:目标字符串_句柄_<br>R3:整数 | IN(R2:SHDL,R3:int):int<br>r2所代表字符串的内容被修改
            case 4:
                stringPool.getResource(r2.asInt())// R2 所代表的字符串
                        .set(String.valueOf(r3.asInt()));// 被修改
                o.set(r3.asInt());// 返回 R3
                break;
            // 5 | 复制字符串 | r3的值 | R2:源字符串句柄<br>R3:目标字符串句柄 | IN(R2:SHDL,R3:SHDL):SHDL<br>r3所代表字符串的内容被修改
            case 5:
            {
                stringHandle(r3.asInt()).set(string(r2.asInt()));
                o.set(r3.asInt());// 返回 R3
            }
            break;
            // 6 | 连接字符串 | r3的值 | R2:源字符串<br>R3:目标字符串 | IN(R2:SHDL,R3:SHDL):SHDL<br>r3所代表字符串的内容被修改 R3+R2
            case 6:
            {
                stringHandle(r3.asInt()).concat(string(r2.asInt()));
                o.set(r3.asInt());// 返回 R3
            }
            break;
            // 7 | 获取字符串长度 | 字符串长度 | R3:字符串 | IN(R3:SHDL):int
            case 7:
            {
                r3.set(string(r3.asInt()).length());
                o.set(r3.asInt());// 返回 R3
            }
            break;
            // 8 | 释放字符串句柄 | r3的值 | R3:字符串句柄 | IN(R3:SHDL):SHDL
            case 8:
            {
                stringPool.release(r3.asInt());
                o.set(r3.asInt());// 返回 R3
            }
            break;
            // 9 | 比较字符串 | 两字符串的差值 相同为0，大于为1,小于为-1 | R2:基准字符串<br>R3:比较字符串 | IN(R2:SHDL,R3:SHDL):int
            case 9:
            {
                int v = string(r2.asInt()).compareTo(string(r3.asInt()));
                // 对结果进行约束
                if (v > 0)
                    v = 1;
                else if (v < 0)
                    v = -1;
                o.set(v);
            }
            break;
            // 10 | 整数转换为浮点数再转换为字符串 | r3的值 | R2:目标字符串<br>R3:整数 | r2所代表字符串的内容被修改
            case 10:
            {
                stringHandle(r2.asInt()).set(String.format("%.6f", (float) r3.asInt()));

                o.set(r3.asInt());
            }
            break;
            // 11 | 字符串转换为浮点数 | 浮点数 | R3:字符串 |
            case 11:
            {
                o.set(Bins.int32(Float.parseFloat(string(r3.asInt()))));
            }
            break;
            // 12 | 获取字符的ASCII码 | ASCII码 | R2:字符位置<br>R3:字符串 |
            // 备注: 返回的结果范围为有符号的 8bit值,因此对中文操作时返回负数
            case 12:
            {
                byte b = string(r3.asInt()).getBytes(DEFAULT_CHARSET)[r2.asInt()];

                o.set((int) b);
            }
            break;
            // 13 | 将给定字符串中指定索引的字符替换为给定的ASCII代表的字符 |
            // r3的值 | R1:ASCII码<br>R2:字符位置<br>R3:目标字符串 |
            // r3所代表字符串的内容被修改, 要求r3是句柄才能修改r3的值,给出的ASCII会进行模256的处理
            case 13:
            {
                StringHandle handle = stringHandle(r3.asInt());
                char[] chars = null;
                int ascii = r1.asInt() % 256;
                char[] c = Character.toChars(ascii);
                chars[r2.asInt()] = c[0];

                handle.set(new String(chars));

                o.set(r3.asInt());
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
            // 16 | 求正弦值 | X!的正弦值 | R3:X! |
            case 16:
            {
                o.set(Bins.int32((float) Math.sin(Bins.float32(r3.asInt()))));
            }
            break;
            // 17 | 求余弦值 | X!的余弦值 | R3:X! |
            case 17:
            {
                o.set(Bins.int32((float) Math.cos(Bins.float32(r3.asInt()))));
            }
            break;
            // 18 | 求正切值 | X!的正切值 | R3:X! |
            case 18:
            {
                o.set(Bins.int32((float) Math.tan(Bins.float32(r3.asInt()))));
            }
            break;
            // 19 | 求平方根值 | X!的平方根值 | R3:X! |
            case 19:
            {
                o.set(Bins.int32((float) Math.sqrt(Bins.float32(r3.asInt()))));
            }
            break;
            // 20 | 求绝对值 | X%的绝对值 | R3:X% |
            case 20:
            {
                o.set(Math.abs(r3.asInt()));
            }
            break;
            // 21 | 求绝对值 | X!的绝对值 | R3:X! |
            case 21:
            {
                o.set(Bins.int32(Math.abs(Bins.float32(r3.asInt()))));
            }
            break;
            // 22 | 重定位数据指针 | r3的值 | R2:数据位置 | r3中为任意值
            case 22:
            {
                dataPtr = r2.asInt();
                o.set(r3.asInt());
            }
            break;
            // 23 | 读内存数据 | 地址内容 | R3:地址 |
            case 23:
            {
                // 由于端口 23 和 24 在虚拟机内不能使用,所以无法测试,
                o.set(Bins.int32l(memory, r3.asInt()));
            }
            break;
            // 24 | 写内存数据 | r3的值 | R2:待写入数据<br>R3:待写入地址 |
            case 24:
            {
                Bins.int32l(memory, r3.asInt(), r2.asInt());

                o.set(r3.asInt());
            }
            break;
            // 25 | 获取环境值 | 环境值 |  |
            case 25:
            {
                o.set(envType.asInt());
            }
            break;
            // 32 | 整数转换为字符串 | r3的值 | R1:整数<br>R3:目标字符串 | r3所代表字符串的内容被修改
            case 32:
            {
                stringHandle(r3.asInt()).set(String.valueOf(r1.asInt()));
                o.set(r3.asInt());
            }
            break;
            // 33 | 字符串转换为整数 | 整数 | R3:字符串 |
            case 33:
            {
                o.set((int) Float.parseFloat(string(r3.asInt())));
            }
            break;
            // 34 | 获取字符的ASCII码 | ASCII码 | R3:字符串 |
            case 34:
            {
                o.set(string(r3.asInt()).codePointAt(0));
            }
            break;
            // 35 | 左取字符串 | r3的值 | R1:截取长度<br>R2:源字符串<br>R3:目标字符串 | r3所代表字符串的内容被修改 （此端口似乎不正常）
            case 35:
            {
                stringHandle(r3.asInt()).set(string(r2.asInt()).substring(0, r1.asInt()));

                o.set(r3.asInt());
            }
            break;
            // 36 | 右取字符串 | r3的值 | R1:截取长度<br>R2:源字符串<br>R3:目标字符串 | r3所代表字符串的内容被修改
            case 36:
            {
                Integer len = r1.asInt();
                String str = string(r2.asInt());
                int start = str.length() - len;
                int end = str.length();
                stringHandle(r3.asInt()).set(str.substring(start, end));

                o.set(r3.asInt());
            }
            break;
            // 37 | 中间取字符串 | R0:截取长度 | R0:截取长度<br>R1:截取位置<br>R2:源字符串<br>R3:目标字符串 | r3所代表字符串的内容被修改
            case 37:
            {
                Integer len = r0.asInt();
                String str = string(r2.asInt());
                int start = r1.asInt();
                int end = start + len;

                stringHandle(r3.asInt()).set(str.substring(start, end));

                o.set(r0.asInt());
            }
            break;
            // 38 | 查找字符串 | 位置 | R1:起始位置<br>R2:子字符串<br>R3:父字符串 |
            case 38:
            {
                int i = string(r3.asInt()).indexOf(string(r2.asInt()), r1.asInt());
                // FIXME PC 虚拟机中有这个BUG,不知道小机中有这个BUG没
                if (i < 0)
                    i = 0;
                o.set(i);
            }
            break;
            // 39 | 获取字符串长度 | 字符串长度 | R3:字符串 |
            case 39:
            {
                o.set(string(r3.asInt()).length());
            }
            break;
            default:
                return false;
        }
        return true;
    }

    protected boolean beforeIn(InstructionContext ctx)
    {
        return false;
    }

    protected boolean afterIn(InstructionContext ctx)
    {
        return false;
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
//            return stringPool.getResource(o).asInt();
            return "";
        return Bins.zString(memory, o, DEFAULT_CHARSET);
    }

    protected StringHandle stringHandle(Integer o)
    {
        if (o < 0)
            return stringPool.getResource(o);
        return StringHandle.valueOf(Bins.zString(memory, o, DEFAULT_CHARSET));
    }

    protected void log(Object... objects)
    {
        System.out.println(logString(true, objects));
    }

    protected String logString(boolean debug, Object... objects)
    {
        StringBuilder builder = new StringBuilder();
        boolean lastIsOperand = false;
        for (Object object : objects)
        {
            if (lastIsOperand)
            {
                builder.append(", ");
            }
            builder.append(object).append(" ");
            lastIsOperand = object instanceof Operand;
        }
        if (debug)
        {
            builder.append("\n;")
                    .append(String.format("R0= %s, R1= %s, R2= %s, R3= %s, RS= %s, RB= %s, rp= %s, RF= %s",
                            r0.asInt(), r1.asInt(), r2.asInt(), r3.asInt(), rs.asInt(), rb.asInt(), rp.asInt(), rf.asInt()));
        }
        return builder.toString();
    }

    @Override
    public void exit()
    {
        running = false;
    }

    /**
     * 获取寄存器
     * <pre>
     * RP | 0x0 | 程序计数器
     * RF | 0x1 |
     * RS | 0x2 | 栈顶位置
     * RB | 0x3 | 栈底位置
     * R0 | 0x4 | #0 寄存器
     * R1 | 0x5 | #1 寄存器
     * R2 | 0x6 | #2 寄存器
     * R3 | 0x7 | #3 寄存器
     * </pre>
     */
    @Override
    public Reg getRegister(int reg)
    {
        return getRegister(Values.fromValue(RegisterType.class, reg));
    }

    /**
     * 获取寄存器
     * <pre>
     * RP | 0x0 | 程序计数器
     * RF | 0x1 |
     * RS | 0x2 | 栈顶位置
     * RB | 0x3 | 栈底位置
     * R0 | 0x4 | #0 寄存器
     * R1 | 0x5 | #1 寄存器
     * R2 | 0x6 | #2 寄存器
     * R3 | 0x7 | #3 寄存器
     * </pre>
     */
    public Reg getRegister(RegisterType r)
    {
        switch (r)
        {
            case RP:
                return rp;
            case RF:
                return rf;
            case RS:
                return rs;
            case RB:
                return rb;
            case R0:
                return r0;
            case R1:
                return r1;
            case R2:
                return r2;
            case R3:
                return r3;
            default:
                throw new IllegalArgumentException("未知的寄存器 :" + r);
        }
    }

}
