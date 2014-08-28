package me.wener.bbvm.core;

public abstract class AbstractDevice<SCR extends Screen<PAGE>,PAGE extends Page,PIC extends Picture> implements Device
{
    public abstract AbstractHandlePool<PAGE> getPagePool();
    public abstract AbstractHandlePool<PIC> getPicturePool();
    public abstract SCR getScreen();
    @Override
    public DeviceFunction getFunction()
    {
        return null;
    }

    public abstract int waitkey();

    public abstract boolean isKeyPressed(int keycode);

    /**
     * @param file 资源文件
     * @param index 资源索引
     * @return 返回图片资源句柄
     */
    public abstract int loadPicture(String file, int index);

    public abstract void setScreenSize(int width, int height);
}
