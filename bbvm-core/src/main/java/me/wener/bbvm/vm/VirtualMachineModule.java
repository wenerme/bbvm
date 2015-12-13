package me.wener.bbvm.vm;

import com.google.inject.AbstractModule;
import com.google.inject.Provides;
import com.google.inject.multibindings.OptionalBinder;
import me.wener.bbvm.vm.res.StringManager;

import javax.inject.Singleton;

/**
 * @author wener
 * @since 15/12/13
 */
public class VirtualMachineModule extends AbstractModule {

    private final Config config;

    public VirtualMachineModule(Config config) {
        this.config = config;
    }

    @Override
    protected void configure() {
        OptionalBinder.newOptionalBinder(binder(), SystemInvokeManager.class)
                .setDefault().to(SystemInvokeManagerImpl.class).in(Singleton.class);
//        OptionalBinder.newOptionalBinder(binder(), StringManager.class)
//                .setDefault().to(StringManager.class).in(Singleton.class);
        bind(StringManager.class).in(Singleton.class);
        bind(VM.class).in(Singleton.class);
    }

    @Provides
    @Singleton
    public Config config() {
        return config;
    }
}
