package me.wener.bbvm.util;

import com.google.common.base.Predicate;
import com.google.common.base.Predicates;
import com.google.common.base.Throwables;
import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import com.google.common.collect.UnmodifiableIterator;
import com.google.common.reflect.ClassPath;
import com.google.common.util.concurrent.Service;
import com.google.common.util.concurrent.ServiceManager;
import com.google.inject.*;
import com.google.inject.multibindings.MapBinder;
import com.google.inject.multibindings.Multibinder;

import javax.inject.Named;
import javax.inject.Singleton;
import java.io.IOException;
import java.lang.reflect.Modifier;
import java.util.ArrayList;
import java.util.Map;
import java.util.Set;

/**
 * @author wener
 * @since 15/12/14
 */
@SuppressWarnings("unused")
public class MoreModules {
    /**
     * Don't support nested class
     */
    public static Module serviceModule(String pkgName, ClassLoader classLoader, Predicate<Map.Entry<String, Class<? extends Service>>> predicate) {
        return new ServiceModule(pkgName, classLoader, predicate);
    }

    public static Module serviceModule(String pkgName, ClassLoader classLoader) {
        return serviceModule(pkgName, classLoader, Predicates.<Map.Entry<String, Class<? extends Service>>>alwaysTrue());
    }


    /**
     * Don't support nested class
     */
    public static Module pluggingModule(String pkgName, ClassLoader classLoader, Predicate<Map.Entry<String, Class<? extends Module>>> predicate) {
        return new PluggingModule(pkgName, classLoader, predicate);
    }

    @SuppressWarnings("unchecked")
    private static <T> Map<String, Class<? extends T>> scanNamed(ClassLoader cl, String pkg, Class<T> base) throws IOException {
        UnmodifiableIterator<ClassPath.ClassInfo> iterator = ClassPath.from(cl).getTopLevelClassesRecursive(pkg).iterator();
        Map<String, Class<? extends T>> scanned = Maps.newHashMap();
        while (iterator.hasNext()) {
            Class<?> cls = iterator.next().load();
            if (!Modifier.isAbstract(cls.getModifiers()) && base.isAssignableFrom(cls)) {
                {
                    Named named = cls.getAnnotation(Named.class);
                    if (named != null) {
                        scanned.put(named.value(), (Class<? extends T>) cls);
                    }
                }
                {
                    com.google.inject.name.Named named = cls.getAnnotation(com.google.inject.name.Named.class);
                    if (named != null) {
                        scanned.put(named.value(), (Class<? extends T>) cls);
                    }
                }
            }
        }
        return scanned;
    }

    public static Module pluggingModule(String name, ClassLoader classLoader) {
        return pluggingModule(name, classLoader, Predicates.<Map.Entry<String, Class<? extends Module>>>alwaysTrue());
    }

    private static class ServiceModule extends AbstractModule {
        private final String pkgName;
        private final ClassLoader classLoader;
        private final Predicate<Map.Entry<String, Class<? extends Service>>> predicate;

        private ServiceModule(String pkgName, ClassLoader classLoader, Predicate<Map.Entry<String, Class<? extends Service>>> predicate) {
            this.pkgName = pkgName;
            this.classLoader = classLoader;
            this.predicate = predicate;
        }

        @Override
        protected void configure() {
            MapBinder<String, Class<? extends Service>> services = MapBinder.newMapBinder(binder(), new TypeLiteral<String>() {
            }, new TypeLiteral<Class<? extends Service>>() {
            });

            Multibinder.newSetBinder(binder(), Service.class);

            if (pkgName != null) {
                try {
                    for (Map.Entry<String, Class<? extends Service>> entry : scanNamed(classLoader, pkgName, Service.class).entrySet()) {
                        services.addBinding(entry.getKey()).toInstance(entry.getValue());
                    }
                } catch (IOException e) {
                    Throwables.propagate(e);
                }
            }
        }

        @Provides
        @Singleton
        public ServiceManager serviceManager(Injector injector, Set<Service> services, Map<String, Class<? extends Service>> serviceClasses) {
            ArrayList<Service> list = Lists.newArrayList(services);
            for (Map.Entry<String, Class<? extends Service>> entry : Maps.filterEntries(serviceClasses, predicate).entrySet()) {
                list.add(injector.getInstance(entry.getValue()));
            }
            return new ServiceManager(list);
        }
    }

    private static class PluggingModule extends AbstractModule {
        private final String pkgName;
        private final ClassLoader classLoader;
        private final Predicate<Map.Entry<String, Class<? extends Module>>> predicate;

        private PluggingModule(String pkgName, ClassLoader classLoader, Predicate<Map.Entry<String, Class<? extends Module>>> predicate) {
            this.pkgName = pkgName;
            this.classLoader = classLoader;
            this.predicate = predicate;
        }

        @Override
        protected void configure() {
            try {
                for (Map.Entry<String, Class<? extends Module>> entry : scanNamed(classLoader, pkgName, Module.class).entrySet()) {
                    install(entry.getValue().newInstance());
                }
            } catch (Exception e) {
                Throwables.propagate(e);
            }
        }

    }
}
