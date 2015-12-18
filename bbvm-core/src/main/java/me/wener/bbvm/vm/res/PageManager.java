package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/18
 */
public interface PageManager extends ResourceManager<PageManager, PageResource> {

    PageResource screen();

    int getWidth();

    int getHeight();
}
