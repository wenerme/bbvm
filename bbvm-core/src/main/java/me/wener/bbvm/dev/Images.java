package me.wener.bbvm.dev;

import com.google.common.base.MoreObjects;
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
    }

    public static List<ImageInfo> load(String file) throws IOException {
        Path path = Paths.get(file);
        for (Map.Entry<String, ImageCodec> entry : CODEC.entrySet()) {
            ImageCodec codec = entry.getValue();
            if (codec.accept(path)) {
                return codec.load(path);
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


    private static List<ImageInfo> loadInfos(Path file, String type, boolean hasName, ByteOrder order) throws IOException {
        try (FileChannel ch = FileChannel.open(file, StandardOpenOption.READ)) {
            int n = ch.map(FileChannel.MapMode.READ_ONLY, 0, 4).order(order).getInt();
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
                list.add(new ImageInfo(i, name, offset, file.toString(), type));

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
            return loadInfos(file, getType(), false, ByteOrder.LITTLE_ENDIAN);
        }

        @Override
        public String getType() {
            return "LIB_RGB565";
        }

        @Override
        public boolean accept(Path path) {
            File file = path.toFile();
            return !file.isDirectory() && Files.getFileExtension(file.getName()).equalsIgnoreCase("RLB");
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

    static class Lib2BitGrayImageCodec implements ImageCodec {

        @Override
        public List<ImageInfo> load(Path file) throws IOException {
            try {
                return loadInfos(file, getType() + "_LE", false, ByteOrder.LITTLE_ENDIAN);
            } catch (IOException e) {
                return loadInfos(file, getType() + "_BE", false, ByteOrder.BIG_ENDIAN);
            }
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
            if (info.getType().endsWith("_LE")) {
                return read0(info, ByteOrder.LITTLE_ENDIAN);
            } else {
                return read0(info, ByteOrder.BIG_ENDIAN);
            }
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
            return loadInfos(path, getType(), true, ByteOrder.LITTLE_ENDIAN);
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

    static class ImageInfo {
        private final int index;
        private final String name;
        private final int offset;
        private final String filename;
        private final String type;

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

        @Override
        public String toString() {
            return MoreObjects.toStringHelper(this)
                    .add("index", index)
                    .add("name", name)
                    .add("offset", offset)
                    .add("filename", filename)
                    .add("type", type)
                    .toString();
        }
    }
}
