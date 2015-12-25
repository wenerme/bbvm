package me.wener.bbvm.vm.res.swing;

import com.google.inject.AbstractModule;
import com.google.inject.Module;
import me.wener.bbvm.dev.FontType;
import me.wener.bbvm.exception.ResourceMissingException;
import me.wener.bbvm.util.IntEnums;
import me.wener.bbvm.vm.res.ImageManager;
import me.wener.bbvm.vm.res.InputManager;
import me.wener.bbvm.vm.res.PageManager;
import me.wener.bbvm.vm.res.ResourceManager;

import java.awt.*;
import java.awt.event.KeyEvent;
import java.awt.event.KeyListener;
import java.awt.event.MouseEvent;
import java.awt.event.MouseListener;
import java.util.concurrent.ArrayBlockingQueue;

/**
 * @author wener
 * @since 15/12/18
 */
public class Swings {
    static {
        IntEnums.cache(FontType.class);
    }

    public static SwingContext createContext() {
        return new SwingContextImpl();
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
        Input input = new Input();
        component.addKeyListener(input);
        component.addMouseListener(input);
        return input;
    }

    static class SwingModule extends AbstractModule {
        @Override
        protected void configure() {
            bind(PageManager.class).to(SwingPageManager.class);
            bind(ImageManager.class).to(SwingImageManager.class);
        }
    }

    static class Input extends QueuedInputManager implements MouseListener, KeyListener {


        public Input() {
            super(new ArrayBlockingQueue<>(64));
        }

        @Override
        public void mouseClicked(MouseEvent e) {
            try {
                events.put(new InputEvent(InputEvent.Type.CLICK, e.getX() | e.getY()));
            } catch (InterruptedException e1) {
                e1.printStackTrace();
            }
        }

        @Override
        public void mousePressed(MouseEvent e) {

        }

        @Override
        public void mouseReleased(MouseEvent e) {

        }

        @Override
        public void mouseEntered(MouseEvent e) {

        }

        @Override
        public void mouseExited(MouseEvent e) {

        }

        @Override
        public void keyTyped(KeyEvent e) {

        }

        @Override
        public void keyPressed(KeyEvent e) {
            try {
                if (e.getKeyChar() == KeyEvent.CHAR_UNDEFINED) {
                    events.put(new InputEvent(InputEvent.Type.DOWN, e.getKeyCode()));
                } else {
                    events.put(new InputEvent(InputEvent.Type.DOWN, e.getKeyChar()));
                }
            } catch (InterruptedException e1) {
                e1.printStackTrace();
            }
        }

        @Override
        public void keyReleased(KeyEvent e) {
            try {
                if (e.getKeyChar() == KeyEvent.CHAR_UNDEFINED) {
                    events.put(new InputEvent(InputEvent.Type.UP, e.getKeyCode()));
                } else {
                    events.put(new InputEvent(InputEvent.Type.UP, e.getKeyChar()));
                }
            } catch (InterruptedException e1) {
                e1.printStackTrace();
            }
        }
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
