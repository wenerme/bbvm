package me.wener.bbvm.core.spi;


import java.util.ServiceLoader;
import me.wener.bbvm.core.Device;

public abstract class DeviceProvider
{
    public static DeviceProvider getProvider()
    {
        ServiceLoader<DeviceProvider> loader = ServiceLoader.load(DeviceProvider.class);
        //noinspection LoopStatementThatDoesntLoop
        for (DeviceProvider provider : loader)
            return provider;
        throw new RuntimeException("No DeviceProvider.");
    }

    public Device createDevice(int width, int height)
    {
        return null;
    }
}
