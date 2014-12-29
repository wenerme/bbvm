package me.wener.bbvm.swing;

import java.awt.Image;
import me.wener.bbvm.api.Picture;

public class SwingPicture implements Picture, IsImage
{
    private final Image image;

    public SwingPicture(Image image)
    {
        this.image = image;
    }

    @Override
    public Image asImage()
    {
        return image;
    }

    @Override
    public int getWidth()
    {
        return image.getWidth(null);
    }

    @Override
    public int getHeight()
    {
        return image.getHeight(null);
    }
}
