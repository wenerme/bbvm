package me.wener.bbvm.bbasm;

import com.google.common.collect.BiMap;
import java.util.List;
import java.util.Map;

public interface Compiler
{
    List<Statement> statements();

    Map<String, Integer> labels();

    BiMap<Integer, Integer> offsetAndLine();

    byte[] toBinary();

    String toAssembly();

    void compile();
}
