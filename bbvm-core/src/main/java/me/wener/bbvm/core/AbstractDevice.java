package me.wener.bbvm.core;

import me.wener.bbvm.core.constant.Device;

public abstract class AbstractDevice
        <SCR extends Screen<PAGE>, PAGE extends Page, PIC extends Picture, FILE extends FileHandle>
        implements Device
{
    protected final ResourceHandlePool<PAGE> pagePool;
    protected final ResourceHandlePool<PIC> picturePool;
    protected final ResourceHandlePool<FILE> filePool;

    protected AbstractDevice()
    {
        pagePool = getPagePool0();
        picturePool = getPicturePool0();
        filePool = getFilePool0();
    }

    public final ResourceHandlePool<PAGE> getPagePool()
    {
        return pagePool;
    }

    public final ResourceHandlePool<PIC> getPicturePool()
    {
        return picturePool;
    }

    public final ResourceHandlePool<FILE> getFilePool()
    {
        return filePool;
    }

    protected abstract ResourceHandlePool<PAGE> getPagePool0();

    protected abstract ResourceHandlePool<PIC> getPicturePool0();

    protected abstract ResourceHandlePool<FILE> getFilePool0();

    public abstract SCR getScreen();

    @Override
    public DeviceFunction getFunction()
    {
        return null;
    }

    public abstract int waitkey();

    public abstract boolean isKeyPressed(int keycode);

    /**
     * @param file  资源文件
     * @param index 资源索引
     * @return 返回图片资源句柄
     */
    public abstract int loadPicture(String file, int index);

    /**
     * @param file  资源文件句柄
     * @param index 资源索引
     * @return 返回图片资源句柄
     */
    public abstract int loadPicture(int file, int index);

    public abstract void setScreenSize(int width, int height);
}
