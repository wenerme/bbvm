package me.wener.bbvm.core;

public interface Screen<T extends Page> extends IsPage<T>
{


    void showPage(T resource);

}
