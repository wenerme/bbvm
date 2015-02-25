package me.wener.bbvm.system.internal;

import lombok.extern.slf4j.Slf4j;
import me.wener.bbvm.system.CPU;
import me.wener.bbvm.system.InstructionIntercepter;
import me.wener.bbvm.system.OpStatus;
import me.wener.bbvm.system.Opcode;
import me.wener.bbvm.system.Operand;

@Slf4j
public class IO2Log implements InstructionIntercepter
{
    @Override
    public void intercept(CPU cpu)
    {
        OpStatus op = cpu.opstatus();
        if (op.opcode() != Opcode.OUT || !log.isDebugEnabled())
        {
            return;
        }
        String msg;
        Operand b = op.b();
        switch (op.a().get())
        {
            // 0 | 显示整数 | 整数 |  | 会换行
            // 1 | 显示字符串 | 字符串 |  | 会换行
            // 2 | 显示字符串 | 字符串 |  |
            // 3 | 显示整数 | 整数 |  |
            // 4 | 显示字符 | 字符ASCII码 |  |
            // 5 | 显示浮点数 | 浮点数 |  |
            case 0:
                msg = b.get().toString() + "/n";
                break;
            case 1:
                msg = b.asString() + "/n";
                break;
            case 2:
                msg = b.asString();
                break;
            case 3:
                msg = b.get().toString();
                break;
            case 4:
                msg = Character.highSurrogate(b.get()) + "";
                break;
            case 6:
                msg = b.asFloat() + "";
                break;
            default:
                return;
        }
        log.debug(msg);
    }
}
