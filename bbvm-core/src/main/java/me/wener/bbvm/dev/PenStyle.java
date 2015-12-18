package me.wener.bbvm.dev;

import me.wener.bbvm.util.IsInt;

public enum PenStyle implements IsInt {
    PEN_SOLID(0), PEN_DASH(1);
    private final int value;

    PenStyle(int value) {
        this.value = value;
    }

    public int asInt() {
        return value;
    }
}
