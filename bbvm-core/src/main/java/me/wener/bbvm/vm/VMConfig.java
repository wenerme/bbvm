package me.wener.bbvm.vm;

import com.google.common.base.Predicate;
import com.google.common.base.Throwables;
import com.google.common.collect.Lists;
import com.google.inject.Guice;
import com.google.inject.Module;
import com.typesafe.config.Config;
import me.wener.bbvm.dev.DeviceConstants;
import me.wener.bbvm.vm.invoke.InputInvoke;
import me.wener.bbvm.vm.invoke.OutputInvoke;

import java.nio.charset.Charset;
import java.util.List;
import java.util.Scanner;

/**
 * @author wener
 * @since 15/12/13
 */
public class VMConfig {
    final Charset charset;
    final Predicate<Exception> errorHandler;
    private final Config config;
    private final List<Object> invokeHandlers;
    private final List<Object> modules;
    private final int envType;

    private VMConfig(Builder builder) {
        charset = builder.charset;
        errorHandler = builder.errorHandler;
        config = builder.config;
        invokeHandlers = builder.invokeHandlers;
        modules = builder.modules;
        envType = builder.envType;
    }

    public static Builder newBuilder() {
        return new Builder();
    }

    public List<Object> getModules() {
        return modules;
    }

    public List<Object> getInvokeHandlers() {
        return invokeHandlers;
    }

    public Charset getCharset() {
        return charset;
    }

    public Config getConfig() {
        return config;
    }

    public boolean isServiceEnabled(String name) {
        String path = "service." + name + ".enabled";
        return config.hasPath(path) && config.getBoolean(path);
    }

    public boolean isModuleEnabled(String name) {
        String path = "module." + name + ".enabled";
        return config.hasPath(path) && config.getBoolean(path);
    }

    public Config getModuleConfig(String name) {
        return null;
    }

    public Config getServiceConfig(String name) {
        return null;
    }

    public Predicate<Exception> getErrorHandler() {
        return errorHandler;
    }

    public int getEnvType() {
        return envType;
    }

    public static final class Builder {
        private final List<Object> modules = Lists.newArrayList();
        private int envType = DeviceConstants.ENV_SIM;
        private List<Object> invokeHandlers = Lists.newArrayList();
        private Charset charset = Charset.forName("UTF-8");
        private Predicate<Exception> errorHandler = e -> {
            throw Throwables.propagate(e);
        };
        private Config config;

        public Builder() {
        }

        public Builder(VMConfig copy) {
            this.charset = copy.charset;
            this.errorHandler = copy.errorHandler;
            this.config = copy.config;
            this.invokeHandlers = copy.invokeHandlers;
        }

        public Builder charset(Charset val) {
            charset = val;
            return this;
        }

        /**
         * @return Handle the exception, return true for exit
         */
        public Builder errorHandler(Predicate<Exception> val) {
            errorHandler = val;
            return this;
        }

        public Builder config(Config val) {
            config = val;
            return this;
        }

        public Builder invokeHandlers(List<Object> val) {
            invokeHandlers = val;
            return this;
        }

        public Builder envType(int val) {
            envType = val;
            return this;
        }

        public Builder withModule(Class<? extends Module> module) {
            modules.add(module);
            return this;
        }

        public List<Object> modules() {
            return modules;
        }

        public Builder withModule(Module module) {
            modules.add(module);
            return this;
        }

        public Builder exitOnError() {
            errorHandler = e -> true;
            return this;
        }

        public Builder invokeWithSystemInput() {
            Scanner scanner = new Scanner(System.in);
            return invokeWith(new InputInvoke().setSupplier(scanner::nextLine));
        }

        public Builder invokeWithSystemOutput() {
            return invokeWith(new OutputInvoke(System.out::print));
        }

        public Builder invokeWith(Object handler) {
            invokeHandlers.add(handler);
            return this;
        }


        public VMConfig build() {
            return new VMConfig(this);
        }

        public VM create() {
            return Guice.createInjector(new VirtualMachineModule(build())).getInstance(VM.class);
        }
    }
}
