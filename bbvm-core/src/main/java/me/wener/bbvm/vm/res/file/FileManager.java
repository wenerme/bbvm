package me.wener.bbvm.vm.res.file;

import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.vm.res.ResourceManager;

/**
 * @author wener
 * @since 15/12/17
 */
public class FileManager implements ResourceManager<FileManager, FileResource> {
    public static final int FILE_NUMBER = 10;
    private final FileResource[] resources = new FileResource[10];

    public FileManager() {
        for (int i = 0; i < FILE_NUMBER; i++) {
            resources[i] = create(i);
        }
    }

    private FileResource create(int i) {
        return null;
    }

    @Override
    public FileResource getResource(int handler) {
        if (handler < 0 || handler > FILE_NUMBER) {
            throw new ExecutionException("No file resource for handler " + handler);
        }
        return resources[handler];
    }

    @Override
    public FileManager reset() {
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
        return "File";
    }
}
