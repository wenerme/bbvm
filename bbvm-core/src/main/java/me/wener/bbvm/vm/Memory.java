package me.wener.bbvm.vm;

import com.google.common.base.MoreObjects;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;

import java.nio.ByteOrder;
import java.nio.charset.Charset;

/**
 * @author wener
 * @since 15/12/10
 */
public class Memory {
    public static final int DEFAULT_MEMORY = 1024 * 1024 * 4;
    public static final int DEFAULT_STACK_SIZE = 1000;
    Register rs;
    Register rb;
    private ByteBuf mem;
    private int memorySize;
    private int stackSize;
    private VM vm;

    private Memory() {
        this(DEFAULT_MEMORY, DEFAULT_STACK_SIZE);
    }

    private Memory(int memorySize, int stackSize) {
        mem = Unpooled.buffer(memorySize + stackSize, memorySize + stackSize).order(ByteOrder.LITTLE_ENDIAN);
        this.stackSize = stackSize;
        this.memorySize = memorySize;
    }

    public static Memory load(ByteBuf buf) {
        Memory memory = new Memory(buf.readableBytes(), DEFAULT_STACK_SIZE);
        memory.mem.writeBytes(buf);
        return memory;
    }

    /**
     * Fill memory with zero
     */
    public Memory clear() {
        mem.clear();
        byte[] bytes = mem.array();
        for (int i = 0; i < bytes.length; i++) {
            bytes[i] = 0;
        }
        return this;
    }

    public VM getVm() {
        return vm;
    }

    public Memory setVm(VM vm) {
        this.vm = vm;
        this.rs = vm.rs;
        this.rb = vm.rb;
        return this;
    }

    public int pop() {
        rs.add(4);
        return mem.getInt(rs.get());
    }

    public void push(int v) {
        // TODO Stack bounds check
//        if (rb.intValue() - rs.intValue() < 4) {
//            throw new RuntimeException("Stack overflow");
//        }
        mem.setInt(rs.get(), v);
        rs.subtract(4);
    }

    public ByteBuf getByteBuf() {
        return mem;
    }

    public int getMemorySize() {
        return memorySize;
    }

    public int getStackSize() {
        return stackSize;
    }

    @Override
    public String toString() {
        return MoreObjects.toStringHelper(this)
            .add("mem", mem)
            .add("rs", rs)
            .add("rb", rb)
            .add("memorySize", memorySize)
            .add("stackSize", stackSize)
            .toString();
    }

    public int read(int addr) {
        return mem.getInt(addr);
    }

    public Memory write(int addr, int v) {
        mem.setInt(addr, v);
        return this;
    }

    public String getString(int i, Charset charset) {
        int end = i;
        byte[] bytes = mem.array();
        //noinspection StatementWithEmptyBody
        while (bytes[end++] != 0) {
            // Ignored
        }
        int length = end - i - 1;
        if (length == 0) {
            return "";
        }
        return new String(bytes, i, length, charset);
    }
}
