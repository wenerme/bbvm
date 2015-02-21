package me.wener.bbvm.system.api;

import java.io.Closeable;
import me.wener.bbvm.utils.val.ValueHolder;

/**
 * 资源内容持有类
 */
public interface Resource extends ValueHolder<Object>, Closeable
{
    <T> T as();

    /**
     * @return 句柄号
     */
    int handler();

    boolean isNull();
}
