package me.wener.bbvm.vm.invoke;

import me.wener.bbvm.dev.InputManager;
import me.wener.bbvm.vm.Register;
import me.wener.bbvm.vm.SystemInvoke;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import javax.inject.Named;

/**
 * @author wener
 * @since 15/12/20
 */
public class KeyInvoke {
    private final static Logger log = LoggerFactory.getLogger(KeyInvoke.class);
    private final Register r3;
    private final Register r2;
    private final Register r1;
    private final Register r0;
    private InputManager inputManager;

    @Inject
    public KeyInvoke(@Named("R3") Register r3, @Named("R2") Register r2, @Named("R1") Register r1, @Named("R0") Register r0, InputManager inputManager) {
        this.r3 = r3;
        this.r2 = r2;
        this.r1 = r1;
        this.r0 = r0;
        this.inputManager = inputManager;
    }

    /*
10 | 键入整数 | 0 |  | r3的值变为键入的整数
11 | 键入字符串 | 0 | r3:目标字符串句柄 | r3所指字符串的内容变为键入的字符串
12 | 键入浮点数 | 0 |  | r3的值变为键入的浮点数

34 | 判定某键是否按下 | 0;r3 | r3:KEY |  KEYPRESS(KEY)
39 | 等待按键 | r3:按键 | - |  WAITKEY()
45 | 获取按键的字符串 | 0 | r3:字符串句柄,用于存储结果 |  InKey$
46 | 获取按键的ASCII码 | 0 | r3:KEYPRESS |  INKEY()
     */
    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 34, b = 0)
    public void keyPress() {
        r3.set(inputManager.isKeyPressed(r3.get()) ? 1 : 0);
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 39, b = 0)
    public void waitKey() {
        r3.set(inputManager.waitKey());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 45, b = 0)
    public void keyString() {
        r3.set(inputManager.inKeyString());
    }

    @SystemInvoke(type = SystemInvoke.Type.OUT, a = 46, b = 0)
    public void keyCode() {
        r3.set(inputManager.inKey());
    }

}
