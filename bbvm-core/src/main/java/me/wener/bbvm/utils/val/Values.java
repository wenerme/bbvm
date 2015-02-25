package me.wener.bbvm.utils.val;

import com.google.common.collect.HashBasedTable;
import com.google.common.collect.Table;
import java.util.EnumSet;
import me.wener.bbvm.utils.val.impl.ReadonlyHolder;
import me.wener.bbvm.utils.val.impl.SimpleStringHolder;
import me.wener.bbvm.utils.val.impl.SimpleValue;
import me.wener.bbvm.utils.val.impl.SimpleValueHolder;

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

    public static <T extends Enum<T> & IsInteger> EnumSet<T> asEnumSet(long flag, Class<T> clazz)
    {
        EnumSet<T> set = EnumSet.noneOf(clazz);
        for (int i = 0; i < 32; i++)
        {
            int v = 1 << i;
            if ((flag & v) > 0)
            {
                T e = fromValue(clazz, v);
                if (e == null)
                    continue;
                set.add(e);
            }
        }

        return set;
    }

    public static int or(IsInteger... integers)
    {
        int result = 0;
        for (IsInteger integer : integers)
        {
            result |= integer.get();
        }
        return result;
    }

    public static int or(Iterable<? extends IsInteger> integers)
    {
        int result = 0;
        for (IsInteger integer : integers)
        {
            result |= integer.get();
        }
        return result;
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
