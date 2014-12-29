package me.wener.bbvm.api;

public interface ValueHolder<T>
{
    T get();
    void set(T v);
}
