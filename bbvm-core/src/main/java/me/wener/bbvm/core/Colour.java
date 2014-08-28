package me.wener.bbvm.core;

import java.beans.ConstructorProperties;

public class Colour
{
    public final static Colour white = new Colour(255, 255, 255);

    public final static Colour lightGray = new Colour(192, 192, 192);

    public final static Colour gray = new Colour(128, 128, 128);

    public final static Colour darkGray = new Colour(64, 64, 64);

    public final static Colour black = new Colour(0, 0, 0);


    public final static Colour red = new Colour(255, 0, 0);


    public final static Colour pink = new Colour(255, 175, 175);


    public final static Colour orange = new Colour(255, 200, 0);


    public final static Colour yellow = new Colour(255, 255, 0);


    public final static Colour green = new Colour(0, 255, 0);


    public final static Colour magenta = new Colour(255, 0, 255);


    public final static Colour cyan = new Colour(0, 255, 255);

    public final static Colour blue = new Colour(0, 0, 255);

    private final int value;


    public Colour(int r, int g, int b)
    {
        this(r, g, b, 255);
    }

    @ConstructorProperties({"red", "green", "blue", "alpha"})
    public Colour(int r, int g, int b, int a)
    {
        value = ((a & 0xFF) << 24) |
                ((r & 0xFF) << 16) |
                ((g & 0xFF) << 8) |
                ((b & 0xFF) << 0);
        testColorValueRange(r, g, b, a);
    }

    public Colour(int rgb)
    {
        value = 0xff000000 | rgb;
    }

    public Colour(int rgba, boolean hasalpha)
    {
        if (hasalpha)
        {
            value = rgba;
        } else
        {
            value = 0xff000000 | rgba;
        }
    }


    /**
     * @param color ARGB 格式
     */
    public static Colour fromARGB(int color)
    {
        int b = color & 0xFF;
        color >>= 8;
        int g = color & 0xFF;
        color >>= 8;
        int r = color & 0xFF;
        color >>= 8;
        int a = color & 0xFF;

        return new Colour(r, g, b, a);
    }

    private static void testColorValueRange(int r, int g, int b, int a)
    {
        boolean rangeError = false;
        String badComponentString = "";

        if (a < 0 || a > 255)
        {
            rangeError = true;
            badComponentString = badComponentString + " Alpha";
        }
        if (r < 0 || r > 255)
        {
            rangeError = true;
            badComponentString = badComponentString + " Red";
        }
        if (g < 0 || g > 255)
        {
            rangeError = true;
            badComponentString = badComponentString + " Green";
        }
        if (b < 0 || b > 255)
        {
            rangeError = true;
            badComponentString = badComponentString + " Blue";
        }
        if (rangeError)
        {
            throw new IllegalArgumentException("Color parameter outside of expected range:"
                    + badComponentString);
        }
    }

    public int getRGB()
    {
        return value;
    }

    public int getRed()
    {
        return (getRGB() >> 16) & 0xFF;
    }

    public int getGreen()
    {
        return (getRGB() >> 8) & 0xFF;
    }

    public int getBlue()
    {
        return (getRGB() >> 0) & 0xFF;
    }

    public int getAlpha()
    {
        return (getRGB() >> 24) & 0xff;
    }

    public String toString()
    {
        return getClass().getName() + "[r=" + getRed() + ",g=" + getGreen() + ",b=" + getBlue() + "]";
    }
}
