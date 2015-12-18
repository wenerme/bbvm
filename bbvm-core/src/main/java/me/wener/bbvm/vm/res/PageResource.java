package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/18
 */
public interface PageResource extends Resource, Drawable {
    @Override
    PageManager getManager();

    PageResource show(Drawable o, int dx, int dy, int w, int h, int x, int y, int mode);

    /**
     * Display on screen
     */
    PageResource display();

    PageResource show(PageResource resource);

    PageResource fill(int x, int y, int w, int h, int color);

    PageResource pixel(int x, int y, int color);

    int pixel(int x, int y);

    PageResource clear();

    PageResource show(PageResource resource, int x, int y);

    PageResource show(PageResource resource, int x, int y, int w, int h, int cx, int cy);
}
