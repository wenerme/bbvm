package me.wener.bbvm.event;

public class VMStateEvent extends BBVmEvent
{
    VMState state;

    public enum VMState
    {
        START, PAUSE, INIT, RESET, RESUME, STOP
    }
}
