package me.wener.bbvm.java;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.RandomAccessFile;
import me.wener.bbvm.core.spi.AbstractFileHandle;

public class JavaFileHandle extends AbstractFileHandle
{
    protected File file;
    protected RandomAccessFile accessFile;

    @Override
    public void open(String path)
    {
        file = new File(path);
        try
        {
            accessFile = new RandomAccessFile(file, "rw");
        } catch (FileNotFoundException e)
        {
            e.printStackTrace();
        }
    }

    @Override
    public void close()
    {
        try
        {
            accessFile.close();
        } catch (IOException e)
        {
            e.printStackTrace();
        }
    }

    @Override
    public byte readByte()
    {
        try
        {
            return accessFile.readByte();
        } catch (IOException e)
        {
            throw new RuntimeException(e);
        }
    }

    @Override
    public int readBytes(byte[] bytes, int index, int len)
    {
        try
        {
            return accessFile.read(bytes, index, len);
        } catch (IOException e)
        {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void writeByte(byte v)
    {
        try
        {
            accessFile.writeByte(v);
        } catch (IOException e)
        {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void writeBytes(byte[] bytes, int index, int len)
    {
        try
        {
            accessFile.write(bytes, index, len);
        } catch (IOException e)
        {
            throw new RuntimeException(e);
        }
    }

    @Override
    public boolean isEOF()
    {
        return length() == offset();
    }

    @Override
    public int length()
    {
        return (int) file.length();
    }

    @Override
    public int offset()
    {
        try
        {
            return (int) accessFile.getFilePointer();
        } catch (IOException e)
        {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void offset(int address)
    {
        try
        {
            accessFile.seek(address);
        } catch (IOException e)
        {
            throw new RuntimeException(e);
        }
    }
}
