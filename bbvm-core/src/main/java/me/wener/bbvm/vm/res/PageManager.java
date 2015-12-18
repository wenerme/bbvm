package me.wener.bbvm.vm.res;

/**
 * @author wener
 * @since 15/12/18
 */
public interface PageManager extends ResourceManager<PageManager, PageResource> {
    @Override
    default String getType() {
        return "graph";
    }

    Color color(int next);

    PageResource screen();

}
