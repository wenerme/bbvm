grammar BBAsm;

// start point
prog
    : stat+
    ;

stat
    : NEWLINE
    | expr NEWLINE
    | comment
    ;

comment
	: LINE_COMMENT NEWLINE
	| COMMENT
	;


expr
    : instruction
    | LABEL // 定义标签的语句
    ;

// 所有可能的指令集
instruction
    : NoOperandIns
    | OneOperandIns operand
    | TowOperandIns operand COMMA operand
	| INS_CAL DataType CalculateOperator operand COMMA operand // CAL int ADD r0,12
	| INS_LD DataType operand COMMA operand  // ld int r1, 1067320848
	| INS_JMP CompareOperator operand // jpc a some-where
	| INS_BLOCK IntegerLiteral IntegerLiteral // .block 1 10
    ;

operand
    : Register						// 使用寄存器
    | LBRACK Register RBRACK
    | ConstFormula					// 使用常量表达式
    | LBRACK ConstFormula RBRACK
    | Identifier					// 使用标识符
    | LBRACK Identifier RBRACK
    ;

// 常量表达式,在编译期间可以进行求值
ConstFormula
	: IntegerLiteral
	;
	
// 无操作数的指令
NoOperandIns
    : INS_NOP
    | INS_RET
    | INS_EXIT
    ;
// 一个操作数的指令
OneOperandIns
    : INS_JMP
    | INS_CALL
    | INS_PUSH
    | INS_POP
    ;
// 两个操作数的指令
TowOperandIns
    : INS_IN
    | INS_OUT
    ;


literal
    :   IntegerLiteral
    |   FloatingPointLiteral
    |   CharacterLiteral
    |   StringLiteral
    |   BooleanLiteral
    |   'null'
    ;


Register
    : RP
    | RF
    | RS
    | RB
    | R0
    | R1
    | R2
    | R3
    ;

// 寄存器类型
fragment RP : R P ;
fragment RF : R F ;
fragment RS : R S ;
fragment RB : R B ;
fragment R0 : R '0' ;
fragment R1 : R '1' ;
fragment R2 : R '2' ;
fragment R3 : R '3' ;

// 比较操作
CompareOperator
	: CMP_Z	
	| CMP_B 	
	| CMP_BE	
	| CMP_A
	| CMP_AE
	| CMP_NZ
	;
	
fragment CMP_Z	: Z	  ;
fragment CMP_B 	: B   ;
fragment CMP_BE	: B E ;
fragment CMP_A 	: A   ;
fragment CMP_AE	: A E ;
fragment CMP_NZ	: N Z ;

// 计算操作
CalculateOperator
	: CAL_ADD
	| CAL_SUB
	| CAL_MUL
	| CAL_DIV
	| CAL_MOD
	;

fragment CAL_ADD : A D D ;
fragment CAL_SUB : S U B ;
fragment CAL_MUL : M U L ;
fragment CAL_DIV : D I V ;
fragment CAL_MOD : M O D ;

// 数据类型
DataType
	: DataType_DWORD
	| DataType_WORD
	| DataType_BYTE
	| DataType_FLOAT
	| DataType_INT
	;

fragment DataType_DWORD :D W O R D ;
fragment DataType_WORD  :W O R D   ;
fragment DataType_BYTE  :B Y T E   ;
fragment DataType_FLOAT :F L O A T ;
fragment DataType_INT   :I N T     ;

// 所有的指令
INS_NOP    : N O P;
INS_LD     : L D;
INS_PUSH   : P U S H;
INS_POP    : P O P;
INS_IN     : I N;
INS_OUT    : O U T;
INS_JMP    : J M P;
INS_JPC    : J P C;
INS_CALL   : C A L L;
INS_RET    : R E T;
INS_CMP    : C M P;
INS_CAL    : C A L;
INS_EXIT   : E X I T;
INS_DATA   : D A T A;
INS_BLOCK  : DOT? B L O C K;
// 标签格式
LABEL : Identifier ':' ;

// 常量算式,在编译期能够得到结果的

// 整数字面值
IntegerLiteral
    :   DecimalIntegerLiteral
    |   HexIntegerLiteral
    |   OctalIntegerLiteral
    |   BinaryIntegerLiteral
    ;

