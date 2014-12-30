package me.wener.bbvm.event;

import java.lang.annotation.Documented;
import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;
import javax.sql.rowset.Predicate;

@Documented
@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
public @interface InstFilter
{
    Class<Predicate>[] value();

    int[] instructs() default {};

    int[] op1() default {};

    int[] op2() default {};

}
