package me.wener.bbvm.dev;

import me.wener.bbvm.util.val.IsInt;

public enum DrawMode implements IsInt {
    KEY_COLOR(1);
    private final int value;

    DrawMode(int value) {
        this.value = value;
    }

    public int asInt() {
        return value;
    }
}