fragment
DecimalIntegerLiteral
    :   DecimalNumeral IntegerTypeSuffix?
    ;

fragment
HexIntegerLiteral
    :   HexNumeral IntegerTypeSuffix?
    ;

fragment
OctalIntegerLiteral
    :   OctalNumeral IntegerTypeSuffix?
    ;

fragment
BinaryIntegerLiteral
    :   BinaryNumeral IntegerTypeSuffix?
    ;

fragment
IntegerTypeSuffix
    :   [lL]
    ;

fragment
DecimalNumeral
    :   '0'
    |   NonZeroDigit (Digits? | Underscores Digits)
    ;

fragment
Digits
    :   Digit (DigitOrUnderscore* Digit)?
    ;

fragment
Digit
    :   '0'
    |   NonZeroDigit
    ;

fragment
NonZeroDigit
    :   [1-9]
    ;

fragment
DigitOrUnderscore
    :   Digit
    |   '_'
    ;

fragment
Underscores
    :   '_'+
    ;

fragment
HexNumeral
    :   '0' [xX] HexDigits
    ;

fragment
HexDigits
    :   HexDigit (HexDigitOrUnderscore* HexDigit)?
    ;

fragment
HexDigit
    :   [0-9a-fA-F]
    ;

fragment
HexDigitOrUnderscore
    :   HexDigit
    |   '_'
    ;

fragment
OctalNumeral
    :   '0' Underscores? OctalDigits
    ;

fragment
OctalDigits
    :   OctalDigit (OctalDigitOrUnderscore* OctalDigit)?
    ;

fragment
OctalDigit
    :   [0-7]
    ;

fragment
OctalDigitOrUnderscore
    :   OctalDigit
    |   '_'
    ;

fragment
BinaryNumeral
    :   '0' [bB] BinaryDigits
    ;

fragment
BinaryDigits
    :   BinaryDigit (BinaryDigitOrUnderscore* BinaryDigit)?
    ;

fragment
BinaryDigit
    :   [01]
    ;

fragment
BinaryDigitOrUnderscore
    :   BinaryDigit
    |   '_'
    ;

// 浮点数字面值
FloatingPointLiteral
    :   DecimalFloatingPointLiteral
    |   HexadecimalFloatingPointLiteral
    ;

fragment
DecimalFloatingPointLiteral
    :   Digits '.' Digits? ExponentPart? FloatTypeSuffix?
    |   '.' Digits ExponentPart? FloatTypeSuffix?
    |   Digits ExponentPart FloatTypeSuffix?
    |   Digits FloatTypeSuffix
    ;

fragment
ExponentPart
    :   ExponentIndicator SignedInteger
    ;

fragment
ExponentIndicator
    :   [eE]
    ;

fragment
SignedInteger
    :   Sign? Digits
    ;

fragment
Sign
    :   [+-]
    ;

fragment
FloatTypeSuffix
    :   [fFdD]
    ;

fragment
HexadecimalFloatingPointLiteral
    :   HexSignificand BinaryExponent FloatTypeSuffix?
    ;

fragment
HexSignificand
    :   HexNumeral '.'?
    |   '0' [xX] HexDigits? '.' HexDigits
    ;

fragment
BinaryExponent
    :   BinaryExponentIndicator SignedInteger
    ;

fragment
BinaryExponentIndicator
    :   [pP]
    ;

// 布尔字面值
BooleanLiteral
    :   'true'
    |   'false'
    ;

// 字符字面值
CharacterLiteral
    :   '\'' SingleCharacter '\''
    |   '\'' EscapeSequence '\''
    ;

