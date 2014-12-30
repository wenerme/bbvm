package me.wener.bbvm.utils.val;

/**
 * 可变值对象
 *
 * @param <T> 值类型
 */
public interface ValueHolder<T> extends IsValue<T>
{
    T get();

    void set(T v);
}
