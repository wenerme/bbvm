package me.wener.bbvm.dev;

/**
 * @author wener
 * @since 15/12/19
 */
public interface InputManager {

    boolean isKeyPressed(int key);

    int waitKey();

    int getLastKeyCode();

    String getLastKeyString();

    String readText();
}
