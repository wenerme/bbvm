package me.wener.bbvm.api;

public interface BBVm
{
    void reset();

    void start();

    void push(int v);

    int pop();

    void exit();
}
