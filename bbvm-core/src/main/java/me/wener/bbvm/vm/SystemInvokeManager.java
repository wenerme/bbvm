package me.wener.bbvm.vm;

import com.google.common.base.Function;
import com.google.common.base.MoreObjects;
import com.google.common.collect.HashBasedTable;
import com.google.common.collect.Maps;
import com.google.common.collect.Table;
import com.google.inject.Injector;

import javax.annotation.Nullable;
import javax.inject.Inject;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.util.Map;

/**
 * @author wener
 * @since 15/12/13
 */
public class SystemInvokeManager {
    final Map<SystemInvoke.Type, Table<Integer, Integer, Function<Instruction, Object>>> tables = Maps.newConcurrentMap();
    @Inject
    private Injector injector;

    @Inject
    public SystemInvokeManager() {
        for (SystemInvoke.Type type : SystemInvoke.Type.values()) {
            tables.put(type, HashBasedTable.create());
        }
    }

    public void register(Object o) {
        for (Method m : o.getClass().getMethods()) {
            SystemInvokes invokes = m.getAnnotation(SystemInvokes.class);
            if (invokes != null) {
                for (SystemInvoke invoke : invokes.value()) {
                    register(invoke, wrap(o, m));
                }
                continue;
            }
            SystemInvoke invoke = m.getAnnotation(SystemInvoke.class);
            if (invoke != null) {
                register(invoke, wrap(o, m));
            }
        }
    }

    private void register(SystemInvoke invoke, InvokeHandler handler) {
        register(invoke.type(), invoke.a(), invoke.b(), handler);
    }


    public void register(SystemInvoke.Type type, int a, int b, Function<Instruction, Object> handler) {
        Table<Integer, Integer, Function<Instruction, Object>> handlers = tables.get(type);
        Function<Instruction, Object> function = handlers.get(a, b);
        if (function == null) {
            handlers.put(a, b, handler);
        } else {
            throw new RuntimeException("Can not register " + handler + ",handler already exists " + a + "," + b + " -> " + function);
        }
    }

    public void invoke(Instruction inst) {
        int a = inst.getA().get();
        int b = inst.getB().get();
        Function<Instruction, Object> handler;
        switch (inst.opcode) {
            case IN:
                handler = getHandler(SystemInvoke.Type.IN, a, b);
                break;
            case OUT:
                handler = getHandler(SystemInvoke.Type.OUT, a, b);
                break;
            default:
                throw new RuntimeException(inst.opcode + " is not a system invoke");
        }
        handler.apply(inst);
    }

    public Function<Instruction, Object> getHandler(SystemInvoke.Type type, int a, int b) {
        Table<Integer, Integer, Function<Instruction, Object>> handlers = tables.get(type);
        Function<Instruction, Object> handler = handlers.get(a, b);
        if (handler == null) {
            handler = handlers.get(a, SystemInvoke.ANY);
        }
        if (handler == null) {
            handler = handlers.get(SystemInvoke.ANY, b);
        }
        if (handler == null) {
            handler = handlers.get(SystemInvoke.ANY, SystemInvoke.ANY);
        }
        return handler;
    }

    private InvokeHandler wrap(Object target, Method m) {
        return new InvokeHandler(target, m);
    }

    class InvokeHandler implements Function<Instruction, Object> {
        Object target;
        Method method;
        Object[] args;
        Function<Instruction, Object>[] mapper;

        @SuppressWarnings("unchecked")
        public InvokeHandler(Object target, Method m) {
            this.target = target;
            this.method = m;
            Parameter[] parameters = method.getParameters();
            args = new Object[parameters.length];
            mapper = new Function[parameters.length];
            int operandNth = 0;

            for (int i = 0; i < parameters.length; i++) {
                Parameter parameter = parameters[i];
                Class<?> type = parameter.getType();
                if (type == Instruction.class) {
                    mapper[i] = new Function<Instruction, Object>() {
                        @Nullable
                        @Override
                        public Object apply(@Nullable Instruction input) {
                            return input;
                        }
                    };
                } else if (type == Operand.class) {
                    switch (operandNth++) {
                        case 0:
                            mapper[i] = new Function<Instruction, Object>() {
                                @Nullable
                                @Override
                                public Object apply(Instruction input) {
                                    return input.a;
                                }
                            };
                            break;
                        case 1:
                            mapper[i] = new Function<Instruction, Object>() {
                                @Nullable
                                @Override
                                public Object apply(Instruction input) {
                                    return input.b;
                                }
                            };
                            break;
                        default:
                            throw new RuntimeException("Got third operand!");
                    }
                } else {
                    mapper[i] = new Function<Instruction, Object>() {
                        @Nullable
                        @Override
                        public Object apply(@Nullable Instruction input) {
                            return injector.getInstance(type);
                        }
                    };
                }
            }
        }


        @Nullable
        @Override
        public Object apply(Instruction input) {
            for (int i = 0; i < mapper.length; i++) {
                args[i] = mapper[i].apply(input);
            }

            try {
                return method.invoke(target, args);
            } catch (IllegalAccessException e) {
                throw new RuntimeException(e);
            } catch (InvocationTargetException e) {
                throw new RuntimeException(e);
            }
        }

        @Override
        public String toString() {
            return MoreObjects.toStringHelper(this)
                    .add("target", target)
                    .add("method", method)
                    .toString();
        }
    }
}
