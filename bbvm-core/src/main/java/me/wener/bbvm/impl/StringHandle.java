package me.wener.bbvm.impl;

import me.wener.bbvm.utils.val.SimpleValueHolder;
import me.wener.bbvm.utils.val.StringHolder;

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
        set(get()+o);
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
