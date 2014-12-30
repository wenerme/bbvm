package me.wener.bbvm.utils.val;

class SimpleValue<T> implements IsValue<T>
{
    private T value;

    public SimpleValue(T value)
    {
        this.value = value;
    }

    @Override
    public T get()
    {
        return value;
    }
}
