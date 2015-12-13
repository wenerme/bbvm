package me.wener.bbvm.vm;

import com.google.common.base.Predicate;
import me.wener.bbvm.exception.ExecutionException;

import java.nio.charset.Charset;

/**
 * @author wener
 * @since 15/12/13
 */
public class Config {
    final Charset charset;
    final Predicate<ExecutionException> errorHandler;

    private Config(Builder builder) {
        charset = builder.charset;
        errorHandler = builder.errorHandler;
    }

    public Charset getCharset() {
        return charset;
    }

    public Predicate<ExecutionException> getErrorHandler() {
        return errorHandler;
    }

    public static final class Builder {
        private Charset charset = Charset.forName("UTF-8");
        private Predicate<ExecutionException> errorHandler = e -> true;

        public Builder() {
        }

        public Builder(Config copy) {
            this.charset = copy.charset;
            this.errorHandler = copy.errorHandler;
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

        public Builder exitOnError() {
            errorHandler = e -> true;
            return this;
        }


        public Config build() {
            return new Config(this);
        }
    }
}
