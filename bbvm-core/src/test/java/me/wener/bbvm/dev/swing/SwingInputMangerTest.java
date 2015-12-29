package me.wener.bbvm.dev.swing;

import org.assertj.swing.core.MouseClickInfo;
import org.assertj.swing.edt.GuiActionRunner;
import org.assertj.swing.edt.GuiQuery;
import org.assertj.swing.fixture.FrameFixture;
import org.assertj.swing.junit.testcase.AssertJSwingJUnitTestCase;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.awt.*;
import java.awt.image.BufferedImage;
import java.util.concurrent.atomic.AtomicInteger;

import static java.awt.event.KeyEvent.*;
import static org.assertj.core.api.Assertions.assertThat;

/**
 * @author wener
 * @since 15/12/27
 */
public class SwingInputMangerTest extends AssertJSwingJUnitTestCase {
    private final static Logger log = LoggerFactory.getLogger(SwingInputMangerTest.class);
    private SwingInputManger inputManger;
    private FrameFixture window;
    private SwingPage page;
    private MainFrame frame;

    public static int charToKeyCode(int c) {
        if (Character.isJavaIdentifierPart(c)) {
            if (Character.isUpperCase(c)) {
                throw new RuntimeException("Not allowed");
            }
            return Character.toUpperCase(c);
        }
        return c;
    }

    public static int[] keyCodes(String str) {
        return str.chars().map(SwingInputMangerTest::charToKeyCode).toArray();
    }

    @Override
    protected void onSetUp() {
        frame = GuiActionRunner.execute(new GuiQuery<MainFrame>() {
            protected MainFrame executeInEDT() {
                BufferedImage image = new BufferedImage(240, 320, BufferedImage.TYPE_INT_RGB);
                MainFrame frame = new MainFrame(() -> image);
                inputManger = new SwingInputManger();
                inputManger.bindKeyEvent(frame);
                inputManger.bindMouseEvent(frame.getImagePanel());
                inputManger.page = page;
                return frame;
            }
        });
        // IMPORTANT: note the call to 'robot()'
        // we must use the Robot from AssertJSwingJUnitTestCase
        window = new FrameFixture(robot(), frame);
        window.show(); // shows the frame to test
    }

    @Test(timeout = 60000)
    public void testReadText() throws Exception {
        Thread thread = new Thread(() -> {
            assertThat(inputManger.readText()).isEqualTo("Wener is grate");
            assertThat(inputManger.readText()).isEqualTo("so\tgood");
        });
        thread.start();
        window.requireFocused()
                .pressKey(VK_SHIFT)
                .pressKey(VK_W)
                .click(MouseClickInfo.leftButton().times(5))
                .releaseKey(VK_W)
                .releaseKey(VK_SHIFT)
                .click()
                .pressAndReleaseKeys(VK_F1)// Will ignore
                .pressAndReleaseKeys(keyCodes("ener is grate\n"))
                // Not works with VK_TAB
                .pressAndReleaseKeys(keyCodes("so\tgood\n"))
        ;
        thread.join();
    }

    @Test(timeout = 60000)
    public void testWaitKey() throws Exception {
        Thread thread = new Thread(() -> {
            assertThat(inputManger.waitKey()).isEqualTo('A');
            assertThat(inputManger.waitKey()).isEqualTo('B');
            assertThat(inputManger.waitKey()).isEqualTo(VK_SHIFT);
            assertThat(inputManger.waitKey()).isEqualTo(VK_TAB);
        });
        thread.start();

        window.requireFocused()
                .pressAndReleaseKeys(VK_A, VK_B, VK_SHIFT, VK_TAB);

        thread.join();
    }

    @Test(timeout = 60000)
    public void testMouseClick() throws InterruptedException {
        Thread thread = new Thread(() -> {
            int x = 0;
            int y = 0;
            do {
                int pos = inputManger.inKey();
                if (pos == 0) {
                    continue;
                }
                x = inputManger.getX(pos);
                y = inputManger.getY(pos);
            } while (x != 20 && y != 20);
        });
        thread.start();

        while (thread.isAlive()) {
            robot().click(frame.getImagePanel(), new Point(20, 20));
        }
    }

    @Test(timeout = 60000)
    public void testKeyPressed() throws InterruptedException {
        AtomicInteger step = new AtomicInteger(0);
        Thread thread = new Thread(() -> {
            while (true) {
                switch (step.get()) {
                    case 0:
                        if (inputManger.isKeyPressed('A')) {
                            step.incrementAndGet();
                            log.debug("Step {}",step);
                        }
                        break;
                    case 1:
                        if (inputManger.isKeyPressed('B')) {
                            step.incrementAndGet();
                            log.debug("Step {}",step);
                        }
                        break;
                    case 2:
                        try {
                            Thread.sleep(100);
                        } catch (InterruptedException ignored) {
                        }
                        if (!inputManger.isKeyPressed('B')) {
                            step.incrementAndGet();
                            log.debug("Step {}",step);
                        }
                        break;
                    default:
                        log.debug("Step {} complete",step);
                        return;
                }
            }
        });
        thread.start();
        while (true) {
            switch (step.get()) {
                case 0:
                    window.pressAndReleaseKeys('A');
                    break;
                case 1:
                    window.releaseKey('B');
                    window.pressKey('B');
                    break;
                case 2:
                    Thread.sleep(10);
                    break;
                default:
                    thread.join();
                    return;
            }
        }
    }
}
