package me.wener.bbvm.def;

import me.wener.bbvm.utils.val.IsInteger;

public enum KeyCode implements IsInteger
{
    KEY_UP(38),
    KEY_DOWN(40),
    KEY_LEFT(37),
    KEY_RIGHT(39),
    KEY_SPACE(32),
    KEY_ESCAPE(27),
    KEY_ENTER(13);
    private final int value;

    KeyCode(int value)
    {
        this.value = value;
    }

    public Integer get()
    {
        return value;
    }
}
