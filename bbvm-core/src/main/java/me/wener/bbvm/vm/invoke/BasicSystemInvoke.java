package me.wener.bbvm.vm.invoke;

import me.wener.bbvm.vm.Operand;
import me.wener.bbvm.vm.Register;
import me.wener.bbvm.vm.SystemInvoke;
import me.wener.bbvm.vm.res.StringManager;

import javax.inject.Named;

/**
 * @author wener
 * @since 15/12/13
 */
public class BasicSystemInvoke {
    /*
;0 | 浮点数转换为整数 | 整数 | r3:浮点数 | int(r3.float)
;1 | 整数转换为浮点数 | 浮点数 | r3:整数 | float(r3.int)
;2 | 申请字符串句柄 | 申请到的句柄 |  |  strPool.acquire
;3 | 字符串转换为整数 | 整数 | r3:字符串句柄,__地址__ | float(r3.str);若r3的值不是合法的字符串句柄则返回r3的值
;4 | 整数转换为字符串 | 返回的值为r3:整数 | r2:目标字符串_句柄_<br>r3:整数 | r2.str=str(r3.int);return r3.int;r2所代表字符串的内容被修改
     */
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 0)
    public void float2int(@Named("A") Operand o, @Named("R3") Register r3) {
        o.set((int) r3.getFloat());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 1)
    public void int2float(@Named("A") Operand o, @Named("R3") Register r3) {
        o.set((float) r3.get());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 2)
    public void acquireStringResource(@Named("A") Operand o, StringManager stringManager) {
        o.set(stringManager.create().getHandler());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 3)
    public void string2int(StringManager stringManager, @Named("A") Operand o, @Named("R3") Register r3) {
        o.set(Integer.parseInt(r3.getString()));
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 4)
    public void int2string(StringManager stringManager, @Named("A") Operand o, @Named("R3") Register r3, @Named("R2") Register r2) {
        o.set(r3.get());
        r2.set(String.valueOf(r3.get()));
    }

    /*
;5 | 复制字符串 | r3的值 | r2:源字符串句柄<br>r3:目标字符串句柄 | r3.str=r2.str;return r3
;6 | 连接字符串 | r3的值 | r2:源字符串<br>r3:目标字符串 | r3.str=r3.str+r2.str
;7 | 获取字符串长度 | 字符串长度 | r3:字符串 | strlen(r3.str)
;8 | 释放字符串句柄 | r3的值 | r3:字符串句柄 | strPool.release(r3);return r3
     */
    @SystemInvoke(type = SystemInvoke.Type.IN, b = 5)
    public void stringCopy(@Named("A") Operand o, @Named("R3") Register r3, @Named("R2") Register r2) {
        o.set(r3.get());
        r3.set(r2.getString());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 6)
    public void stringConcat(@Named("A") Operand o, @Named("R3") Register r3, @Named("R2") Register r2) {
        o.set(r3.get());
        r3.set(r3.getString() + r2.getString());
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 7)
    public void stringLength(@Named("A") Operand o, @Named("R3") Register r3, @Named("R2") Register r2) {
        o.set(r3.getString().length());//TODO Char length or bytes length ?
    }

    @SystemInvoke(type = SystemInvoke.Type.IN, b = 8)
    public void releaseString(StringManager stringManager, @Named("A") Operand o, @Named("R3") Register r3) {
        o.set(r3.get());
        stringManager.getResource(r3.get()).close();
    }
}
