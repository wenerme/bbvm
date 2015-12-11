package me.wener.bbvm.asm;

import com.google.common.collect.HashMultimap;
import com.google.common.collect.Multimap;
import me.wener.bbvm.vm.Operand;

/**
 * @author wener
 * @since 15/12/10
 */
public class BaseBBAsmParser {

    Multimap<String, Operand> labels = HashMultimap.create();

    void jjtreeOpenNodeScope(Node n) {
    }

    void jjtreeCloseNodeScope(Node n) {
    }

}
