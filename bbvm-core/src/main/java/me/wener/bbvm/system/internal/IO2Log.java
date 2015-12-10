package me.wener.bbvm.system.internal;

import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.system.*;
import me.wener.bbvm.system.Operand;

@Slf4j
public class IO2Log implements InstructionIntercepter
{
    @Override
    public void intercept(CPU cpu)
    {
        OpState op = cpu.opstatus();
        if (op.opcode() != Opcode.OUT || !log.isDebugEnabled())
        {
            return;
        }
        String msg;
        Operand b = op.b();
        switch (op.a().asInt())
        {
            // 0 | 显示整数 | 整数 |  | 会换行
            // 1 | 显示字符串 | 字符串 |  | 会换行
            // 2 | 显示字符串 | 字符串 |  |
            // 3 | 显示整数 | 整数 |  |
            // 4 | 显示字符 | 字符ASCII码 |  |
            // 5 | 显示浮点数 | 浮点数 |  |
            case 0:
                msg = b.asInt() + "/n";
                break;
            case 1:
                msg = b.asString() + "/n";
                break;
            case 2:
                msg = b.asString();
                break;
            case 3:
                msg = String.valueOf(b.asInt());
                break;
            case 4:
                msg = Character.highSurrogate(b.asInt()) + "";
                break;
            case 5:
                msg = b.asFloat() + "";
                break;
            default:
                return;
        }
        log.debug(msg);
    }
}
