package me.wener.bbvm.vm.res.swing;

import me.wener.bbvm.vm.res.ImageManager;
import me.wener.bbvm.vm.res.InputManager;
import me.wener.bbvm.vm.res.PageManager;

import javax.swing.*;

/**
 * @author wener
 * @since 15/12/26
 */
public interface SwingContext {
    JFrame getFrame();

    PageManager getPageManager();

    ImageManager getImageManager();

    InputManager getInputManager();
}
