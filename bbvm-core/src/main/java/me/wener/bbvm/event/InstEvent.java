package me.wener.bbvm.event;

import lombok.Data;
import lombok.EqualsAndHashCode;
import me.wener.bbvm.def.Instruction;

@EqualsAndHashCode(callSuper = true)
@Data
public class InstEvent extends BBVmEvent
{
    Instruction instruction;
}
