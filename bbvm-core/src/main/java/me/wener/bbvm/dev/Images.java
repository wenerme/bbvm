package me.wener.bbvm.dev;

import com.google.common.base.MoreObjects;
import com.google.common.collect.ImmutableList;
import com.google.common.collect.Maps;
import com.google.common.io.Files;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;

import javax.imageio.ImageIO;
import java.awt.image.BufferedImage;
import java.awt.image.DataBufferByte;
import java.awt.image.DataBufferUShort;
import java.awt.image.IndexColorModel;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.channels.FileChannel;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 * @author wener
 * @since 15/12/19
 */
public class Images {
    private static final Map<String, ImageCodec> CODEC = Maps.newHashMap();

    static {
        register(new LibRGB565ImageCodec());
        register(new Lib2BitGrayBEImageCodec());
        register(new Lib2BitGrayLEImageCodec());
        register(new RLBImageCodec());
    }

    private static void register(ImageCodec codec) {
        CODEC.put(codec.getType(), codec);
    }

    public static List<ImageInfo> load(String file) throws IOException {
        Path path = Paths.get(file);
        for (Map.Entry<String, ImageCodec> entry : CODEC.entrySet()) {
            ImageCodec codec = entry.getValue();
            if (codec.accept(path)) {
                try {
                    List<ImageInfo> infos = codec.load(path);
                    if (infos != null) {
                        return infos;
                    }
                } catch (IOException e) {
                    e.printStackTrace();
                } catch (UnsupportedOperationException e) {
//                    e.printStackTrace();
                }
            }
        }
        throw new RuntimeException("Can not load image file " + file);
    }

    public static BufferedImage read(ImageInfo info) throws IOException {
        return CODEC.get(info.getType()).read(info);
    }

    public static BufferedImage read(String file, int index) throws IOException {
        return read(load(file).get(index));
    }


    private static List<ImageInfo> loadInfos(Path file, String type, boolean hasName, ByteOrder order, int expectedBits) throws IOException {
        try (FileChannel ch = FileChannel.open(file, StandardOpenOption.READ)) {
            int n = ch.map(FileChannel.MapMode.READ_ONLY, 0, 4).order(order).getInt();
            // Never reach this size
            if (n > 65535) {
                throw new UnsupportedOperationException();
            }
            ByteBuf buf = Unpooled.wrappedBuffer(ch.map(FileChannel.MapMode.READ_ONLY, 4, n * (4 + 32))).order(order);
            List<ImageInfo> list = new ArrayList<>(n);
            for (int i = 0; i < n; i++) {
                int offset = buf.readInt();
                String name = null;
                if (hasName) {
                    int length = buf.bytesBefore(32, (byte) 0);
                    if (length > 0) {
                        name = buf.toString(buf.readerIndex(), length, StandardCharsets.UTF_8);
                    }
                    buf.skipBytes(32);
                }
                if (name == null) {
                    name = "NO-" + i;
                }
                ImageInfo info = new ImageInfo(i, name, offset, file.toString(), type);

                buf = Unpooled.wrappedBuffer(ch.map(FileChannel.MapMode.READ_ONLY, offset, 12)).order(order);
                info.setSize(buf.readInt() - 12)
                        .setWidth(buf.readUnsignedShort())
                        .setHeight(buf.readUnsignedShort());

                if (expectedBits > 0 && info.getBits() != expectedBits) {
                    throw new UnsupportedOperationException("Expected bits " + expectedBits + " but got " + info.getBits());
                }

                list.add(info);

            }
            return list;
        }
    }


    interface ImageCodec {
        List<ImageInfo> load(Path file) throws IOException;

        default BufferedImage read(Path path, int index) throws IOException {
            return read(load(path).get(index));
        }

        String getType();

        boolean accept(Path path);

        BufferedImage read(ImageInfo info) throws IOException;
    }

    static class LibRGB565ImageCodec implements ImageCodec {

        @Override
        public List<ImageInfo> load(Path file) throws IOException {
            return loadInfos(file, getType(), false, ByteOrder.LITTLE_ENDIAN, 16);
        }

        @Override
        public String getType() {
            return "LIB_RGB565";
        }

        @Override
        public boolean accept(Path path) {
            File file = path.toFile();
            return !file.isDirectory() && Files.getFileExtension(file.getName()).equalsIgnoreCase("LIB");
        }

        @Override
        public BufferedImage read(ImageInfo info) throws IOException {
            try (FileChannel ch = FileChannel.open(Paths.get(info.getFilename()), StandardOpenOption.READ)) {
                // + 4 to skip length
                ByteBuf buf = Unpooled.wrappedBuffer(ch.map(FileChannel.MapMode.READ_ONLY, info.getOffset() + 4, 8)).order(ByteOrder.LITTLE_ENDIAN);
                int w = buf.readUnsignedShort(), h = buf.readUnsignedShort();
                ByteBuffer buffer = ByteBuffer.allocate(w * h * 2).order(ByteOrder.LITTLE_ENDIAN);
                ch.read(buffer);
                BufferedImage image = new BufferedImage(w, h, BufferedImage.TYPE_USHORT_565_RGB);
                // Fast way to load the data
                buffer.flip();
                buffer.asShortBuffer().get(((DataBufferUShort) image.getRaster().getDataBuffer()).getData());
                return image;
            }
        }
    }


    static class Lib2BitGrayLEImageCodec extends Lib2BitGrayImageCodec {

        Lib2BitGrayLEImageCodec() {
            super(ByteOrder.LITTLE_ENDIAN);
        }

        @Override
        public String getType() {
            return super.getType() + "_LE";
        }
    }

    static class Lib2BitGrayBEImageCodec extends Lib2BitGrayImageCodec {

