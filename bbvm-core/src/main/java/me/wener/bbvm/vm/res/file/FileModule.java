package me.wener.bbvm.vm.res.file;

import com.google.common.base.Function;
import com.google.inject.AbstractModule;
import com.google.inject.TypeLiteral;
import com.google.inject.multibindings.OptionalBinder;

import javax.inject.Inject;

/**
 * @author wener
 * @since 15/12/17
 */
public class FileModule extends AbstractModule {
    @Override
    protected void configure() {
        OptionalBinder.newOptionalBinder(binder(), new TypeLiteral<Function<Integer, FileResource>>() {
        }).setDefault().to(JavaFileResourceFunction.class);
        bind(FileManager.class).to(FileManagerImpl.class).asEagerSingleton();
    }

    private static class JavaFileResourceFunction implements Function<Integer, FileResource> {
        @Inject
        FileManager fileManager;

        @Override
        public FileResource apply(Integer input) {
            return new JavaFileResourceImpl(input, fileManager);
        }
    }
}
