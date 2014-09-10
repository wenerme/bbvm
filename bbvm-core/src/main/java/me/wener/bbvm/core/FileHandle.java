package me.wener.bbvm.core;

public interface FileHandle extends AutoCloseable
{
    void open(String path);

    void close();

    byte readByte();

    /**
     * @return 返回读取的实际长度
     */
    int readBytes(byte[] bytes, int index, int len);

    short readShort();

    int readInt();

    float readFloat();

    void writeByte(byte v);

    void writeBytes(byte[] bytes, int index, int len);

    void writeShort(short v);

    void writeInt(int v);

    void writeFloat(float v);

    /**
     * 判断文件是否 EOF
     */
    boolean isEOF();

    /**
     * 获取文件长度
     */
    int length();

    /**
     * 获取文件指针偏移量
     */
    int offset();

    /**
     * 设置文件指针偏移量
     */
    void offset(int address);

}
