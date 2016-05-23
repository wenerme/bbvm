'use strict';
class EnumInstance {
    constructor({name, value}) {
        this.name = name;
        this.value = value;
    }

    toString() {
        return this.name
    }
}

class EnumType {
    constructor(enums) {
        Object.assign(this, enums)
    }
}

function Enum(spec, type = EnumInstance) {
    let enums = Object
        .keys(spec)
        .reduce((o, k)=> {
            o[k] = new type({name: k, value: spec[k]});
            return o
        }, {});

    return Object.freeze(new EnumType(enums));
}

/**
 * <pre>
 * 表示     | 字节码 | 说明
 * --------|-------|----
 * rx      | 0x0   | 寄存器寻址
 * [rx]    | 0x1   | 寄存器间接寻址
 * n       | 0x2   | 立即数寻址
 * [n]     | 0x3   | 直接寻址
 *
 * op1/op2 |rx      | [rx]   | n      | [n]
 * --------|--------|--------|--------|----
 * rx      | 0x0    | 0x1    | 0x2    | 0x3
 * [rx]    | 0x4    | 0x5    | 0x6    | 0x7
 * n       | 0x8    | 0x9    | 0xa    | 0xb
 * [n]     | 0xc    | 0xd    | 0xe    | 0xf
 */
const AddrType = Enum({
    /**
     * 寄存器寻址
     */
    REGISTER: 0,
    /**
     * 寄存器间接寻址
     */
    REGISTER_DEFERRED: 1,
    /**
     * 立即数
     */
    IMMEDIATE: 2,
    /**
     * 直接寻址
     */
    DIRECT: 3
});

const DataType = Enum({
    DWORD: 0,
    WORD: 1,
    BYTE: 2,
    FLOAT: 3,
    INT: 4
});


const CalcType = Enum({
    ADD: 0,
    SUB: 1,
    MUL: 2,
    DIV: 3,
    MOD: 4
});
/**
 * <pre>
 * Z    | 0x1 | 等于
 * B    | 0x2 | Below,小于
 * BE   | 0x3 | 小于等于
 * A    | 0x4 | Above,大于
 * AE   | 0x5 | 大于等于
 * NZ   | 0x6 | 不等于
 */
const CmpType = Enum({
    Z: 1,
    B: 2,
    BE: 3,
    A: 4,
    AE: 5,
    NZ: 6,
});
const Opcode = Enum({
    NOP: 0,
    LD: 1,
    PUSH: 2,
    POP: 3,
    IN: 4,
    OUT: 5,
    JMP: 6,
    JPC: 7,
    CALL: 8,
    RET: 9,
    CMP: 0xA,
    CAL: 0xB,
    EXIT: 0xF,
}, class OpcodeType extends EnumInstance {
    constructor(args) {
        super(args)
    }

    get length() {
        switch (this) {
            case Opcode.NOP:
            case Opcode.RET:
            case Opcode.EXIT:
                return 1;
            case Opcode.PUSH:
            case Opcode.POP:
            case Opcode.JMP:
                return 5;
            case Opcode.JPC:
                return 6;
            case Opcode.IN:
            case Opcode.LD:
            case Opcode.CMP:
            case Opcode.CAL:
            case Opcode.OUT:
                return 10
        }
    }
});
/**
 * 寄存器类型
 * <pre>
 * RP | 0x0 | 程序计数器
 * RF | 0x1 | 比较标示符
 * RS | 0x2 | 栈顶位置
 * RB | 0x3 | 栈底位置
 * R0 | 0x4 | #0 寄存器
 * R1 | 0x5 | #1 寄存器
 * R2 | 0x6 | #2 寄存器
 * R3 | 0x7 | #3 寄存器
 * </pre>
 */
const RegType = Enum({
    RP: 0,
    RF: 1,
    RS: 2,
    RB: 3,
    R0: 4,
    R1: 5,
    R2: 6,
    R3: 7,
});

module.exports = {
    RegType, Opcode, AddrType, DataType, CalcType, CmpType
};
