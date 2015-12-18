package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/18
 */
public interface ImageResource extends Resource, Drawable {
    @Override
    ImageManager getManager();

    int getWidth();

    int getHeight();
}
