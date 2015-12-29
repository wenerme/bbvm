package me.wener.bbvm.dev;

import java.io.IOException;

/**
 * @author wener
 * @since 15/12/17
 */
public interface FileResource extends Resource {

    @Override
    FileManager getManager();

    FileResource open(String string) throws IOException;

    default int readInt(int address) throws IOException {
        return seek(address).readInt();
    }

    default String readString(int address) throws IOException {
        return seek(address).readString();
    }

    int readInt() throws IOException;

    String readString() throws IOException;

    FileResource writeInt(int address, int v) throws IOException;

    FileResource writeString(int address, String v) throws IOException;

    FileResource writeInt(int v) throws IOException;

    FileResource writeString(String v) throws IOException;

    boolean isEof() throws IOException;

    int length() throws IOException;

    int tell() throws IOException;

    FileResource seek(int i) throws IOException;

    @Override
    void close();
}
