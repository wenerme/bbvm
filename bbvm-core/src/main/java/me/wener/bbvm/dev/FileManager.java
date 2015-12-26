package me.wener.bbvm.dev;

/**
 * @author wener
 * @since 15/12/17
 */
public interface FileManager extends ResourceManager<FileManager, FileResource> {
    @Override
    default String getType() {
        return "file";
    }
}
