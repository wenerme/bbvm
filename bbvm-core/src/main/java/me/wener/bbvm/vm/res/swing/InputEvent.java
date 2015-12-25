package me.wener.bbvm.vm.res.swing;

/**
 * @author wener
 * @since 15/12/26
 */
class InputEvent {
    private Type type;
    private int keyCode;
    private char keyChar;

    public InputEvent() {
    }

    public InputEvent(Type type, int keyCode) {
        this.type = type;
        this.keyCode = keyCode;
        this.keyChar = '#';
    }

    public InputEvent(Type type, int keyCode, char keyChar) {
        this.type = type;
        this.keyCode = keyCode;
        this.keyCode = keyChar;
    }

    public int getKeyCode() {
        return keyCode;
    }

    public char getKeyChar() {
//                Character.isC
        return keyChar;
    }

    public enum Type {
        UP, DOWN, CLICK
    }
}
