package me.wener.bbvm.impl.spi;

import me.wener.bbvm.api.FileHandle;
import me.wener.bbvm.utils.Bins;

public abstract class AbstractFileHandle implements FileHandle
{
    // vm 环境默认是 le
    boolean bigEndian = false;

    public boolean isBigEndian()
    {
        return bigEndian;
    }

    public void setBigEndian(boolean bigEndian)
    {
        this.bigEndian = bigEndian;
    }

    public short readShort()
    {
        byte[] bytes = {0, 0};
        readBytes(bytes, 0, bytes.length);

        if (bigEndian)
            return Bins.int16b(bytes, 0);
        else
            return Bins.int16l(bytes, 0);
    }


    public int readInt()
    {
        byte[] bytes = {0, 0, 0, 0};
        readBytes(bytes, 0, bytes.length);

        if (bigEndian)
            return Bins.int32b(bytes, 0);
        else
            return Bins.int32l(bytes, 0);
    }

    public float readFloat()
    {
        return Bins.float32(readInt());
    }

    public void writeShort(short v)
    {
        byte[] bytes = {0, 0, 0, 0};

        if (bigEndian)
            Bins.int16b(bytes, 0, v);
        else
            Bins.int16l(bytes, 0, v);

        writeBytes(bytes, 0, bytes.length);
    }

    public void writeInt(int v)
    {
        byte[] bytes = {0, 0, 0, 0};

        if (bigEndian)
            Bins.int32b(bytes, 0, v);
        else
            Bins.int32l(bytes, 0, v);

        writeBytes(bytes, 0, bytes.length);
    }

    public void writeFloat(float v)
    {
        writeInt(Bins.int32(v));
    }
}
