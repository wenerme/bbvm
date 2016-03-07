package me.wener.bbvm.asm;

import com.google.common.base.Preconditions;
import io.netty.buffer.ByteBuf;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * @author wener
 * @since 15/12/11
 */
public class PseudoBlock extends AbstractAssembly implements Assembly {
    private final static Logger log = LoggerFactory.getLogger(PseudoBlock.class);
    private int length;
    private int value;
    private int line;

    public PseudoBlock() {
    }

    public PseudoBlock(int length, int value) {
        this.length = length;
        this.value = value;
        Preconditions.checkArgument(length >= 0);
        if (value < 0 || value > 0xff) {
            log.warn("BLOCK value {} will cast to {} ", value, value & 0xff);
            this.value = value & 0xff;
        }
    }

    @Override
    public Type getType() {
        return Type.PSEUDO;
    }

    @Override
    public String toAssembly() {
        return ".BLOCK " + length + " " + value;
    }

    @Override
    public void write(ByteBuf buf) {
        for (int i = 0; i < length; i++) {
            buf.writeByte(value);
        }
    }

    @Override
    public int getLine() {
        return line;
    }

    public PseudoBlock setLine(int line) {
        this.line = line;
        return this;
    }

    @Override
    public int length() {
        return length;
    }
}
