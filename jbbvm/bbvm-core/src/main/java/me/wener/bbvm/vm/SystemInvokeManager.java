package me.wener.bbvm.vm;

import com.google.common.base.Function;
import com.google.inject.ImplementedBy;

/**
 * @author wener
 * @since 15/12/13
 */
@ImplementedBy(SystemInvokeManagerImpl.class)
public interface SystemInvokeManager {

    /**
     * Register all method with annotation {@link SystemInvoke}, if register a class will use {@linkplain com.google.inject.Injector Injector} to create new instance.
     * If there is already have a handler for type,a,b will throw an exception.
     *
     * @param handlers Object or Class
     */
    void register(Object... handlers);

//    void register(boolean override,Object... handlers);

    void register(SystemInvoke.Type type, int a, int b, Function<Instruction, Object> handler);

    Object invoke(Instruction inst);

    Function<Instruction, Object> getHandler(SystemInvoke.Type type, int a, int b);
}
