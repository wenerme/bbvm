package me.wener.bbvm.swing;

import me.wener.bbvm.core.FileHandle;
import me.wener.bbvm.core.Page;
import me.wener.bbvm.core.Picture;
import me.wener.bbvm.core.ResourceHandlePool;
import me.wener.bbvm.core.spi.AbstractDevice;
import me.wener.bbvm.core.DeviceFunction;

public class SwingDevice extends AbstractDevice
{

    @Override
    public ResourceHandlePool<Page> getPagePool0()
    {
        return null;
    }

    @Override
    public ResourceHandlePool<Picture> getPicturePool0()
    {
        return null;
    }

    @Override
    public ResourceHandlePool<FileHandle> getFilePool0()
    {
        return null;
    }

    @Override
    public DeviceFunction getFunction()
    {
        return null;
    }

    @Override
    public SwingScreen getScreen()
    {
        return null;
    }

    @Override
    public int waitkey()
    {
        return 0;
    }

    @Override
    public boolean isKeyPressed(int keycode)
    {
        return false;
    }

    @Override
    public int loadPicture(String file, int index)
    {
        return 0;
    }

    @Override
    public int loadPicture(int file, int index)
    {
        return 0;
    }

    @Override
    public void setScreenSize(int width, int height)
    {

    }
}
