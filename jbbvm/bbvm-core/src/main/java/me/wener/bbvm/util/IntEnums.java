package me.wener.bbvm.util;

import com.google.common.collect.HashBasedTable;
import com.google.common.collect.Table;

/**
 * @author wener
 * @since 15/12/10
 */
public class IntEnums {
    private static final Table<Class, Integer, Object> cache = HashBasedTable.create();

    @SafeVarargs
    public static <T extends Enum & IsInt> void cache(Class<? extends T>... type) {
        for (Class<? extends T> t : type) {
            cache(t);
        }
    }

    private static <T extends Enum & IsInt> void cache(Class<? extends T> type) {
        for (T item : type.getEnumConstants()) {
            cache.put(type, item.asInt(), item);
        }
    }

    @SuppressWarnings("unchecked")
    public static <T extends Enum & IsInt> T fromInt(Class<T> type, int value) {
        return (T) cache.get(type, value);
    }
}
