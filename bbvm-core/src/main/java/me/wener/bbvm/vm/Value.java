package me.wener.bbvm.vm;

import me.wener.bbvm.dev.Resource;
import me.wener.bbvm.dev.ResourceManager;
import me.wener.bbvm.util.IsInt;

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

    /**
     * @param def Default value if getString return null
     */
    default String getString(String def) {
        String s = getString();
        return s == null ? def : s;
    }

    default <M extends ResourceManager<M, R>, R extends Resource> R get(M manager) {
        return manager.getResource(get());
    }

    @SuppressWarnings("unchecked")
    default T set(String v) {
        getVm().getStringManager().getResource(get()).setValue(v);
        return (T) this;
    }

    VM getVm();

    default T add(int v) {
        return set(get() + v);
    }

    default T subtract(int v) {
        return set(get() - v);
    }
}
