package me.wener.bbvm.core;

/**
 * BB 虚拟机相关的函数
 */
@SuppressWarnings("unused")
public interface BBVMFunc
{

    /**
     * 打开一个待操作的文件<br>
     * <b>语法</b>OPEN 文件名 FOR BINARY AS #文件号
     * <ul>
     * <li>OPEN是打开文件的关键字。</li>
     * <li>文件名是要打开的文件。</li>
     * <li>FOR BINARY是将文件以二进制方式打开。由于BB只支持二进制，这里打开其他类型的文件不在赘叙。</li>
     * <li>AS #文件号，是指定打开的文件所占据的文件号，后边对这个文件的操作都是通过操作这个文件号进行的。</li>
     * </ul>
     * <p/>
     * <pre>
     *     例如：OPEN "hszj.sav" FOR BIANRY AS #1
     *     这句是出现在《幻兽战记》代码中的一行代码，意思是使用1号通道，打开幻兽战记的存盘文件"hszj.sav"，以便对其进行操作。
     * </pre>
     */
    void OPEN(String filename, int mode, int fnum);

    /**
     * 关闭一个已经打开的文件<br>
     * <b>语法：</b>CLOSE #文件号<br>
     * CLOSE语句后面的文件号，就是在OPEN语句中，指定的文件号。
     * 所有的文件，在对其操作完毕后，都要使用CLOSE语句来关闭它。
     * QB中支持使用不带文件号的CLOSE关闭所有打开的文件，目前在BB中尚未支持这个功能。
     *
     * @param fnum 文件号
     */
    void CLOSE(int fnum);

    /**
     * 向一个打开的文件中写入数据<br>
     * <b>语法：</b>PUT #文件号,常量或变量
     *
     * @param fnum  文件号
     * @param value 写入的值
     */
    void PUT(int fnum, Object value);

    /**
     * 从一个打开的文件中读取数据<br>
     * <b>语法：</b>GET #文件号,变量名
     *
     * @param fnum 文件号
     * @return 获取到的值
     */
    Object GET(int fnum);

    /**
     * 判断文件指针是否已到文件结束位置
     *
     * @param fnum 文件号
     * @return 是否已到文件结束位置
     */
    boolean EOF(int fnum);

    /**
     * 获取文件的字节长度。
     *
     * @param fnum 文件号
     * @return 文件的字节长度
     */
    int LOF(int fnum);

    /**
     * 获取文件指针的当前位置。
     *
     * @param fnum 文件号
     * @return 文件指针的当前位置
     */
    int LOC(int fnum);

    /**
     * 设置随机数种子
     *
     * @param SEED 随机数种子
     */
    void RANDOMIZE(int SEED);

    /**
     * 获取随机数
     *
     * @param RANGE 随机数范围
     * @return 随机数
     */
    int RND(int RANGE);

    /**
     * 获取滴答数
     *
     * @return 运行到现在的毫秒数
     */
    int GETTICK();

    /**
     * 等待指定时间
     *
     * @param MSEC 毫秒
     */
    void MSDELAY(int MSEC);

    int PEEK(int ADDRESS);

    /**
     * 获取当前环境
     *
     * @return 当前环境
     */
    int GETENV();

    /**
     * 测试是否在VM下
     */
    boolean VMTEST();
}

