package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/18
 */
public interface PageResource extends Resource, Drawable {
    @Override
    PageManager getManager();

    PageResource draw(Drawable o, int dx, int dy, int w, int h, int x, int y, int mode);

    /**
     * Display on screen
     */
    PageResource display();

    PageResource draw(PageResource resource);

    PageResource fill(int x, int y, int w, int h, int rgb);

    PageResource pixel(int x, int y, int rgb);

    int pixel(int x, int y);

    /**
     * Fill page with pen color
     */
    PageResource fill();

    PageResource draw(PageResource resource, int x, int y);

    PageResource draw(PageResource resource, int x, int y, int w, int h, int cx, int cy);

    PageResource pen(int width, int style, int rgb);

    PageResource circle(int cx, int cy, int r);

    PageResource rectangle(int left, int top, int right, int bottom);

    PageResource line(int x, int y);

    PageResource move(int x, int y);
}