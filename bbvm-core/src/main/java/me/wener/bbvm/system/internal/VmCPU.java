package me.wener.bbvm.system.internal;

import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import lombok.Getter;
import lombok.Setter;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.system.*;
import me.wener.bbvm.util.Bins;
import me.wener.bbvm.util.val.IntEnums;

import java.io.IOException;
import java.nio.ByteBuffer;
import java.util.LinkedList;
import java.util.Map;

@Accessors(chain = true, fluent = true)
@Slf4j
public class VmCPU extends OpStates.DefaultOpState implements CPU, VmStatus, Defines
{
    @Getter
    private final Register.MonitoredRegister rp = Register.monitor(new Register("RP"));
    @Getter
    private final Register rb = new Register("RB");
    @Getter
    private final Register rs = new Register("RS");
    @Getter
    private final Register rf = new Register("RF");
    @Getter
    private final Register r0 = new Register("R0");
    @Getter
    private final Register r1 = new Register("R1");
    @Getter
    private final Register r2 = new Register("R2");
    @Getter
    private final Register r3 = new Register("R3");
    private final LinkedList<Integer> stack = Lists.newLinkedList();
    @Getter
    private VmMemory memory = new VmMemory();
    @Getter
    @Setter
    private boolean ignoreProcess;
    @Getter
    @Accessors(fluent = false)// isExit
    private boolean exit;
    private Map<String, ResourcePool> resources = Maps.newConcurrentMap();
    private boolean isJumped = false;


    public VmCPU()
    {
        // 初始化资源池
        // 字符串的句柄为负
        resources.put(RES_STRING, new NegativeHandlerResourcePool());
        resources.put(RES_FILE, new ResourcePool());
        resources.put(RES_PAGE, new ResourcePool());
        resources.put(RES_RES, new ResourcePool());
        rp.listeners().add(new Register.RegisterChangeListener()
        {
            @Override
            public void onChange(me.wener.bbvm.system.Register register, Integer val)
            {
                isJumped = true;
            }
        });
        a = new Operand().cpu(this);
        b = new Operand().cpu(this);
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
        exit = memory.hasRemaining(rp.asInt());
        if (!exit)
        {
            isJumped = false;
            readInstruction();
            log.trace("[{}] {}", rp.asInt(), toAssembly());

            if (!ignoreProcess)
            {
                process();
            }
            // 如果没有跳转,则正常增加
            if (!isJumped)
            {
                rp.set(rp.asInt() + opcode.length());
            }
        }
        return !exit;
    }

    private void readInstruction()
    {
        ByteBuffer buf = memory.buffer();
        buf.position(rp.asInt());
        OpStates.readBinary(this, buf);
    }

    private void readOperand(me.wener.bbvm.system.Operand o)
    {
        int v = memory.buffer().getInt();
        o.value(v);
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
            default:
                throw new UnsupportedOperationException();
        }
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
                a.set(b.asInt());
                break;
            case T_BYTE:
                a.set((a.asInt() & 0xffffff00) | (b.asInt() & 0xff));
                break;
            case T_WORD:
                a.set((a.asInt() & 0xffff0000) | (b.asInt() & 0xffff));
                break;
            default:
                throw new AssertionError("未知的数据类型:" + dataType);
        }
    }

    void PUSH()
    {
        push(a.asInt());
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
        rp.set(a.asInt());
    }

    void JPC()
    {
        // JPC 的数据类型为比较操作
        CompareType org = compareType;
        CompareType flag = IntEnums.fromInt(CompareType.class, rf.asInt());


        if (CompareType.isMatch(org, flag))
            rp.set(a.asInt());
    }

    void CALL()
    {
        push(rp.asInt() + opcode.length());
        rp.set(a.asInt());
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
            oa = Bins.float32(a.asInt());
            ob = Bins.float32(b.asInt());
        } else
        {
            oa = a.asInt();
            ob = b.asInt();
        }
        float oc = oa - ob;
        if (oc > 0)
            rf.set(CompareType.A.asInt());
        else if (oc < 0)
            rf.set(CompareType.B.asInt());
        else
            rf.set(CompareType.Z.asInt());
    }

    void CAL()
    {
        // 返回结果为 R0
        double oa;
        double ob;
        double oc;

        if (dataType == DataType.T_FLOAT)
        {
            oa = Bins.float32(a.asInt());
            ob = Bins.float32(b.asInt());
        } else
        {
            oa = a.asInt();
            ob = b.asInt();
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
        stack.push(v);
//        memory.writeInt(RF.get(), v);
//        RF.set(RF.get() + 4);
    }

    public int pop()
    {
//        RF.set(RF.get() - 4);
//        return memory.readInt(RF.get());
        return stack.pop();
    }

    @Override
    public OpState opstatus()
    {
        return this;
    }

    @Override
    public VmStatus vmstatus()
    {
        return this;
    }

    @Override
    public me.wener.bbvm.system.Register register(RegisterType type)
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
        log.debug("Load content, size {}", content.length);
        memory.load(content);
    }
}
