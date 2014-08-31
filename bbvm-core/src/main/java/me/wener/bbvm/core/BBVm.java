package me.wener.bbvm.core;

import java.nio.charset.Charset;
import java.util.logging.Level;
import java.util.logging.Logger;
import me.wener.bbvm.core.asm.CalOP;
import me.wener.bbvm.core.asm.CmpOP;
import me.wener.bbvm.core.asm.DataType;
import me.wener.bbvm.core.asm.Instruction;
import me.wener.bbvm.core.asm.RegType;
import me.wener.bbvm.core.constant.Device;
import me.wener.bbvm.utils.Bins;

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
public class BBVm
{
    private final Device device;
    private final Logger log = Logger.getLogger(BBVm.class.toString());
    private final boolean logInst = true;//log.getLevel() == Level.INFO;
    private final DeviceFunction deviceFunction;
    private final StringHandlePool stringPool = new StringHandlePool();
    private final byte[] stack = new byte[1024];
    private final Reg rp = new Reg("rp");
    private final Reg rb = new Reg("rb");
    private final Reg rs = new Reg("rs");
    private final Reg rf = new Reg("rf");
    private final Reg r0 = new Reg("r0");
    private final Reg r1 = new Reg("r1");
    private final Reg r2 = new Reg("r2");
    private final Reg r3 = new Reg("r3");
    private byte[] memory;
    private InstructionContext context;
    private int stackSize = 1000;

    static
    {
        Values.cache(CalOP.class);
        Values.cache(CmpOP.class);
        Values.cache(DataType.class);
        Values.cache(Instruction.class);
        Values.cache(RegType.class);
    }

    private int romSize;

    protected BBVm(Device device)
    {

        this.device = device;
        deviceFunction = null;
        //deviceFunction = device.getFunction();
    }

    byte[] getMemory()
    {
        return memory;
    }

    public void load(byte[] bytes)
    {
        memory = new byte[bytes.length + stackSize];
        this.romSize = bytes.length;
        System.arraycopy(bytes, 0, memory, 0, bytes.length);
    }

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

    public void start()
    {

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
            rp.set(pc + Instruction.length(context.getInstruction()));

        return true;
    }

    protected void doInstruction(InstructionContext ctx)
    {
        final Operand op1 = ctx.getOp1();
        final Operand op2 = ctx.getOp2();
        final Integer opv1 = op1.get();
        final Integer opv2 = op2.get();
        final DataType dataType = ctx.getDataType();
        final Instruction instruction = ctx.getInstruction();

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
                    case T_DWORD:
                    case T_FLOAT:
                    case T_INT:
                        op1.set(op2.get());
                        break;
                    case T_BYTE:
                        op1.set(op1.get() & 0xffffff00 | (op2.get() & 0xff));
                        break;
                    case T_WORD:
                        op1.set(op1.get() & 0xffff0000 | (op2.get() & 0xffff));
                        break;
                    default:
                        throw unsupport("未知的数据类型: %s", dataType);
                }
                break;
            case PUSH:
                if (logInst)
                    log(instruction, op1);

                push(op1.get());
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

