package me.wener.bbvm.swing;

import java.awt.Toolkit;
import java.io.IOException;
import me.wener.bbvm.api.DeviceFunction;
import me.wener.bbvm.api.FileHandle;
import me.wener.bbvm.api.Page;
import me.wener.bbvm.api.Picture;
import me.wener.bbvm.impl.ResourceHandlePool;
import me.wener.bbvm.impl.spi.AbstractDevice;
import me.wener.bbvm.java.JavaFileHandle;
import me.wener.bbvm.swing.image.ImageFactory;

public class SwingDevice extends AbstractDevice
{
    final static Toolkit kit = Toolkit.getDefaultToolkit();

    private final SwingDeviceFunction function;
    private final int height;
    private final int width;

    public SwingDevice(int width, int height)
    {
        super(new SwingScreen(new SwingPage(width, height)));
        this.width = width;
        this.height = height;

        function = new SwingDeviceFunction(this);
    }

    @Override
    public ResourceHandlePool<Page> getPagePool0()
    {
        return new AutoResourceHandlePool<Page>(10, true)
        {
            @Override
            public Page createResource()
            {
                return new SwingPage(width, height);
            }
        };
    }

    @Override
    public ResourceHandlePool<Picture> getPicturePool0()
    {
        return new AutoResourceHandlePool<Picture>(10, true, SwingPicture.class);
    }

    @Override
    public ResourceHandlePool<FileHandle> getFilePool0()
    {
        AutoResourceHandlePool<FileHandle> pool = new AutoResourceHandlePool<FileHandle>(10, true, JavaFileHandle.class);
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
