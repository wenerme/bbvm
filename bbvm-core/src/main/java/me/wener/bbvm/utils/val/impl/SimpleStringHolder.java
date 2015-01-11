package me.wener.bbvm.utils.val.impl;

import me.wener.bbvm.utils.val.StringHolder;

public class SimpleStringHolder
        extends SimpleValueHolder<String>
        implements StringHolder, CharSequence
{

    public SimpleStringHolder(String value)
    {
        super(value);
    }

    public SimpleStringHolder()
    {
    }

    @Override
    public int length()
    {
        return get().length();
    }

    @Override
    public char charAt(int index)
    {
        return get().charAt(index);
    }

    @Override
    public CharSequence subSequence(int start, int end)
    {
        return get().subSequence(start, end);
    }
}
