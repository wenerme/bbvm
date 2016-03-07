package me.wener.bbvm.dev.swing;

import me.wener.bbvm.dev.AbstractFileManager;
import me.wener.bbvm.dev.FileManager;
import me.wener.bbvm.dev.FileResource;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Singleton;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.channels.FileChannel;
import java.nio.charset.Charset;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;

/**
 * @author wener
 * @since 15/12/26
 */
@Singleton
public class JavaFileManager extends AbstractFileManager {

    @Override
    protected FileResource createResource(int i) {
        return new JavaFileResource(i, this, charset);
    }

    static class JavaFileResource implements FileResource {
        private final static Logger log = LoggerFactory.getLogger(FileResource.class);
        private final int handler;
        private final FileManager manager;
        private final Charset charset;
        private Path path;
        private FileChannel channel;

        protected JavaFileResource(int handler, FileManager manager, Charset charset) {
            this.handler = handler;
            this.manager = manager;
            this.charset = charset;
        }

        @Override
        public FileResource open(String string) throws IOException {
            log.info("Open file #{} {}", handler, string);
            path = Paths.get(string);
            channel = FileChannel.open(path, StandardOpenOption.READ, StandardOpenOption.WRITE, StandardOpenOption.CREATE);
            return this;
        }

        @Override
        public FileResource writeInt(int address, int v) throws IOException {
            return seek(address).writeInt(v);
        }

        @Override
        public FileResource writeString(int address, String v) throws IOException {
            return seek(address).writeString(v);
        }

        @Override
        public FileResource writeInt(int v) throws IOException {
            ByteBuffer buffer = ByteBuffer.allocate(4);
            buffer.order(ByteOrder.LITTLE_ENDIAN);
            buffer.putInt(v).flip();
            ch().write(buffer);
            return this;
        }

        @Override
        public FileResource writeString(String v) throws IOException {
            byte[] bytes = v.getBytes(charset);
            ByteBuffer buffer = ByteBuffer.allocate(bytes.length + 1);
            buffer.order(ByteOrder.LITTLE_ENDIAN);
            buffer.put(bytes).put((byte) 0).flip();
            ch().write(buffer);
            return this;
        }

        @Override
        public boolean isEof() throws IOException {
//            if (channel == null) {
//                log.info("File handler #{} not open yet", handler);
//                return true;
//            }
//            ByteBuffer buffer = ByteBuffer.allocate(1);
//            boolean eof = ch().read(buffer) < 0;
//            if (!eof) {
//                // Reverse position if we advance.
//                channel.position(channel.position() - 1);
//            }
//            return eof;
            return ch().position() == channel.size();
        }

        private FileChannel ch() throws IOException {
            if (channel == null) {
                throw new IOException(String.format("File handler #%s not open yet", handler));
            }
            return channel;
        }

        @Override
        public int length() throws IOException {
            return (int) ch().size();
        }

        @Override
        public int readInt() throws IOException {
            ByteBuffer buffer = ByteBuffer.allocate(4);
            buffer.order(ByteOrder.LITTLE_ENDIAN);
            ch().read(buffer);
            buffer.flip();
            return buffer.getInt();
        }

        @Override
        public String readString() throws IOException {
            ByteArrayOutputStream os = new ByteArrayOutputStream();
            ByteBuffer buffer = ByteBuffer.allocate(1);
            while (ch().read(buffer) > 0) {
                buffer.flip();
                byte b = buffer.get();
                if (b == 0) {
                    break;
                }
                os.write(b);
                buffer.clear();
            }
            return os.toString(charset.name());
        }

        @Override
        public int tell() throws IOException {
            return (int) ch().position();
        }

        @Override
        public FileResource seek(int i) throws IOException {
            ch().position(i);
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
}
