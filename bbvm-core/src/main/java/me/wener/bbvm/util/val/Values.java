package me.wener.bbvm.util.val;

import com.google.common.collect.HashBasedTable;
import com.google.common.collect.Table;
import me.wener.bbvm.util.val.impl.ReadonlyHolder;
import me.wener.bbvm.util.val.impl.SimpleStringHolder;
import me.wener.bbvm.util.val.impl.SimpleValue;
import me.wener.bbvm.util.val.impl.SimpleValueHolder;

public class Values
{
    private static final Table<Class, Object, Object> cache = HashBasedTable.create();

    private Values() {}

    public static <V, T extends Enum & IsValue<V>> void cache(Class<T> type)
    {
        for (T item : type.getEnumConstants())
        {
            cache.put(type, item.get(), item);
        }
    }

    @SafeVarargs
    public static <T extends Enum & IsValue> void cache(Class<? extends T>... type)
    {
        for (Class<? extends T> t : type)
        {
            cache(t);
        }
    }

    public static int or(IsInt... integers)
    {
        int result = 0;
        for (IsInt integer : integers)
        {
            result |= integer.asInt();
        }
        return result;
    }

    public static int or(Iterable<? extends IsInt> integers)
    {
        int result = 0;
        for (IsInt integer : integers)
        {
            result |= integer.asInt();
        }
        return result;
    }

    public static <V, T extends Enum & IsInt> T fromValue(Class<T> type, int v) {
        return IntEnums.fromInt(type, v);
    }

    /**
     * @return null if not found
     */
    @SuppressWarnings("unchecked")
    public static <V, T extends Enum & IsValue<V>> T fromValue(Class<T> type, V v)
    {
        return (T) cache.get(type, v);
    }

    @SuppressWarnings("unchecked")
    public static <V, T extends Enum & IsValue<V>> T fromValue(Class<T> type, V v, T forNull)
    {
        T val = (T) cache.get(type, v);
        return val == null ? forNull : val;
    }

    public static <T> ValueHolder<T> hold(T value)
    {
        return new SimpleValueHolder<>(value);
    }

    public static StringHolder hold(String value)
    {
        return new SimpleStringHolder(value);
    }

    public static <T> ValueHolder<T> readonly(T value)
    {
        return new ReadonlyHolder<>(value);
    }

    public static <T> IsValue<T> valueOf(T value)
    {
        return new SimpleValue<>(value);
    }
}
