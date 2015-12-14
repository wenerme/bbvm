package me.wener.bbvm.vm;

import java.lang.annotation.Documented;
import java.lang.annotation.Repeatable;
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
@Repeatable(SystemInvokes.class)
public @interface SystemInvoke {
    int ANY = Integer.MIN_VALUE;

    Type type();

    int a() default ANY;

    int b() default ANY;

    enum Type {
        IN, OUT
    }
}
