package me.wener.bbvm.vm.res.file;

import me.wener.bbvm.exception.ExecutionException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.channels.FileChannel;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;

/**
 * @author wener
 * @since 15/12/17
 */
public class JavaFileResourceImpl implements FileResource {
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
