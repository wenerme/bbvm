package me.wener.bbvm.system;

import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import lombok.Data;
import lombok.experimental.Accessors;


@Data
@Accessors(chain = true)
public class Memory implements Resettable
{
    private static final Charset DEFAULT_CHARSET = Charset.forName("GBK");
    private ByteBuffer mem;
    private Charset charset = DEFAULT_CHARSET;

    byte read(short address)
    {
        return 0;
    }

    void write(short address, byte val)
    {
    }

    public int readInt(int pos)
    {
        return mem.getInt(pos);
    }

    public String readString(int pos)
    {
        return readString(pos, charset);
    }

    public String readString(int pos, Charset charset)
    {
        return null;
    }

    public void writeInt(int pos, int value)
    {
        mem.putInt(pos, value);
    }

    public void load(byte[] content)
    {
        mem.put(content);
    }

    @Override
    public void reset()
    {
        mem.clear();
    }
}
