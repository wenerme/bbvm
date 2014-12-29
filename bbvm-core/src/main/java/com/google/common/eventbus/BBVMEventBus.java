package com.google.common.eventbus;

import java.util.concurrent.Executor;

public class BBVMEventBus extends AsyncEventBus
{
    public BBVMEventBus(Executor executor, SubscriberExceptionHandler subscriberExceptionHandler)
    {
        super(executor, subscriberExceptionHandler);
    }

    @Override
    void dispatch(Object event, EventSubscriber subscriber)
    {
        super.dispatch(event, subscriber);
    }
}
