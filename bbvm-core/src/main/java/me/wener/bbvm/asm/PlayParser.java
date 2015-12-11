package me.wener.bbvm.asm;

import me.wener.bbvm.vm.Instruction;

import java.io.StringReader;

/**
 * @author wener
 * @since 15/12/11
 */
public class PlayParser {
    public static void main(String[] args) throws ParseException {

        BBAsmParser parser = new BBAsmParser(new StringReader("" +
                "ld int r1, 1\n\n\n" +
                "ld word   rp, r1\n" +
                "ld int [123],123\n" +
                "ld int [r1],r0\n" +
                "HERE:\n" +
                "ld int [HERE],123\n" +
                "THERE :\n" +
                "jmp THERE" +
                "\n"));
        for (Instruction i : parser.start()) {
            System.out.println(i);
            System.out.println(i.toAssembly());
        }
    }
}
