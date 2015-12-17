package me.wener.bbvm.vm.res.file;

import me.wener.bbvm.vm.res.AbstractResource;

/**
 * @author wener
 * @since 15/12/17
 */
public class FileResource extends AbstractResource {
    private final FileManager manager;

    protected FileResource(int handler, FileManager manager) {
        super(handler);
        this.manager = manager;
    }

    @Override
    public FileManager getManager() {
        return manager;
    }

    @Override
    public void close() {
        manager.close(this);
    }

    public FileResource open(String string) {
        return this;
    }

    public int readInt(int address) {
        return 0;
    }

    public float readFloat(int address) {
        return 0;
    }

    public String readString(int address) {
        return "";
    }

    public FileResource writeInt(int address, int v) {
        return this;
    }

    public FileResource writeFloat(int address, float v) {
        return this;
    }

    public FileResource writeString(int address, String v) {
        return this;
    }

    public boolean isEof() {
        return false;
    }

    public int length() {
        return 0;
    }

    public int tell() {
        return 0;
    }

    public FileResource seek(int i) {
        return this;
    }
}
