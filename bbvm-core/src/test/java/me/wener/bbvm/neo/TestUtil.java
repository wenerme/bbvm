package me.wener.bbvm.neo;


import com.google.common.base.Predicate;
import com.google.common.base.Strings;
import com.google.common.collect.Iterables;
import com.google.common.collect.Lists;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import java.io.IOException;
import java.nio.ByteOrder;
import java.util.List;
import java.util.regex.Pattern;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.io.HexDump;

@Slf4j
public class TestUtil
{
    public static final Pattern MATCH_HEX_DATA = Pattern
            .compile("\\s{2,}([^ \r\n]*[^0-9a-fA-F]+[^ \r\n]*)$", Pattern.MULTILINE);
    public static final Pattern MATCH_OFFSET = Pattern
            .compile("^[^\\s]+\\s", Pattern.MULTILINE);
    public static final Predicate<String> NON_NULL_OR_EMPTY = new Predicate<String>()
    {
        @Override
        public boolean apply(String input)
        {
            return !Strings.isNullOrEmpty(input);
        }
    };
    protected static boolean logDump = false;

    public static ByteBuf fromDumpBytes(String dump)
    {
        ByteBuf buf = Unpooled.buffer(20);
        String origin = dump;
        // 删除偏移值
        if (dump.startsWith("00000000 "))
            dump = MATCH_OFFSET.matcher(dump).replaceAll("");
        dump = MATCH_HEX_DATA.matcher(dump).replaceAll("");
        String[] lines = dump.split("[\n\r]+");
        for (String line : lines)
        {
            Iterable<String> iterable = Iterables.filter(Lists.newArrayList(line.split("\\s+")), NON_NULL_OR_EMPTY);
            List<String> split = Lists.newArrayList(iterable);
            // 一行最多16个
            for (int i = 0; i < split.size() && i < 16; i++)
            {
                String b = split.get(i);
                buf.writeByte(Integer.parseInt(b, 16));
            }
        }

        if (log.isTraceEnabled() || logDump)
            try
            {
                System.out.println("原始内容");
                System.out.println(origin);
                System.out.println("解析结果 长度:" + buf.readableBytes());
                HexDump.dump(buf.copy().array(), 0, System.out, 0);
            } catch (IOException e)
            {
                e.printStackTrace();
            }

        return buf.order(ByteOrder.LITTLE_ENDIAN);
    }

    public static byte[] readableToBytes(ByteBuf buf)
    {
        byte[] bytes = new byte[buf.readableBytes()];
        buf.readBytes(bytes);
        return bytes;
    }
}