        Lib2BitGrayBEImageCodec() {
            super(ByteOrder.BIG_ENDIAN);
        }

        @Override
        public String getType() {
            return super.getType() + "_BE";
        }
    }

    static class GenericImageCodec implements ImageCodec {

        @Override
        public List<ImageInfo> load(Path file) throws IOException {
            BufferedImage image = ImageIO.read(file.toFile());
            if (image == null) {
                throw new UnsupportedOperationException(file + " is not generic image");
            }
            GenericImage im = new GenericImage(0, "NO-0", 0, file.toString(), getType(), image);
            im.setHeight(im.getHeight()).setWidth(im.getWidth());
            return ImmutableList.of(im);
        }

        @Override
        public String getType() {
            return "GENERIC";
        }

        @Override
        public boolean accept(Path path) {
            return true;
        }

        @Override
        public BufferedImage read(ImageInfo info) throws IOException {
            if (info instanceof GenericImage) {
                return ((GenericImage) info).getImage();
            }
            return ImageIO.read(new File(info.getFilename()));
        }
    }

    static class Lib2BitGrayImageCodec implements ImageCodec {

        private final ByteOrder order;

        Lib2BitGrayImageCodec(ByteOrder order) {
            this.order = order;
        }

        @Override
        public List<ImageInfo> load(Path file) throws IOException {
            return loadInfos(file, getType(), false, order, 2);
        }

        @Override
        public String getType() {
            return "LIB_2BIT_GRAY";
        }

        @Override
        public boolean accept(Path path) {
            File file = path.toFile();
            return !file.isDirectory() && Files.getFileExtension(file.getName()).equalsIgnoreCase("LIB");
        }

        @Override
        public BufferedImage read(ImageInfo info) throws IOException {
            return read0(info, order);
        }

        private BufferedImage read0(ImageInfo info, ByteOrder order) throws IOException {
            try (FileChannel ch = FileChannel.open(Paths.get(info.getFilename()), StandardOpenOption.READ)) {
                ByteBuf buf = Unpooled.wrappedBuffer(ch.map(FileChannel.MapMode.READ_ONLY, info.getOffset(), 8)).order(order);
                final int len = buf.readInt() - 12;
                int w = buf.readUnsignedShort(), h = buf.readUnsignedShort();
                ByteBuffer buffer = ch.map(FileChannel.MapMode.READ_ONLY, info.getOffset() + 12, len).order(order);
                byte[] gray;
                if (order == ByteOrder.BIG_ENDIAN) {
                    gray = new byte[]{-127, 127, 64, 0};
                } else {
                    gray = new byte[]{0, 64, 127, -127};
                }
                BufferedImage image = new BufferedImage(w, h, BufferedImage.TYPE_BYTE_BINARY, new IndexColorModel(2, 4, gray, gray, gray));

                byte[] data = ((DataBufferByte) image.getRaster().getDataBuffer()).getData();
                buffer.get(data);
                return image;
            }
        }
    }

    static class RLBImageCodec implements ImageCodec {
        @Override
        public List<ImageInfo> load(Path path) throws IOException {
            return loadInfos(path, getType(), true, ByteOrder.LITTLE_ENDIAN, -1);
        }

        @Override
        public BufferedImage read(Path path, int index) throws IOException {
            return read(load(path).get(index));
        }

        @Override
        public String getType() {
            return "RLB";
        }

        @Override
        public boolean accept(Path path) {
            File file = path.toFile();
            return !file.isDirectory() && Files.getFileExtension(file.getName()).equalsIgnoreCase("RLB");
        }

        @Override
        public BufferedImage read(ImageInfo info) throws IOException {
            FileInputStream is = new FileInputStream(new File(info.getFilename()));
            long skipped = is.skip(info.getOffset() + 4);
            assert skipped == info.getOffset() + 4;
            return ImageIO.read(is);
        }
    }

    static class GenericImage extends ImageInfo {
        private BufferedImage image;

        public GenericImage(int index, String name, int offset, String filename, String type, BufferedImage image) {
            super(index, name, offset, filename, type);
            this.image = image;
        }

        public BufferedImage getImage() {
            return image;
        }

        public GenericImage setImage(BufferedImage image) {
            this.image = image;
            return this;
        }
    }

    static class ImageInfo {
        private int index;
        private String name;
        private int offset;
        private String filename;
        private String type;
        private int width = -1;
        private int height = -1;
        private int size = -1;

        public ImageInfo() {
        }

        public ImageInfo(int index, String name, int offset, String filename, String type) {
            this.index = index;
            this.name = name;
            this.offset = offset;
            this.filename = filename;
            this.type = type;
        }

        public int getOffset() {
            return offset;
        }

        public int getIndex() {
            return index;
        }

        public String getType() {
            return type;
        }

        public String getName() {
            return name;
        }

        public String getFilename() {
            return filename;
        }

        public int getWidth() {
            return width;
        }

        public ImageInfo setWidth(int width) {
            this.width = width;
            return this;
        }

        public int getHeight() {
            return height;
        }

        public ImageInfo setHeight(int height) {
            this.height = height;
            return this;
        }

        /**
         * Pixel pre size
         */
        public int getBits() {
            return size * 8 / (width * height);
        }

        public int getSize() {
            return size;
        }

        public ImageInfo setSize(int size) {
            this.size = size;
            return this;
        }

        @Override
        public String toString() {
            return MoreObjects.toStringHelper(this)
                    .add("index", index)
                    .add("name", name)
                    .add("offset", offset)
                    .add("filename", filename)
                    .add("type", type)
                    .add("width", width)
                    .add("height", height)
                    .add("size", size)
                    .toString();
        }
    }
}
