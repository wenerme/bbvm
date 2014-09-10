package me.wener.bbvm.swing;

import java.awt.Color;
import java.awt.Graphics2D;
import java.awt.image.BufferedImage;
import me.wener.bbvm.core.AbstractPage;
import me.wener.bbvm.core.Colour;
import me.wener.bbvm.core.constant.DrawMode;
import me.wener.bbvm.core.Page;
import me.wener.bbvm.core.Picture;

public class SwingPage
        extends AbstractPage<SwingPage>
        implements IsImage<BufferedImage>
{
    private final int width;
    private final int height;
    private final BufferedImage image;
    private final Graphics2D g;

    public SwingPage(int width, int height)
    {
        this(new BufferedImage(width, height, BufferedImage.TYPE_3BYTE_BGR));
    }

    public SwingPage(BufferedImage image)
    {
        this.image = image;
        width = image.getWidth();
        height = image.getHeight();
        g = (Graphics2D) image.getGraphics();
    }

    @Override
    public void drawLine(int x, int y, int dx, int dy)
    {
        g.setColor(VMUtils.color(penColor));
        g.drawLine(x, y, dx, dy);
    }

    @Override
    public void draw(Page src)
    {
        SwingPage page = (SwingPage) src;
        g.drawImage(page.asImage(), 0, 0, page.getWidth(), page.getHeight(), null);
    }

    @Override
    public void drawString(String content, int x, int y)
    {
        g.setColor(VMUtils.color(foreground));
        g.drawString(content, x, y);
    }

    @Override
    public void draw(Picture picture, int dx, int dy, int width, int height, int sx, int sy, DrawMode drawMode)
    {
        SwingPicture pic = (SwingPicture) picture;
        g.drawImage(pic.asImage(), dx, dy, dx + width, dy + height, sx, sy, sx + width, sy + height, null);
    }

    @Override
    public void rectangle(int left, int top, int right, int bottom)
    {
        g.setColor(VMUtils.color(penColor));
        g.drawRect(left, top, right - left, bottom - top);
    }

    @Override
    public void circle(int cx, int cy, int cr)
    {
        g.setColor(VMUtils.color(penColor));
        g.drawOval(cx - cr/2, cy - cr/2, cr, cr);
    }

    @Override
    public Colour pixel(int x, int y)
    {
        return new Colour(image.getRGB(x, y));
    }

    @Override
    public void pixel(int x, int y, Colour color)
    {
        image.setRGB(x, y, color.getRGB());
    }
    public void pixel(int x, int y, Color color)
    {
        image.setRGB(x, y, color.getRGB());
    }

    @Override
    public void fill(int x, int y, int width, int height, Colour color)
    {
        g.setColor(VMUtils.color(color));
        g.fillRect(x, y, width, height);
    }

    @Override
    public void draw(Page src, int x, int y, int width, int height, int sx, int sy)
    {
        SwingPage page = (SwingPage) src;
        g.drawImage(page.asImage(), x, y, x + width, y + height, sx, sy, sx + width, sy + height, null);
    }

    @Override
    public void draw(Page src, int x, int y)
    {
        SwingPage page = (SwingPage) src;
        g.drawImage(page.asImage(), x, y, null);
    }

    public int getWidth()
    {
        return width;
    }

    public int getHeight()
    {
        return height;
    }

    @Override
    public BufferedImage asImage()
    {
        return image;
    }
}
