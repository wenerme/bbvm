package me.wener.bbvm.impl.subscriber;

import com.google.common.base.Function;
import java.util.List;
import me.wener.bbvm.impl.InstructionContext;

public class InOutSubscriber
{
    List<Function<InstructionContext, Boolean>> functions;
}
