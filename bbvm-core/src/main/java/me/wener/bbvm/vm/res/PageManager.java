package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/18
 */
public interface PageManager extends ResourceManager<PageManager, PageResource> {

    PageResource getScreen();

    int getWidth();

    int getHeight();

    PageManager setSize(int w, int h);

    @Override
    default String getType() {
        return "page";
    }
}
