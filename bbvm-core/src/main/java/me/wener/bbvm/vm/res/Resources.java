package me.wener.bbvm.vm.res;

import com.google.common.base.Function;
import com.google.common.collect.Maps;
import com.google.common.eventbus.EventBus;
import com.google.common.eventbus.Subscribe;
import com.google.inject.AbstractModule;
import com.google.inject.Module;
import com.google.inject.TypeLiteral;
import com.google.inject.multibindings.OptionalBinder;
import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.exception.ResourceMissingException;
import me.wener.bbvm.vm.SystemInvokeManager;
import me.wener.bbvm.vm.event.ResetEvent;
import me.wener.bbvm.vm.invoke.FileInvoke;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import javax.inject.Singleton;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.channels.FileChannel;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;
import java.util.Map;

/**
 * @author wener
 * @since 15/12/18
 */
public class Resources {
    public static Module stringModule() {
        return new AbstractModule() {
            @Override
            protected void configure() {
                bind(StringManager.class).to(StringManagerImpl.class);
            }
        };
    }

    public static Class<? extends StringManager> stringManager() {
        return StringManagerImpl.class;
    }

    public static Module fileModule() {
        return new FileModule();
    }

    /**
     * @author wener
     * @since 15/12/13
     */
    private static class StringResourceImpl implements StringResource {
        final private StringManagerImpl manager;
        final private int handler;
        private String value;

        public StringResourceImpl(StringManagerImpl manager, int handler) {
            this.manager = manager;
            this.handler = handler;
        }

        @Override
        public int getHandler() {
            return handler;
        }

        @Override
        public StringManager getManager() {
            return manager;
        }

        @Override
        public void close() {
            manager.close(this);
        }

        public String getValue() {
            return value;
        }

        public StringResourceImpl setValue(String v) {
            value = v;
            return this;
        }

    }

    /**
     * Standard string resource implementation
     *
     * @author wener
     * @since 15/12/13
     */
    private static class StringManagerImpl implements StringManager {
        private final static Logger log = LoggerFactory.getLogger(StringManagerImpl.class);
        int handler = -1;
        // TODO Reuse handler ?
        Map<Integer, StringResource> resources = Maps.newHashMap();

        @Override
        public StringResource getResource(int handler) {
            StringResource resource = resources.get(handler);
            if (resource == null) {
                throw new ResourceMissingException(getType(), handler);
            }
            return resource;
        }

        @Override
        public StringManagerImpl reset() {
            handler = -1;
            resources.clear();
            return this;
        }

        void close(StringResource resource) {
            resources.remove(resource.getHandler());
        }

        @Override
        public StringResource create() {
            log.debug("Create string resource {}", handler);
            StringResource resource = new StringResourceImpl(this, handler--);
            resources.put(resource.getHandler(), resource);
            return resource;
        }

        @Inject
        void init(EventBus eventBus) {
            eventBus.register(this);
        }

        @Subscribe
        public void onReset(ResetEvent resetEvent) {
            log.debug("Reset all string resources");
            reset();
        }
    }

    /**
     * @author wener
     * @since 15/12/17
     */
    private static class JavaFileResourceImpl implements FileResource {
        private final static Logger log = LoggerFactory.getLogger(FileResource.class);
        private final int handler;
        private final FileManager manager;
        private Path path;
        private FileChannel channel;

        protected JavaFileResourceImpl(int handler, FileManager manager) {
            this.handler = handler;
            this.manager = manager;
        }

