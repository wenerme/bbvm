package me.wener.bbvm.bbasm;

public interface Compilable
{
    /**
     * @return 编译结果
     */
    byte[] toBinary();

    /**
     * @return 是否延迟求值
     */
    boolean isLazy();

    /**
     * @return 结果长度
     */
    int length();
}
