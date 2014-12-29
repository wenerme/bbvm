package me.wener.bbvm.swing.test;

import java.awt.Color;
import java.awt.Graphics;
import java.awt.event.KeyEvent;
import java.awt.image.BufferedImage;
import java.util.Timer;
import java.util.TimerTask;
import javax.swing.JFrame;
import javax.swing.WindowConstants;
import me.wener.bbvm.impl.plaf.Colour;
import me.wener.bbvm.impl.spi.DeviceProvider;
import me.wener.bbvm.swing.KeyStatus;
import me.wener.bbvm.swing.SwingDeviceProvider;
import me.wener.bbvm.swing.SwingPage;
import org.junit.Test;

public class TestFrame
{
    public static void main(String[] args)
    {
        new Timer().schedule(new TimerTask() {
            @Override
            public void run()
            {
                System.out.println("A pressed ? "+ KeyStatus.isPressed(KeyEvent.VK_A));
            }
        }, 0, 1000);

        final BufferedImage image = new BufferedImage(240, 320, BufferedImage.TYPE_3BYTE_BGR);
        JFrame frame = new JFrame()
        {
            @Override
            public void paint(Graphics g)
            {
                super.paint(g);
                g.drawImage(image, 0, 0, null);
            }
        };
        SwingPage page = new SwingPage(image);
        frame.setSize(240, 320);
        frame.setUndecorated(true);
        frame.setDefaultCloseOperation(WindowConstants.EXIT_ON_CLOSE);
        frame.setPreferredSize(frame.getSize());
        frame.pack();
        frame.setLocationRelativeTo(null);
        frame.setVisible(true);
        page.drawLine(0, 0, 240, 320);
        page.rectangle(60, 80, 180, 240);
        page.circle(0, 0, 60);
        page.circle(240,0,60);
        {
            for (int i = 0; i < 20; i++)
            {
                for (int j = 0; j < 20; j++)
                {
                    page.pixel(110+j,150+i, Color.green);
                }
            }
        }
        page.fill(0,260,60,60, Colour.gray);
        page.fill(180,260,60,60, Colour.gray);

        page.drawString("Demo here", 0, 40);
        page.drawString("测试示例", 0, 60);
    }

    @Test
    public void test()
    {
        assert DeviceProvider.getProvider().getClass() == SwingDeviceProvider.class;
    }
}
