package com.google.common.eventbus;

import java.lang.reflect.Method;
import java.util.concurrent.Executor;
import me.wener.bbvm.event.InstEvent;
import me.wener.bbvm.event.InstFilter;

public class BBVMEventBus extends AsyncEventBus
{
    public BBVMEventBus(Executor executor, SubscriberExceptionHandler subscriberExceptionHandler)
    {
        super(executor, subscriberExceptionHandler);
    }

    @Override
    void dispatch(Object event, EventSubscriber subscriber)
    {
        if (event instanceof InstEvent)
        {
            Method method = subscriber.getMethod();
            InstFilter instFilter = method.getAnnotation(InstFilter.class);
            if (instFilter != null)
            {
                // 实际逻辑处理
            }
        }
        super.dispatch(event, subscriber);
    }
}
