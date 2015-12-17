package me.wener.bbvm.vm.res;

/**
 * Standard string resource implementation
 *
 * @author wener
 * @since 15/12/13
 */
public interface StringManager extends ResourceManager<StringManager, StringResource> {

    @Override
    default String getType() {
        return "string";
    }
}
