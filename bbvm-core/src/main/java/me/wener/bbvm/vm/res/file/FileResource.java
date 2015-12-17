package me.wener.bbvm.vm.res.file;

import me.wener.bbvm.vm.res.Resource;

/**
 * @author wener
 * @since 15/12/17
 */
public interface FileResource extends Resource {

    FileResource open(String string);

    int readInt(int address);

    float readFloat(int address);

    String readString(int address);

    FileResource writeInt(int address, int v);

    FileResource writeFloat(int address, float v);

    FileResource writeString(int address, String v);

    boolean isEof();

    int length();

    int tell();

    FileResource seek(int i);

    @Override
    void close();
}
