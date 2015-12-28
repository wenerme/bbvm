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
import java.awt.event.InputEvent;
import java.awt.event.KeyEvent;
import java.awt.event.MouseEvent;
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
class SwingInputManger implements InputManager {
    private final static Logger log = LoggerFactory.getLogger(SwingInputManger.class);
    final protected BlockingQueue<InputEvent> events;
    private final PrimitiveIterator.OfInt keyCodeIterator;
    private final PrimitiveIterator.OfInt charIterator;
    private final KeyStateTracker tracker;
    SwingPage page;
    @Inject
    private SwingPageManager pageManager;

    SwingInputManger(BlockingQueue<InputEvent> events) {
        this.events = events;
        keyCodeIterator = getKeyCodeStream().iterator();
        charIterator = getCharStream().iterator();
        tracker = new KeyStateTracker();
    }

    @Inject
    public SwingInputManger() {
        this(new SynchronousQueue<>());
    }

    @Override
    public boolean isKeyPressed(int key) {
        return tracker.isPressed(key);
    }

    @Override
    public int waitKey() {
        return keyCodeIterator.next();
    }

    @Override
    public int inKey() {
        try {
            InputEvent e = events.poll(5, TimeUnit.MILLISECONDS);
            if (e == null) {
                return 0;
            }
            switch (e.getID()) {
                case KeyEvent.KEY_PRESSED:
                    return ((KeyEvent) e).getKeyCode();
                case MouseEvent.MOUSE_CLICKED:
                    return makeClickKey(((MouseEvent) e).getX(), ((MouseEvent) e).getY());
            }
        } catch (InterruptedException e) {
            throw new ExecutionException(e);
        }
        return 0;
    }

    @Override
    public String inKeyString() {
        try {
            InputEvent e = events.poll(5, TimeUnit.MILLISECONDS);
            if (e == null) {
                return "";
            }
            if (e.getID() == KeyEvent.KEY_PRESSED) {
                int c = ((KeyEvent) e).getKeyCode();
                if (c == KeyEvent.CHAR_UNDEFINED) {
                    return "";
                }
                return String.valueOf((char) c);
            }
        } catch (InterruptedException e) {
            throw new ExecutionException(e);
        }
        return "";
    }

    /**
     * This method is thread safe
     */
    public void offer(InputEvent e) {
        tracker.offer(e);
        try {
            boolean offer = events.offer(e, 5, TimeUnit.MILLISECONDS);
            if (log.isTraceEnabled()) {
                log.trace("{} {}", offer ? "Got" : "Drooped", e);
            }
        } catch (InterruptedException ignored) {
        }
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
                .filter(e -> e.getID() == KeyEvent.KEY_PRESSED || e.getID() == MouseEvent.MOUSE_CLICKED)
                .mapToInt(e -> {
                    switch (e.getID()) {
                        case KeyEvent.KEY_PRESSED:
                            return ((KeyEvent) e).getKeyCode();
                        case MouseEvent.MOUSE_CLICKED:
                            return makeClickKey(((MouseEvent) e).getX(), ((MouseEvent) e).getY());
                        default:
                            throw new AssertionError();
                    }
                });
    }

    @NotNull
    private Stream<InputEvent> generator() {
        return Stream
                .generate(() -> {
                    try {
                        return events.take();
                    } catch (InterruptedException e) {
                        throw new ExecutionException(e);
                    }
                });
    }
}
