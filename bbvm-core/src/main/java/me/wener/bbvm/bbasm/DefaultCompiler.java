package me.wener.bbvm.bbasm;

import me.wener.bbvm.system.OpState;
import me.wener.bbvm.system.internal.OpStates;

public abstract class DefaultCompiler extends OpStates.DefaultOpState implements OpState, Compiler
{
    String current;
    int line;
    String content;

    protected void readStatement()
    {

    }

    protected void isComment() {}

    protected void skipLine() {}

    protected void readOpcode() {}

    protected void readDataType()
    {
    }

    protected void readString() {}

    protected void readNumber() {}

    protected void readInteger() {}


    protected void readComma() {}

    protected void readOperand()
    {

    }

    protected boolean hasMore()
    {
        return false;
    }

    @Override
    public void compile()
    {
        while (hasMore())
        {
            readStatement();
        }
    }
}
