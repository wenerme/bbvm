package me.wener.bbvm.utils;

import java.nio.charset.Charset;

@SuppressWarnings("unused")
public class Bins
{
    public static byte int8(byte[] bytes, int offset)
    {
        return bytes[offset];
    }

    public static short uint8(byte[] bytes, int offset)
    {
        return (short) (int8(bytes, offset) & 0xFF);
    }

    /**
     * 获取一个两字节的 UTF-16 字符
     */
    public static char char16(byte[] bytes, int offset)
    {
        return (char)int16(bytes, offset);
    }

    public static short int16(byte[] bytes, int offset)
    {
        return (short) (bytes[offset] << 8 | bytes[offset + 1] & 0xFF);
    }

    public static int uint16(byte[] bytes, int offset)
    {
        return int16(bytes, offset) & 0xFFFF;
    }

    public static int int24(byte[] bytes, int offset)
    {
        int value = uint24(bytes, offset);
        if ((value & 0x800000) != 0)
        {
            value |= 0xff000000;
        }
        return value;
    }

    public static int uint24(byte[] bytes, int offset)
    {
        return (bytes[offset] & 0xff) << 16 |
                (bytes[offset + 1] & 0xff) << 8 |
                bytes[offset + 2] & 0xff;
    }

    public static int int32(byte[] bytes, int offset)
    {
        return (bytes[offset] & 0xff) << 24 |
                (bytes[offset + 1] & 0xff) << 16 |
                (bytes[offset + 2] & 0xff) << 8 |
                bytes[offset + 3] & 0xff;
    }

    public static int int32(float v)
    {
        return Float.floatToRawIntBits(v);
    }

    public static long uint32(byte[] bytes, int offset)
    {
        return int32(bytes, offset) & 0xFFFFFFFFL;
    }

    public static long int64(byte[] bytes, int offset)
    {
        return ((long) bytes[offset] & 0xff) << 56 |
                ((long) bytes[offset + 1] & 0xff) << 48 |
                ((long) bytes[offset + 2] & 0xff) << 40 |
                ((long) bytes[offset + 3] & 0xff) << 32 |
                ((long) bytes[offset + 4] & 0xff) << 24 |
                ((long) bytes[offset + 5] & 0xff) << 16 |
                ((long) bytes[offset + 6] & 0xff) << 8 |
                (long) bytes[offset + 7] & 0xff;
    }

    public static long uint64(byte[] bytes, int offset) { return 0;}

    public static float float32(byte[] bytes, int offset)
    {
        return float32(int32(bytes, offset));
    }

    public static float float32(int v)
    {
        return Float.intBitsToFloat(v);
    }

    public static long int64(long v)
    {
        return Double.doubleToRawLongBits(v);
    }

    public static double double64(byte[] bytes, int offset)
    {
        return Double.longBitsToDouble(int64(bytes, offset));
    }

    public static void rBytes(byte[] bytes, int offset, byte[] to, int toOffset, int len)
    {
        System.arraycopy(bytes, offset, to, toOffset, len);
    }

    public static void wBytes(byte[] bytes, int offset, byte[] from, int fromOffset, int len)
    {
        System.arraycopy(from, fromOffset, bytes, offset, len);
    }

    public static byte[] bytes(String string, Charset charset)
    {
        return string.getBytes(charset);
    }

    public static byte[] bytes(char[] chars, Charset charset)
    {
        return new String(chars).getBytes(charset);
    }

    public static char[] chars(String string)
    {
        return string.toCharArray();
    }

    public static String string(byte[] bytes, Charset charset)
    {
        /*
         CharBuffer cBuffer = ByteBuffer.wrap(bytes).asCharBuffer();
         return cBuffer.toString();
         */
        return new String(bytes, charset);
    }

    public static short reverse(short x)
    {
        return Short.reverseBytes(x);
    }

    public static char reverse(char x)
    {
        return Character.reverseBytes(x);
    }

    public static int reverse(int x)
    {
        return Integer.reverseBytes(x);
    }

    public static long reverse(long x)
    {
        return Long.reverseBytes(x);
    }
}
