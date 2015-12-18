package me.wener.bbvm.vm;

import com.google.common.base.Function;
import com.google.common.base.MoreObjects;
import com.google.common.base.Throwables;
import com.google.common.collect.HashBasedTable;
import com.google.common.collect.ImmutableMap;
import com.google.common.collect.Maps;
import com.google.common.collect.Table;
import com.google.inject.Injector;
import com.google.inject.Key;
import com.google.inject.name.Names;
import me.wener.bbvm.exception.ExecutionException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.annotation.Nullable;
import javax.inject.Inject;
import javax.inject.Named;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.util.Map;

/**
 * @author wener
 * @since 15/12/13
 */
class SystemInvokeManagerImpl implements SystemInvokeManager {
    final static Map<String, Function<Instruction, Object>> MAPPER_MAP;
    private final static Logger log = LoggerFactory.getLogger(SystemInvokeManager.class);

    static {
        ImmutableMap.Builder<String, Function<Instruction, Object>> builder = ImmutableMap.builder();
        builder.put("A", Instruction::getA);
        builder.put("B", Instruction::getB);
        builder.put("INSTRUCTION", i -> i);
        MAPPER_MAP = builder.build();
    }

    final Map<SystemInvoke.Type, Table<Integer, Integer, Function<Instruction, Object>>> tables = Maps.newConcurrentMap();
    @Inject
    private Injector injector;

    @Inject
    public SystemInvokeManagerImpl() {
        for (SystemInvoke.Type type : SystemInvoke.Type.values()) {
            tables.put(type, HashBasedTable.create());
        }
    }

    public void register(Object handler) {
        Object o = handler;
        if (o instanceof Class) {
            o = injector.getInstance((Class) o);
        } else {
            injector.injectMembers(o);
        }
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

    public void register(SystemInvoke invoke, InvokeHandler handler) {
        if (log.isTraceEnabled()) {
            log.trace("Register {} {}, {} -> {}"
                    , invoke.type(), invoke.a() == SystemInvoke.ANY ? "ANY" : invoke.a()
                    , invoke.b() == SystemInvoke.ANY ? "ANY" : invoke.b(), handler);
        }
        register(invoke.type(), invoke.a(), invoke.b(), handler);
    }

    @Override
    public void register(Object... handlers) {
        for (Object handler : handlers) {
            register(handler);
        }
    }

    @Override
    public void register(SystemInvoke.Type type, int a, int b, Function<Instruction, Object> handler) {
        Table<Integer, Integer, Function<Instruction, Object>> handlers = tables.get(type);
        Function<Instruction, Object> function = handlers.get(a, b);
        if (function == null) {
            handlers.put(a, b, handler);
        } else {
            throw new RuntimeException("Duplicated invoke register " + handler + ",handler already exists " + a + "," + b + " -> " + function);
        }
    }

    @Override
    public Object invoke(Instruction inst) {
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
        if (handler == null) {
            throw new ExecutionException(String.format("No handler for system invoke %s %s,%s", inst.opcode, a, b));
        }
        return handler.apply(inst);
    }

    @Override
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
                    mapper[i] = MAPPER_MAP.get("INSTRUCTION");
                } else if (type == Operand.class) {
                    String n = getName(parameter);
                    if (n == null) {
                        switch (operandNth++) {
                            case 0:
                                n = "A";
                                break;
                            case 1:
                                n = "B";
                                break;
                            default:
                                throw new ExecutionException(String.format("Got three operand for %s", parameter));
                        }
                    }
                    switch (n) {
                        case "A":
                            mapper[i] = MAPPER_MAP.get("A");
                            break;
                        case "B":
                            mapper[i] = MAPPER_MAP.get("B");
                            break;
                        default:
                            throw new ExecutionException(String.format("No operand %s for %s", n, parameter));
                    }
                } else {
                    String n = getName(parameter);
                    if (n == null) {
                        args[i] = injector.getInstance(type);
                    } else {
                        args[i] = injector.getInstance(Key.get(Register.class, Names.named(n)));
                    }
                }
            }
        }

        private String getName(Parameter parameter) {
            String n = null;

            Named named = parameter.getAnnotation(Named.class);
            if (named == null) {
                com.google.inject.name.Named gNamed = parameter.getAnnotation(com.google.inject.name.Named.class);
                if (gNamed != null) {
                    n = gNamed.value();
                }
            } else {
                n = named.value();
            }
            return n;
        }


        @Nullable
        @Override
        public Object apply(Instruction input) {
            for (int i = 0; i < mapper.length; i++) {
                if (mapper[i] != null) {
                    args[i] = mapper[i].apply(input);
                }
            }

            try {
                return method.invoke(target, args);
            } catch (IllegalAccessException e) {
                throw new RuntimeException(e);
            } catch (InvocationTargetException e) {
                throw Throwables.propagate(e.getCause());
//                throw new RuntimeException(e);
            } catch (Exception e) {
                log.warn("Call {} failed with {}", method, args);
                throw Throwables.propagate(e);
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
