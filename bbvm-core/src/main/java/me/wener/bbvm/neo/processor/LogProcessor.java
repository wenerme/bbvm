package me.wener.bbvm.neo.processor;

import lombok.extern.slf4j.Slf4j;

@Slf4j
public class LogProcessor implements Processor
{
    @Override
    public ProcessResult apply(ProcessContext input)
    {
        log.info(input.instruction().toString());
        return Results.keepGoing();
    }
}
