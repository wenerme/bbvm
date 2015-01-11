package me.wener.bbvm.neo;

import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.experimental.Accessors;
import me.wener.bbvm.utils.val.IntegerHolder;
import me.wener.bbvm.utils.val.impl.SimpleValueHolder;

public class Registers
{
    public static Register create(String name)
    {
        return new DefaultRegister(name);
    }

    public static Register create(String name, IntegerHolder holder)
    {
        return new ProxyRegister(name, holder);
    }

    @EqualsAndHashCode(callSuper = true)
    @Data
    @Accessors(chain = true)
    public static class DefaultRegister extends SimpleValueHolder<Integer> implements Register
    {

        private final String name;

        public DefaultRegister(String name)
        {
            super(0);
            this.name = name;
        }

        @Override
        public String toString()
        {
            return name;
        }

        @Override
        public String name()
        {
            return name;
        }
    }

    private static class ProxyRegister implements Register
    {
        private final IntegerHolder internal;
        private final String name;

        private ProxyRegister(String name, IntegerHolder internal)
        {
            this.internal = internal;
            this.name = name;
        }

        @Override
        public String name()
        {
            return null;
        }

        @Override
        public Integer get()
        {
            return null;
        }

        @Override
        public void set(Integer v)
        {

        }
    }
}
