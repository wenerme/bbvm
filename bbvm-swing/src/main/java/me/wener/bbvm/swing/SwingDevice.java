package me.wener.bbvm.swing;

import java.awt.Toolkit;
import java.io.IOException;
import me.wener.bbvm.core.DeviceFunction;
import me.wener.bbvm.core.FileHandle;
import me.wener.bbvm.core.Page;
import me.wener.bbvm.core.Picture;
import me.wener.bbvm.core.ResourceHandlePool;
import me.wener.bbvm.core.spi.AbstractDevice;
import me.wener.bbvm.swing.image.ImageFactory;

public class SwingDevice extends AbstractDevice
{
    final static Toolkit kit = Toolkit.getDefaultToolkit();

    private final SwingDeviceFunction function;

    public SwingDevice(int width, int height)
    {
        super(new SwingScreen(new SwingPage(width, height)));
        function = new SwingDeviceFunction(this);
    }

    @Override
    public ResourceHandlePool<Page> getPagePool0()
    {
        return new AutoResourceHandlePool<>(10, true);
    }

    @Override
    public ResourceHandlePool<Picture> getPicturePool0()
    {
        return new AutoResourceHandlePool<>(10, true);
    }

    @Override
    public ResourceHandlePool<FileHandle> getFilePool0()
    {
        AutoResourceHandlePool<FileHandle> pool = new AutoResourceHandlePool<>(10, true);
        pool.prepare();
        return pool;
    }

    @Override
    public DeviceFunction getFunction()
    {
        return function;
    }


    @Override
    public int waitkey()
    {
        return 0;
    }

    @Override
    public boolean isKeyPressed(int keycode)
    {
        return KeyStatus.isPressed(keycode);
    }

    @Override
    public int loadPicture(String file, int index)
    {
        try
        {
            return picturePool.addResource(new SwingPicture(ImageFactory.loadImage(file, index)));
        } catch (IOException e)
        {
            e.printStackTrace();
            return -1;
        }
    }

    @Override
    public void setScreenSize(int width, int height)
    {

    }
}
