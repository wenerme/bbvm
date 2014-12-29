package me.wener.bbvm.def;

import me.wener.bbvm.utils.val.IsInteger;

public enum BackgroundMode implements IsInteger
{
    /**
     * 透明显示，即字体的背景颜色无效。
     */
    TRANSPARENT(1),
    /**
     * 不透明显示，即字体的背景颜色有效
     */
    OPAQUE(2);
    private final int value;

    BackgroundMode(int value)
    {
        this.value = value;
    }

    public Integer asValue()
    {
        return value;
    }
}
