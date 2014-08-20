package me.wener.bbvm.core;

/**
 * 基本操作函数
 */
public interface BBBasicFunc
{
    float SIN(float X);

    float COS(float X);

    float TAN(float X);

    float SQR(float X);

    float ABS(float X);

    int ABS(int X);

    /**
     * 求取字符串长度
     *
     * @param X$ 字符串
     * @return 长度
     */
    int LEN(String X$);

    /**
     * 将数值转换成字符串
     *
     * @param V 数值
     * @return 字符串
     */
    String STR$(int V);

    /**
     * 将字符串转变成数值
     *
     * @param X$ 字符串
     * @return 数值
     */
    int VAL(String X$);

    /**
     * 把ASCII值转换成字符串
     *
     * @param X ASCII值
     */
    String CHR$(int X);

    /**
     * 把字符串转换成ASCII值
     * 这个函数可以用来求取字符的内码。
     *
     * @param X$ 字符串
     * @return ASCII值
     */
    int ASC(String X$);

    /**
     * 左取子字符串
     * 函数返回源字符串前N个字符的子字符串，N为截取长度。当N<0或N>源字符串长度时，将返回整个源字符串。
     *
     * @param X$ 字符串
     * @param N  截取长度
     * @return 结果字符串
     */
    String LEFT$(String X$, int N);

    /**
     * 右取子字符串
     * 该函数返回源字符串后N个字符的子字符串，N为截取长度。当N<0或N>源字符串长度时，将返回整个源字符串。
     *
     * @param X$ 源字符串
     * @param N  截取长度
     * @return 结果字符串
     */
    String RIGHT$(String X$, int N);

    /**
     * 从中间指定开始位置取指定长度的字符串
     * 该函数返回从指定位置开始的指定长度子字符串。当N<0或N>源字符串长度时，将返回整个源字符串；当S + N大于源字符串的长度时，将截取从S开始到字符串结束的字符串。
     *
     * @param X$ 源字符串
     * @param S  开始位置
     * @param N  截取长度
     * @return 结果字符串
     */
    String MID$(String X$, int S, int N);
}
