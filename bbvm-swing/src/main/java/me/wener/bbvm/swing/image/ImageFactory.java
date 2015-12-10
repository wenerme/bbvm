package me.wener.bbvm.swing.image;

import com.google.common.io.LittleEndianDataInputStream;
import me.wener.bbvm.util.Bins;

import javax.imageio.ImageIO;
import javax.swing.*;
import java.awt.*;
import java.awt.image.BufferedImage;
import java.awt.image.WritableRaster;
import java.io.*;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.EnumSet;

public class ImageFactory
{
    private static BufferedImage loadLibRGB565(DataInput is) throws IOException
    {

        int len = is.readInt();
        int w = is.readUnsignedShort();
        int h = is.readUnsignedShort();
        BufferedImage image = new BufferedImage(w, h, BufferedImage.TYPE_USHORT_565_RGB);
        WritableRaster raster = image.getRaster();

        is.skipBytes(8);// 跳过无用的8字节
        //image.getColorModel().getComponents()
        int[] color = new int[3];
        for (int y = 0; y < h; y++)
        {
            for (int x = 0; x < w; x++)
            {
                raster.setPixel(x, y, RGB565ToRGB(is.readUnsignedShort(), color));
            }
        }

        return image;
    }


    private static BufferedImage loadLibBit2Gray(DataInput is, final ImageType type) throws IOException
    {
        int len = is.readInt();
        int w = is.readUnsignedShort();
        int h = is.readUnsignedShort();
        BufferedImage image = new BufferedImage(w, h, BufferedImage.TYPE_BYTE_GRAY);
        WritableRaster raster = image.getRaster();

        is.skipBytes(8);// 跳过无用的8字节
        //image.getColorModel().getComponents()
        int[] color = new int[3];
        for (int y = 0; y < h; y++)
        {
            for (int x = 0; x < w; )
            {
                int c = is.readUnsignedByte();
                for (int i = 0; i < 4 && x < w; i++)
                {
                    raster.setPixel(x, y, Bit2GrayToRGB(c, color, i, type == ImageType.LIB_Bit2Gray_BE));
                    x++;
                }
            }
        }

        return image;
    }


    public static BufferedImage loadImage(String file, int index) throws IOException
    {
        try
        {
            byte[] bytes = Files.readAllBytes(Paths.get(file));
            ImageType type = detectType(file, bytes);
            if (type == null)
                return ImageIO.read(new ByteArrayInputStream(bytes));

            return loadLibrary(new ByteArrayInputStream(bytes), type)[index];
        } catch (Exception e)
        {
            e.printStackTrace();
            return null;
        }
    }

    public static BufferedImage[] loadLibrary(String file, ImageType type) throws IOException
    {
        return loadLibrary(new FileInputStream(file), type);
    }

    public static BufferedImage[] loadLibrary(InputStream is, ImageType type) throws IOException
    {
        DataInput in;
        if (type.isLittleEndian())
            in = new LittleEndianDataInputStream(is);
        else in = new DataInputStream(is);


        int gap = 0;
        if (type == ImageType.RLB)
            gap = 32;
        int[] offsets = readOffset(in, gap);
        int number = offsets.length;
        BufferedImage[] images = new BufferedImage[number];

        // 读取所有Image
        for (int i = 0; i < number; i++)
        {
            switch (type)
            {
                case LIB_RGB565:
                    images[i] = loadLibRGB565(in);
                    break;
                case LIB_Bit2Gray_LE:
                case LIB_Bit2Gray_BE:
                    images[i] = loadLibBit2Gray(in, type);
                    break;
                case RLB:
                    images[i] = loadRlb(is);
                    break;
                default:
                    throw new RuntimeException("不支持的格式.");
            }
        }

        return images;
    }

    private static int[] readOffset(DataInput in, int gap) throws IOException
    {
        int n = in.readInt();
        int[] offsets = new int[n];

        // 读取所有偏移量
        for (int i = 0; i < n; i++)
        {
            offsets[i] = in.readInt();
            in.skipBytes(gap);
        }
        return offsets;
    }

    private static BufferedImage loadRlb(InputStream in) throws IOException
    {
        //noinspection ResultOfMethodCallIgnored
        in.skip(4);// 前面 4 位为长度
        return ImageIO.read(in);
    }

    private static Image[] loadRlb(String file) throws IOException
    {
        byte[] bytes = Files.readAllBytes(Paths.get(file));
        int offset = 0;
        int n = Bins.int32l(bytes, offset);
        offset += 4;

        Image[] images = new Image[n];
        int[] offsets = new int[n];

        // 读取所有偏移量
        for (int i = 0; i < n; i++)
        {
            offsets[i] = Bins.int32l(bytes, offset);
            offset += 4 + 32;// 跳过名字
        }

        // 读取所有Image
        for (int i = 0; i < n; i++)
        {
            int len = Bins.int32l(bytes, offset);
            offset += 4;
            images[i] = ImageIO.read(new ByteArrayInputStream(bytes, offset, bytes.length - offset));
            offset += len;
        }

        return images;
    }

