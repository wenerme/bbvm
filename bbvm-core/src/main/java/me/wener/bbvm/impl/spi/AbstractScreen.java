package me.wener.bbvm.impl.spi;

import me.wener.bbvm.api.Page;
import me.wener.bbvm.api.Screen;

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
