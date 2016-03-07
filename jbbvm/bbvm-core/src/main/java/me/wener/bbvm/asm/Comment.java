package me.wener.bbvm.asm;

import io.netty.buffer.ByteBuf;

/**
 * @author wener
 * @since 15/12/11
 */
public class Comment implements Assembly {
    Token token;

    public Comment(Token specialToken) {
        this.token = specialToken;
    }

    public static Comment createFor(Token token) {
        if (token.specialToken != null) {
            return new Comment(token.specialToken);
        }
        return null;
    }

    @Override
    public Type getType() {
        return Type.COMMENT;
    }

    @Override
    public String toAssembly() {
        return token.toString();
    }

    @Override
    public Comment getComment() {
        return this;
    }

    @Override
    public void setComment(Comment comment) {
        throw new UnsupportedOperationException();
    }

    @Override
    public boolean hasComment() {
        return true;
    }

    @Override
    public int length() {
        return 0;
    }

    @Override
    public void write(ByteBuf buf) {

    }

    @Override
    public int getLine() {
        return token.beginLine;
    }

    @Override
    public String toString() {
        return "Comment{" + token.toString() + "}";
    }
}
