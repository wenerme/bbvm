package me.wener.bbvm.asm;

import io.netty.buffer.ByteBuf;

/**
 * @author wener
 * @since 15/12/11
 */
public interface Assembly {
    Type getType();

    String toAssembly();

    /**
     * @return {@link Comment} for this {@link Assembly} or {@code null}
     */
    Comment getComment();

    void setComment(Comment comment);

    boolean hasComment();

    /**
     * @return The length of this assembly
     */
    int length();

    void write(ByteBuf buf);

    int getLine();

    enum Type {
        LABEL, COMMENT, PSEUDO, INST
    }
}
