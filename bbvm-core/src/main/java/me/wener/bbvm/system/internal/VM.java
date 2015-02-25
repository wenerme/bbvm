package me.wener.bbvm.system.internal;

/**
 * 工厂类
 */
public class VM
{
    public static me.wener.bbvm.system.Register register(String name)
    {
        return new Register(name);
    }
}
