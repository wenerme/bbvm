package me.wener.bbvm.core;

public class BasicFunctionImpl implements BasicFunction
{
    @Override
    public float SIN(float X)
    {
        return (float) Math.sin(X);
    }

    @Override
    public float COS(float X)
    {
        return (float) Math.cos(X);
    }

    @Override
    public float TAN(float X)
    {
        return (float) Math.tan(X);
    }

    @Override
    public float SQRT(float X)
    {
        return (float) Math.sqrt(X);
    }

    @Override
    public float ABS(float X)
    {
        return Math.abs(X);

    }

    @Override
    public int ABS(int X)
    {
        return Math.abs(X);

    }

    @Override
    public int LEN(String X$)
    {
        return X$.length();
    }

    @Override
    public String STR$(int V)
    {
        return String.valueOf(V);
    }

    @Override
    public int VAL(String X$)
    {
        return Integer.parseInt(X$);
    }

    @Override
    public String CHR$(int X)
    {
        return String.valueOf((char)X);
    }

    @Override
    public int ASC(String X$)
    {
        return X$.charAt(0);
    }

    @Override
    public String LEFT$(String X$, int N)
    {
        if (N < 0 || N > X$.length())
            return X$;
        return X$.substring(0,N);
    }

    @Override
    public String RIGHT$(String X$, int N)
    {
        if (N < 0 || N > X$.length())
            return X$;
        return X$.substring(X$.length()-N);
    }

    @Override
    public String MID$(String X$, int S, int N)
    {
        if (N < 0 || N > X$.length())
            return X$;
        return X$.substring(S, S + N);
    }

    @Override
    public int INSTR(int index, String sub, String B$)
    {
        return B$.indexOf(sub, index);
    }
}
