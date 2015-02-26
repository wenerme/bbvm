package me.wener.bbvm.system;

public interface CPU extends Resettable
{
    OpState opstatus();

    VmStatus vmstatus();

    /**
     * @return 是否 EXIT
     */
    boolean step();
}