        @Override
        public FileResource open(String string) {
            log.info("Open file #{} {}", handler, string);
            try {
                path = Paths.get(string);
                channel = FileChannel.open(path, StandardOpenOption.READ, StandardOpenOption.WRITE, StandardOpenOption.CREATE);
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
            return this;
        }

        @Override
        public FileResource writeInt(int address, int v) {
            try {
                ByteBuffer buf = ch().map(FileChannel.MapMode.READ_WRITE, address, 4);
                buf.putInt(v);
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
            return this;
        }

        @Override
        public FileResource writeFloat(int address, float v) {
            try {
                ByteBuffer buf = ch().map(FileChannel.MapMode.READ_WRITE, address, 4);
                buf.putFloat(v);
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
            return this;
        }

        @Override
        public FileResource writeString(int address, String v) {
            try {
                byte[] bytes = v.getBytes();
                ByteBuffer buf = ch().map(FileChannel.MapMode.READ_WRITE, address, bytes.length + 1);
                buf.put(bytes).put((byte) 0);
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
            return this;
        }

        @Override
        public boolean isEof() {
            if (channel == null) {
                log.info("File handler #{} not open yet", handler);
                return true;
            }
            try {
                ByteBuffer buffer = ByteBuffer.allocate(1);
                boolean eof = channel.read(buffer) < 0;
                if (!eof) {
                    channel.position(channel.position() - 1);
                }
                return eof;
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
        }

        private FileChannel ch() {
            if (channel == null) {
                throw new ExecutionException(String.format("File handler #%s not open yet", handler));
            }
            return channel;
        }

        @Override
        public int length() {
            try {
                return (int) ch().size();
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
        }

        @Override
        public int readInt(int address) {
            try {
                ByteBuffer buf = ch().map(FileChannel.MapMode.READ_WRITE, address, 4);
                return buf.getInt();
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
        }

        @Override
        public float readFloat(int address) {
            try {
                ByteBuffer buf = ch().map(FileChannel.MapMode.READ_WRITE, address, 4);
                return buf.getFloat();
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
        }

        @Override
        public String readString(final int address) {
            ByteArrayOutputStream os = new ByteArrayOutputStream();
            try {
                FileChannel ch = ch().position(address);
                ByteBuffer buffer = ByteBuffer.allocate(1);
                while (ch.read(buffer) > 0) {
                    buffer.flip();
                    byte b = buffer.get();
                    if (b == 0) {
                        break;
                    }
                    os.write(b);
                    buffer.clear();
                }
                // TODO Charset
                return os.toString("UTF-8");
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
        }

        @Override
        public int tell() {
            try {
                return (int) ch().position();
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
        }

        @Override
        public FileResource seek(int i) {
            try {
                ch().position(i);
            } catch (IOException e) {
                throw new ExecutionException(e);
            }
            return this;
        }

        @Override
        public int getHandler() {
            return handler;
        }

        @Override
        public FileManager getManager() {
            return manager;
        }

        @Override
        public void close() {
            if (channel != null) {
                log.info("Close file #{} {}", handler, path);
                try {
                    channel.close();
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
            path = null;
            channel = null;
        }
    }

    /**
     * @author wener
     * @since 15/12/17
     */
    @Singleton
    public static class FileManagerImpl implements FileManager {
        static final int HANDLER_NUMBER = 10;
        static final String TYPE = "File";
        private final static Logger log = LoggerFactory.getLogger(FileManager.class);
        private final FileResource[] resources = new FileResource[10];
        private final Function<Integer, FileResource> creator;

        @Inject
        public FileManagerImpl(Function<Integer, FileResource> creator) {
            this.creator = creator;
            for (int i = 0; i < HANDLER_NUMBER; i++) {
                resources[i] = createNew(i);
            }
        }

        @Inject
        private void init(SystemInvokeManager systemInvokeManager) {
            systemInvokeManager.register(FileInvoke.class);
        }

        private FileResource createNew(int i) {
            return creator.apply(i);
        }

        @Override
        public FileResource getResource(int handler) {
            if (handler < 0 || handler > HANDLER_NUMBER) {
                throw new ExecutionException("No file resource for handler " + handler);
            }
            return resources[handler];
        }

        @Override
        public FileManagerImpl reset() {
            for (FileResource resource : resources) {
                resource.close();
            }
            return this;
        }

        @Override
        public FileResource create() {
            throw new ExecutionException("Resource for file is fixed");
        }

        @Override
        public String getType() {
            return TYPE;
        }

        @Inject
        void init(EventBus eventBus) {
            eventBus.register(this);
        }

        @Subscribe
        public void onReset(ResetEvent resetEvent) {
            log.debug("Reset all file resources");
            reset();
        }
    }

    /**
     * @author wener
     * @since 15/12/17
     */
    private static class FileModule extends AbstractModule {
        @Override
        protected void configure() {
            OptionalBinder.newOptionalBinder(binder(), new TypeLiteral<Function<Integer, FileResource>>() {
            }).setDefault().to(JavaFileResourceFunction.class);
            bind(FileManager.class).to(FileManagerImpl.class).asEagerSingleton();
        }

        private static class JavaFileResourceFunction implements Function<Integer, FileResource> {
            @Inject
            FileManager fileManager;

            @Override
            public FileResource apply(Integer input) {
                return new JavaFileResourceImpl(input, fileManager);
            }
        }
    }
}
