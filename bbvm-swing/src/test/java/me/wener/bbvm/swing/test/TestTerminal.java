package me.wener.bbvm.swing.test;

import com.googlecode.lanterna.graphics.TextGraphics;
import com.googlecode.lanterna.gui.GUIScreen;
import com.googlecode.lanterna.gui.dialog.DialogButtons;
import com.googlecode.lanterna.gui.dialog.MessageBox;
import com.googlecode.lanterna.screen.TerminalScreen;
import com.googlecode.lanterna.terminal.DefaultTerminalFactory;
import com.googlecode.lanterna.terminal.swing.SwingTerminalFrame;
import org.junit.Test;

import java.io.IOException;

public class TestTerminal
{
    @Test
    public void test() throws InterruptedException, IOException
    {
        DefaultTerminalFactory factory = new DefaultTerminalFactory();
        SwingTerminalFrame terminal = (SwingTerminalFrame) factory.createTerminal();
        TerminalScreen screen = new TerminalScreen(terminal);
        screen.startScreen();
        terminal.setTitle("BeBasicVirtualMachine");
        GUIScreen gui = new GUIScreen(screen);
        MessageBox.showMessageBox(gui,"t","c",DialogButtons.OK);
        TextGraphics graphics = terminal.newTextGraphics();

        terminal.enterPrivateMode();
        terminal.clearScreen();
        terminal.putCharacter('我');
        terminal.setCursorPosition(10, 10);
        graphics.putString(70, 10, "我就是我,你就是你");

        terminal.setCursorPosition(70, 11);
        for (char c : "我就是我,你就是你".toCharArray())
        {
            terminal.putCharacter(c);
        }
        terminal.setCursorVisible(true);
        System.out.println(terminal.getFont());
        Thread.sleep(10000);
    }
}
