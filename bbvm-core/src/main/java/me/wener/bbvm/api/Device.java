package me.wener.bbvm.api;

import me.wener.bbvm.impl.ResourceHandlePool;

public interface Device
{
    DeviceFunction getFunction();

    Screen getScreen();

    ResourceHandlePool<Page> getPagePool();

    ResourceHandlePool<Picture> getPicturePool();

    ResourceHandlePool<FileHandle> getFilePool();
}
