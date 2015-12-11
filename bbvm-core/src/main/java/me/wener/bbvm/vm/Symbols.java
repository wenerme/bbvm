package me.wener.bbvm.vm;

import com.google.common.collect.Maps;

import java.util.Collection;
import java.util.Map;
import java.util.NavigableMap;
import java.util.TreeMap;

/**
 * @author wener
 * @since 15/12/11
 */
public class Symbols {
    public static SymbolTable table(Collection<? extends Symbol> symbols) {
        Map<String, Symbol> nameMap = Maps.newHashMap();
        NavigableMap<Integer, Symbol> addressMap = new TreeMap<>();
        for (Symbol symbol : symbols) {
            nameMap.put(symbol.getName(), symbol);
            addressMap.put(symbol.getAddress(), symbol);
        }

        return new Table(nameMap, addressMap);
    }

    private static class Table implements SymbolTable {
        private Map<String, Symbol> nameMap;
        private NavigableMap<Integer, Symbol> addressMap;

        public Table(Map<String, Symbol> nameMap, NavigableMap<Integer, Symbol> addressMap) {
            this.nameMap = nameMap;
            this.addressMap = addressMap;
        }

        @Override
        public Symbol getSymbol(String name) {
            return nameMap.get(name);
        }

        @Override
        public Symbol getSymbol(int address) {
            return addressMap.get(address);
        }

        @Override
        public Map<String, Symbol> getNameMap() {
            return nameMap;
        }

        @Override
        public NavigableMap<Integer, Symbol> getAddressMap() {
            return addressMap;
        }
    }
}
