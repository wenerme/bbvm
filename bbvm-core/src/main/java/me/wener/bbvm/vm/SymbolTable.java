package me.wener.bbvm.vm;

import java.util.Map;

/**
 * @author wener
 * @since 15/12/11
 */
public interface SymbolTable {
    Symbol getSymbol(String name);

    Map<String, Symbol> getSymbols();
}
