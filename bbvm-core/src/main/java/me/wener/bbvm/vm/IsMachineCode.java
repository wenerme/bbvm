package me.wener.bbvm.vm;

import me.wener.bbvm.util.val.IsInt;

/**
 * @author wener
 * @since 15/12/10
 */
public interface IsMachineCode extends IsInt {
    /**
     * @return The machine code of this target
     */
    @Override
    int asInt();
}