    public static void main(String[] args) throws IOException
    {
        String file = "D:\\dev\\projects\\bbvm\\doc\\testsuit\\out\\wener.bmp";
        String rlb = "D:\\dev\\projects\\bbvm\\doc\\testsuit\\out\\wener.rlb";
        String lib9688 = "D:\\dev\\projects\\bbvm\\doc\\testsuit\\out\\9688-wener.lib";
        String lib9288 = "D:\\dev\\projects\\bbvm\\doc\\testsuit\\out\\9288-wener.lib";
        String lib9188 = "D:\\dev\\projects\\bbvm\\doc\\testsuit\\out\\9188-wener.lib";

        //Image image = loadRlb(rlb)[0];
//        BufferedImage image = loadLibrary(lib9188, ImageType.LIB_Bit2Gray_BE)[0];
        BufferedImage image = loadImage(rlb, 0);
        ImageIO.write(image, "BMP", new File("D:\\dev\\projects\\bbvm\\doc\\testsuit\\out\\tmp.bmp"));

        assert detectType(lib9188, Files.readAllBytes(Paths.get(lib9188))) == ImageType.LIB_Bit2Gray_BE;
        assert detectType(lib9288, Files.readAllBytes(Paths.get(lib9288))) == ImageType.LIB_Bit2Gray_LE;
        assert detectType(lib9688, Files.readAllBytes(Paths.get(lib9688))) == ImageType.LIB_RGB565;
        assert detectType(rlb, Files.readAllBytes(Paths.get(rlb))) == ImageType.RLB;

        JFrame frame = new BackgroundImageJFrame(image);
        frame.setVisible(true);
        frame.repaint();


    }

    static int[] Bit2GrayToRGB(int color, int[] bytes, int i)
    {
        return Bit2GrayToRGB(color, bytes, i, false);
    }

    static int[] Bit2GrayToRGB(int color, int[] bytes, int i, boolean inverse)
    {
        i = 3 - i;
        color = color >> (i * 2) & 0b11;
        color = color << 6;
        if (inverse)
            color = 255 - color;
        bytes[0] = color;
        bytes[1] = color;
        bytes[2] = color;

        return bytes;
    }

    static int[] BGR565ToRGB(int color, int[] bytes)
    {
        // bbbbb|ggggg|rrrrr
        bytes[0] = color & 0b11111;
        bytes[1] = color >> 5 & 0b111111;
        bytes[2] = color >> 11 & 0b11111;

        return bytes;
    }

    static int[] RGB565ToRGB(int color, int[] bytes)
    {
        // bbbbb|ggggg|rrrrr
        bytes[2] = color & 0b11111;
        bytes[1] = color >> 5 & 0b111111;
        bytes[0] = color >> 11 & 0b11111;

        return bytes;
    }

    static int[] BGR555ToRGB(int color, int[] bytes)
    {
        // bbbbb|ggggg|rrrrr
        bytes[0] = color & 0b11111;
        bytes[1] = color >> 5 & 0b11111;
        bytes[2] = color >> 10 & 0b11111;

        return bytes;
    }

    /**
     * @return 无法检测类型返回 {@code null}
     */
    public static ImageType detectType(String fn, byte[] bytes)
    {
        int indexOf = fn.lastIndexOf('.');

        String ext = null;
        if (indexOf >= 0)
            ext = fn.substring(indexOf + 1).toUpperCase();

        switch (ext)
        {
            case "DLX":
                return ImageType.DLX;
            case "RLB":
                return ImageType.RLB;
            case "LIB":
            default:
                return detectType(bytes);
        }
    }

    private static ImageType detectType(byte[] bytes)
    {

        if (bytes[0] == 'B' && bytes[1] == 'M')
        {
            return null;
        }
        if (bytes[0] == 'D' && bytes[1] == 'L' && bytes[1] == 'X')
        {
            return ImageType.DLX;
        }
        int n = Bins.int32l(bytes, 0);
        boolean be = false;
        int maxImages = 0xffff;
        if (n > maxImages)
        {
            n = Bins.int32b(bytes, 0);
            be = true;
        }

        int offset = Bins.int32(bytes, 4, be);
        int len = Bins.int32(bytes, offset, be);
        int w, h;

        w = Bins.int16(bytes, offset + 4, be);
        h = Bins.int16(bytes, offset + 6, be);

        if (w * h / 4 == len - 12)
        {
            if (be)
                return ImageType.LIB_Bit2Gray_BE;
            else
                return ImageType.LIB_Bit2Gray_LE;
        }

        if (w * h * 2 == len - 12)
        {
            return ImageType.LIB_RGB565;
        }

        return null;
    }

    public enum ImageType
    {
        DLX,
        RLB,
        LIB_RGB565,
        LIB_Bit2Gray_LE,
        LIB_Bit2Gray_BE;
        private static EnumSet<ImageType> be = EnumSet.of(LIB_Bit2Gray_BE);

        public boolean isLittleEndian()
        {
            return !be.contains(this);
        }
    }

    public static class BackgroundImageJFrame extends JFrame
    {
        private final Image image;
        JButton b1;
        JLabel l1;

        public BackgroundImageJFrame(Image image)
        {
            setTitle("Background Color for JFrame");
            setSize(400, 400);
            setLocationRelativeTo(null);
            setDefaultCloseOperation(EXIT_ON_CLOSE);
            setVisible(true);

            // Another way
            setLayout(new BorderLayout());
            this.image = image;
            setContentPane(new JLabel(new ImageIcon(this.image)));
            setLayout(new FlowLayout());
            l1 = new JLabel("Here is a button");
            b1 = new JButton("I am a button");
//            add(l1);
//            add(b1);
            // Just for refresh :) Not optional!
//            setSize(399, 399);
//            setSize(400, 400);
        }

    }
}
