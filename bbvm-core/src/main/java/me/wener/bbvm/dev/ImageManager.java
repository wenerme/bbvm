package me.wener.bbvm.dev;

import java.util.List;

/**
 * @author wener
 * @since 15/12/18
 */
public interface ImageManager extends ResourceManager<ImageManager, ImageResource> {
    /**
     * @param file  Resource name
     * @param index Resource index start from 0
     */
    ImageResource load(String file, int index);

    /**
     * @return Mutable directories to search the resource
     */
    List<String> getResourceDirectory();

    @Override
    default String getType() {
        return "image";
    }
}
