package me.wener.bbvm.system;

import com.google.common.base.Preconditions;
import com.google.common.collect.Lists;
import com.google.common.primitives.Bytes;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.Charset;
import java.util.List;
import lombok.Getter;
import lombok.Setter;
import lombok.experimental.Accessors;
import me.wener.bbvm.system.api.Memory;

/**
 * VM 内存内容,栈在最后1K
 */
@Accessors(chain = true, fluent = true)
public class VmMemory implements Memory
{
    private static final Charset GBK_CHARSET = Charset.forName("GBK");
    @Getter
    private final ByteBuffer buffer;// 4m
    @Getter
    @Setter
    private Charset charset = GBK_CHARSET;

    @Getter
    @Setter
    private VmCPU cpu;

    /**
     * 栈大小,默认为 1K
     */
    @Getter
    @Setter
    private int stackSize = 1024;

    @Getter
    private int length;

    public VmMemory()
    {
        // BB 的位序
        buffer = ByteBuffer.allocate(1024 * 1024 * 4);
        buffer.order(ByteOrder.LITTLE_ENDIAN);
    }

    @Override
    public byte read(int address)
    {
        return buffer.get(address);
    }

    private void checkPos(int pos)
    {
        Preconditions.checkArgument(pos <= length);
    }

    @Override
    public void write(int address, byte val)
    {
        buffer.put(address, val);
    }

    @Override
    public int readInt(int pos)
    {
        return buffer.getInt(pos);
    }

    @Override
    public String readString(int pos)
    {
        return readString(pos, charset);
    }

    public byte[] readBytesUntil(int pos, int until)
    {
        int start = pos;
        byte b;
        List<Byte> bytes = Lists.newArrayList();
        for (; ; pos++)
        {
            b = buffer.get(pos);
            if (b == until)
            {
                break;
            }
            bytes.add(b);
        }
        return Bytes.toArray(bytes);
    }

    @Override
    public String readString(int pos, Charset charset)
    {
        return new String(readBytesUntil(pos, 0), charset);
    }

    @Override
    public void writeInt(int pos, int value)
    {
        buffer.putInt(pos, value);
    }

    @Override
    public void load(byte[] content)
    {
        length = content.length;
        buffer.limit(length + stackSize);
        buffer.put(content);
    }

    public void push(int v)
    {

    }

    public int pop()
    {
        return 0;
    }

    @Override
    public void reset()
    {
        buffer.clear();
    }

    public boolean hasRemaining(int pos)
    {
        return length <= pos;
    }
}
