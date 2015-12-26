package me.wener.bbvm.dev.swing;

import com.google.inject.Exposed;
import com.google.inject.PrivateModule;
import com.google.inject.Provides;
import me.wener.bbvm.dev.*;

import javax.inject.Singleton;

/**
 * Provide a {@link SwingContext}
 *
 * @author wener
 * @since 15/12/26
 */
public class SwingModule extends PrivateModule {
    @Override
    protected void configure() {

    }

    @Exposed
    @Provides
    @Singleton
    InputManager inputManager(SwingInputManger swingInputManger) {
        return swingInputManger;
    }

    @Exposed
    @Provides
    @Singleton
    FileManager fileManager(JavaFileManager javaFileManager) {
        return javaFileManager;
    }

    @Exposed
    @Provides
    @Singleton
    ImageManager imageManager(SwingImageManager swingImageManager) {
        return swingImageManager;
    }

    @Exposed
    @Provides
    @Singleton
    PageManager pageManager(SwingPageManager swingPageManager) {
        return swingPageManager;
    }

    @Exposed
    @Provides
    @Singleton
    DeviceContext deviceContext(SwingContextImpl swingContext) {
        return swingContext;
    }

    @Exposed
    @Provides
    @Singleton
    SwingContext swingContext(SwingContextImpl swingContext) {
        return swingContext;
    }
}
