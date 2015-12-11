package me.wener.bbvm.asm;

import io.netty.buffer.ByteBuf;

/**
 * @author wener
 * @since 15/12/11
 */
public class Value {
    int integerValue;
    long longValue;
    String string;
    byte[] bytes;
    Token token;
    Type type;
    String assembly;

    public static Value forInteger(String v) {
        Value val = new Value();
        val.longValue = Long.parseLong(v);
        val.type = Type.INTEGER;
        val.assembly = v;
        return val;
    }

    public static Value forInteger(Token token) {
        Value val = forInteger(token.image);
        val.token = token;
        return val;
    }

    public static Value forFloat(Token token) {
        return null;
    }

    public static Value forString(Token token) {
        Value val = new Value().setType(Type.STRING).setToken(token);
        val.string = token.image.substring(1, token.image.length() - 1);
        return val;
    }

    public Value setType(Type type) {
        this.type = type;
        return this;
    }

    public Value setToken(Token token) {
        this.token = token;
        return this;
    }

    void write(ByteBuf buf) {

    }

    String toAssembly() {
        return token.toString();
    }

    public int asInt() {
        return 0;
    }

    enum Type {
        INTEGER, LONG, FLOAT, DOUBLE, BYTES, STRING
    }
}
