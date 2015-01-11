package me.wener.bbvm.neo;

import io.netty.buffer.ByteBuf;
import me.wener.bbvm.neo.codec.InstCodec;
import org.junit.Test;

public class CodecTest extends TestUtil
{
    @Test
    public void testReadInst()
    {
        /*
CALL LABEL0 ; For initialization
LD dword r2,CD_INITDATA
IN r2,22
PUSH [CSTRING_3]
POP r2
IN r3,2
IN r2,5
PUSH r3
POP r3
OUT 2,r3
IN r3,8
OUT 4,10
EXIT
LABEL0:
LD dword [CSTRING_3],CS_CSTRING_3
RET
LABEL1:
EXIT
LABEL2:
DATA CSTRING_3 dword 0
DATA CS_CSTRING_3 char "ABC",0
DATA CD_INITDATA bin %%
         */
        String dump = "00000000  42 42 45 00 00 00 00 40  00 00 00 00 00 00 00 00  |BBE....@........|\n" +
                "00000010  82 60 00 00 00 10 02 06  00 00 00 77 00 00 00 40  |.`.........w...@|\n" +
                "00000020  02 06 00 00 00 16 00 00  00 23 6c 00 00 00 30 06  |.........#l...0.|\n" +
                "00000030  00 00 00 40 02 07 00 00  00 02 00 00 00 40 02 06  |...@.........@..|\n" +
                "00000040  00 00 00 05 00 00 00 20  07 00 00 00 30 07 00 00  |....... ....0...|\n" +
                "00000050  00 50 08 02 00 00 00 07  00 00 00 40 02 07 00 00  |.P.........@....|\n" +
                "00000060  00 08 00 00 00 50 0a 04  00 00 00 0a 00 00 00 f0  |.....P..........|\n" +
                "00000070  10 0e 6c 00 00 00 70 00  00 00 90 f0 00 00 00 00  |..l...p.........|\n" +
                "00000080  41 42 43 00 00 00 00                              |ABC....|";
        ByteBuf mem = fromDumpBytes(dump);
        mem.skipBytes(16);
        BBVMContext ctx = new BBVMContext(mem);

        for (int i = 0; i < 16; i++)
        {
            System.out.println(InstCodec.read(ctx));
        }
    }
}
