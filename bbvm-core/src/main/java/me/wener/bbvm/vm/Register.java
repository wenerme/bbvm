package me.wener.bbvm.vm;

/**
 * @author wener
 * @since 15/12/10
 */
public interface Register extends Value<Register> {
    RegisterType getType();
}
