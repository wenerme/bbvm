package me.wener.bbvm.dev.swing;

import javax.inject.Inject;
import javax.inject.Singleton;
import java.awt.event.KeyEvent;
import java.awt.event.KeyListener;
import java.awt.event.MouseEvent;
import java.awt.event.MouseListener;
import java.util.concurrent.ArrayBlockingQueue;

/**
 * @author wener
 * @since 15/12/26
 */
@Singleton
class SwingInputManger extends QueuedInputManager implements MouseListener, KeyListener {


    @Inject
    public SwingInputManger() {
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

    @Override
    public String readText() {
        return null;
    }
}