                rp.set(op1.get());
                break;
            case JPC:
            {
                // JPC 的数据类型为比较操作
                CmpOP org = Values.fromValue(CmpOP.class, (int) Bins.int4(ctx.getFirstByte(), 2));
                CmpOP flag = Values.fromValue(CmpOP.class, rf.get());
                boolean valid = false;

                if (logInst)
                    log(instruction, org, op1);


                switch (flag)
                {
                    case A:
                        if (org == CmpOP.AE || org == CmpOP.A || org == CmpOP.NZ)
                            valid = true;
                        break;
                    case B:
                        if (org == CmpOP.BE || org == CmpOP.B || org == CmpOP.NZ)
                            valid = true;
                        break;
                    case Z:
                        if (org == CmpOP.Z || org == CmpOP.AE || org == CmpOP.BE)
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
                push(rp.get() + Instruction.length(instruction));
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
                if (dataType == DataType.T_FLOAT)
                {
                    a = Bins.float32(opv1);
                    b = Bins.float32(opv2);
                }
                float c = a - b;
                if (c > 0)
                    rf.set(CmpOP.A.asValue());
                else if (c < 0)
                    rf.set(CmpOP.B.asValue());
                else
                    rf.set(CmpOP.Z.asValue());
            }
            break;
            case CAL:
            {
                if (logInst)
                    log(instruction, op1);

                CalOP op = Values.fromValue(CalOP.class, this.context.getSpecialByte());
                // 返回结果为 r0
                double a = opv1;
                double b = opv2;
                if (dataType.equals(DataType.T_FLOAT))
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
                    case T_FLOAT:
                        ret = Bins.int32((float) c);
                        break;
                    case T_BYTE:
                        ret &= 0xff;
                        break;
                    case T_WORD:
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

    protected void push(int v)
    {
        Bins.int32l(memory, rs.get(), v);
        rs.set(rs.get() - 4);
    }

    protected int pop()
    {
        rs.set(rs.get() + 4);
        return Bins.int32l(memory, rs.get());
    }

    /**
     * 处理 out 端口操作
     *
     * @return 如果被处理了, 返回 true 否则 false
     */
    protected boolean out(InstructionContext ctx)
    {
        Integer opv2 = ctx.getOp2().get();
        switch (ctx.getOp1().get())
        {
            case 0:
                System.out.println(opv2);
                break;
            case 1:
                System.out.println(string(opv2));
                break;
            case 2:
                System.out.printf("%s", string(opv2));
                break;
            case 3:
                System.out.printf("%s", opv2);
                break;
            case 4:
                System.out.printf("%c", Character.toChars(opv2)[0]);
                break;
            case 5:
                System.out.printf("%.6f", Bins.float32(opv2));
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
        Operand op1 = ctx.getOp1();
        Operand op2 = ctx.getOp2();
        Integer opv2 = op2.get();
        switch (opv2)
        {
            //0 | 浮点数转换为整数 | 整数 | r3:浮点数
            case 0:
                op1.set((int)Bins.float32(r3.get()));
                break;
            //1 | 整数转换为浮点数 | 浮点数 | r3:整数
            case 1:
                op1.set(Bins.int32((float)r3.get()));
                break;
            //2 | 申请字符串句柄 | 申请到的句柄 |  |  IN():SHDL
            case 2:
                op1.set(stringPool.acquire());
                break;
            //3 | 字符串转换为整数 | 整数 | r3:字符串句柄,__地址__ | IN(r3:SHDL):int<br>若r3的值不是合法的字符串句柄则返回r3的值
            case 3:
                String s = stringPool.getResource(r3.get()).get();
                try
                {
                    op1.set(Integer.parseInt(s));
                } catch (NumberFormatException e)
                {
                    op1.set(r3.get());
                }
                break;
            //4 | 整数转换为字符串 | 返回的值为r3:整数 | r2:目标字符串_句柄_<br>r3:整数 | IN(r2:SHDL,r3:int):int<br>r2所代表字符串的内容被修改
            case 4:
                stringPool.getResource(r2.get())// r2 所代表的字符串
                          .set(r3.get().toString());// 被修改
                op1.set(r3.get());// 返回 r3
                break;
            // 5 | 复制字符串 | r3的值 | r2:源字符串句柄<br>r3:目标字符串句柄 | IN(r2:SHDL,r3:SHDL):SHDL<br>r3所代表字符串的内容被修改
            case 5:
            {
                StringHandle src = stringPool.getResource(r2.get());
                StringHandle dest = stringPool.getResource(r3.get());
                dest.set(src.get());
                op1.set(r3.get());// 返回 r3
            }
                break;
            // 6 | 连接字符串 | r3的值 | r2:源字符串<br>r3:目标字符串 | IN(r2:SHDL,r3:SHDL):SHDL<br>r3所代表字符串的内容被修改 r3+r2
            case 6:
            {
                StringHandle src = stringPool.getResource(r2.get());
                StringHandle dest = stringPool.getResource(r3.get());
                dest.set(src.get()+dest.get());
                op1.set(r3.get());// 返回 r3
            }
            break;
            // 7 | 获取字符串长度 | 字符串长度 | r3:字符串 | IN(r3:SHDL):int
            case 7:
            {
                StringHandle src = stringPool.getResource(r2.get());
                StringHandle dest = stringPool.getResource(r3.get());
                dest.set(src.get()+dest.get());
                op1.set(r3.get());// 返回 r3
            }
            break;
            // 8 | 释放字符串句柄 | r3的值 | r3:字符串句柄 | IN(r3:SHDL):SHDL
            case 8:
            {
                stringPool.release(r3.get());
                op1.set(r3.get());// 返回 r3
            }
            break;




            default:
                return false;
        }
        return true;
    }

    /**
     * 获取内存中的字符串
     *
     * @param offset 字符串偏移量
     */
    protected String string(Integer offset)
    {
        return Bins.zString(memory, offset, Charset.forName("GBK"));
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
                   .append(String.format("r0= %s, r1= %s, r2= %s, r3= %s, rs= %s, rb= %s, rp= %s, rf= %s",
                           r0.get(),r1.get(),r2.get(),r3.get(),rs.get(),rb.get(),rp.get(),rf.get()));
        }
        return builder.toString();
    }

    private void exit()
    {
        System.exit(0);
    }

    /**
     * 获取寄存器
     * <pre>
     * rp | 0x0 | 程序计数器
     * rf | 0x1 |
     * rs | 0x2 | 栈顶位置
     * rb | 0x3 | 栈底位置
     * r0 | 0x4 | #0 寄存器
     * r1 | 0x5 | #1 寄存器
     * r2 | 0x6 | #2 寄存器
     * r3 | 0x7 | #3 寄存器
     * </pre>
     */
    IntHolder register(int reg)
    {
        RegType r = Values.fromValue(RegType.class, reg);
        switch (r)
        {
            case rp:
                return rp;
            case rf:
                return rf;
            case rs:
                return rs;
            case rb:
                return rb;
            case r0:
                return r0;
            case r1:
                return r1;
            case r2:
                return r2;
            case r3:
                return r3;
            default:
                throw new IllegalArgumentException("未知的寄存器 :" + reg);
        }
    }

}
