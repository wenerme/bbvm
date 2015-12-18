package me.wener.bbvm.dev;

import me.wener.bbvm.util.IsInt;

public enum BackgroundMode implements IsInt {
    /**
     * 透明显示，即字体的背景颜色无效。
     */
    TRANSPARENT(1),
    /**
     * 不透明显示，即字体的背景颜色有效
     */
    OPAQUE(2);
    private final int value;

    BackgroundMode(int value) {
        this.value = value;
    }

    public int asInt() {
        return value;
    }
}
