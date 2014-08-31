package me.wener.bbvm.core;

/**
 * 字符串句柄资源池需要使用负数
 */
public class StringHandlePool extends ResourceNegativeHandlePool<StringHandle>
{
    protected StringHandlePool()
    {
        super(-1);
    }

    @Override
    public StringHandle createResource()
    {
        return new StringHandle();
    }

    @Override
    public void freeResource(StringHandle resource)
    {

    }

    public static void main(String[] args)
    {
        StringHandlePool pool = new StringHandlePool();
        int h = pool.acquire();
        pool.getResource(h).set("wener");
        System.out.println(h);// -1
        System.out.println(pool.acquire());// -2
        pool.release(h);
        System.out.println(pool.acquire());// -1
        System.out.println(pool.acquire());// -3
    }
}
