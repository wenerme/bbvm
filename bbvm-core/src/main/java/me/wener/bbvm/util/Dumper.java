package me.wener.bbvm.util;

import com.google.common.base.Throwables;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import org.apache.commons.io.HexDump;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.util.Arrays;

/**
 * 与 Stringer 类似,但是会输出更多调试相关的信息
 */
public class Dumper
{
    // region 数据 dump


    public static String dump(ByteBuf buf)
    {
        return buf == null ? null : ByteBufUtil.hexDump(buf);
    }

    public static String hexDumpReadable(ByteBuf buf)
    {
        if (buf.readableBytes() == 0)
        {
            return "00000000 \n";
        }

        return hexDump(Arrays.copyOfRange(buf.array(), buf.readerIndex(), buf.readableBytes()), buf.readerIndex());
    }

    public static String hexDump(ByteBuf buf)
    {
        return hexDump(buf.array());
    }

    public static String hexDumpOut(byte[] buf)
    {
        String dump = hexDump(buf);
        System.out.println(dump);
        return dump;
    }

    public static String hexDumpOut(ByteBuf buf)
    {
        String dump = hexDump(buf.array());
        System.out.println(dump);
        return dump;
    }

    public static String hexDump(byte[] bytes)
    {
        return hexDump(bytes, 0);
    }

    public static String hexDump(byte[] bytes, int index)
    {
        ByteArrayOutputStream os = new ByteArrayOutputStream();
        try
        {
            HexDump.dump(bytes, 0, os, index);
        } catch (IOException e)
        {
            Throwables.propagate(e);
        }
        return new String(os.toByteArray());
    }

    // endregion

}
