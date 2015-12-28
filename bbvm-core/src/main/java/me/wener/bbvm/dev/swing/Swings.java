package me.wener.bbvm.dev.swing;

import com.google.inject.Module;
import me.wener.bbvm.dev.DeviceConstants;
import me.wener.bbvm.dev.ResourceManager;
import me.wener.bbvm.exception.ResourceMissingException;
import me.wener.bbvm.util.IntEnums;

/**
 * @author wener
 * @since 15/12/18
 */
public class Swings implements DeviceConstants {
    static {
        IntEnums.cache(FontType.class);
    }

    public static Module module() {
        return new SwingModule();
    }

    static <T> T checkMissing(ResourceManager mgr, int handler, T v) {
        if (v == null) {
            throw new ResourceMissingException(String.format("%s #%s not exists", mgr.getType(), handler), handler);
        }
        return v;
    }
}
