package me.wener.bbvm.dev;

/**
 * @author wener
 * @since 15/12/19
 */
public interface InputManager {

    boolean isKeyPressed(int key);

    /**
     * @return Wait any key or mouse click.
     */
    int waitKey();

    /**
     * Test is there any key press or mouse click.
     */
    int inKey();

    String inKeyString();

    /**
     * @return read a string, use Enter of confirm.
     */
    String readText();

    default int makeClickKey(int x, int y) {
        return x | 0x80000000 | (y << 16);
    }

    default int getX(int key) {
        return key & 0xFFFF;
    }

    default int getY(int key) {
        return key >>> 16 ^ 0x8000;
    }
}
