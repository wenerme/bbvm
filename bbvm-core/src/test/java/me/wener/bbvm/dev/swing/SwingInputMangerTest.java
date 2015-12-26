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
                JFrame frame = new JFrame(SwingInputMangerTest.class.getSimpleName());
                ImageIcon icon = new ImageIcon(image);
                JLabel label = new JLabel(icon);
                frame.setLocationRelativeTo(null);
                frame.setFocusTraversalKeysEnabled(false);// Make VK_TAB works
                label.setLocation(0, 0);
                frame.getContentPane().add(label);
                frame.setDefaultCloseOperation(WindowConstants.EXIT_ON_CLOSE);
                frame.setResizable(false);
//                frame.setVisible(true);
                frame.pack();
                frame.repaint();
//                page = new SwingPage(1, null, image);
//                Graphics2D d = page.g;
//                d.setColor(Color.DARK_GRAY);
                inputManger = new SwingInputManger();
                frame.addKeyListener(inputManger);
                frame.addMouseListener(inputManger);
                inputManger.page = page;
//                page.fill();
//                page.draw("Hello\n");
//                new Thread(() -> {
//                    try {
//                        while (true) {
//                            label.repaint();
//                            Thread.sleep(1000 / 60);
//                        }
//                    } catch (Exception e) {
//                        Throwables.propagate(e);
//                    }
//                }).start();
                return frame;
            }
        });
        // IMPORTANT: note the call to 'robot()'
        // we must use the Robot from AssertJSwingJUnitTestCase
        window = new FrameFixture(robot(), frame);
        window.show(); // shows the frame to test
    }

    @Test(timeout = 20000)
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

    @Test(timeout = 10000)
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
