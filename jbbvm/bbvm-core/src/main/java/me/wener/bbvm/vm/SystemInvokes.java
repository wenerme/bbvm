package me.wener.bbvm.vm;

import java.lang.annotation.Documented;
import java.lang.annotation.Retention;
import java.lang.annotation.Target;

import static java.lang.annotation.ElementType.METHOD;
import static java.lang.annotation.RetentionPolicy.RUNTIME;

/**
 * @author wener
 * @since 15/12/13
 */
@Documented
@Target(METHOD)
@Retention(RUNTIME)
public @interface SystemInvokes {
    SystemInvoke[] value();
}
