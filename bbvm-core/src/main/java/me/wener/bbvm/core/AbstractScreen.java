package me.wener.bbvm.core;

public abstract class AbstractScreen<P extends Page> implements Screen<P>
{

    private final P page = getPage();

    protected abstract P getPage();

    @Override
    public void showPage(P page)
    {
        this.page.draw(page);
    }

    @Override
    public P asPage()
    {
        return page;
    }
}
