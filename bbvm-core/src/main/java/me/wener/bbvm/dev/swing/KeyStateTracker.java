package me.wener.bbvm.dev.swing;

import com.google.common.cache.Cache;
import com.google.common.cache.CacheBuilder;

import java.awt.event.InputEvent;
import java.awt.event.KeyEvent;
import java.util.concurrent.TimeUnit;

/**
 * @author wener
 * @since 15/12/27
 */
class KeyStateTracker {
    private final Cache<Integer, Boolean> cache = CacheBuilder
            .newBuilder()
            .expireAfterWrite(100, TimeUnit.MILLISECONDS)
            .build();

    public boolean isPressed(int keyCode) {
        return cache.getIfPresent(keyCode) != null;
    }

    public KeyStateTracker offer(InputEvent e) {
        switch (e.getID()) {
            case KeyEvent.KEY_PRESSED:
                cache.put(((KeyEvent) e).getKeyCode(), Boolean.TRUE);
                break;
            case KeyEvent.KEY_RELEASED:
                cache.invalidate(((KeyEvent) e).getKeyCode());
                break;
        }
        return this;
    }
}
