package me.wener.bbvm.utils.val;

import com.google.common.collect.HashBasedTable;
import com.google.common.collect.Table;
import me.wener.bbvm.def.Instruction;

public class Values
{
    private static final Table<Class,Object,Object> cache = HashBasedTable.create();

    private Values() {}

    public static <V, T extends Enum & IsValue<V>> void cache(Class<T> type)
    {
        for (T item : type.getEnumConstants())
        {
            cache.put(type,item.asValue(), item);
        }
    }
    @SuppressWarnings("unchecked")
    public static <V, T extends Enum & IsValue<V>> T fromValue(Class<T> type, V v)
    {
        return (T) cache.get(type, v);
    }

    public static void main(String[] args)
    {
        cache(Instruction.class);
        assert fromValue(Instruction.class, Instruction.CAL.asValue()).equals(Instruction.CAL);

        long start = System.currentTimeMillis();
        int n = 100000;
        for (int i = 0; i < n; i++)
        {
            fromValue(Instruction.class, 0);
            fromValue(Instruction.class, 1);
            fromValue(Instruction.class, 2);
            fromValue(Instruction.class, 3);
            fromValue(Instruction.class, 4);
            fromValue(Instruction.class, 5);
            fromValue(Instruction.class, 6);
            fromValue(Instruction.class, 7);
            fromValue(Instruction.class, 8);
            fromValue(Instruction.class, 9);
            fromValue(Instruction.class, 0xA);
            fromValue(Instruction.class, 0XB);
        }
        long esplase = System.currentTimeMillis() - start;
        System.out.println(n*12 + " used "+esplase+" ms");

    }
}
