package me.wener.bbvm.neo.codec;

import com.google.common.base.Function;
import com.google.common.collect.Maps;
import io.netty.buffer.ByteBuf;
import java.util.Map;
import me.wener.bbvm.neo.BBVMContext;
import me.wener.bbvm.neo.Operand;
import me.wener.bbvm.neo.Operands;
import me.wener.bbvm.neo.define.Flags;
import me.wener.bbvm.neo.inst.*;

/**
 * 指令编码
 */
public class InstCodec
{


    public static final Function<Integer, Inst> CREATE_FACTORY = new Function<Integer, Inst>()
    {
        @Override
        public Inst apply(Integer input)
        {
            return create(input);
        }
    };
    private Map<Integer, Inst> cache = Maps.newConcurrentMap();
    private ByteBuf memory;


    private InstCodec(ByteBuf memory)
    {
        this.memory = memory;
    }

    public static InstCodec decode(ByteBuf memory)
    {
        return new InstCodec(memory);
    }


    public static Inst create(int code)
    {
        switch (code)
        {
            case Flags.NOP:
                return new NOP();
            case Flags.LD:
                return new LD();
            case Flags.PUSH:
                return new PUSH();
            case Flags.POP:
                return new POP();
            case Flags.IN:
                return new IN();
            case Flags.OUT:
                return new OUT();
            case Flags.JMP:
                return new JMP();
            case Flags.JPC:
                return new JPC();
            case Flags.CALL:
                return new CALL();
            case Flags.RET:
                return new RET();
            case Flags.CMP:
                return new CMP();
            case Flags.CAL:
                return new CAL();
            case Flags.EXIT:
                return new EXIT();

            default:
                throw new UnsupportedOperationException("不能创建指令 " + code);
        }
    }

    public static <T extends Inst> T read(BBVMContext ctx)
    {
        return read(ctx, CREATE_FACTORY);
    }

    @SuppressWarnings({"ConstantConditions", "unchecked"})
    public static <T extends Inst> T read(BBVMContext ctx, Function<Integer, Inst> factory)
    {
        ByteBuf memory = ctx.memory();
        /*
            指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数 + 第二个操作数
         0x 0       0         0           0        00000000     00000000
        */
        short b = memory.readUnsignedByte();
        int opcode = b >> 4;


        Inst inst = factory.apply(opcode);
        switch (opcode)
        {
            case Flags.JPC:
            case Flags.POP:
            case Flags.PUSH:
            case Flags.CALL:
            {
                /*
无操作数 1byte
   指令码 + 无用
0x 0       0
一个操作数 5byte
   指令码 + 寻址方式 + 第一个操作数
0x 0       0        00000000
两个操作数 10byte
   指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数 + 第二个操作数
0x 0       0         0           0        00000000     00000000
JPC指令 6byte
   指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数
0x 0       0         0           0        00000000
        */
                int addressingMode = b & 0xf;
                // 一个操作数
                ((OneOperandInst) inst).a = readOperand(ctx, addressingMode);
            }
            break;
            case Flags.RET:
            case Flags.NOP:
            case Flags.EXIT:
            {
                // 没有操作数
            }
            break;
            case Flags.LD:
            case Flags.IN:
            case Flags.OUT:
            case Flags.JMP:
            case Flags.CMP:
            case Flags.CAL:
            {

                int dataType = b & 0xf;
                b = memory.readUnsignedByte();
                int special = b >> 4;
                int addressingMode = b & 0xf;

                TowOperandInst tow = (TowOperandInst) inst;
                tow.a = readOperand(ctx, addressingMode / 4);
                tow.b = readOperand(ctx, addressingMode % 4);
            }
            break;
        }
        return (T) inst;
    }

    public static Operand readOperand(BBVMContext ctx, int addressingMode)
    {
        ByteBuf mem = ctx.memory();
        int v = mem.readInt();
        Operand op;
        switch (addressingMode)
        {
            case Flags.REGISTER:
                op = Operands.wrap(ctx.register(v));
                break;
            case Flags.REGISTER_DEFERRED:
                op = Operands.direct(ctx.register(v), mem);
                break;
            case Flags.DIRECT:
                op = Operands.direct(v, mem);
                break;
            case Flags.IMMEDIATE:
                op = Operands.immediate(v);
                break;
            default:
                throw new AssertionError();
        }
        return op.addressingMode(addressingMode);
    }

    private static CALL readInst(ByteBuf buf, CALL inst)
    {

        return null;
    }

    public Inst retrieve(int code)
    {
        Inst inst = cache.get(code);
        // cacheaware
        if (inst != null)
            return inst;

        inst = create(code);
        cache.put(code, inst);
        return inst;
    }

    public <T extends Inst> T read()
    {
        return null;
    }
}