fragment
SingleCharacter
    :   ~['\\]
    ;

// 字符串字面值
StringLiteral
    :   '"' StringCharacters? '"'
    ;

fragment
StringCharacters
    :   StringCharacter+
    ;

fragment
StringCharacter
    :   ~["\\]
    |   EscapeSequence
    ;

// 字符和字符串的转义序列
fragment
EscapeSequence
    :   '\\' [btnfr"'\\]
    |   OctalEscape
    |   UnicodeEscape
    ;

fragment
OctalEscape
    :   '\\' OctalDigit
    |   '\\' OctalDigit OctalDigit
    |   '\\' ZeroToThree OctalDigit OctalDigit
    ;

fragment
UnicodeEscape
    :   '\\' 'u' HexDigit HexDigit HexDigit HexDigit
    ;

fragment
ZeroToThree
    :   [0-3]
    ;

// The Null Literal

NullLiteral
    :   'null'
    ;

// Separators

LPAREN          : '(';
RPAREN          : ')';
LBRACE          : '{';
RBRACE          : '}';
LBRACK          : '[';
RBRACK          : ']';
SEMI            : ';';
COMMA           : ',';
DOT             : '.';

// Operators

ASSIGN          : '=';
GT              : '>';
LT              : '<';
BANG            : '!';
TILDE           : '~';
QUESTION        : '?';
COLON           : ':';
EQUAL           : '==';
LE              : '<=';
GE              : '>=';
NOTEQUAL        : '!=';
AND             : '&&';
OR              : '||';
INC             : '++';
DEC             : '--';
ADD             : '+';
SUB             : '-';
MUL             : '*';
DIV             : '/';
BITAND          : '&';
BITOR           : '|';
CARET           : '^';
MOD             : '%';

ADD_ASSIGN      : '+=';
SUB_ASSIGN      : '-=';
MUL_ASSIGN      : '*=';
DIV_ASSIGN      : '/=';
AND_ASSIGN      : '&=';
OR_ASSIGN       : '|=';
XOR_ASSIGN      : '^=';
MOD_ASSIGN      : '%=';
LSHIFT_ASSIGN   : '<<=';
RSHIFT_ASSIGN   : '>>=';
URSHIFT_ASSIGN  : '>>>=';

// 标识符 必须在所有关键词之后
Identifier
    :   JavaLetter JavaLetterOrDigit*
    ;

fragment
JavaLetter
    :   [a-zA-Z$_] // these are the "java letters" below 0xFF
    |   // covers all characters above 0xFF which are not a surrogate
        ~[\u0000-\u00FF\uD800-\uDBFF]
        {Character.isJavaIdentifierStart(_input.LA(-1))}?
    |   // covers UTF-16 surrogate pairs encodings for U+10000 to U+10FFFF
        [\uD800-\uDBFF] [\uDC00-\uDFFF]
        {Character.isJavaIdentifierStart(Character.toCodePoint((char)_input.LA(-2), (char)_input.LA(-1)))}?
    ;

fragment
JavaLetterOrDigit
    :   [a-zA-Z0-9$_] // these are the "java letters or digits" below 0xFF
    |   // covers all characters above 0xFF which are not a surrogate
        ~[\u0000-\u00FF\uD800-\uDBFF]
        {Character.isJavaIdentifierPart(_input.LA(-1))}?
    |   // covers UTF-16 surrogate pairs encodings for U+10000 to U+10FFFF
        [\uD800-\uDBFF] [\uDC00-\uDFFF]
        {Character.isJavaIdentifierPart(Character.toCodePoint((char)_input.LA(-2), (char)_input.LA(-1)))}?
    ;

//
// Additional symbols not defined in the lexical specification
//

AT : '@';
ELLIPSIS : '...';

//
// Whitespace and comments
//

WS  :  [ \t\u000C]+ -> skip
    ;

COMMENT
    :   '/*' .*? '*/'
    ;

LINE_COMMENT
    :   ('//' | ';' | '\'') ~[\r\n]*
    ;

NEWLINE:'\r'? '\n' ; // return newlines to parser (is end-statement signal)

// 便于做大小写无关的语法
fragment A:('a'|'A');
fragment B:('b'|'B');
fragment C:('c'|'C');
fragment D:('d'|'D');
fragment E:('e'|'E');
fragment F:('f'|'F');
fragment G:('g'|'G');
fragment H:('h'|'H');
fragment I:('i'|'I');
fragment J:('j'|'J');
fragment K:('k'|'K');
fragment L:('l'|'L');
fragment M:('m'|'M');
fragment N:('n'|'N');
fragment O:('o'|'O');
fragment P:('p'|'P');
fragment Q:('q'|'Q');
fragment R:('r'|'R');
fragment S:('s'|'S');
fragment T:('t'|'T');
fragment U:('u'|'U');
fragment V:('v'|'V');
fragment W:('w'|'W');
fragment X:('x'|'X');
fragment Y:('y'|'Y');
fragment Z:('z'|'Z');