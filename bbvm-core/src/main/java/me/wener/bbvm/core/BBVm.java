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
    private final boolean logInst = log.getLevel() == Level.INFO;
    private final DeviceFunction deviceFunction;
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

        if (pc >= memory.length)
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
        System.out.println(logString(objects));
    }

    protected String logString(Object... objects)
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
        return builder.toString();
    }

    /**
     * 处理 in 端口操作
     *
     * @return 如果被处理了, 返回 true 否则 false
     */
    protected boolean in(InstructionContext ctx)
    {
        return true;
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
