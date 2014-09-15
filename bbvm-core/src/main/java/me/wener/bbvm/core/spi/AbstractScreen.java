package me.wener.bbvm.core.spi;

import me.wener.bbvm.core.Page;
import me.wener.bbvm.core.Screen;

public abstract class AbstractScreen implements Screen
{

    private final Page page = getPage();

    protected abstract Page getPage();

    @Override
    public void showPage(Page page)
    {
        this.page.draw(page);
    }

    @Override
    public Page asPage()
    {
        return page;
    }
}
