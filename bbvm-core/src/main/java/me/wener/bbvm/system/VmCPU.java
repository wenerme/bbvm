package me.wener.bbvm.system;

import static me.wener.bbvm.neo.inst.def.CompareTypes.*;

import com.google.common.collect.Maps;
import java.io.IOException;
import java.io.Writer;
import java.nio.ByteBuffer;
import java.util.Map;
import lombok.Getter;
import lombok.Setter;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.system.api.*;
import me.wener.bbvm.utils.Bins;
import me.wener.bbvm.utils.val.Values;

@Accessors(chain = true, fluent = true)
@Slf4j
public class VmCPU extends OpStatusImpl implements CPU, VmStatus, Defines
{
    @Getter
    private final RegisterImpl.MonitoredRegister rp = RegisterImpl.monitor(new RegisterImpl("rp"));
    @Getter
    private final RegisterImpl rb = new RegisterImpl("rb");
    @Getter
    private final RegisterImpl rs = new RegisterImpl("rs");
    @Getter
    private final RegisterImpl rf = new RegisterImpl("rf");
    @Getter
    private final RegisterImpl r0 = new RegisterImpl("r0");
    @Getter
    private final RegisterImpl r1 = new RegisterImpl("r1");
    @Getter
    private final RegisterImpl r2 = new RegisterImpl("r2");
    @Getter
    private final RegisterImpl r3 = new RegisterImpl("r3");

    @Getter
    private VmMemory memory = new VmMemory();
    @Getter
    @Setter
    private Writer asmWriter;
    @Getter
    private boolean exit;
    private Map<String, ResourcePool> resources = Maps.newConcurrentMap();
    private boolean isJumped = false;

    public VmCPU()
    {
        // 初始化资源池
        // 字符串的句柄为负
        resources.put(RES_STRING, new ResourcePool()
        {
            @Override
            protected int next()
            {
                // 第一个句柄为 -1
                return handler.decrementAndGet();
            }
        });

        resources.put(RES_FILE, new ResourcePool());
        resources.put(RES_PAGE, new ResourcePool());
        resources.put(RES_RES, new ResourcePool());
        rp.listeners().add(new RegisterImpl.RegisterChangeListener()
        {
            @Override
            public void onChange(Register register, Integer val)
            {
                isJumped = true;
            }
        });

        reset();
    }

    public void reset()
    {
        for (ResourcePool pool : resources.values())
        {
            try
            {
                pool.close();
            } catch (IOException e)
            {
                log.warn("Close resource pool failed.", e);
            }
        }
        memory.reset();
        rp.set(0);
        rb.set(0);
        rs.set(0);
        rf.set(0);
        r0.set(0);
        r1.set(0);
        r2.set(0);
        r3.set(0);
        exit = false;
    }

    public boolean step()
    {
        if (!exit)
        {
            isJumped = false;
            readInstruction();
            if (asmWriter != null)
            {
                try
                {
                    asmWriter.write(toAssembly());
                } catch (IOException e)
                {
                    log.info("Write asm failed.", e);
                }
            }
            process();
            // 如果没有跳转,则正常增加
            if (!isJumped)
            {
                rp.set(rp.get() + opcode.length());
            }
        }

        return exit;
    }

    private void readInstruction()
    {
        /*
            指令码 + 数据类型 + 特殊用途字节 + 寻址方式 + 第一个操作数 + 第二个操作数
         0x 0       0         0           0        00000000     00000000
        */
        ByteBuffer buf = memory.buffer();
        buf.position(rp.get()).mark();
        short first = Bins.unsigned(buf.get());
        opcode = Values.fromValue(Opcode.class, first >> 4);
        switch (opcode)
        {
            case RET:
            case NOP:
            case EXIT:
            {
                // 无操作数
            }
            break;
            case LD:
            case IN:
            case OUT:
            case CMP:
            {
                // 两个操作数
//                int dataType = first & 0xf;
//                short second = memory.readUnsignedByte();
//                int special = second >> 4;
//                int addressingMode = second & 0xf;
//
//                TowOperandInst tow = (TowOperandInst) inst;
//                tow.a = readOperand(ctx, addressingMode / 4);
//                tow.b = readOperand(ctx, addressingMode % 4);
//                tow.dataType = dataType;
                dataType = Values.fromValue(DataType.class, first & 0xf);
                short second = Bins.unsigned(buf.get());
                int special = second >> 4;
                int addressingMode = second & 0xf;
                a.addressingMode(Values.fromValue(AddressingMode.class, addressingMode / 4));
                b.addressingMode(Values.fromValue(AddressingMode.class, addressingMode % 4));
                readOperand(a);
                readOperand(b);
            }
            break;
            default:
                throw new UnsupportedOperationException();
        }
        buf.reset();
    }

