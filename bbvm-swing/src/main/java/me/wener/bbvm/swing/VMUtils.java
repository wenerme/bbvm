package me.wener.bbvm.swing;

import java.awt.Color;
import me.wener.bbvm.impl.plaf.Colour;

public class VMUtils
{
    public static Color color(Colour colour)
    {
        return new Color(colour.getRGB());
    }
}
