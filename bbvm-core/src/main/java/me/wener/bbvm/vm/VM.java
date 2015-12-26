package me.wener.bbvm.vm;

import com.google.common.base.Preconditions;
import com.google.common.base.Throwables;
import com.google.common.eventbus.EventBus;
import com.google.inject.Guice;
import io.netty.buffer.ByteBuf;
import me.wener.bbvm.dev.StringManager;
import me.wener.bbvm.exception.ExecutionException;
import me.wener.bbvm.exception.ResourceMissingException;
import me.wener.bbvm.util.IntEnums;
import me.wener.bbvm.vm.event.ResetEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.annotation.Nonnull;
import javax.inject.Inject;
import java.util.Iterator;
import java.util.Map;
import java.util.NavigableMap;

import static com.google.common.base.Preconditions.checkState;

/**
 * @author wener
 * @since 15/12/10
 */
public class VM {
    private final static Logger log = LoggerFactory.getLogger(VM.class);
    //    Injector injector;
    final Register r0 = new Reg(RegisterType.R0, this);
    final Register r1 = new Reg(RegisterType.R1, this);
    final Register r2 = new Reg(RegisterType.R2, this);
    final Register r3 = new Reg(RegisterType.R3, this);
    final Register rs = new Reg(RegisterType.RS, this);
    final Register rf = new Reg(RegisterType.RF, this);
    final Register rb = new Reg(RegisterType.RB, this);
    final Register rp = new Reg(RegisterType.RP, this);
    Memory memory;
    SymbolTable symbolTable;
    NavigableMap<Integer, Integer> addressTable;
    @Inject
    StringManager stringManager;
    @Inject
    SystemInvokeManager systemInvokeManager;
    private boolean exit = false;
    @Inject
    private VMConfig config;
    private Throwable lastError;
    @Inject
    private EventBus eventBus;

    @Inject
    private VM() {
    }

    public static VM create() {
        return Guice.createInjector(new VirtualMachineModule(new VMConfig.Builder().build())).getInstance(VM.class);
    }

    private static double cal(CalculateType calculateType, DataType dataType, Operand a, Operand b) {
        int ia = a.get(), ib = b.get();
        switch (dataType) {
            case FLOAT: {
                float fa = a.getFloat(), fb = b.getFloat();
                switch (calculateType) {
                    case ADD:
                        fa += fb;
                        break;
                    case SUB:
                        fa -= fb;
                        break;
                    case MUL:
                        fa *= fb;
                        break;
                    case DIV:
                        fa /= fb;
                        break;
                    case MOD:
                        fa %= fb;
                        break;
                }
                return fa;
            }
            case WORD:
                ia &= 0xffff;
                ib &= 0xffff;
                break;
            case BYTE:
                ia &= 0xff;
                ib &= 0xff;
                break;
        }
        switch (calculateType) {
            case ADD:
                ia += ib;
                break;
            case SUB:
                ia -= ib;
                break;
            case MUL:
                ia *= ib;
                break;
            case DIV:
                ia /= ib;
                break;
            case MOD:
                ia %= ib;
                break;
        }
        return ia;
    }

    private static double cal(CalculateType calculateType, DataType dataType, Operand a, Operand b, Value out) {
        double vc = cal(calculateType, dataType, a, b);
        switch (dataType) {
            case FLOAT:
                out.set((float) vc);
                break;
            case DWORD:
            case INT:
                // Allowed overflow
                out.set((int) (((long) vc)));
                break;
            case WORD:
                // TODO is unsigned ok ?
                // Same on the PC vm, not sure in the 9688
                out.set(((short) vc) & 0xffff);
                break;
            case BYTE:
                out.set(((byte) vc) & 0xff);
                break;
        }
        return vc;
    }

    public NavigableMap<Integer, Integer> getAddressTable() {
        return addressTable;
    }

    public VM setAddressTable(NavigableMap<Integer, Integer> addressTable) {
        this.addressTable = addressTable;
        return this;
    }

    /**
     * @return Tick clock
     */
    public int getTick() {
        // TODO Is this ok ?
        return (int) (System.currentTimeMillis() & 0xfffffff);
    }

    @Inject
    private void init(SystemInvokeManager systemInvokeManager) {
        config.getInvokeHandlers().forEach(systemInvokeManager::register);
    }

    public StringManager getStringManager() {
        return stringManager;
    }

    public SymbolTable getSymbolTable() {
        return symbolTable;
    }

    public VM setSymbolTable(SymbolTable symbolTable) {
        this.symbolTable = symbolTable;
        return this;
    }

    boolean hasRemaining() {
        return rp.get() < memory.getMemorySize();
    }

    public Throwable getLastError() {
        return lastError;
    }

    public void run() {
        Preconditions.checkState(!exit);
        Instruction instruction = new Instruction().setVm(this);
        ByteBuf buf = this.memory.getByteBuf();
        int last;
        while (hasRemaining()) {
            last = rp.get();
            instruction.reset().read(buf, last);
            try {
                run(instruction);
            } catch (Exception e) {
                log.warn("Catch exception when > {} ' {} @ {}", instruction, debugAsm(), getLine(instruction.address), e);
                if (e instanceof ExecutionException) {
                    ((ExecutionException) e).setVm(this);
                }
                lastError = e;
                if (config.getErrorHandler().apply(e)) {
                    exit();
                }
                Throwables.propagateIfInstanceOf(e, ExecutionException.class);
                throw new ExecutionException(e);
            }
            if (exit) {
                return;
            }
            if (rp.get() == last) {
                rp.add(instruction.getOpcode().length());
            }
        }
    }

