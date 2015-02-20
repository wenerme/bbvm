package me.wener.bbvm.system;

import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import lombok.Getter;
import lombok.Setter;
import lombok.experimental.Accessors;


@Accessors(chain = true, fluent = true)
public class VmMemory implements Memory
{
    private static final Charset GBK_CHARSET = Charset.forName("GBK");
    @Getter
    private ByteBuffer buffer = ByteBuffer.allocate(1024);
    @Getter
    @Setter
    private Charset charset = GBK_CHARSET;

    byte read(int address)
    {
        return buffer.get(address);
    }

    void write(int address, byte val)
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

    @Override
    public String readString(int pos, Charset charset)
    {
        return null;
    }

    @Override
    public void writeInt(int pos, int value)
    {
        buffer.putInt(pos, value);
    }

    public void load(byte[] content)
    {
        buffer.put(content);
    }

    @Override
    public void reset()
    {
        buffer.clear();
    }
}
