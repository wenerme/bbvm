package me.wener.bbvm.swing;

import me.wener.bbvm.core.AbstractDevice;
import me.wener.bbvm.core.AdapterResourceHandlePool;
import me.wener.bbvm.core.ResourceHandlePool;
import me.wener.bbvm.java.JavaFileHandle;

public class SwingDevice extends AbstractDevice<SwingScreen, SwingPage, SwingPicture, JavaFileHandle>
{

    @Override
    public ResourceHandlePool<SwingPage> getPagePool0()
    {
        return new AdapterResourceHandlePool<>(-1, new AdapterResourceHandlePool.Adapter<SwingPage>()
        {
            @Override
            public SwingPage createResource()
            {
                return null;
            }

            @Override
            public void freeResource(SwingPage resource)
            {

            }
        });
    }

    @Override
    public ResourceHandlePool<SwingPicture> getPicturePool0()
    {
        return new AdapterResourceHandlePool<>(-1, new AdapterResourceHandlePool.Adapter<SwingPicture>()
        {
            @Override
            public SwingPicture createResource()
            {
                return null;
            }

            @Override
            public void freeResource(SwingPicture resource)
            {

            }
        });
    }

    @Override
    public ResourceHandlePool<JavaFileHandle> getFilePool0()
    {
        return new AdapterResourceHandlePool<>(-1, new AdapterResourceHandlePool.Adapter<JavaFileHandle>()
        {
            @Override
            public JavaFileHandle createResource()
            {
                return new JavaFileHandle();
            }

            @Override
            public void freeResource(JavaFileHandle resource)
            {
                resource.close();
            }
        });
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
