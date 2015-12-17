package me.wener.bbvm.vm.res.file;

import com.google.common.base.Function;
import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.vm.SystemInvokeManager;

import javax.inject.Inject;
import javax.inject.Singleton;

/**
 * @author wener
 * @since 15/12/17
 */
@Singleton
public class FileManagerImpl implements FileManager {
    static final int HANDLER_NUMBER = 10;
    static final String TYPE = "File";
    private final FileResource[] resources = new FileResource[10];
    private final Function<Integer, FileResource> creator;

    @Inject
    public FileManagerImpl(Function<Integer, FileResource> creator) {
        this.creator = creator;
        for (int i = 0; i < HANDLER_NUMBER; i++) {
            resources[i] = create(i);
        }
    }

    @Inject
    private void init(SystemInvokeManager systemInvokeManager) {
        systemInvokeManager.register(FileInvoke.class);
    }

    private FileResource create(int i) {
        return creator.apply(i);
    }

    @Override
    public FileResource getResource(int handler) {
        if (handler < 0 || handler > HANDLER_NUMBER) {
            throw new ExecutionException("No file resource for handler " + handler);
        }
        return resources[handler];
    }

    @Override
    public FileManagerImpl reset() {
        for (FileResource resource : resources) {
            resource.close();
        }
        return this;
    }

    @Override
    public FileResource create() {
        throw new ExecutionException("Resource for file is fixed");
    }

    @Override
    public String getType() {
        return TYPE;
    }
}
