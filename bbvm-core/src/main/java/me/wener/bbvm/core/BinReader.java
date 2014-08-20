package me.wener.bbvm.core;

public interface BinReader
{
    byte i16();
    int i32();
    long ui32();
    long i64();
    float f32();
    double f64();
    int offset();
    BinReader bytes(byte[] bytes, int offset, int len);
    byte[] bytes(int len);

    BinReader i16(byte v);
    BinReader i32(int v);
    BinReader ui32(long v);
    BinReader i64(long v);
    BinReader f32(float v);
    BinReader f64(double v);
    BinReader offset(int v);
    BinReader next(int n);
}
