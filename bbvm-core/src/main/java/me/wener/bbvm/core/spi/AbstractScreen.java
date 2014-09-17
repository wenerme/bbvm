package me.wener.bbvm.core.spi;

import me.wener.bbvm.core.Page;
import me.wener.bbvm.core.Screen;

public abstract class AbstractScreen implements Screen
{
    protected final Page page;

    protected AbstractScreen(Page page) {this.page = page;}

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
