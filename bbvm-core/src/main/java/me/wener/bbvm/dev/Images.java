package me.wener.bbvm.dev;

import com.google.common.base.MoreObjects;
import com.google.common.collect.Maps;
import com.google.common.io.Files;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;

import javax.imageio.ImageIO;
import java.awt.image.BufferedImage;
import java.awt.image.DataBufferUShort;
import java.awt.image.DirectColorModel;
import java.awt.image.WritableRaster;
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
    private static final int DCM_565_RED_MASK = 0xf800;
    private static final int DCM_565_GRN_MASK = 0x07E0;
    private static final int DCM_565_BLU_MASK = 0x001F;

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


    private static List<ImageInfo> loadInfos(Path file, String type, boolean hasName) throws IOException {
        try (FileChannel ch = FileChannel.open(file, StandardOpenOption.READ)) {
            int n = ch.map(FileChannel.MapMode.READ_ONLY, 0, 4).order(ByteOrder.LITTLE_ENDIAN).getInt();
            ByteBuf buf = Unpooled.wrappedBuffer(ch.map(FileChannel.MapMode.READ_ONLY, 4, n * (4 + 32))).order(ByteOrder.LITTLE_ENDIAN);
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
            return loadInfos(file, getType(), false);
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
                // Fast way to load the data
                DirectColorModel colorModel = new DirectColorModel(16, DCM_565_RED_MASK, DCM_565_GRN_MASK, DCM_565_BLU_MASK);
                WritableRaster raster = colorModel.createCompatibleWritableRaster(w, h);

                buffer.flip();
                buffer.asShortBuffer().get(((DataBufferUShort) raster.getDataBuffer()).getData());
                return new BufferedImage(colorModel, raster, false, null);
            }
        }
    }

    static class RLBImageCodec implements ImageCodec {

        @Override
        public List<ImageInfo> load(Path path) throws IOException {
            return loadInfos(path, getType(), true);
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
