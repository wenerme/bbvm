package me.wener.bbvm.system;

import com.google.common.base.Preconditions;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.Charset;
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
    private final ByteBuffer buffer = ByteBuffer.allocate(1024 * 1024 * 4);// 4m
    @Getter
    @Setter
    private Charset charset = GBK_CHARSET;

    @Getter
    private int length;

    public VmMemory()
    {
        buffer.order(ByteOrder.LITTLE_ENDIAN);
    }

    @Override
    public byte read(int address)
    {
        checkPos(address);
        return buffer.get(address);
    }

    private void checkPos(int pos)
    {
        Preconditions.checkArgument(pos <= length);
    }

    @Override
    public void write(int address, byte val)
    {
        checkPos(address);
        buffer.put(address, val);
    }

    @Override
    public int readInt(int pos)
    {
        checkPos(pos);
        return buffer.getInt(pos);
    }

    @Override
    public String readString(int pos)
    {
        return readString(pos, charset);
    }

    @Override
    public String readString(int pos, Charset charset)
    {
        checkPos(pos);
        return null;
    }

    @Override
    public void writeInt(int pos, int value)
    {
        checkPos(pos);
        buffer.putInt(pos, value);
    }

    @Override
    public void load(byte[] content)
    {
        length = content.length;
        buffer.put(content);
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
