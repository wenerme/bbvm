package me.wener.bbvm.vm.res.file;

import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.vm.res.ResourceManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.MappedByteBuffer;
import java.nio.channels.FileChannel;
import java.nio.file.Path;

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
        log.info("Open file {}", string);
        try {
            channel = FileChannel.open(path);
        } catch (IOException e) {
            throw new ExecutionException(e);
        }
        return this;
    }

    @Override
    public FileResource writeInt(int address, int v) {
        try {
            ByteBuffer buf = channel.map(FileChannel.MapMode.READ_WRITE, address, 4);
            buf.putInt(v);
        } catch (IOException e) {
            throw new ExecutionException(e);
        }
        return this;
    }

    @Override
    public FileResource writeFloat(int address, float v) {
        try {
            ByteBuffer buf = channel.map(FileChannel.MapMode.READ_WRITE, address, 4);
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
            ByteBuffer buf = channel.map(FileChannel.MapMode.READ_WRITE, address, bytes.length + 1);
            buf.put(bytes).put((byte) 0);
        } catch (IOException e) {
            throw new ExecutionException(e);
        }
        return this;
    }

    @Override
    public boolean isEof() {
        try {
            MappedByteBuffer buf = channel.map(FileChannel.MapMode.READ_ONLY, channel.position(), 1);
            return buf.get() < 0;
        } catch (IOException e) {
            throw new ExecutionException(e);
        }
    }

    @Override
    public int length() {
        try {
            return (int) channel.size();
        } catch (IOException e) {
            throw new ExecutionException(e);
        }
    }

    @Override
    public int readInt(int address) {
        try {
            ByteBuffer buf = channel.map(FileChannel.MapMode.READ_WRITE, address, 4);
            return buf.getInt();
        } catch (IOException e) {
            throw new ExecutionException(e);
        }
    }

    @Override
    public float readFloat(int address) {
        try {
            ByteBuffer buf = channel.map(FileChannel.MapMode.READ_WRITE, address, 4);
            return buf.getFloat();
        } catch (IOException e) {
            throw new ExecutionException(e);
        }
    }

    @Override
    public String readString(final int address) {
        int pos = address;
        ByteArrayOutputStream os = new ByteArrayOutputStream();
        try {
            byte b = -1;
            while (b != 0) {
                ByteBuffer buf = channel.map(FileChannel.MapMode.READ_WRITE, pos += 256, 256);
                while ((b = buf.get()) != 0) {
                    os.write(b);
                }
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
            return (int) channel.position();
        } catch (IOException e) {
            throw new ExecutionException(e);
        }
    }

    @Override
    public FileResource seek(int i) {
        try {
            channel.position(i);
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
    public ResourceManager getManager() {
        return manager;
    }

    @Override
    public void close() {
        if (channel != null) {
            log.info("Close file {}", path);
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