    public VM exit() {
        Preconditions.checkState(!exit);
        log.info("Exit vm");
        exit = true;
        return this;
    }

    public Iterable<Instruction> instructions(final Instruction instruction, final int position) {
        return new InstructionIterable(position, instruction);
    }

    public Memory getMemory() {
        return memory;
    }

    public VM setMemory(@Nonnull Memory memory) {
        this.memory = memory.setVm(this);
        return this;
    }

    /**
     * Reset virtual machine state, should be this before every rerun
     */
    public VM reset() {
        r0.set(0);
        r1.set(0);
        r2.set(0);
        r3.set(0);
        rs.set(0);
        rf.set(0);
        rb.set(0);
        rp.set(0);
        exit = false;
        if (memory != null) {
            rs.set(memory.getMemorySize() + memory.getStackSize() - 4);
            rb.set(memory.getMemorySize() - 4);
        }
        log.debug("VM Reset {}", debugAsm());
        eventBus.post(new ResetEvent(this));
        return this;
    }

    public int getLine(int address) {
        if (addressTable == null) {
            return -1;
        }
        Map.Entry<Integer, Integer> entry = addressTable.ceilingEntry(address);
        if (entry == null) {
            return -1;
        }
        return entry.getValue();
    }

    public void run(Instruction inst) {
        checkState(!exit, "Exited");
        if (log.isTraceEnabled()) {
            log.trace("{}", inst);
        }
        if (log.isDebugEnabled()) {
            log.debug("{} ' A={} B={} {} @{} {}",
                    inst.toAssembly(),
                    inst.hasA() ? inst.getA().get() : "NaN",
                    inst.hasB() ? inst.getB().get() : "NaN",
                    debugAsm(), getLine(inst.getAddress()), inst.getAddress());
        }
        run(inst, inst.opcode, inst.a, inst.b);
    }

    String debugAsm() {
        return String.format("RP=%s RF=%s RS=%s RB=%s R0=%s R1=%s R2=%s R3=%s"
                , rp.get(), rf.get(), rs.get(), rb.get(), r0.get(), r1.get(), r2.get(), r3.get());
    }

    private void run(Instruction inst, Opcode opcode, Operand a, Operand b) {
        switch (opcode) {
            case NOP:
                break;
            case LD:
                // TODO Data type overflow check
                a.set(b.get());
                break;
            case PUSH:
                push(a.get());
                break;
            case POP:
                a.set(pop());
                break;
            case OUT:
            case IN:
                systemInvokeManager.invoke(inst);
                break;
            case JMP:
                jmp(a.get());
                break;
            case JPC:
                if (IntEnums.fromInt(CompareType.class, rf.get()).isMatch(inst.compareType)) {
                    jmp(a.get());
                }
                break;
            case CALL:
                push(rp.get() + inst.getOpcode().length());
                jmp(a.get());
                break;
            case RET:
                ret();
                break;
            case CMP: {
                double vc = cal(CalculateType.SUB, inst.getDataType(), a, b);
                if (vc > 0)
                    rf.set(CompareType.A);
                else if (vc < 0)
                    rf.set(CompareType.B);
                else
                    rf.set(CompareType.Z);
            }
            break;
            case CAL: {
                cal(inst.getCalculateType(), inst.getDataType(), a, b, a);
            }
            break;
            case EXIT:
                exit = true;
                break;
        }
    }

    public boolean isExit() {
        return exit;
    }

    void jmp(int i) {
        rp.set(i);
    }

    void ret() {
        rp.set(pop());
    }

    VM push(int v) {
        memory.push(v);
        return this;
    }

    int pop() {
        return memory.pop();
    }

    /**
     * If i >= 0 then read the string from memory, if i < 0 will get string from string resource.
     * If i > memory size or ResourceMissing will return null
     */
    public String getString(int i) {
        if (i >= 0) {
            if (i > memory.getMemorySize()) {
                return null;
            }
            return memory.getString(i, config.getCharset());
        }
        try {
            return stringManager.getResource(i).getValue();
        } catch (ResourceMissingException e) {
            return null;
        }
    }

    public Register getRegister(RegisterType type) {
        switch (type) {
            case RP:
                return rp;
            case RF:
                return rf;
            case RS:
                return rs;
            case RB:
                return rb;
            case R0:
                return r0;
            case R1:
                return r1;
            case R2:
                return r2;
            case R3:
                return r3;
        }
        throw new UnsupportedOperationException();
    }

    public Symbol getSymbol(int address) {
        return symbolTable != null ? symbolTable.getSymbol(address) : null;
    }

    /**
     * @author wener
     * @since 15/12/15
     */
    private static class Reg implements Register {
        final RegisterType type;
        final VM vm;
        int value;

        public Reg(RegisterType type, VM vm) {
            this.type = type;
            this.vm = vm;
        }

        public VM getVm() {
            return vm;
        }

        @Override
        public int get() {
            return value;
        }

        public Reg set(int value) {
            this.value = value;
            return this;
        }

        @Override
        public RegisterType getType() {
            return type;
        }
    }

    private class InstructionIterable implements Iterable<Instruction> {
        private final int position;
        private final Instruction instruction;

        public InstructionIterable(int position, Instruction instruction) {
            this.position = position;
            this.instruction = instruction;
        }

        @Override
        public Iterator<Instruction> iterator() {
            return new Iterator<Instruction>() {
                int pos = position;

                @Override
                public boolean hasNext() {
                    return pos < memory.getMemorySize();
                }

                @Override
                public Instruction next() {
                    instruction.reset().read(memory.getByteBuf(), pos);
                    pos += instruction.opcode.length();
                    return instruction;
                }
            };
        }
    }
}
