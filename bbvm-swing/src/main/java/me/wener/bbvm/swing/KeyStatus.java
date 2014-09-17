package me.wener.bbvm.swing;

import com.google.common.collect.Maps;
import java.awt.KeyEventDispatcher;
import java.awt.KeyboardFocusManager;
import java.awt.event.KeyEvent;
import java.util.Map;

public class KeyStatus
{
    private static final Map<Integer, Boolean> status = Maps.newConcurrentMap();

    static
    {
        KeyboardFocusManager manager = KeyboardFocusManager.getCurrentKeyboardFocusManager();
        manager.addKeyEventDispatcher(new MyDispatcher());
    }

    private static class MyDispatcher implements KeyEventDispatcher
    {
        @Override
        public boolean dispatchKeyEvent(KeyEvent e)
        {
            if (e.getID() == KeyEvent.KEY_PRESSED)
            {
                status.put(e.getKeyCode(), true);
            } else if (e.getID() == KeyEvent.KEY_RELEASED)
            {
                status.put(e.getKeyCode(), false);
            }
            return false;
        }
    }

    public static boolean isNotPressed(int keyCode)
    {
        return !isPressed(keyCode);
    }
    public static boolean isPressed(int keyCode)
    {
        Boolean pressed = status.get(keyCode);
        return pressed == null? false: pressed;
    }
}
