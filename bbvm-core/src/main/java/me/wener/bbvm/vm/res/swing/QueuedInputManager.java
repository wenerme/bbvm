package me.wener.bbvm.vm.res.swing;

import me.wener.bbvm.vm.res.InputManager;

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

}
