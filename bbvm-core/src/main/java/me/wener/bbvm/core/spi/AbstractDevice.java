package me.wener.bbvm.core.spi;

import me.wener.bbvm.core.Device;
import me.wener.bbvm.core.FileHandle;
import me.wener.bbvm.core.Page;
import me.wener.bbvm.core.Picture;
import me.wener.bbvm.core.ResourceHandlePool;
import me.wener.bbvm.core.Screen;

public abstract class AbstractDevice implements Device
{
    protected final Screen screen;
    protected final ResourceHandlePool<Page> pagePool;
    protected final ResourceHandlePool<Picture> picturePool;
    protected final ResourceHandlePool<FileHandle> filePool;

    @Override
    public Screen getScreen()
    {
        return screen;
    }

    protected AbstractDevice(Screen screen)
    {
        this.screen = screen;
        pagePool = getPagePool0();
        picturePool = getPicturePool0();
        filePool = getFilePool0();
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
