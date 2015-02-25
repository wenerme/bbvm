package me.wener.bbvm.system;

import java.nio.charset.Charset;

public interface Memory extends Resettable
{
    byte read(int address);

    void write(int address, byte val);

    int readInt(int pos);

    String readString(int pos);

    String readString(int pos, Charset charset);

    void writeInt(int pos, int value);

    Charset charset();

    Memory charset(Charset charset);

    void load(byte[] content);
}
