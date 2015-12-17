package me.wener.bbvm.vm;

import com.google.common.base.Function;

/**
 * @author wener
 * @since 15/12/13
 */
public interface SystemInvokeManager {

    void register(Object... handlers);


    void register(SystemInvoke.Type type, int a, int b, Function<Instruction, Object> handler);

    void invoke(Instruction inst);

    Function<Instruction, Object> getHandler(SystemInvoke.Type type, int a, int b);
}
