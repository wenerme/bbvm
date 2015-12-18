package me.wener.bbvm.dev;

import me.wener.bbvm.util.IsInt;

public enum BrushStyle implements IsInt {
    BRUSH_SOLID(0);
    private final int value;

    BrushStyle(int value) {
        this.value = value;
    }

    public int asInt() {
        return value;
    }
}
