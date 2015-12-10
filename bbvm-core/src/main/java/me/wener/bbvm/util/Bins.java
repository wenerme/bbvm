package me.wener.bbvm.util;

import java.nio.charset.Charset;

/**
 * 参考实现 {@link java.nio.Bits}, 和 Netty 的ByteBuf 和 ByteBuffer<br>
 * 后缀 b 和 l 分别代表 big-endian 和 little-endian
 * <p/>
 * 在书写的时候,用了一些额外的 +0 和多余的括号,为了使代码对齐,可能是强迫症,得治.
 */
@SuppressWarnings({"unused", "PointlessArithmeticExpression"})
public class Bins
{
    public static void main(String[] args)
    {
        int i = 0x12345678;
        assert Bins.int4(i, 0) == 8;
        assert Bins.int4(i, 1) == 7;
        assert Bins.int4(i, 2) == 6;
        assert Bins.int4(i, 3) == 5;
        assert Bins.int4(i, 4) == 4;
        assert Bins.int4(i, 5) == 3;
        assert Bins.int4(i, 6) == 2;
        assert Bins.int4(i, 7) == 1;
        assert Bins.int4(i, 8) == 8;
    }

    /**
     * 返回一个范围在 0xf 内的值
     *
     * @param i 索引, 从右到左, 该索引会mod到8
     */
    public static byte int4(int v, int i)
    {
        i %= 8;
        return (byte) (v >> (i * 4) & 0xf);
    }

    public static byte int4(long v, int i)
    {
        i %= 16;
        return (byte) (v >> (i * 4) & 0xf);
    }

    public static void int8(byte[] bytes, int offset, byte v)
    {
        bytes[offset] = v;
    }

    public static byte int8(byte[] bytes, int offset)
    {
        return bytes[offset];
    }

    public static short uint8(byte[] bytes, int offset)
    {
        return (short) (int8(bytes, offset) & 0xFF);
    }

    public static char[] char16(int code)
    {
        return Character.toChars(code);
    }

    public static char char16(short code)
    {
        return Character.toChars(code)[0];
    }

    /**
     * 获取一个两字节的 UTF-16 字符
     */
    public static char char16b(byte[] bytes, int offset)
    {
        return (char) int16b(bytes, offset);
    }

    public static void char16b(byte[] bytes, int offset, char v)
    {
        bytes[offset + 0] = char0(v);
        bytes[offset + 1] = char1(v);
    }

    public static void int16b(byte[] bytes, int offset, short v)
    {
        bytes[offset + 0] = short0(v);
        bytes[offset + 1] = short1(v);
    }

    public static short int16b(byte[] bytes, int offset)
    {
        return (short) (bytes[offset] << 8 | bytes[offset + 1] & 0xFF);
    }

    public static void int16l(byte[] bytes, int offset, short v)
    {
        bytes[offset + 1] = short0(v);
        bytes[offset + 0] = short1(v);
    }

    public static short int16l(byte[] bytes, int offset)
    {
        return (short) (bytes[offset + 1] << 8 | bytes[offset + 0] & 0xFF);
    }

    public static int uint16l(byte[] bytes, int offset)
    {
        return int16l(bytes, offset) & 0xFFFF;
    }

    public static int uint16b(byte[] bytes, int offset)
    {
        return int16b(bytes, offset) & 0xFFFF;
    }

    public static int int16(byte[] bytes, int offset, boolean be)
    {
        if (be)
            return int16b(bytes, offset);
        else
            return int16l(bytes, offset);
    }

    public static int int24b(byte[] bytes, int offset)
    {
        int value = uint24b(bytes, offset);
        if ((value & 0x800000) != 0)
        {
            value |= 0xff000000;
        }
        return value;
    }

    public static int uint24b(byte[] bytes, int offset)
    {
        return ((bytes[offset + 0] & 0xff)) << 16 |
                (bytes[offset + 1] & 0xff) << 8 |
                (bytes[offset + 2] & 0xff);
    }

    public static int uint24l(byte[] bytes, int offset)
    {
        return ((bytes[offset + 2] & 0xff)) << 16 |
                (bytes[offset + 1] & 0xff) << 8 |
                (bytes[offset + 0] & 0xff);
    }

    public static void int32b(byte[] bytes, int offset, int v)
    {
        bytes[offset + 0] = int3(v);
        bytes[offset + 1] = int2(v);
        bytes[offset + 2] = int1(v);
        bytes[offset + 3] = int0(v);
    }

    public static int int32b(byte[] bytes, int offset)
    {
        return ((bytes[offset + 0] & 0xff)) << 24 |
                (bytes[offset + 1] & 0xff) << 16 |
                (bytes[offset + 2] & 0xff) << 8 |
                (bytes[offset + 3] & 0xff);
    }

    public static int int32(byte[] bytes, int offset, boolean be)
    {
        if (be)
            return int32b(bytes, offset);
        else
            return int32l(bytes, offset);
    }

    public static void int32l(byte[] bytes, int offset, int v)
    {
        bytes[offset + 3] = int3(v);
        bytes[offset + 2] = int2(v);
        bytes[offset + 1] = int1(v);
        bytes[offset + 0] = int0(v);
    }

