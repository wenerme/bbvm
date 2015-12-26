package me.wener.bbvm.dev.swing;

import me.wener.bbvm.dev.*;

import javax.inject.Inject;
import javax.inject.Singleton;
import javax.swing.*;

/**
 * @author wener
 * @since 15/12/26
 */
@Singleton
class SwingContextImpl implements SwingContext {
    private JFrame frame;
    @Inject
    private SwingPageManager pageManager;
    @Inject
    private SwingInputManger inputManger;
    @Inject
    private SwingImageManager imageManager;
    @Inject
    private JavaFileManager fileManager;
    @Inject
    private StringManager stringManager;


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
        return pageManager;
    }

    @Override
    public ImageManager getImageManager() {
        return imageManager;
    }

    @Override
    public InputManager getInputManager() {
        return inputManger;
    }

    @Override
    public FileManager getFileManager() {
        return fileManager;
    }

    @Override
    public StringManager getStringManager() {
        return stringManager;
    }

    protected JFrame createFrame() {
        return null;
    }
}
