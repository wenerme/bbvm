package me.wener.bbvm.dev;

import me.wener.bbvm.util.IsInt;

public enum FontType implements IsInt {
    FONT_12SONG(0),
    FONT_12KAI(1),
    FONT_12HEI(2),
    FONT_16SONG(3),
    FONT_16KAI(4),
    FONT_16HEI(5),
    FONT_24SONG(6),
    FONT_24KAI(7),
    FONT_24HEI(8);
    private final int value;

    FontType(int value) {
        this.value = value;
    }

    public int asInt() {
        return value;
    }
}
