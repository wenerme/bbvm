package me.wener.bbvm.system;

import me.wener.bbvm.util.val.ValueHolder;

import java.io.Closeable;

/**
 * 资源内容持有类
 */
public interface Resource extends ValueHolder<Object>, Closeable
{
    <T> T as();

    /**
     * @return 句柄号
     */
    Integer handler();

    boolean isNull();

    Resource handler(Integer handler);
}
