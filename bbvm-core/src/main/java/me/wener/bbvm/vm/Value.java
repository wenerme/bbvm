package me.wener.bbvm.vm;

import me.wener.bbvm.util.val.IsInt;

/**
 * @author wener
 * @since 15/12/15
 */
public interface Value<T extends Value> {
    T set(int v);

    int get();

    default T set(IsInt v) {
        return set(v.asInt());
    }

    default T set(float v) {
        return set(Float.floatToRawIntBits(v));
    }

    default float getFloat() {
        return Float.intBitsToFloat(get());
    }

    default String getString() {
        return getVm().getString(get());
    }

    @SuppressWarnings("unchecked")
    default T set(String v) {
        getVm().getStringManager().getResource(get()).setValue(v);
        return (T) this;
    }

//    public T add()

    VM getVm();

    default T add(int v) {
        return set(get() + v);
    }

    default T subtract(int v) {
        return set(get() - v);
    }
}
