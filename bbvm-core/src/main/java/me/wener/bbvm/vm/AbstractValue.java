package me.wener.bbvm.vm;

import me.wener.bbvm.util.val.IsInt;

/**
 * @author wener
 * @since 15/12/13
 */
@SuppressWarnings("unchecked")
abstract class AbstractValue<T extends AbstractValue> {
    abstract T set(int v);

    abstract int get();

    public T set(IsInt v) {
        return set(v.asInt());
    }

    public T set(float v) {
        return set(Float.floatToRawIntBits(v));
    }

    public float getFloat() {
        return Float.intBitsToFloat(get());
    }

    public String getString() {
        return getVm().getString(get());
    }

    public T set(String v) {
        getVm().getStringManager().getResource(get()).setValue(v);
        return (T) this;
    }

//    public T add()

    abstract VM getVm();

    public T add(int v) {
        return set(get() + v);
    }

    public T subtract(int v) {
        return set(get() - v);
    }
}