    public static int int32l(byte[] bytes, int offset)
    {
        return ((bytes[offset + 3] & 0xff)) << 24 |
                (bytes[offset + 2] & 0xff) << 16 |
                (bytes[offset + 1] & 0xff) << 8 |
                (bytes[offset + 0] & 0xff);
    }

    public static int int32(float v)
    {
        return Float.floatToRawIntBits(v);
    }

    public static long uint32b(byte[] bytes, int offset)
    {
        return int32b(bytes, offset) & 0xFFFFFFFFL;
    }

    public static long uint32l(byte[] bytes, int offset)
    {
        return int32l(bytes, offset) & 0xFFFFFFFFL;
    }

    public static void int64b(byte[] bytes, int offset, long v)
    {
        bytes[offset + 0] = long0(v);
        bytes[offset + 1] = long1(v);
        bytes[offset + 2] = long2(v);
        bytes[offset + 3] = long3(v);
        bytes[offset + 4] = long4(v);
        bytes[offset + 5] = long5(v);
        bytes[offset + 6] = long6(v);
        bytes[offset + 7] = long7(v);
    }

    public static long int64b(byte[] bytes, int offset)
    {
        return (((long) bytes[offset + 0] & 0xff)) << 56 |
                ((long) bytes[offset + 1] & 0xff) << 48 |
                ((long) bytes[offset + 2] & 0xff) << 40 |
                ((long) bytes[offset + 3] & 0xff) << 32 |
                ((long) bytes[offset + 4] & 0xff) << 24 |
                ((long) bytes[offset + 5] & 0xff) << 16 |
                ((long) bytes[offset + 6] & 0xff) << 8 |
                ((long) bytes[offset + 7] & 0xff);
    }

    public static long uint64(byte[] bytes, int offset)
    {
        // TODO
        return 0;
    }

    public static float float32b(byte[] bytes, int offset)
    {
        return float32(int32b(bytes, offset));
    }

    public static void float32b(byte[] bytes, int offset, float v)
    {
        int32b(bytes, offset, int32(v));
    }

    public static float float32l(byte[] bytes, int offset)
    {
        return float32(int32l(bytes, offset));
    }

    public static void float32l(byte[] bytes, int offset, float v)
    {
        int32l(bytes, offset, int32(v));
    }

    public static float float32(int v)
    {
        return Float.intBitsToFloat(v);
    }

    public static long int64(double v)
    {
        return Double.doubleToRawLongBits(v);
    }

    public static double double64b(byte[] bytes, int offset)
    {
        return Double.longBitsToDouble(int64b(bytes, offset));
    }

    public static void double64b(byte[] bytes, int offset, double v)
    {
        int64b(bytes, offset, int64(v));
    }

    /**
     * 读取以0结尾的字符,使用指定编码将其转换为字符串
     */
    public static String zString(byte[] bytes, int offset, Charset charset)
    {
        return string(zByte(bytes, offset), charset);
    }

    /**
     * 返回byte数组,从bytes中读取,直到遇到0
     */
    public static byte[] zByte(byte[] bytes, int offset)
    {
        int len = 0;
        for (int i = offset; i < bytes.length; i++)
        {
            if (bytes[i] == 0)
                break;
            len++;
        }
        byte[] out = new byte[len];
        rBytes(bytes, offset, out, 0, len);
        return out;
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

    // region unsigned
    public static short unsigned(byte v)
    {
        return (short) (v & 0xFF);
    }

    public static int unsigned(short v)
    {
        return v & 0xFFFF;
    }

    public static long unsigned(int v)
    {
        return v & 0xFFFFFFFFL;
    }
    // endregion

    //region 位逆序操作
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

    //endregion

    public static int makeInt(byte b3, byte b2, byte b1, byte b0)
    {
        return (((b3) << 24) |
                ((b2 & 0xff) << 16) |
                ((b1 & 0xff) << 8) |
                ((b0 & 0xff)));
    }

    //region 基本类型的byte获取操作
    public static byte int3(int x)
    {
        return (byte) (x >> 24);
    }

    public static byte int2(int x) { return (byte) (x >> 16); }

    public static byte int1(int x) { return (byte) (x >> 8); }

    public static byte int0(int x) { return (byte) (x); }

    public static byte long7(long x) { return (byte) (x >> 56); }

    public static byte long6(long x) { return (byte) (x >> 48); }

    public static byte long5(long x) { return (byte) (x >> 40); }

    public static byte long4(long x) { return (byte) (x >> 32); }

    public static byte long3(long x) { return (byte) (x >> 24); }

    public static byte long2(long x) { return (byte) (x >> 16); }

    public static byte long1(long x) { return (byte) (x >> 8); }

    public static byte long0(long x) { return (byte) (x); }

    public static byte char1(char x) { return (byte) (x >> 8); }

    public static byte char0(char x) { return (byte) (x); }

    public static byte short1(short x) { return (byte) (x >> 8); }

    public static byte short0(short x) { return (byte) (x); }

    //endregion
}
