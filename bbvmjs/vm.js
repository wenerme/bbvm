// import BBAsm from './bbasm'
var BBAsm = require('./bbasm');
var {RegType, Opcode, AddrType, DataType, CalcType, CmpType} =require('./const');
var util = require('util');
var fs = require('fs');
class BBVM {
    constructor(buffer) {
        this.ram = new DataView(buffer)
    }
}
class Assembly {
    constructor({op, a, b, dt, cmp, cal, comment}) {
        if (op && Opcode[op]) {
            this.op = Opcode[op]
        }
        if (a) {
            this.a = new Operand(a)
        }
        if (b) {
            this.b = new Operand(b)
        }
        if (dt) {
            this.dataType = DataType[dt]
        }
        if (cmp) {
            this.cmpType = CmpType[cmp]
        }
        if (cal) {
            this.calType = CalType[cmp]
        }
        if (comment) {
            this.comment = comment
        }
    }

    toAssembly() {

    }

    get length() {
        return 0
    }

    static from(asm) {
        if (asm.label) {
            return new Label(asm)
        }
        if (!asm.op && asm.comment) {
            return new Comment(asm)
        }
        switch (asm.op) {

        }
        var type = Assembly.opcodes[asm.op];
        if (type) {
            return new type(asm)
        }
        return new Instruction(asm)
    }
}
Assembly.opcodes = {};
class Instruction extends Assembly {
    constructor(args) {
        super(args);
    }

    get length() {
        return this.op.length
    }

    toAssembly() {
        var st = this.op.name;
        switch (this.op) {
            case Opcode.NOP:
            case Opcode.RET:
            case Opcode.EXIT:
                break;
            case Opcode.PUSH:
            case Opcode.POP:
            case Opcode.JMP:
                st += ` ${this.a.toAssembly()}`;
                break;
            case Opcode.JPC:
                st += ` ${this.cmpType.name} ${this.a.toAssembly()}`;
                break;
            case Opcode.IN:
            case Opcode.OUT:
                st += ` ${this.a.toAssembly()}, ${this.b.toAssembly()}`;
                break;
            case Opcode.LD:
                st += ` ${this.dataType.name}  ${this.a.toAssembly()}, ${this.b.toAssembly()}`;
                break;
            case Opcode.CMP:
                st += ` ${this.cmpType.name}  ${this.a.toAssembly()}, ${this.b.toAssembly()}`;
                break;
            case Opcode.CAL:
                st += ` ${this.calType.name} ${this.dataType.name}  ${this.a.toAssembly()}, ${this.b.toAssembly()}`;
                break
        }
        return st
    }
}
class Comment extends Assembly {
    constructor(args) {
        super(args)
    }

    toAssembly() {
        return `; ${this.comment}`
    }
}
class Label extends Assembly {
    constructor(args) {
        super(args);
        this.symbol = args.symbol
    }

    toAssembly() {
        return `${this.symbol}:`
    }
}
class PseudoData extends Assembly {
    constructor(args) {
        super(args);
        this.symbol = args.symbol;
        this.values = args.values.map(v=>new Value(v))
    }

    toAssembly() {
        return `DATA ${this.symbol || ''} ${this.dataType ? this.dataType.name : ''} ${this.values.map(v=>v.toAssembly()).join(", ")}`
    }
}
class PseudoBlock extends Assembly {
    constructor(args) {
        super(args);

        this.size = new Value(args.size);
        this.value = new Value(args.value)
    }

    toAssembly() {
        return `.BLOCK ${this.size.toAssembly()} ${this.value.toAssembly()}`
    }

    get length() {
        return this.size
    }
}
Assembly.opcodes['.BLOCK'] = PseudoBlock;
Assembly.opcodes['DATA'] = PseudoData;
class Value {
    constructor({type, value}) {
        this.type = type;
        this.value = value
    }

    toAssembly() {
        return `${this.value}`
    }
}
class Operand {
    constructor({direct, symbol, register, value}) {
        this.direct = direct || false;
        this.symbol = symbol;
        this.value = value ? new Value(value) : null;
        this.register = register || false
    }

    toAssembly() {
        var st = this.symbol || this.value.toAssembly();
        if (!this.direct) {
            st = `[ ${st} ]`
        }
        return st
    }
}


function normalizerInstruction(v) {
    for (let k of ["op", "dt", "cmp", "cal"]) {
        if (v[k]) {
            v[k] = v[k].toUpperCase()
        }
    }

    switch (v.op) {
        case '.BLOCK':
        case 'DATA':
            v.pseudo = true
    }

    var normalRegisterSymbol = (o)=> {
        if (o && o.symbol && /R[0-3PFSB]/i.test(o.symbol)) {
            o.register = true;
            o.symbol = o.symbol.toUpperCase()
        }
    };
    normalRegisterSymbol(v.a);
    normalRegisterSymbol(v.b);

    return v
}


fs.readFile('./cal.basm', 'utf8', (err, data)=> {
    if (err)throw err;
    // console.log(BBAsm.parse(data).map(normalizerInstruction).map(instructionToString).join('\n'))
    console.log(util.inspect(BBAsm.parse(data).map(normalizerInstruction).map(Assembly.from).map(v=>v.toAssembly()), false, null))
});
