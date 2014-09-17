package me.wener.bbvm.swing;

import me.wener.bbvm.core.Device;
import me.wener.bbvm.core.spi.DeviceProvider;

public class SwingDeviceProvider extends DeviceProvider
{
    @Override
    public Device createDevice(int width, int height)
    {
        return  new SwingDevice(width, height);
    }
}
