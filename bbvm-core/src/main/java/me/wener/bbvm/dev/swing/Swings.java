package me.wener.bbvm.dev.swing;

import com.google.inject.Module;
import me.wener.bbvm.dev.DeviceConstants;
import me.wener.bbvm.dev.InputManager;
import me.wener.bbvm.dev.ResourceManager;
import me.wener.bbvm.exception.ResourceMissingException;
import me.wener.bbvm.util.IntEnums;

import java.awt.*;
import java.awt.event.KeyEvent;

/**
 * @author wener
 * @since 15/12/18
 */
public class Swings implements DeviceConstants {
    static {
        IntEnums.cache(FontType.class);
    }

    public static SwingContext createContext() {
        return new SwingContextImpl();
    }

    public static Module module() {
        return new SwingModule();
    }

    public static Class<? extends SwingContext> swingContext() {
        return SwingContextImpl.class;
    }

    public static Module graphModule() {
        return new SwingModule();
    }

    static Module graphModule(SwingContext context) {
        return new SwingModule();
    }

    static <T> T checkMissing(ResourceManager mgr, int handler, T v) {
        if (v == null) {
            throw new ResourceMissingException(String.format("%s #%s not exists", mgr.getType(), handler), handler);
        }
        return v;
    }

    static InputManager bind(Component component) {
        SwingInputManger input = new SwingInputManger();
        component.addKeyListener(input);
        component.addMouseListener(input);
        return input;
    }

    static class IsKeyPressed {
        private static boolean wPressed = false;

        public static boolean isWPressed() {
            synchronized (IsKeyPressed.class) {
                return wPressed;
            }
        }

        public static void main(String[] args) {
            KeyboardFocusManager.getCurrentKeyboardFocusManager().addKeyEventDispatcher(ke -> {
                synchronized (IsKeyPressed.class) {
                    switch (ke.getID()) {
                        case KeyEvent.KEY_PRESSED:
                            if (ke.getKeyCode() == KeyEvent.VK_W) {
                                wPressed = true;
                            }
                            break;

                        case KeyEvent.KEY_RELEASED:
                            if (ke.getKeyCode() == KeyEvent.VK_W) {
                                wPressed = false;
                            }
                            break;
                    }
                    return false;
                }
            });
        }
    }

}
