package me.wener.bbvm.dev;

import me.wener.bbvm.dev.Images.ImageInfo;
import org.junit.Test;

import javax.imageio.ImageIO;
import java.io.File;
import java.io.IOException;

/**
 * @author wener
 * @since 15/12/20
 */
public class ImagesTest {

    @Test(expected = Exception.class)
    public void testFileNotFound() throws IOException {
        Images.load("../bbvm-test/image/9288-xxxx.lib").get(0);
    }

    @Test
    public void testLib2BitLeGray() throws IOException {
        ImageInfo info = Images.load("../bbvm-test/image/9288.lib").get(0);
        ImageIO.write(Images.read(info), "png", new File(info.getType() + ".png"));
    }

    @Test
    public void testLib2BitBeGray() throws IOException {
        ImageInfo info = Images.load("../bbvm-test/image/9188.lib").get(0);
        ImageIO.write(Images.read(info), "png", new File(info.getType() + ".png"));
    }

    @Test
    public void testLibRGB565() throws IOException {
        ImageInfo info = Images.load("../bbvm-test/image/9688.lib").get(0);
        ImageIO.write(Images.read(info), "BMP", new File(info.getType() + ".bmp"));
    }

    @Test
    public void testRlb() throws IOException {
        ImageInfo info = Images.load("../bbvm-test/image/bmp.rlb").get(0);
        ImageIO.write(Images.read(info), "BMP", new File(info.getType() + ".bmp"));
    }

    @Test
    public void testGeneric() throws IOException {
        ImageInfo info = Images.load("../bbvm-test/image/bmp.bmp").get(0);
        ImageIO.write(Images.read(info), "BMP", new File(info.getType() + ".bmp"));
    }
}
