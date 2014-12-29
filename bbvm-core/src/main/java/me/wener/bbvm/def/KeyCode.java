package me.wener.bbvm.def;

public enum KeyCode implements IsIntValue
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

    public Integer asValue()
    {
        return value;
    }
}
