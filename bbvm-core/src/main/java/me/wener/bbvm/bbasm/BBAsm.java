package me.wener.bbvm.bbasm;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class BBAsm
{
    static
    {

    }

    public static Compiler compile(String asm)
    {
        return null;
    }

    public static void main(String[] args)
    {
        String regex = "((?<opcode>LD)\\s+(?<datatype>\\w+)\\s+(?<a>[^,]+)\\s*,\\s*(?<b>[^,]+))";
        Pattern compile = Pattern.compile(regex, Pattern.MULTILINE | Pattern.CASE_INSENSITIVE | Pattern.COMMENTS);
        Matcher matcher = compile.matcher("ld int [r1 ], 123\nnop");
        while (matcher.find())
        {
//            System.out.println(matcher.group(0));
            System.out.println(matcher.group("opcode"));
            System.out.println(matcher.group("datatype"));
            System.out.println(matcher.group("a"));
            System.out.println(matcher.group("b"));
        }
        // readStatement
        // isComment
        // skipLine
        // readOpcode
        // switch Opcode
        // readDataType
        // readOperand
        // readComma
        // readOperand
    }
}
