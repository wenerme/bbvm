package me.wener.bbvm.core;

public interface ValueHolder<T>
{
    T get();
    void set(T v);
}
