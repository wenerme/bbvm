package me.wener.bbvm.dev;

/**
 * Represent a runtime context
 *
 * @author wener
 * @since 15/12/26
 */
public interface DeviceContext {
    PageManager getPageManager();

    ImageManager getImageManager();

    InputManager getInputManager();

    FileManager getFileManager();

    StringManager getStringManager();
}
