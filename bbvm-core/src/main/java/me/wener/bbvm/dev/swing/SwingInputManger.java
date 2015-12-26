package me.wener.bbvm.dev.swing;

import com.google.common.eventbus.EventBus;
import com.google.common.eventbus.Subscribe;
import me.wener.bbvm.dev.InputManager;
import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.vm.event.ResetEvent;
import org.jetbrains.annotations.NotNull;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import javax.inject.Singleton;
import java.awt.event.*;
import java.util.PrimitiveIterator;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.SynchronousQueue;
import java.util.concurrent.TimeUnit;
import java.util.stream.IntStream;
import java.util.stream.Stream;

/**
 * @author wener
 * @since 15/12/26
 */
@Singleton
class SwingInputManger implements InputManager, MouseListener, KeyListener {
    private final static Logger log = LoggerFactory.getLogger(SwingInputManger.class);
    final protected BlockingQueue<InputEvent> events;
    private final PrimitiveIterator.OfInt keyCodeIterator;
    private final PrimitiveIterator.OfInt charIterator;
    SwingPage page;
    @Inject
    private SwingPageManager pageManager;

    SwingInputManger(BlockingQueue<InputEvent> events) {
        this.events = events;
        keyCodeIterator = getKeyCodeStream().iterator();
        charIterator = getCharStream().iterator();
    }

    @Inject
    public SwingInputManger() {
        this(new SynchronousQueue<>());
    }

    @Override
    public boolean isKeyPressed(int key) {
        return false;
    }

    @Override
    public int waitKey() {
        return keyCodeIterator.next();
    }

    @Override
    public int getLastKeyCode() {
        return 0;
    }

    @Override
    public String getLastKeyString() {
        return null;
    }

    @Override
    public void mouseClicked(MouseEvent e) {
        try {
            events.put(e);
        } catch (InterruptedException e1) {
            e1.printStackTrace();
        }
    }

    @Override
    public void mousePressed(MouseEvent e) {

    }

    @Override
    public void mouseReleased(MouseEvent e) {

    }

    @Override
    public void mouseEntered(MouseEvent e) {

    }

    @Override
    public void mouseExited(MouseEvent e) {

    }

    @Override
    public void keyTyped(KeyEvent e) {
    }

    @Override
    public void keyPressed(KeyEvent e) {
        offer(e);
    }

    private void offer(KeyEvent e) {
        try {
            if (!events.offer(e, 5, TimeUnit.MILLISECONDS) && log.isTraceEnabled()) {
                log.trace("Dropped {}", e);
            }
        } catch (InterruptedException ignored) {
        }
    }

    @Override
    public void keyReleased(KeyEvent e) {
        offer(e);
    }

    @Override
    public String readText() {
        StringBuilder sb = new StringBuilder();

        while (true) {
            int c = charIterator.nextInt();
            if (page != null) {
                page.draw((char) c);
            }
            if (c == '\n') {
                break;
            }
            sb.append((char) c);
        }
        return sb.toString();
    }

    @Inject
    void init(EventBus eventBus) {
        eventBus.register(this);
    }

    @Subscribe
    public void onVmReset(ResetEvent e) {
        page = pageManager.getScreen();
    }

    private IntStream getCharStream() {
        return generator()
                .filter(e -> e.getID() == KeyEvent.KEY_PRESSED && ((KeyEvent) e).getKeyChar() != KeyEvent.CHAR_UNDEFINED)
                .mapToInt(e -> ((KeyEvent) e).getKeyChar());
    }

    private IntStream getKeyCodeStream() {
        return generator()
                .filter(e -> e.getID() == KeyEvent.KEY_PRESSED)
                .mapToInt(e -> ((KeyEvent) e).getKeyCode());
    }

    @NotNull
    private Stream<InputEvent> generator() {
        return Stream
                .generate(() -> {
                    try {
                        InputEvent e = events.take();
                        if (log.isTraceEnabled()) {
                            log.trace("Got {}", e);
                        }
                        return e;
                    } catch (InterruptedException e) {
                        throw new ExecutionException(e);
                    }
                });
    }
}
