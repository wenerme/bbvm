package me.wener.bbvm.system;

import java.io.Closeable;
import java.io.IOException;
import lombok.Getter;
import lombok.Setter;
import lombok.experimental.Accessors;

@Accessors(chain = true, fluent = true)
public class Resource
        implements me.wener.bbvm.system.api.Resource
{
    @Getter
    @Setter
    private int handler = Integer.MAX_VALUE;

    private volatile Object value;

    @SuppressWarnings("unchecked")
    @Override
    public <T> T as()
    {
        return (T) value;
    }


    @Override
    public boolean isNull()
    {
        return value == null;
    }

    @Override
    public Object get()
    {
        return value;
    }

    @Override
    public void set(Object v)
    {
        value = v;
    }

    @Override
    public void close() throws IOException
    {
        if (value instanceof Closeable)
        {
            ((Closeable) value).close();
        }
    }
}
