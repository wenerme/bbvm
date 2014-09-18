package me.wener.bbvm.swing;

import static me.wener.bbvm.swing.AsConvert.as;

import java.awt.Graphics;
import java.awt.image.BufferedImage;
import java.util.Timer;
import java.util.TimerTask;
import javax.swing.JFrame;
import javax.swing.JPanel;
import javax.swing.WindowConstants;
import me.wener.bbvm.core.Device;
import me.wener.bbvm.core.spi.DeviceProvider;
import me.wener.bbvm.swing.image.ImageFactory;

public class SwingDeviceProvider extends DeviceProvider
{
    @Override
    public Device createDevice(int width, int height)
    {
        final SwingDevice device = new SwingDevice(width, height);
        final JPanel panel = new JPanel();
        panel.setSize(width, height);
        panel.setPreferredSize(panel.getSize());
        panel.setLocation(0,0);
        final JFrame frame = new JFrame("BBVM");
        frame.setLayout(null);
        frame.setDefaultCloseOperation(WindowConstants.EXIT_ON_CLOSE);
     //   frame.setResizable(false);
        frame.add(panel);
        frame.setSize(width, height);
        //frame.pack();
        //frame.setSize(frame.getPreferredSize());
        frame.setLocationRelativeTo(null);
        //frame.setVisible(true);
        final SwingPage page = as(device.getScreen().asPage());

//        ImageFactory.BackgroundImageJFrame imageJFrame = new ImageFactory.BackgroundImageJFrame(page.asImage());
//        imageJFrame.setVisible(true);
        new Timer().schedule(new TimerTask()
        {
            @Override
            public void run()
            {
                Graphics graphics = panel.getGraphics().create();
                BufferedImage image = page.asImage();
                graphics.drawImage(image, 0, 0, null);
                graphics.dispose();
            }
        }, 0, 30);

        return device;
    }
}
