package me.wener.bbvm.vm;

/**
 * A symbol represent a address in memory
 *
 * @author wener
 * @since 15/12/11
 */
public interface Symbol {
    String getName();

    int getAddress();
}
