package me.wener.bbvm.core;


public interface HandlePool
{
    /**
     * 句柄池大小
     */
    int size();

    /**
     * 句柄池最大大小, -1 为不限制
     */
    int maxsize();

    /**
     * 请求句柄
     */
    int acquire();

    /**
     * 释放句柄
     */
    void release(int handle);
}
