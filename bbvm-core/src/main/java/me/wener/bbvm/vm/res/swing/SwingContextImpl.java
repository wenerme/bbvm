package me.wener.bbvm.vm.res.swing;

import me.wener.bbvm.vm.res.ImageManager;
import me.wener.bbvm.vm.res.InputManager;
import me.wener.bbvm.vm.res.PageManager;

import javax.swing.*;

/**
 * @author wener
 * @since 15/12/26
 */
class SwingContextImpl implements SwingContext {
    private JFrame frame;

    @Override
    public JFrame getFrame() {
        if (frame == null) {
            synchronized (this) {
                if (frame == null) {
                    frame = createFrame();
                }
            }
        }
        return frame;
    }

    @Override
    public PageManager getPageManager() {
        return null;
    }

    @Override
    public ImageManager getImageManager() {
        return null;
    }

    @Override
    public InputManager getInputManager() {
        return null;
    }

    protected JFrame createFrame() {
        return null;
    }
}
