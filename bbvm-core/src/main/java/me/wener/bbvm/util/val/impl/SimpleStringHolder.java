package me.wener.bbvm.util.val.impl;

import me.wener.bbvm.util.val.StringHolder;

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
        return asInt();
    }

    @Override
    public char charAt(int index)
    {
        return 0;
    }

    @Override
    public CharSequence subSequence(int start, int end)
    {
        return "";
    }
}
