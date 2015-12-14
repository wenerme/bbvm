package me.wener.bbvm.vm;

import com.google.common.base.Predicate;
import com.typesafe.config.Config;
import me.wener.bbvm.exception.ExecutionException;

import java.nio.charset.Charset;

/**
 * @author wener
 * @since 15/12/13
 */
public class VMConfig {
    final Charset charset;
    final Predicate<ExecutionException> errorHandler;
    private final com.typesafe.config.Config config;

    private VMConfig(Builder builder) {
        charset = builder.charset;
        errorHandler = builder.errorHandler;
        config = builder.config;
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

    public com.typesafe.config.Config getModuleConfig(String name) {
        return null;
    }

    public com.typesafe.config.Config getServiceConfig(String name) {
        return null;
    }

    public Predicate<ExecutionException> getErrorHandler() {
        return errorHandler;
    }

    public static final class Builder {
        private Charset charset = Charset.forName("UTF-8");
        private Predicate<ExecutionException> errorHandler = e -> true;
        private com.typesafe.config.Config config;

        public Builder() {
        }

        public Builder(VMConfig copy) {
            this.charset = copy.charset;
            this.errorHandler = copy.errorHandler;
            this.config = copy.config;
        }

        public Builder charset(Charset val) {
            charset = val;
            return this;
        }

        /**
         * @return Handle the exception, return true for exit
         */
        public Builder errorHandler(Predicate<ExecutionException> val) {
            errorHandler = val;
            return this;
        }

        public Builder config(com.typesafe.config.Config val) {
            config = val;
            return this;
        }

        public Builder exitOnError() {
            errorHandler = e -> true;
            return this;
        }


        public VMConfig build() {
            return new VMConfig(this);
        }
    }
}
