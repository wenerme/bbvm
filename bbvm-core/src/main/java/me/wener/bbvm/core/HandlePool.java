package me.wener.bbvm.core;


public interface HandlePool
{
    int size();
    int maxsize();
    int acquire();
    void release(int handle);
}
