package me.wener.bbvm.dev.swing;

import me.wener.bbvm.dev.InputManager;

import java.util.concurrent.BlockingQueue;

/**
 * @author wener
 * @since 15/12/26
 */
class QueuedInputManager implements InputManager {
    final protected BlockingQueue<InputEvent> events;

    QueuedInputManager(BlockingQueue<InputEvent> events) {
        this.events = events;
    }

    @Override
    public boolean isKeyPressed(int key) {
        return false;
    }

    @Override
    public int waitKey() {
        return events.poll().getKeyCode();
    }

    @Override
    public int getLastKeyCode() {
        return 0;
    }

    @Override
    public String getLastKeyString() {
        return null;
    }

    @Override
    public String readText() {
        return null;
    }

}
