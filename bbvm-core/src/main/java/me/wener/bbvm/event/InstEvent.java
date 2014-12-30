package me.wener.bbvm.event;

import lombok.Data;
import lombok.EqualsAndHashCode;
import me.wener.bbvm.impl.InstructionContext;

@EqualsAndHashCode(callSuper = true)
@Data
public class InstEvent extends BBVMEvent
{
    InstructionContext context;
}
