package me.wener.bbvm.neo.processor;

import io.netty.buffer.ByteBuf;
import me.wener.bbvm.neo.BBVMContext;
import me.wener.bbvm.neo.Register;
import me.wener.bbvm.neo.define.RegisterType;

public class VMContext implements VMContextAware
{
    protected BBVMContext vm;
    protected Register rp;
    protected Register rb;
    protected Register rs;
    protected Register rf;
    protected Register r0;
    protected Register r1;
    protected Register r2;
    protected Register r3;
    protected ByteBuf memory;

    public void initialize(BBVMContext vm)
    {
        this.vm = vm;
        memory = vm.memory();
        rp = vm.register(RegisterType.rp);
        rb = vm.register(RegisterType.rb);
        rs = vm.register(RegisterType.rs);
        rf = vm.register(RegisterType.rf);
        r0 = vm.register(RegisterType.r0);
        r1 = vm.register(RegisterType.r1);
        r2 = vm.register(RegisterType.r2);
        r3 = vm.register(RegisterType.r3);
    }


    public BBVMContext push(int value)
    {
        vm.push(value);
        return vm;
    }

    public int pop()
    {
        return vm.pop();
    }
}
