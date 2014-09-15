package me.wener.bbvm.core;

public interface Device
{
    DeviceFunction getFunction();

    Screen getScreen();

    ResourceHandlePool<Page> getPagePool();

    ResourceHandlePool<Picture> getPicturePool();

    ResourceHandlePool<FileHandle> getFilePool();
}
