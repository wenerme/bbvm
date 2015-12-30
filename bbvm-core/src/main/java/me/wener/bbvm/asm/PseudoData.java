package me.wener.bbvm.asm;

import com.google.common.collect.Lists;
import io.netty.buffer.ByteBuf;

import java.util.Iterator;
import java.util.List;

/**
 * @author wener
 * @since 15/12/11
 */
public class PseudoData extends Label implements Assembly {
    Token dataTypeToken;
    private List<Value> values = Lists.newArrayList();

    @Override
    public Type getType() {
        return Type.PSEUDO;
    }

    public void add(Value data) {
        values.add(data);
    }

    @Override
    public void write(ByteBuf buf) {
        for (Value value : values) {
            value.write(buf);
        }
    }

    @Override
    public String toAssembly() {
        StringBuilder sb = new StringBuilder();
        sb.append("DATA ")
            .append(name)
            .append(' ');
        if (dataTypeToken != null) {
            sb.append(dataTypeToken.image).append(' ');
        }
        for (Iterator<Value> iterator = values.iterator(); iterator.hasNext(); ) {
            sb.append(iterator.next().toAssembly());
            if (iterator.hasNext()) {
                sb.append(", ");
            }
        }
        return sb.toString();
    }

    @Override
    public int length() {
        int len = 0;
        for (Value v : values) {
            len += v.length();
        }
        return len;
    }
}
