package me.wener.bbvm.dev.swing;

import org.assertj.swing.core.MouseClickInfo;
import org.assertj.swing.edt.GuiActionRunner;
import org.assertj.swing.edt.GuiQuery;
import org.assertj.swing.fixture.FrameFixture;
import org.assertj.swing.junit.testcase.AssertJSwingJUnitTestCase;
import org.junit.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.swing.*;
import java.awt.event.KeyAdapter;
import java.awt.event.KeyEvent;
import java.awt.event.MouseAdapter;
import java.awt.event.MouseEvent;
import java.awt.image.BufferedImage;

import static java.awt.event.KeyEvent.*;
import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.Assert.assertEquals;

/**
 * @author wener
 * @since 15/12/27
 */
public class SwingInputMangerTest extends AssertJSwingJUnitTestCase {
    private final static Logger log = LoggerFactory.getLogger(SwingInputMangerTest.class);
    private SwingInputManger inputManger;
    private FrameFixture window;
    private SwingPage page;

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

    @Test
    public void testKeyCode() {
        assertEquals(VK_ENTER, charToKeyCode('\n'));
        assertEquals(VK_TAB, charToKeyCode('\t'));
    }

    @Override
    protected void onSetUp() {
        JFrame frame = GuiActionRunner.execute(new GuiQuery<JFrame>() {
            protected JFrame executeInEDT() {
                BufferedImage image = new BufferedImage(240, 320, BufferedImage.TYPE_INT_RGB);
                MainFrame frame = new MainFrame(() -> image);
                inputManger = new SwingInputManger();
                frame.addKeyListener(new KeyAdapter() {
                    @Override
                    public void keyPressed(KeyEvent e) {
                        inputManger.offer(e);
                    }

                    @Override
                    public void keyReleased(KeyEvent e) {
                        inputManger.offer(e);
                    }
                });
                frame.getImagePanel().addMouseListener(new MouseAdapter() {
                    @Override
                    public void mouseClicked(MouseEvent e) {
                        inputManger.offer(e);
                    }
                });
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
}
