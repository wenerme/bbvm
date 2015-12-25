package me.wener.bbvm.vm.res.swing;

import com.google.common.base.MoreObjects;
import me.wener.bbvm.vm.res.ImageManager;
import me.wener.bbvm.vm.res.ImageResource;

import java.awt.image.BufferedImage;

/**
 * @author wener
 * @since 15/12/26
 */
class SwingImage extends Draw implements ImageResource {
    private final int handler;
    private final SwingImageManager manager;
    public String name;

    SwingImage(int handler, SwingImageManager manager, BufferedImage image) {
        super(image);
        this.handler = handler;
        this.manager = manager;
    }

    @Override
    public int getHandler() {
        return handler;
    }

    @Override
    public ImageManager getManager() {
        return manager;
    }

    @Override
    public void close() {
        manager.close(this);
    }

    @Override
    public String toString() {
        return MoreObjects.toStringHelper(this)
                .add("handler", handler)
                .add("width", getWidth())
                .add("height", getHeight())
                .add("name", name)
                .toString();
    }
}
