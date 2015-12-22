package me.wener.bbvm.dev;

import me.wener.bbvm.util.IsInt;

public enum FontType implements IsInt {
    FONT_12SONG(0, 12),
    FONT_12KAI(1, 12),
    FONT_12HEI(2, 12),
    FONT_16SONG(3, 16),
    FONT_16KAI(4, 16),
    FONT_16HEI(5, 16),
    FONT_24SONG(6, 24),
    FONT_24KAI(7, 24),
    FONT_24HEI(8, 24);
    private final int value;
    private final int size;

    FontType(int value, int size) {
        this.value = value;
        this.size = size;
    }

    public int asInt() {
        return value;
    }

    public int getSize() {
        return size;
    }
}
