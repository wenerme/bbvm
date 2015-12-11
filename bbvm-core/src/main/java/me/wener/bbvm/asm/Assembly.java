package me.wener.bbvm.asm;

/**
 * @author wener
 * @since 15/12/11
 */
public interface Assembly {
    Type getType();

    String toAssembly();

    enum Type {
        LABEL, COMMENT, PSEUDO, INST
    }
}
