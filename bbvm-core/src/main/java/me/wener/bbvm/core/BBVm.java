package me.wener.bbvm.core;

import me.wener.bbvm.core.asm.CalOP;
import me.wener.bbvm.core.asm.CmpOP;
import me.wener.bbvm.core.asm.DataType;
import me.wener.bbvm.core.asm.Instruction;
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
public abstract class BBVm
{
    private final Device device;
    private final DeviceFunction deviceFunction;
    private final byte[] stack = new byte[1024];
    private byte[] memory;
    InstructionContext ctx = new InstructionContext(memory);
    private Reg rp = new Reg("rp");
    private Reg rb = new Reg("rb");
    private Reg rs = new Reg("rs");
    private Reg rf = new Reg("rf");
    private Reg r0 = new Reg("r0");
    private Reg r1 = new Reg("r1");
    private Reg r2 = new Reg("r2");
    private Reg r3 = new Reg("r3");


    protected BBVm(Device device)
    {

        this.device = device;
        deviceFunction = device.getFunction();
    }

    public void load(byte[] bytes)
    {
        memory = bytes;
    }

    public void reset()
    {

    }

    private void loop()
    {
        final int pc = rp.get();// 记录位置

        ctx.read(pc);

        try
        {
            doInstruction(ctx);
        } catch (Exception e)
        {
            e.printStackTrace();
        }

        // 如果 rp 没变,则自增
        if (pc == rp.get())
            rp.set(pc + Instruction.length(ctx.instruction));
    }

    protected void doInstruction(InstructionContext ctx)
    {
        final Operand op1 = ctx.op1;
        final Operand op2 = ctx.op2;
        final Integer opv1 = ctx.op1.get();
        final Integer opv2 = ctx.op2.get();
        final DataType dataType = ctx.dataType;
        final Instruction instruction = ctx.instruction;

        switch (instruction)
        {
            case NOP:
                break;
            case LD:
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
                push(op1.get());
                break;
            case POP:
                op1.set(pop());
                break;
            case IN:
                in(ctx);
                break;
            case OUT:
                out(ctx);
                break;
            case JMP:
                rp.set(op1.get());
                break;
            case JPC:
            {
                CmpOP org = Values.fromValue(CmpOP.class, ctx.getSpecialByte());
                CmpOP flag = Values.fromValue(CmpOP.class, rf.get());
                boolean valid = false;

                switch (flag)
                {
                    case A:
                        if (org == CmpOP.AE || org == CmpOP.A)
                            valid = true;
                        break;
                    case B:
                        if (org == CmpOP.BE || org == CmpOP.B)
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
                push(rp.get());
                rp.set(opv1);
                break;
            case RET:
                rp.set(pop());
                break;
            case CMP:
            {
                float a = opv1;
                float b = opv2;
                if (dataType.equals(DataType.T_FLOAT))
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
                CalOP op = Values.fromValue(CalOP.class, this.ctx.specialByte);
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
            case EXIT:
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
    {}

    protected int pop()
    {
        return 0;
    }

    protected void out(InstructionContext port)
    {}

    protected void in(InstructionContext port)
    {}

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
        switch (reg)
        {
            case 0:
                return rp;
            case 1:
                return rf;
            case 2:
                return rs;
            case 3:
                return rb;
            case 4:
                return r0;
            case 5:
                return r1;
            case 6:
                return r2;
            case 7:
                return r3;
            default:
                throw new IllegalArgumentException("未知的寄存器 :" + reg);
        }
    }

    class InstructionContext
    {
        private final byte[] memory;
        private Instruction instruction;
        private Operand op1;
        private Operand op2;
        private DataType dataType;
        private int specialByte;
        private int addressingType;

        InstructionContext(byte[] memory) {this.memory = memory;}

        public Instruction getInstruction()
        {

            return instruction;
        }

        public Operand getOp1()
        {
            return op1;
        }

        public Operand getOp2()
        {
            return op2;
        }

        public DataType getDataType()
        {
            return dataType;
        }

        public int getSpecialByte()
        {
            return specialByte;
        }

        public int getAddressingType()
        {
            return addressingType;
        }

        void read(int pc)
        {
           /*
                指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数 + 第二个操作数
             0x 0        0          0             0          0000           0000
            */
            int code = Bins.uint16(memory, pc);
            int opcode = code >> 12;// 指令码
            //dataType = (code & 0x0F00) >> 8;
            dataType = Values.fromValue(DataType.class, (code & 0x0F00) >> 8);

            specialByte = (code & 0x00F0) >> 4;
            addressingType = code & 0x000F;
            instruction = Values.fromValue(Instruction.class, opcode);

            op1 = operand(addressingType % 4, Bins.int32(memory, pc + 2));
            op2 = operand(addressingType / 4, Bins.int32(memory, pc + 6));
        }

        /**
         * <pre>
         * rx		| 0x0 | 寄存器寻址
         * [rx]	| 0x1 | 寄存器间接寻址
         * n		| 0x2 | 立即数寻址
         * [n]	| 0x3 | 间接寻址
         * </pre>
         */
        private Operand operand(int type, int op)
        {
            switch (type)
            {
                case 0:
                    return Operand.holder(register(op));
                case 1:
                    return Operand.indirect(register(op), memory);
                case 2:
                    return Operand.value(op);
                case 3:
                    return Operand.address(op, memory);
                default:
                    throw unsupport("未知的寻址类型: %s 操作数为: %s", type, op);
            }
        }

    }
}
