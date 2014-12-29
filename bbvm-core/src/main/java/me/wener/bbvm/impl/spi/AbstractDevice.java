package me.wener.bbvm.impl.spi;

import me.wener.bbvm.api.Device;
import me.wener.bbvm.api.FileHandle;
import me.wener.bbvm.api.Page;
import me.wener.bbvm.api.Picture;
import me.wener.bbvm.api.Screen;
import me.wener.bbvm.impl.ResourceHandlePool;

public abstract class AbstractDevice implements Device
{
    protected final Screen screen;
    protected final ResourceHandlePool<Page> pagePool;
    protected final ResourceHandlePool<Picture> picturePool;
    protected final ResourceHandlePool<FileHandle> filePool;

    protected AbstractDevice(Screen screen)
    {
        this.screen = screen;
        pagePool = getPagePool0();
        picturePool = getPicturePool0();
        filePool = getFilePool0();
    }

    @Override
    public Screen getScreen()
    {
        return screen;
    }

    @Override
    public final ResourceHandlePool<Page> getPagePool()
    {
        return pagePool;
    }

    public final ResourceHandlePool<Picture> getPicturePool()
    {
        return picturePool;
    }

    public final ResourceHandlePool<FileHandle> getFilePool()
    {
        return filePool;
    }

    protected abstract ResourceHandlePool<Page> getPagePool0();

    protected abstract ResourceHandlePool<Picture> getPicturePool0();

    protected abstract ResourceHandlePool<FileHandle> getFilePool0();

    public abstract int waitkey();

    public abstract boolean isKeyPressed(int keycode);

    /**
     * @param file  资源文件
     * @param index 资源索引
     * @return 返回图片资源句柄
     */
    public abstract int loadPicture(String file, int index);

    public abstract void setScreenSize(int width, int height);
}
