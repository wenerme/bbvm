package me.wener.bbvm.dev.swing;

import me.wener.bbvm.dev.DeviceContext;

import javax.swing.*;

/**
 * @author wener
 * @since 15/12/26
 */
public interface SwingContext extends DeviceContext {
    JFrame getFrame();
}
