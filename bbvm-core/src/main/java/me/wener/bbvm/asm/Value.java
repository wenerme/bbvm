package me.wener.bbvm.asm;

import com.google.common.base.Preconditions;
import com.google.common.io.BaseEncoding;
import io.netty.buffer.ByteBuf;
import org.apache.commons.lang3.mutable.MutableDouble;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;

/**
 * @author wener
 * @since 15/12/11
 */
public class Value {

    final MutableDouble number = new MutableDouble();
    String string;
    byte[] bytes;
    Token token;
    Type type;
    String assembly;
    Charset charset = StandardCharsets.UTF_8;

    public static Value forNumber(Token v) {
        return forNumber(v.image).setToken(v);
    }

    public static Value forNumber(String v) {
        Type type = null;
        double doubleValue = Double.parseDouble(v);

        switch (v.substring(v.length() - 1)) {
            case "l":
            case "L":
                type = Type.LONG;
                break;
            case "f":
            case "F":
                type = Type.FLOAT;
                break;
            case "d":
            case "D":
                type = Type.DOUBLE;
                break;
        }

        if (type == null) {
            if (doubleValue == (int) doubleValue) {
                type = Type.INTEGER;
            } else if (doubleValue == (float) doubleValue) {
                type = Type.FLOAT;
            } else if (doubleValue == (long) doubleValue) {
                type = Type.LONG;
            } else {
                type = Type.DOUBLE;
            }
        }

        Value value = new Value();
        value.number.setValue(doubleValue);
        return value.setAssembly(v).setType(type);
    }

    public static Value forNumber(String v, Type type) {
        return new Value().setAssembly(v).setValue(Double.parseDouble(v)).setType(type);
    }

    public static Value forNumber(Token token, Type type) {
        return forNumber(token.image, type).setToken(token);
    }

    public static Value forHexBytes(Token token) {
        String c = token.image;
        // Currently only support %
        switch (c.codePointAt(0)) {
            case '%':
                c = c.substring(1, c.length() - 1);
                if (c.length() % 2 != 0) {
                    throw new RuntimeException(String.format("%s:%s Hex bytes length value (%s)'%s'", token.beginLine, token.beginColumn, c.length(), c));
                }
                return new Value().setToken(token).setValue(BaseEncoding.base16().decode(c));
            default:
                throw new AssertionError();
        }
    }

    public static Value forString(Token token) {
        Value val = new Value().setAssembly(token.image).setType(Type.STRING).setToken(token);
        // Strip the quote
        val.string = token.image.substring(1, token.image.length() - 1);
        return val;
    }

    public Value setCharset(Charset charset) {
        this.charset = charset;
        return this;
    }

    public int length() {
        switch (type) {
            case INTEGER:
            case FLOAT:
                return 4;
            case LONG:
            case DOUBLE:
                return 8;
            case STRING:
                return string.getBytes(charset).length;
            case BYTES:
                return bytes.length;
            default:
                throw new AssertionError();
        }
    }

    public Value setAssembly(String assembly) {
        this.assembly = assembly;
        return this;
    }

    public Type getType() {
        return type;
    }

    public Value setType(Type type) {
        this.type = type;
        return this;
    }

    public Value setToken(Token token) {
        this.token = token;
        assembly = token.image;
        return this;
    }

    void write(ByteBuf buf) {
        switch (type) {
            case INTEGER:
                buf.writeInt(number.intValue());
                break;
            case LONG:
                buf.writeLong(number.longValue());
                break;
            case FLOAT:
                buf.writeFloat(number.floatValue());
                break;
            case DOUBLE:
                buf.writeDouble(number.doubleValue());
                break;
            case BYTES:
                buf.writeBytes(bytes);
                break;
            case STRING:
                buf.writeBytes(string.getBytes(charset));
                break;
            default:
                throw new AssertionError();
        }
    }

    public String toAssembly() {
        return assembly;
    }

    public int checkedInteger() {
        if (type != Type.INTEGER) {
            throw new RuntimeException("Need integer but got " + type);
        }
        return number.intValue();
    }

    public Number asNumber() {
        if (type.isNumber()) {
            return number;
        }
        throw new RuntimeException(this + " is not a number");
    }

    public void add(Value o) {
        if (o.type == Type.STRING && type == Type.STRING) {
            string += o.string;
            return;
        }

        type = type.higherNumberType(o.type);
        number.add(o.number);
    }

    public Object getValue() {
        switch (type) {
            case INTEGER:
                return number.intValue();
            case LONG:
                return number.longValue();
            case FLOAT:
                return number.floatValue();
            case DOUBLE:
                return number.doubleValue();
            case BYTES:
                return bytes;
            case STRING:
                return string;
        }
        throw new AssertionError();
    }

    public Value setValue(Object value) {
        if (value instanceof Integer) {
            number.setValue((int) value);
            type = Type.INTEGER;
        } else if (value instanceof Long) {
            number.setValue((long) value);
            type = Type.LONG;
        } else if (value instanceof Float) {
            number.setValue((float) value);
            type = Type.FLOAT;
        } else if (value instanceof Double) {
            number.setValue((double) value);
            type = Type.DOUBLE;
        } else if (value instanceof CharSequence) {
            string = value.toString();
            type = Type.STRING;
        } else if (value instanceof byte[]) {
            bytes = (byte[]) value;
            type = Type.BYTES;
        } else {
            throw new UnsupportedOperationException("Unsupported type " + value.getClass() + " : " + value);
        }
        return this;
    }

    enum Type {
        INTEGER, LONG, FLOAT, DOUBLE, BYTES, STRING;

        public boolean isNumber() {
            switch (this) {
                case INTEGER:
                case LONG:
                case FLOAT:
                case DOUBLE:
                    return true;
            }
            return false;
        }

        public Type higherNumberType(Type o) {
            Preconditions.checkArgument(isNumber() && o.isNumber(), this + " & " + o);
            if (ordinal() > o.ordinal()) {
                return this;
            } else {
                return o;
            }
        }
    }
}
