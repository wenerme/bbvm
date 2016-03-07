package me.wener.bbvm.vm;

import java.util.Map;
import java.util.NavigableMap;

/**
 * @author wener
 * @since 15/12/11
 */
public interface SymbolTable {
    Symbol getSymbol(String name);

    Symbol getSymbol(int address);

    Map<String, Symbol> getNameMap();

    NavigableMap<Integer, Symbol> getAddressMap();
}
