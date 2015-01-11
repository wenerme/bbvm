package me.wener.bbvm.neo.processor;

import me.wener.bbvm.neo.BBVMContext;

public interface VMContextAware
{
    void initialize(BBVMContext ctx);
}
