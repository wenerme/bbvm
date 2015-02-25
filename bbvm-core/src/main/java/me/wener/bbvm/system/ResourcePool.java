package me.wener.bbvm.system;

import java.io.Closeable;
import java.util.Map;

public interface ResourcePool extends Closeable
{
    Resource request();

    Map<Integer, Resource> resources();

    Resource get(int handler);
}
