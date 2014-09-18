package me.wener.bbvm.swing;

import static com.google.common.base.Preconditions.*;

import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import me.wener.bbvm.core.ResourceHandlePool;

public class AutoResourceHandlePool<T> extends ResourceHandlePool<T>
{
    private final Class<? extends T> type;
    private final boolean autoClose;

    public AutoResourceHandlePool(int maxsize, boolean autoClose, Class<? extends T> type)
    {
        super(maxsize);
        this.type = type;
        this.autoClose = autoClose;
    }
    @SuppressWarnings("unchecked")
    public AutoResourceHandlePool(int maxsize, boolean autoClose)
    {
        super(maxsize);
        this.type = (Class<T>) capture();
        this.autoClose = autoClose;
    }

    /**
     * Returns the captured type.
     */
    @SuppressWarnings({"ConstantConditions"})
    final Type capture()
    {
        Type superclass = getClass().getGenericSuperclass();
        checkArgument(superclass instanceof ParameterizedType,
                "%s isn't parameterized", superclass);
        return ((ParameterizedType) superclass).getActualTypeArguments()[0];
    }

    @Override
    public T createResource()
    {
        try
        {
            return type.newInstance();
        } catch (Exception e)
        {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void freeResource(T resource)
    {
        if (autoClose && resource instanceof  AutoCloseable)
        {
            try
            {
                ((AutoCloseable) resource).close();
            } catch (Exception e)
            {
                e.printStackTrace();
            }
        }
    }
}
