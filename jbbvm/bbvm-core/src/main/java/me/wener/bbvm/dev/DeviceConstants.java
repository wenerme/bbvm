package me.wener.bbvm.dev;

import me.wener.bbvm.util.IsInt;
import org.intellij.lang.annotations.MagicConstant;

/**
 * Predefined device relate constants
 *
 * @author wener
 * @since 15/12/26
 */
public interface DeviceConstants {
    //    Brush Style
    int BRUSH_SOLID = 0;

    //    Pen style
    int PEN_SOLID = 0;
    int PEN_DASH = 1;

    // Env type
    int ENV_SIM = 0;
    int ENV_9288 = 9288;
    int ENV_9188 = 9188;
    int ENV_9288T = 9287;
    int ENV_9288S = 9286;
    int ENV_9388 = 9388;
    int ENV_9688 = 9688;

    /**
     * Draw image with key color
     */
    int DRAW_KEY_COLOR = 1;
    /**
     * No font background
     */
    int BACKGROUND_TRANSPARENT = 1;

    /**
     * With font background
     */
    int BACKGROUND_OPAQUE = 2;

    enum FontType implements IsInt {
        FONT_12SONG(0, 12),
        FONT_12KAI(1, 12),
        FONT_12HEI(2, 12),
        FONT_16SONG(3, 16),
        FONT_16KAI(4, 16),
        FONT_16HEI(5, 16),
        FONT_24SONG(6, 24),
        FONT_24KAI(7, 24),
        FONT_24HEI(8, 24);
        private final int value;
        private final int size;

        FontType(int value, int size) {
            this.value = value;
            this.size = size;
        }

        public int asInt() {
            return value;
        }

        public int getSize() {
            return size;
        }
    }

//    enum BackgroundMode implements IsInt {
//        /**
//         * 透明显示，即字体的背景颜色无效。
//         */
//        TRANSPARENT(1),
//        /**
//         * 不透明显示，即字体的背景颜色有效
//         */
//        OPAQUE(2);
//        private final int value;
//
//        BackgroundMode(int value) {
//            this.value = value;
//        }
//
//        public int asInt() {
//            return value;
//        }
//    }

    @MagicConstant(intValues = {ENV_SIM, ENV_9188, ENV_9288, ENV_9288S, ENV_9288T, ENV_9388, ENV_9688})
    @interface Env {
    }
}
