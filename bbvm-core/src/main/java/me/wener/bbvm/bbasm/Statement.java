package me.wener.bbvm.bbasm;

public interface Statement extends Compilable
{
    int line();

    String toAssembly();
}
