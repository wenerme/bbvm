package me.wener.bbvm.system;

import java.io.Closeable;
import java.util.Map;

/**
 * 资源池
 */
public interface ResourcePool extends Closeable
{
    Resource request();

    Map<Integer, Resource> resources();

    Resource get(int handler);
}