    private void readOperand(OperandImpl o)
    {
        int v = memory.buffer().getInt();
        switch (o.addressingMode())
        {
            case REGISTER:
            case REGISTER_DEFERRED:
                o.indirect(register(Values.fromValue(RegisterType.class, v)));
                break;
            case IMMEDIATE:
            case DIRECT:
                o.value(v);
                break;
        }
    }

    private void process()
    {
        switch (opcode)
        {
            case NOP:
                NOP();
                break;
            case LD:
                LD();
                break;
            case PUSH:
                PUSH();
                break;
            case POP:
                POP();
                break;
            case IN:
                IN();
                break;
            case OUT:
                OUT();
                break;
            case JMP:
                JMP();
                break;
            case JPC:
                JPC();
                break;
            case CALL:
                CALL();
                break;
            case RET:
                RET();
                break;
            case CMP:
                CMP();
                break;
            case CAL:
                CAL();
                break;
            case EXIT:
                EXIT();
                break;
        }
        throw new UnsupportedOperationException();
    }

    void NOP()
    {
    }

    void LD()
    {
        switch (dataType)
        {
            case T_DWORD:
            case T_FLOAT:
            case T_INT:
                a.set(b.get());
                break;
            case T_BYTE:
                a.set((a.get() & 0xffffff00) | (b.get() & 0xff));
                break;
            case T_WORD:
                a.set((a.get() & 0xffff0000) | (b.get() & 0xffff));
                break;
            default:
                throw new AssertionError("未知的数据类型:" + dataType);
        }
    }

    void PUSH()
    {
        push(a.get());
    }

    void POP()
    {
        a.set(pop());
    }

    void IN()
    {
    }

    void OUT()
    {
    }

    void JMP()
    {
        rp.set(a.get());
    }

    void JPC()
    {
        // JPC 的数据类型为比较操作
        CompareType org = compareType;
        CompareType flag = Values.fromValue(CompareType.class, rf.get());


        if (CompareType.isMatch(org, flag))
            rp.set(a.get());
    }

    void CALL()
    {
        push(rp.get() + opcode.length());
        rp.set(a.get());
    }

    void RET()
    {
        rp.set(pop());
    }

    void CMP()
    {
        float oa;
        float ob;

        if (dataType == DataType.T_FLOAT)
        {
            oa = Bins.float32(a.get());
            ob = Bins.float32(b.get());
        } else
        {
            oa = a.get();
            ob = b.get();
        }
        float oc = oa - ob;
        if (oc > 0)
            rf.set(A);
        else if (oc < 0)
            rf.set(B);
        else
            rf.set(Z);
    }

    void CAL()
    {
        // 返回结果为 r0
        double oa;
        double ob;
        double oc;

        if (dataType == DataType.T_FLOAT)
        {
            oa = Bins.float32(a.get());
            ob = Bins.float32(b.get());
        } else
        {
            oa = a.get();
            ob = b.get();
        }
        switch (calculateType)
        {
            case ADD:
                oc = oa + ob;
                break;
            case DIV:
                oc = oa / ob;
                break;
            case MOD:
                oc = oa % ob;
                break;
            case MUL:
                oc = oa * ob;
                break;
            case SUB:
                oc = oa - ob;
                break;
            default:
                throw new AssertionError("未知计算操作: " + calculateType);
        }
        int ret = (int) oc;
        // 值返回归约
        switch (dataType)
        {
            case T_FLOAT:
                ret = Bins.int32((float) oc);
                break;
            case T_BYTE:
                ret &= 0xff;
                break;
            case T_WORD:
                ret &= 0xffff;
                break;
        }
        a.set(ret);
    }

    void EXIT()
    {
        exit = true;
    }

    public void push(int v)
    {
        // FIXME
        memory.writeInt(rf.get(), v);
        rf.set(rf.get() + 4);
    }

    public int pop()
    {
        rf.set(rf.get() - 4);
        return memory.readInt(rf.get());
    }

    @Override
    public OpStatus opstatus()
    {
        return this;
    }

    @Override
    public VmStatus vmstatus()
    {
        return this;
    }

    @Override
    public Register register(RegisterType type)
    {
        switch (type)
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
        }
        throw new UnsupportedOperationException();
    }

    @Override
    public ResourcePool resources(String resourceName)
    {
        return resources.get(resourceName);
    }

    @Override
    public Map<String, ResourcePool> resources()
    {
        return resources;
    }

    public void load(byte[] content)
    {
        memory.load(content);
    }
}
