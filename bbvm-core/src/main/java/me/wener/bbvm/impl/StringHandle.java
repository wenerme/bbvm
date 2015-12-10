package me.wener.bbvm.impl;

import me.wener.bbvm.util.val.StringHolder;
import me.wener.bbvm.util.val.impl.SimpleValueHolder;

/**
 * 字符串句柄
 */
public class StringHandle extends SimpleValueHolder<String> implements StringHolder
{
    public static StringHandle valueOf(String v)
    {
        return new ReadOnlyStringHandle(v);
    }

    public StringHandle concat(StringHandle o)
    {
        set(get()+o.get());
        return this;
    }
    public StringHandle concat(String o)
    {
        set(asInt() + o);
        return this;
    }
    static class ReadOnlyStringHandle extends StringHandle
    {
        ReadOnlyStringHandle(String v)
        {
            super.set(v);
        }

        @Override
        public void set(String v)
        {
            throw new UnsupportedOperationException();
        }
    }
}
