package me.wener.bbvm.asm;

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

    enum Type {
        LABEL, COMMENT, PSEUDO, INST
    }
}
