package me.wener.bbvm.system.internal;

import java.io.Closeable;
import java.io.IOException;
import lombok.Getter;
import lombok.Setter;
import lombok.experimental.Accessors;

@Accessors(chain = true, fluent = true)
class DefaultResource
        implements me.wener.bbvm.system.Resource
{
    @Getter
    @Setter
    private Integer handler = null;

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
        return handler == null || value == null;
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
