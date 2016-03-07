package me.wener.bbvm.dev;

import com.google.inject.ImplementedBy;

/**
 * Standard string resource implementation
 *
 * @author wener
 * @since 15/12/13
 */
@ImplementedBy(DefaultStringManager.class)
public interface StringManager extends ResourceManager<StringManager, StringResource> {

    @Override
    default String getType() {
        return "string";
    }
}
