package me.wener.bbvm.utils;

import com.google.common.base.Throwables;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.util.EnumSet;
import me.wener.bbvm.neo.Stringer;
import me.wener.bbvm.utils.val.IsInteger;
import me.wener.bbvm.utils.val.Values;
import org.apache.commons.io.HexDump;

/**
 * 与 Stringer 类似,但是会输出更多调试相关的信息
 */
public class Dumper
{
    // region 数据 dump
    public static <T extends Enum<T> & IsInteger> String dump(long flags, Class<T> type)
    {
        EnumSet<T> set = Values.asEnumSet(flags, type);
        return String.format("%s -> %s", Long.toBinaryString(flags), set);
    }


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

        return hexDump(buf.array(), buf.readerIndex());
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
        return Stringer.string(os.toByteArray());
    }

    // endregion

}
