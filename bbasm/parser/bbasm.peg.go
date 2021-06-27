package parser

import (
	"fmt"
	"github.com/wenerme/bbvm/bbasm"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleStart
	ruleAssembly
	ruleComment
	ruleLabel
	ruleInst
	rulePseudo
	rulePseudoDataValue
	rulePSEUDO_DATA_TYPE
	ruleOperand
	ruleSpacing
	ruleSpace
	ruleIdentifier
	ruleLetter
	ruleLetterOrDigit
	ruleEXIT
	ruleRET
	ruleNOP
	ruleCALL
	rulePUSH
	rulePOP
	ruleJMP
	ruleIN
	ruleOUT
	ruleCAL
	ruleLD
	ruleCMP
	ruleJPC
	ruleBLOCK
	ruleDATA
	ruleCAL_OP
	ruleCMP_OP
	ruleDATA_TYPE
	ruleLBRK
	ruleRBRK
	ruleCOMMA
	ruleSEMICOLON
	ruleCOLON
	ruleMINUS
	ruleNL
	ruleEOT
	ruleLiteral
	ruleIntegerLiteral
	ruleDecimalNumeral
	ruleHexNumeral
	ruleBinaryNumeral
	ruleOctalNumeral
	ruleFloatLiteral
	ruleDecimalFloat
	ruleExponent
	ruleHexFloat
	ruleHexSignificand
	ruleBinaryExponent
	ruleDigits
	ruleHexDigits
	ruleHexDigit
	ruleCharLiteral
	ruleStringLiteral
	ruleEscape
	ruleOctalEscape
	ruleUnicodeEscape
	ruleAction0
	ruleAction1
	ruleAction2
	rulePegText
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30
	ruleAction31
	ruleAction32
	ruleAction33
	ruleAction34
	ruleAction35

	rulePre
	ruleIn
	ruleSuf
)

var rul3s = [...]string{
	"Unknown",
	"Start",
	"Assembly",
	"Comment",
	"Label",
	"Inst",
	"Pseudo",
	"PseudoDataValue",
	"PSEUDO_DATA_TYPE",
	"Operand",
	"Spacing",
	"Space",
	"Identifier",
	"Letter",
	"LetterOrDigit",
	"EXIT",
	"RET",
	"NOP",
	"CALL",
	"PUSH",
	"POP",
	"JMP",
	"IN",
	"OUT",
	"CAL",
	"LD",
	"CMP",
	"JPC",
	"BLOCK",
	"DATA",
	"CAL_OP",
	"CMP_OP",
	"DATA_TYPE",
	"LBRK",
	"RBRK",
	"COMMA",
	"SEMICOLON",
	"COLON",
	"MINUS",
	"NL",
	"EOT",
	"Literal",
	"IntegerLiteral",
	"DecimalNumeral",
	"HexNumeral",
	"BinaryNumeral",
	"OctalNumeral",
	"FloatLiteral",
	"DecimalFloat",
	"Exponent",
	"HexFloat",
	"HexSignificand",
	"BinaryExponent",
	"Digits",
	"HexDigits",
	"HexDigit",
	"CharLiteral",
	"StringLiteral",
	"Escape",
	"OctalEscape",
	"UnicodeEscape",
	"Action0",
	"Action1",
	"Action2",
	"PegText",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"Action31",
	"Action32",
	"Action33",
	"Action34",
	"Action35",

	"Pre_",
	"_In_",
	"_Suf",
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (node *node32) Print(buffer string) {
	node.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: ruleIn, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: ruleSuf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens32) Expand(index int) {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
}

type BBAsm struct {
	parser

	Buffer string
	buffer []rune
	rules  [98]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokens32
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *BBAsm
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *BBAsm) PrintSyntaxTree() {
	p.tokens32.PrintSyntaxTree(p.Buffer)
}

func (p *BBAsm) Highlighter() {
	p.PrintSyntax()
}

func (p *BBAsm) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.line++
		case ruleAction1:
			p.AddAssembly()
		case ruleAction2:
			p.AddAssembly()
			p.AddComment()
		case ruleAction3:
			p.Push(&Comment{})
			p.Push(text)
		case ruleAction4:
			p.Push(&Label{})
		case ruleAction5:
			p.Push(lookup(bbasm.INT, text))
		case ruleAction6:
			p.Push(lookup(bbasm.ADD, text))
		case ruleAction7:
			p.Push(lookup(bbasm.INT, text))
		case ruleAction8:
			p.Push(lookup(bbasm.INT, text))
		case ruleAction9:
			p.Push(lookup(bbasm.A, text))
		case ruleAction10:
			p.AddPseudoDataValue()
		case ruleAction11:
			p.AddPseudoDataValue()
		case ruleAction12:
			p.AddPseudoDataValue()
		case ruleAction13:
			p.Push(text)
			p.AddPseudoDataValue()
		case ruleAction14:
			p.AddOperand(true)
		case ruleAction15:
			p.AddOperand(false)
		case ruleAction16:
			p.AddOperand(true)
		case ruleAction17:
			p.AddOperand(false)
		case ruleAction18:
			p.Push(text)
		case ruleAction19:
			p.PushInst(bbasm.EXIT)
		case ruleAction20:
			p.PushInst(bbasm.RET)
		case ruleAction21:
			p.PushInst(bbasm.NOP)
		case ruleAction22:
			p.PushInst(bbasm.CALL)
		case ruleAction23:
			p.PushInst(bbasm.PUSH)
		case ruleAction24:
			p.PushInst(bbasm.POP)
		case ruleAction25:
			p.PushInst(bbasm.JMP)
		case ruleAction26:
			p.PushInst(bbasm.IN)
		case ruleAction27:
			p.PushInst(bbasm.OUT)
		case ruleAction28:
			p.PushInst(bbasm.CAL)
		case ruleAction29:
			p.PushInst(bbasm.LD)
		case ruleAction30:
			p.PushInst(bbasm.CMP)
		case ruleAction31:
			p.PushInst(bbasm.JPC)
		case ruleAction32:
			p.Push(&PseudoBlock{})
		case ruleAction33:
			p.Push(&PseudoData{})
		case ruleAction34:
			p.Push(text)
			p.AddInteger()
		case ruleAction35:
			p.Push(text)

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *BBAsm) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
		p.buffer = append(p.buffer, endSymbol)
	}

	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	var max token32
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		tree.Expand(tokenIndex)
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position, depth}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Start <- <((Spacing Assembly? NL Action0)* EOT Literal?)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
			l2:
				{
					position3, tokenIndex3, depth3 := position, tokenIndex, depth
					if !_rules[ruleSpacing]() {
						goto l3
					}
					{
						position4, tokenIndex4, depth4 := position, tokenIndex, depth
						if !_rules[ruleAssembly]() {
							goto l4
						}
						goto l5
					l4:
						position, tokenIndex, depth = position4, tokenIndex4, depth4
					}
				l5:
					if !_rules[ruleNL]() {
						goto l3
					}
					if !_rules[ruleAction0]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position3, tokenIndex3, depth3
				}
				if !_rules[ruleEOT]() {
					goto l0
				}
				{
					position6, tokenIndex6, depth6 := position, tokenIndex, depth
					if !_rules[ruleLiteral]() {
						goto l6
					}
					goto l7
				l6:
					position, tokenIndex, depth = position6, tokenIndex6, depth6
				}
			l7:
				depth--
				add(ruleStart, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Assembly <- <((Label / ((&('.' | 'D' | 'd') Pseudo) | (&(';') Comment) | (&('C' | 'E' | 'I' | 'J' | 'L' | 'N' | 'O' | 'P' | 'R' | 'c' | 'e' | 'i' | 'j' | 'l' | 'n' | 'o' | 'p' | 'r') Inst))) Action1 (Comment Action2)?)> */
		func() bool {
			position8, tokenIndex8, depth8 := position, tokenIndex, depth
			{
				position9 := position
				depth++
				{
					position10, tokenIndex10, depth10 := position, tokenIndex, depth
					if !_rules[ruleLabel]() {
						goto l11
					}
					goto l10
				l11:
					position, tokenIndex, depth = position10, tokenIndex10, depth10
					{
						switch buffer[position] {
						case '.', 'D', 'd':
							if !_rules[rulePseudo]() {
								goto l8
							}
							break
						case ';':
							if !_rules[ruleComment]() {
								goto l8
							}
							break
						default:
							if !_rules[ruleInst]() {
								goto l8
							}
							break
						}
					}

				}
			l10:
				if !_rules[ruleAction1]() {
					goto l8
				}
				{
					position13, tokenIndex13, depth13 := position, tokenIndex, depth
					if !_rules[ruleComment]() {
						goto l13
					}
					if !_rules[ruleAction2]() {
						goto l13
					}
					goto l14
				l13:
					position, tokenIndex, depth = position13, tokenIndex13, depth13
				}
			l14:
				depth--
				add(ruleAssembly, position9)
			}
			return true
		l8:
			position, tokenIndex, depth = position8, tokenIndex8, depth8
			return false
		},
		/* 2 Comment <- <(SEMICOLON <(!NL .)*> Action3)> */
		func() bool {
			position15, tokenIndex15, depth15 := position, tokenIndex, depth
			{
				position16 := position
				depth++
				if !_rules[ruleSEMICOLON]() {
					goto l15
				}
				{
					position17 := position
					depth++
				l18:
					{
						position19, tokenIndex19, depth19 := position, tokenIndex, depth
						{
							position20, tokenIndex20, depth20 := position, tokenIndex, depth
							if !_rules[ruleNL]() {
								goto l20
							}
							goto l19
						l20:
							position, tokenIndex, depth = position20, tokenIndex20, depth20
						}
						if !matchDot() {
							goto l19
						}
						goto l18
					l19:
						position, tokenIndex, depth = position19, tokenIndex19, depth19
					}
					depth--
					add(rulePegText, position17)
				}
				if !_rules[ruleAction3]() {
					goto l15
				}
				depth--
				add(ruleComment, position16)
			}
			return true
		l15:
			position, tokenIndex, depth = position15, tokenIndex15, depth15
			return false
		},
		/* 3 Label <- <(Action4 Identifier Spacing COLON)> */
		func() bool {
			position21, tokenIndex21, depth21 := position, tokenIndex, depth
			{
				position22 := position
				depth++
				if !_rules[ruleAction4]() {
					goto l21
				}
				if !_rules[ruleIdentifier]() {
					goto l21
				}
				if !_rules[ruleSpacing]() {
					goto l21
				}
				if !_rules[ruleCOLON]() {
					goto l21
				}
				depth--
				add(ruleLabel, position22)
			}
			return true
		l21:
			position, tokenIndex, depth = position21, tokenIndex21, depth21
			return false
		},
		/* 4 Inst <- <(((PUSH / ((&('J' | 'j') JMP) | (&('P' | 'p') POP) | (&('C' | 'c') CALL))) Operand) / (CAL <DATA_TYPE> Action5 <CAL_OP> Action6 Operand COMMA Operand) / ((&('J' | 'j') (JPC <CMP_OP> Action9 Operand)) | (&('C' | 'c') (CMP <DATA_TYPE> Action8 Operand COMMA Operand)) | (&('L' | 'l') (LD <DATA_TYPE> Action7 Operand COMMA Operand)) | (&('N' | 'n') NOP) | (&('R' | 'r') RET) | (&('E' | 'e') EXIT) | (&('I' | 'O' | 'i' | 'o') ((IN / OUT) Operand COMMA Operand))))> */
		func() bool {
			position23, tokenIndex23, depth23 := position, tokenIndex, depth
			{
				position24 := position
				depth++
				{
					position25, tokenIndex25, depth25 := position, tokenIndex, depth
					{
						position27, tokenIndex27, depth27 := position, tokenIndex, depth
						if !_rules[rulePUSH]() {
							goto l28
						}
						goto l27
					l28:
						position, tokenIndex, depth = position27, tokenIndex27, depth27
						{
							switch buffer[position] {
							case 'J', 'j':
								if !_rules[ruleJMP]() {
									goto l26
								}
								break
							case 'P', 'p':
								if !_rules[rulePOP]() {
									goto l26
								}
								break
							default:
								if !_rules[ruleCALL]() {
									goto l26
								}
								break
							}
						}

					}
				l27:
					if !_rules[ruleOperand]() {
						goto l26
					}
					goto l25
				l26:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					if !_rules[ruleCAL]() {
						goto l30
					}
					{
						position31 := position
						depth++
						if !_rules[ruleDATA_TYPE]() {
							goto l30
						}
						depth--
						add(rulePegText, position31)
					}
					if !_rules[ruleAction5]() {
						goto l30
					}
					{
						position32 := position
						depth++
						if !_rules[ruleCAL_OP]() {
							goto l30
						}
						depth--
						add(rulePegText, position32)
					}
					if !_rules[ruleAction6]() {
						goto l30
					}
					if !_rules[ruleOperand]() {
						goto l30
					}
					if !_rules[ruleCOMMA]() {
						goto l30
					}
					if !_rules[ruleOperand]() {
						goto l30
					}
					goto l25
				l30:
					position, tokenIndex, depth = position25, tokenIndex25, depth25
					{
						switch buffer[position] {
						case 'J', 'j':
							if !_rules[ruleJPC]() {
								goto l23
							}
							{
								position34 := position
								depth++
								if !_rules[ruleCMP_OP]() {
									goto l23
								}
								depth--
								add(rulePegText, position34)
							}
							if !_rules[ruleAction9]() {
								goto l23
							}
							if !_rules[ruleOperand]() {
								goto l23
							}
							break
						case 'C', 'c':
							if !_rules[ruleCMP]() {
								goto l23
							}
							{
								position35 := position
								depth++
								if !_rules[ruleDATA_TYPE]() {
									goto l23
								}
								depth--
								add(rulePegText, position35)
							}
							if !_rules[ruleAction8]() {
								goto l23
							}
							if !_rules[ruleOperand]() {
								goto l23
							}
							if !_rules[ruleCOMMA]() {
								goto l23
							}
							if !_rules[ruleOperand]() {
								goto l23
							}
							break
						case 'L', 'l':
							if !_rules[ruleLD]() {
								goto l23
							}
							{
								position36 := position
								depth++
								if !_rules[ruleDATA_TYPE]() {
									goto l23
								}
								depth--
								add(rulePegText, position36)
							}
							if !_rules[ruleAction7]() {
								goto l23
							}
							if !_rules[ruleOperand]() {
								goto l23
							}
							if !_rules[ruleCOMMA]() {
								goto l23
							}
							if !_rules[ruleOperand]() {
								goto l23
							}
							break
						case 'N', 'n':
							if !_rules[ruleNOP]() {
								goto l23
							}
							break
						case 'R', 'r':
							if !_rules[ruleRET]() {
								goto l23
							}
							break
						case 'E', 'e':
							if !_rules[ruleEXIT]() {
								goto l23
							}
							break
						default:
							{
								position37, tokenIndex37, depth37 := position, tokenIndex, depth
								if !_rules[ruleIN]() {
									goto l38
								}
								goto l37
							l38:
								position, tokenIndex, depth = position37, tokenIndex37, depth37
								if !_rules[ruleOUT]() {
									goto l23
								}
							}
						l37:
							if !_rules[ruleOperand]() {
								goto l23
							}
							if !_rules[ruleCOMMA]() {
								goto l23
							}
							if !_rules[ruleOperand]() {
								goto l23
							}
							break
						}
					}

				}
			l25:
				depth--
				add(ruleInst, position24)
			}
			return true
		l23:
			position, tokenIndex, depth = position23, tokenIndex23, depth23
			return false
		},
		/* 5 Pseudo <- <((BLOCK IntegerLiteral IntegerLiteral) / (DATA Identifier PSEUDO_DATA_TYPE? PseudoDataValue (COMMA PseudoDataValue)*))> */
		func() bool {
			position39, tokenIndex39, depth39 := position, tokenIndex, depth
			{
				position40 := position
				depth++
				{
					position41, tokenIndex41, depth41 := position, tokenIndex, depth
					if !_rules[ruleBLOCK]() {
						goto l42
					}
					if !_rules[ruleIntegerLiteral]() {
						goto l42
					}
					if !_rules[ruleIntegerLiteral]() {
						goto l42
					}
					goto l41
				l42:
					position, tokenIndex, depth = position41, tokenIndex41, depth41
					if !_rules[ruleDATA]() {
						goto l39
					}
					if !_rules[ruleIdentifier]() {
						goto l39
					}
					{
						position43, tokenIndex43, depth43 := position, tokenIndex, depth
						if !_rules[rulePSEUDO_DATA_TYPE]() {
							goto l43
						}
						goto l44
					l43:
						position, tokenIndex, depth = position43, tokenIndex43, depth43
					}
				l44:
					if !_rules[rulePseudoDataValue]() {
						goto l39
					}
				l45:
					{
						position46, tokenIndex46, depth46 := position, tokenIndex, depth
						if !_rules[ruleCOMMA]() {
							goto l46
						}
						if !_rules[rulePseudoDataValue]() {
							goto l46
						}
						goto l45
					l46:
						position, tokenIndex, depth = position46, tokenIndex46, depth46
					}
				}
			l41:
				depth--
				add(rulePseudo, position40)
			}
			return true
		l39:
			position, tokenIndex, depth = position39, tokenIndex39, depth39
			return false
		},
		/* 6 PseudoDataValue <- <((&('%') (<('%' HexDigits '%')> Spacing Action13)) | (&('"') (StringLiteral Action12)) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (IntegerLiteral Action10)) | (&('$' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (Identifier Action11)))> */
		func() bool {
			position47, tokenIndex47, depth47 := position, tokenIndex, depth
			{
				position48 := position
				depth++
				{
					switch buffer[position] {
					case '%':
						{
							position50 := position
							depth++
							if buffer[position] != rune('%') {
								goto l47
							}
							position++
							if !_rules[ruleHexDigits]() {
								goto l47
							}
							if buffer[position] != rune('%') {
								goto l47
							}
							position++
							depth--
							add(rulePegText, position50)
						}
						if !_rules[ruleSpacing]() {
							goto l47
						}
						if !_rules[ruleAction13]() {
							goto l47
						}
						break
					case '"':
						if !_rules[ruleStringLiteral]() {
							goto l47
						}
						if !_rules[ruleAction12]() {
							goto l47
						}
						break
					case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if !_rules[ruleIntegerLiteral]() {
							goto l47
						}
						if !_rules[ruleAction10]() {
							goto l47
						}
						break
					default:
						if !_rules[ruleIdentifier]() {
							goto l47
						}
						if !_rules[ruleAction11]() {
							goto l47
						}
						break
					}
				}

				depth--
				add(rulePseudoDataValue, position48)
			}
			return true
		l47:
			position, tokenIndex, depth = position47, tokenIndex47, depth47
			return false
		},
		/* 7 PSEUDO_DATA_TYPE <- <(DATA_TYPE / (((('c' / 'C') ('h' / 'H') ('a' / 'A') ('r' / 'R')) / (('b' / 'B') ('i' / 'I') ('n' / 'N'))) Space))> */
		func() bool {
			position51, tokenIndex51, depth51 := position, tokenIndex, depth
			{
				position52 := position
				depth++
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if !_rules[ruleDATA_TYPE]() {
						goto l54
					}
					goto l53
				l54:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
					{
						position55, tokenIndex55, depth55 := position, tokenIndex, depth
						{
							position57, tokenIndex57, depth57 := position, tokenIndex, depth
							if buffer[position] != rune('c') {
								goto l58
							}
							position++
							goto l57
						l58:
							position, tokenIndex, depth = position57, tokenIndex57, depth57
							if buffer[position] != rune('C') {
								goto l56
							}
							position++
						}
					l57:
						{
							position59, tokenIndex59, depth59 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l60
							}
							position++
							goto l59
						l60:
							position, tokenIndex, depth = position59, tokenIndex59, depth59
							if buffer[position] != rune('H') {
								goto l56
							}
							position++
						}
					l59:
						{
							position61, tokenIndex61, depth61 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l62
							}
							position++
							goto l61
						l62:
							position, tokenIndex, depth = position61, tokenIndex61, depth61
							if buffer[position] != rune('A') {
								goto l56
							}
							position++
						}
					l61:
						{
							position63, tokenIndex63, depth63 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l64
							}
							position++
							goto l63
						l64:
							position, tokenIndex, depth = position63, tokenIndex63, depth63
							if buffer[position] != rune('R') {
								goto l56
							}
							position++
						}
					l63:
						goto l55
					l56:
						position, tokenIndex, depth = position55, tokenIndex55, depth55
						{
							position65, tokenIndex65, depth65 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l66
							}
							position++
							goto l65
						l66:
							position, tokenIndex, depth = position65, tokenIndex65, depth65
							if buffer[position] != rune('B') {
								goto l51
							}
							position++
						}
					l65:
						{
							position67, tokenIndex67, depth67 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l68
							}
							position++
							goto l67
						l68:
							position, tokenIndex, depth = position67, tokenIndex67, depth67
							if buffer[position] != rune('I') {
								goto l51
							}
							position++
						}
					l67:
						{
							position69, tokenIndex69, depth69 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l70
							}
							position++
							goto l69
						l70:
							position, tokenIndex, depth = position69, tokenIndex69, depth69
							if buffer[position] != rune('N') {
								goto l51
							}
							position++
						}
					l69:
					}
				l55:
					if !_rules[ruleSpace]() {
						goto l51
					}
				}
			l53:
				depth--
				add(rulePSEUDO_DATA_TYPE, position52)
			}
			return true
		l51:
			position, tokenIndex, depth = position51, tokenIndex51, depth51
			return false
		},
		/* 8 Operand <- <(((LBRK Identifier RBRK Action15) / ((&('[') (LBRK IntegerLiteral RBRK Action17)) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (IntegerLiteral Action16)) | (&('$' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (Identifier Action14)))) Spacing)> */
		func() bool {
			position71, tokenIndex71, depth71 := position, tokenIndex, depth
			{
				position72 := position
				depth++
				{
					position73, tokenIndex73, depth73 := position, tokenIndex, depth
					if !_rules[ruleLBRK]() {
						goto l74
					}
					if !_rules[ruleIdentifier]() {
						goto l74
					}
					if !_rules[ruleRBRK]() {
						goto l74
					}
					if !_rules[ruleAction15]() {
						goto l74
					}
					goto l73
				l74:
					position, tokenIndex, depth = position73, tokenIndex73, depth73
					{
						switch buffer[position] {
						case '[':
							if !_rules[ruleLBRK]() {
								goto l71
							}
							if !_rules[ruleIntegerLiteral]() {
								goto l71
							}
							if !_rules[ruleRBRK]() {
								goto l71
							}
							if !_rules[ruleAction17]() {
								goto l71
							}
							break
						case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if !_rules[ruleIntegerLiteral]() {
								goto l71
							}
							if !_rules[ruleAction16]() {
								goto l71
							}
							break
						default:
							if !_rules[ruleIdentifier]() {
								goto l71
							}
							if !_rules[ruleAction14]() {
								goto l71
							}
							break
						}
					}

				}
			l73:
				if !_rules[ruleSpacing]() {
					goto l71
				}
				depth--
				add(ruleOperand, position72)
			}
			return true
		l71:
			position, tokenIndex, depth = position71, tokenIndex71, depth71
			return false
		},
		/* 9 Spacing <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position77 := position
				depth++
			l78:
				{
					position79, tokenIndex79, depth79 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l79
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l79
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l79
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l79
							}
							position++
							break
						}
					}

					goto l78
				l79:
					position, tokenIndex, depth = position79, tokenIndex79, depth79
				}
				depth--
				add(ruleSpacing, position77)
			}
			return true
		},
		/* 10 Space <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))+> */
		func() bool {
			position81, tokenIndex81, depth81 := position, tokenIndex, depth
			{
				position82 := position
				depth++
				{
					switch buffer[position] {
					case '\f':
						if buffer[position] != rune('\f') {
							goto l81
						}
						position++
						break
					case '\r':
						if buffer[position] != rune('\r') {
							goto l81
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l81
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l81
						}
						position++
						break
					}
				}

			l83:
				{
					position84, tokenIndex84, depth84 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l84
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l84
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l84
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l84
							}
							position++
							break
						}
					}

					goto l83
				l84:
					position, tokenIndex, depth = position84, tokenIndex84, depth84
				}
				depth--
				add(ruleSpace, position82)
			}
			return true
		l81:
			position, tokenIndex, depth = position81, tokenIndex81, depth81
			return false
		},
		/* 11 Identifier <- <(<(Letter LetterOrDigit*)> Spacing Action18)> */
		func() bool {
			position87, tokenIndex87, depth87 := position, tokenIndex, depth
			{
				position88 := position
				depth++
				{
					position89 := position
					depth++
					if !_rules[ruleLetter]() {
						goto l87
					}
				l90:
					{
						position91, tokenIndex91, depth91 := position, tokenIndex, depth
						if !_rules[ruleLetterOrDigit]() {
							goto l91
						}
						goto l90
					l91:
						position, tokenIndex, depth = position91, tokenIndex91, depth91
					}
					depth--
					add(rulePegText, position89)
				}
				if !_rules[ruleSpacing]() {
					goto l87
				}
				if !_rules[ruleAction18]() {
					goto l87
				}
				depth--
				add(ruleIdentifier, position88)
			}
			return true
		l87:
			position, tokenIndex, depth = position87, tokenIndex87, depth87
			return false
		},
		/* 12 Letter <- <((&('$' | '_') ('_' / '$')) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position92, tokenIndex92, depth92 := position, tokenIndex, depth
			{
				position93 := position
				depth++
				{
					switch buffer[position] {
					case '$', '_':
						{
							position95, tokenIndex95, depth95 := position, tokenIndex, depth
							if buffer[position] != rune('_') {
								goto l96
							}
							position++
							goto l95
						l96:
							position, tokenIndex, depth = position95, tokenIndex95, depth95
							if buffer[position] != rune('$') {
								goto l92
							}
							position++
						}
					l95:
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l92
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l92
						}
						position++
						break
					}
				}

				depth--
				add(ruleLetter, position93)
			}
			return true
		l92:
			position, tokenIndex, depth = position92, tokenIndex92, depth92
			return false
		},
		/* 13 LetterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position97, tokenIndex97, depth97 := position, tokenIndex, depth
			{
				position98 := position
				depth++
				{
					switch buffer[position] {
					case '$', '_':
						{
							position100, tokenIndex100, depth100 := position, tokenIndex, depth
							if buffer[position] != rune('_') {
								goto l101
							}
							position++
							goto l100
						l101:
							position, tokenIndex, depth = position100, tokenIndex100, depth100
							if buffer[position] != rune('$') {
								goto l97
							}
							position++
						}
					l100:
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l97
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l97
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l97
						}
						position++
						break
					}
				}

				depth--
				add(ruleLetterOrDigit, position98)
			}
			return true
		l97:
			position, tokenIndex, depth = position97, tokenIndex97, depth97
			return false
		},
		/* 14 EXIT <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('t' / 'T') Spacing Action19)> */
		func() bool {
			position102, tokenIndex102, depth102 := position, tokenIndex, depth
			{
				position103 := position
				depth++
				{
					position104, tokenIndex104, depth104 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l105
					}
					position++
					goto l104
				l105:
					position, tokenIndex, depth = position104, tokenIndex104, depth104
					if buffer[position] != rune('E') {
						goto l102
					}
					position++
				}
			l104:
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if buffer[position] != rune('x') {
						goto l107
					}
					position++
					goto l106
				l107:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
					if buffer[position] != rune('X') {
						goto l102
					}
					position++
				}
			l106:
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l109
					}
					position++
					goto l108
				l109:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
					if buffer[position] != rune('I') {
						goto l102
					}
					position++
				}
			l108:
				{
					position110, tokenIndex110, depth110 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l111
					}
					position++
					goto l110
				l111:
					position, tokenIndex, depth = position110, tokenIndex110, depth110
					if buffer[position] != rune('T') {
						goto l102
					}
					position++
				}
			l110:
				if !_rules[ruleSpacing]() {
					goto l102
				}
				if !_rules[ruleAction19]() {
					goto l102
				}
				depth--
				add(ruleEXIT, position103)
			}
			return true
		l102:
			position, tokenIndex, depth = position102, tokenIndex102, depth102
			return false
		},
		/* 15 RET <- <(('r' / 'R') ('e' / 'E') ('t' / 'T') Spacing Action20)> */
		func() bool {
			position112, tokenIndex112, depth112 := position, tokenIndex, depth
			{
				position113 := position
				depth++
				{
					position114, tokenIndex114, depth114 := position, tokenIndex, depth
					if buffer[position] != rune('r') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex, depth = position114, tokenIndex114, depth114
					if buffer[position] != rune('R') {
						goto l112
					}
					position++
				}
			l114:
				{
					position116, tokenIndex116, depth116 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l117
					}
					position++
					goto l116
				l117:
					position, tokenIndex, depth = position116, tokenIndex116, depth116
					if buffer[position] != rune('E') {
						goto l112
					}
					position++
				}
			l116:
				{
					position118, tokenIndex118, depth118 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l119
					}
					position++
					goto l118
				l119:
					position, tokenIndex, depth = position118, tokenIndex118, depth118
					if buffer[position] != rune('T') {
						goto l112
					}
					position++
				}
			l118:
				if !_rules[ruleSpacing]() {
					goto l112
				}
				if !_rules[ruleAction20]() {
					goto l112
				}
				depth--
				add(ruleRET, position113)
			}
			return true
		l112:
			position, tokenIndex, depth = position112, tokenIndex112, depth112
			return false
		},
		/* 16 NOP <- <(('n' / 'N') ('o' / 'O') ('p' / 'P') Spacing Action21)> */
		func() bool {
			position120, tokenIndex120, depth120 := position, tokenIndex, depth
			{
				position121 := position
				depth++
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l123
					}
					position++
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if buffer[position] != rune('N') {
						goto l120
					}
					position++
				}
			l122:
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l125
					}
					position++
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if buffer[position] != rune('O') {
						goto l120
					}
					position++
				}
			l124:
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l127
					}
					position++
					goto l126
				l127:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
					if buffer[position] != rune('P') {
						goto l120
					}
					position++
				}
			l126:
				if !_rules[ruleSpacing]() {
					goto l120
				}
				if !_rules[ruleAction21]() {
					goto l120
				}
				depth--
				add(ruleNOP, position121)
			}
			return true
		l120:
			position, tokenIndex, depth = position120, tokenIndex120, depth120
			return false
		},
		/* 17 CALL <- <(('c' / 'C') ('a' / 'A') ('l' / 'L') ('l' / 'L') Space Action22)> */
		func() bool {
			position128, tokenIndex128, depth128 := position, tokenIndex, depth
			{
				position129 := position
				depth++
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l131
					}
					position++
					goto l130
				l131:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
					if buffer[position] != rune('C') {
						goto l128
					}
					position++
				}
			l130:
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l133
					}
					position++
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if buffer[position] != rune('A') {
						goto l128
					}
					position++
				}
			l132:
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l135
					}
					position++
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					if buffer[position] != rune('L') {
						goto l128
					}
					position++
				}
			l134:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l137
					}
					position++
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					if buffer[position] != rune('L') {
						goto l128
					}
					position++
				}
			l136:
				if !_rules[ruleSpace]() {
					goto l128
				}
				if !_rules[ruleAction22]() {
					goto l128
				}
				depth--
				add(ruleCALL, position129)
			}
			return true
		l128:
			position, tokenIndex, depth = position128, tokenIndex128, depth128
			return false
		},
		/* 18 PUSH <- <(('p' / 'P') ('u' / 'U') ('s' / 'S') ('h' / 'H') Space Action23)> */
		func() bool {
			position138, tokenIndex138, depth138 := position, tokenIndex, depth
			{
				position139 := position
				depth++
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l141
					}
					position++
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					if buffer[position] != rune('P') {
						goto l138
					}
					position++
				}
			l140:
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l143
					}
					position++
					goto l142
				l143:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
					if buffer[position] != rune('U') {
						goto l138
					}
					position++
				}
			l142:
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l145
					}
					position++
					goto l144
				l145:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
					if buffer[position] != rune('S') {
						goto l138
					}
					position++
				}
			l144:
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if buffer[position] != rune('h') {
						goto l147
					}
					position++
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('H') {
						goto l138
					}
					position++
				}
			l146:
				if !_rules[ruleSpace]() {
					goto l138
				}
				if !_rules[ruleAction23]() {
					goto l138
				}
				depth--
				add(rulePUSH, position139)
			}
			return true
		l138:
			position, tokenIndex, depth = position138, tokenIndex138, depth138
			return false
		},
		/* 19 POP <- <(('p' / 'P') ('o' / 'O') ('p' / 'P') Space Action24)> */
		func() bool {
			position148, tokenIndex148, depth148 := position, tokenIndex, depth
			{
				position149 := position
				depth++
				{
					position150, tokenIndex150, depth150 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l151
					}
					position++
					goto l150
				l151:
					position, tokenIndex, depth = position150, tokenIndex150, depth150
					if buffer[position] != rune('P') {
						goto l148
					}
					position++
				}
			l150:
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l153
					}
					position++
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if buffer[position] != rune('O') {
						goto l148
					}
					position++
				}
			l152:
				{
					position154, tokenIndex154, depth154 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l155
					}
					position++
					goto l154
				l155:
					position, tokenIndex, depth = position154, tokenIndex154, depth154
					if buffer[position] != rune('P') {
						goto l148
					}
					position++
				}
			l154:
				if !_rules[ruleSpace]() {
					goto l148
				}
				if !_rules[ruleAction24]() {
					goto l148
				}
				depth--
				add(rulePOP, position149)
			}
			return true
		l148:
			position, tokenIndex, depth = position148, tokenIndex148, depth148
			return false
		},
		/* 20 JMP <- <(('j' / 'J') ('m' / 'M') ('p' / 'P') Space Action25)> */
		func() bool {
			position156, tokenIndex156, depth156 := position, tokenIndex, depth
			{
				position157 := position
				depth++
				{
					position158, tokenIndex158, depth158 := position, tokenIndex, depth
					if buffer[position] != rune('j') {
						goto l159
					}
					position++
					goto l158
				l159:
					position, tokenIndex, depth = position158, tokenIndex158, depth158
					if buffer[position] != rune('J') {
						goto l156
					}
					position++
				}
			l158:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l161
					}
					position++
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					if buffer[position] != rune('M') {
						goto l156
					}
					position++
				}
			l160:
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l163
					}
					position++
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if buffer[position] != rune('P') {
						goto l156
					}
					position++
				}
			l162:
				if !_rules[ruleSpace]() {
					goto l156
				}
				if !_rules[ruleAction25]() {
					goto l156
				}
				depth--
				add(ruleJMP, position157)
			}
			return true
		l156:
			position, tokenIndex, depth = position156, tokenIndex156, depth156
			return false
		},
		/* 21 IN <- <(('i' / 'I') ('n' / 'N') Space Action26)> */
		func() bool {
			position164, tokenIndex164, depth164 := position, tokenIndex, depth
			{
				position165 := position
				depth++
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l167
					}
					position++
					goto l166
				l167:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
					if buffer[position] != rune('I') {
						goto l164
					}
					position++
				}
			l166:
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l169
					}
					position++
					goto l168
				l169:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
					if buffer[position] != rune('N') {
						goto l164
					}
					position++
				}
			l168:
				if !_rules[ruleSpace]() {
					goto l164
				}
				if !_rules[ruleAction26]() {
					goto l164
				}
				depth--
				add(ruleIN, position165)
			}
			return true
		l164:
			position, tokenIndex, depth = position164, tokenIndex164, depth164
			return false
		},
		/* 22 OUT <- <(('o' / 'O') ('u' / 'U') ('t' / 'T') Space Action27)> */
		func() bool {
			position170, tokenIndex170, depth170 := position, tokenIndex, depth
			{
				position171 := position
				depth++
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l173
					}
					position++
					goto l172
				l173:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
					if buffer[position] != rune('O') {
						goto l170
					}
					position++
				}
			l172:
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l175
					}
					position++
					goto l174
				l175:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
					if buffer[position] != rune('U') {
						goto l170
					}
					position++
				}
			l174:
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l177
					}
					position++
					goto l176
				l177:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
					if buffer[position] != rune('T') {
						goto l170
					}
					position++
				}
			l176:
				if !_rules[ruleSpace]() {
					goto l170
				}
				if !_rules[ruleAction27]() {
					goto l170
				}
				depth--
				add(ruleOUT, position171)
			}
			return true
		l170:
			position, tokenIndex, depth = position170, tokenIndex170, depth170
			return false
		},
		/* 23 CAL <- <(('c' / 'C') ('a' / 'A') ('l' / 'L') Space Action28)> */
		func() bool {
			position178, tokenIndex178, depth178 := position, tokenIndex, depth
			{
				position179 := position
				depth++
				{
					position180, tokenIndex180, depth180 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l181
					}
					position++
					goto l180
				l181:
					position, tokenIndex, depth = position180, tokenIndex180, depth180
					if buffer[position] != rune('C') {
						goto l178
					}
					position++
				}
			l180:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l183
					}
					position++
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if buffer[position] != rune('A') {
						goto l178
					}
					position++
				}
			l182:
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l185
					}
					position++
					goto l184
				l185:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
					if buffer[position] != rune('L') {
						goto l178
					}
					position++
				}
			l184:
				if !_rules[ruleSpace]() {
					goto l178
				}
				if !_rules[ruleAction28]() {
					goto l178
				}
				depth--
				add(ruleCAL, position179)
			}
			return true
		l178:
			position, tokenIndex, depth = position178, tokenIndex178, depth178
			return false
		},
		/* 24 LD <- <(('l' / 'L') ('d' / 'D') Space Action29)> */
		func() bool {
			position186, tokenIndex186, depth186 := position, tokenIndex, depth
			{
				position187 := position
				depth++
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l189
					}
					position++
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if buffer[position] != rune('L') {
						goto l186
					}
					position++
				}
			l188:
				{
					position190, tokenIndex190, depth190 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l191
					}
					position++
					goto l190
				l191:
					position, tokenIndex, depth = position190, tokenIndex190, depth190
					if buffer[position] != rune('D') {
						goto l186
					}
					position++
				}
			l190:
				if !_rules[ruleSpace]() {
					goto l186
				}
				if !_rules[ruleAction29]() {
					goto l186
				}
				depth--
				add(ruleLD, position187)
			}
			return true
		l186:
			position, tokenIndex, depth = position186, tokenIndex186, depth186
			return false
		},
		/* 25 CMP <- <(('c' / 'C') ('m' / 'M') ('p' / 'P') Space Action30)> */
		func() bool {
			position192, tokenIndex192, depth192 := position, tokenIndex, depth
			{
				position193 := position
				depth++
				{
					position194, tokenIndex194, depth194 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l195
					}
					position++
					goto l194
				l195:
					position, tokenIndex, depth = position194, tokenIndex194, depth194
					if buffer[position] != rune('C') {
						goto l192
					}
					position++
				}
			l194:
				{
					position196, tokenIndex196, depth196 := position, tokenIndex, depth
					if buffer[position] != rune('m') {
						goto l197
					}
					position++
					goto l196
				l197:
					position, tokenIndex, depth = position196, tokenIndex196, depth196
					if buffer[position] != rune('M') {
						goto l192
					}
					position++
				}
			l196:
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l199
					}
					position++
					goto l198
				l199:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
					if buffer[position] != rune('P') {
						goto l192
					}
					position++
				}
			l198:
				if !_rules[ruleSpace]() {
					goto l192
				}
				if !_rules[ruleAction30]() {
					goto l192
				}
				depth--
				add(ruleCMP, position193)
			}
			return true
		l192:
			position, tokenIndex, depth = position192, tokenIndex192, depth192
			return false
		},
		/* 26 JPC <- <(('j' / 'J') ('p' / 'P') ('c' / 'C') Space Action31)> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if buffer[position] != rune('j') {
						goto l203
					}
					position++
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if buffer[position] != rune('J') {
						goto l200
					}
					position++
				}
			l202:
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l205
					}
					position++
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if buffer[position] != rune('P') {
						goto l200
					}
					position++
				}
			l204:
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l207
					}
					position++
					goto l206
				l207:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
					if buffer[position] != rune('C') {
						goto l200
					}
					position++
				}
			l206:
				if !_rules[ruleSpace]() {
					goto l200
				}
				if !_rules[ruleAction31]() {
					goto l200
				}
				depth--
				add(ruleJPC, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 27 BLOCK <- <('.' ('b' / 'B') ('l' / 'L') ('o' / 'O') ('c' / 'C') ('k' / 'K') Space Action32)> */
		func() bool {
			position208, tokenIndex208, depth208 := position, tokenIndex, depth
			{
				position209 := position
				depth++
				if buffer[position] != rune('.') {
					goto l208
				}
				position++
				{
					position210, tokenIndex210, depth210 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l211
					}
					position++
					goto l210
				l211:
					position, tokenIndex, depth = position210, tokenIndex210, depth210
					if buffer[position] != rune('B') {
						goto l208
					}
					position++
				}
			l210:
				{
					position212, tokenIndex212, depth212 := position, tokenIndex, depth
					if buffer[position] != rune('l') {
						goto l213
					}
					position++
					goto l212
				l213:
					position, tokenIndex, depth = position212, tokenIndex212, depth212
					if buffer[position] != rune('L') {
						goto l208
					}
					position++
				}
			l212:
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if buffer[position] != rune('o') {
						goto l215
					}
					position++
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if buffer[position] != rune('O') {
						goto l208
					}
					position++
				}
			l214:
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l217
					}
					position++
					goto l216
				l217:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
					if buffer[position] != rune('C') {
						goto l208
					}
					position++
				}
			l216:
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if buffer[position] != rune('k') {
						goto l219
					}
					position++
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if buffer[position] != rune('K') {
						goto l208
					}
					position++
				}
			l218:
				if !_rules[ruleSpace]() {
					goto l208
				}
				if !_rules[ruleAction32]() {
					goto l208
				}
				depth--
				add(ruleBLOCK, position209)
			}
			return true
		l208:
			position, tokenIndex, depth = position208, tokenIndex208, depth208
			return false
		},
		/* 28 DATA <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') Space Action33)> */
		func() bool {
			position220, tokenIndex220, depth220 := position, tokenIndex, depth
			{
				position221 := position
				depth++
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l223
					}
					position++
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if buffer[position] != rune('D') {
						goto l220
					}
					position++
				}
			l222:
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l225
					}
					position++
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if buffer[position] != rune('A') {
						goto l220
					}
					position++
				}
			l224:
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l227
					}
					position++
					goto l226
				l227:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
					if buffer[position] != rune('T') {
						goto l220
					}
					position++
				}
			l226:
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l229
					}
					position++
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if buffer[position] != rune('A') {
						goto l220
					}
					position++
				}
			l228:
				if !_rules[ruleSpace]() {
					goto l220
				}
				if !_rules[ruleAction33]() {
					goto l220
				}
				depth--
				add(ruleDATA, position221)
			}
			return true
		l220:
			position, tokenIndex, depth = position220, tokenIndex220, depth220
			return false
		},
		/* 29 CAL_OP <- <(((('m' / 'M') ('u' / 'U') ('l' / 'L')) / ((&('M' | 'm') (('m' / 'M') ('o' / 'O') ('d' / 'D'))) | (&('D' | 'd') (('d' / 'D') ('i' / 'I') ('v' / 'V'))) | (&('S' | 's') (('s' / 'S') ('u' / 'U') ('b' / 'B'))) | (&('A' | 'a') (('a' / 'A') ('d' / 'D') ('d' / 'D'))))) Space)> */
		func() bool {
			position230, tokenIndex230, depth230 := position, tokenIndex, depth
			{
				position231 := position
				depth++
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					{
						position234, tokenIndex234, depth234 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l235
						}
						position++
						goto l234
					l235:
						position, tokenIndex, depth = position234, tokenIndex234, depth234
						if buffer[position] != rune('M') {
							goto l233
						}
						position++
					}
				l234:
					{
						position236, tokenIndex236, depth236 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l237
						}
						position++
						goto l236
					l237:
						position, tokenIndex, depth = position236, tokenIndex236, depth236
						if buffer[position] != rune('U') {
							goto l233
						}
						position++
					}
				l236:
					{
						position238, tokenIndex238, depth238 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l239
						}
						position++
						goto l238
					l239:
						position, tokenIndex, depth = position238, tokenIndex238, depth238
						if buffer[position] != rune('L') {
							goto l233
						}
						position++
					}
				l238:
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					{
						switch buffer[position] {
						case 'M', 'm':
							{
								position241, tokenIndex241, depth241 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l242
								}
								position++
								goto l241
							l242:
								position, tokenIndex, depth = position241, tokenIndex241, depth241
								if buffer[position] != rune('M') {
									goto l230
								}
								position++
							}
						l241:
							{
								position243, tokenIndex243, depth243 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l244
								}
								position++
								goto l243
							l244:
								position, tokenIndex, depth = position243, tokenIndex243, depth243
								if buffer[position] != rune('O') {
									goto l230
								}
								position++
							}
						l243:
							{
								position245, tokenIndex245, depth245 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l246
								}
								position++
								goto l245
							l246:
								position, tokenIndex, depth = position245, tokenIndex245, depth245
								if buffer[position] != rune('D') {
									goto l230
								}
								position++
							}
						l245:
							break
						case 'D', 'd':
							{
								position247, tokenIndex247, depth247 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l248
								}
								position++
								goto l247
							l248:
								position, tokenIndex, depth = position247, tokenIndex247, depth247
								if buffer[position] != rune('D') {
									goto l230
								}
								position++
							}
						l247:
							{
								position249, tokenIndex249, depth249 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l250
								}
								position++
								goto l249
							l250:
								position, tokenIndex, depth = position249, tokenIndex249, depth249
								if buffer[position] != rune('I') {
									goto l230
								}
								position++
							}
						l249:
							{
								position251, tokenIndex251, depth251 := position, tokenIndex, depth
								if buffer[position] != rune('v') {
									goto l252
								}
								position++
								goto l251
							l252:
								position, tokenIndex, depth = position251, tokenIndex251, depth251
								if buffer[position] != rune('V') {
									goto l230
								}
								position++
							}
						l251:
							break
						case 'S', 's':
							{
								position253, tokenIndex253, depth253 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l254
								}
								position++
								goto l253
							l254:
								position, tokenIndex, depth = position253, tokenIndex253, depth253
								if buffer[position] != rune('S') {
									goto l230
								}
								position++
							}
						l253:
							{
								position255, tokenIndex255, depth255 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l256
								}
								position++
								goto l255
							l256:
								position, tokenIndex, depth = position255, tokenIndex255, depth255
								if buffer[position] != rune('U') {
									goto l230
								}
								position++
							}
						l255:
							{
								position257, tokenIndex257, depth257 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l258
								}
								position++
								goto l257
							l258:
								position, tokenIndex, depth = position257, tokenIndex257, depth257
								if buffer[position] != rune('B') {
									goto l230
								}
								position++
							}
						l257:
							break
						default:
							{
								position259, tokenIndex259, depth259 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l260
								}
								position++
								goto l259
							l260:
								position, tokenIndex, depth = position259, tokenIndex259, depth259
								if buffer[position] != rune('A') {
									goto l230
								}
								position++
							}
						l259:
							{
								position261, tokenIndex261, depth261 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l262
								}
								position++
								goto l261
							l262:
								position, tokenIndex, depth = position261, tokenIndex261, depth261
								if buffer[position] != rune('D') {
									goto l230
								}
								position++
							}
						l261:
							{
								position263, tokenIndex263, depth263 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l264
								}
								position++
								goto l263
							l264:
								position, tokenIndex, depth = position263, tokenIndex263, depth263
								if buffer[position] != rune('D') {
									goto l230
								}
								position++
							}
						l263:
							break
						}
					}

				}
			l232:
				if !_rules[ruleSpace]() {
					goto l230
				}
				depth--
				add(ruleCAL_OP, position231)
			}
			return true
		l230:
			position, tokenIndex, depth = position230, tokenIndex230, depth230
			return false
		},
		/* 30 CMP_OP <- <(((('b' / 'B') ('e' / 'E')) / (('a' / 'A') ('e' / 'E')) / ((&('N' | 'n') (('n' / 'N') ('z' / 'Z'))) | (&('A' | 'a') ('a' / 'A')) | (&('Z') 'Z') | (&('z') 'z') | (&('B' | 'b') ('b' / 'B')))) Space)> */
		func() bool {
			position265, tokenIndex265, depth265 := position, tokenIndex, depth
			{
				position266 := position
				depth++
				{
					position267, tokenIndex267, depth267 := position, tokenIndex, depth
					{
						position269, tokenIndex269, depth269 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l270
						}
						position++
						goto l269
					l270:
						position, tokenIndex, depth = position269, tokenIndex269, depth269
						if buffer[position] != rune('B') {
							goto l268
						}
						position++
					}
				l269:
					{
						position271, tokenIndex271, depth271 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l272
						}
						position++
						goto l271
					l272:
						position, tokenIndex, depth = position271, tokenIndex271, depth271
						if buffer[position] != rune('E') {
							goto l268
						}
						position++
					}
				l271:
					goto l267
				l268:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
					{
						position274, tokenIndex274, depth274 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l275
						}
						position++
						goto l274
					l275:
						position, tokenIndex, depth = position274, tokenIndex274, depth274
						if buffer[position] != rune('A') {
							goto l273
						}
						position++
					}
				l274:
					{
						position276, tokenIndex276, depth276 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l277
						}
						position++
						goto l276
					l277:
						position, tokenIndex, depth = position276, tokenIndex276, depth276
						if buffer[position] != rune('E') {
							goto l273
						}
						position++
					}
				l276:
					goto l267
				l273:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
					{
						switch buffer[position] {
						case 'N', 'n':
							{
								position279, tokenIndex279, depth279 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l280
								}
								position++
								goto l279
							l280:
								position, tokenIndex, depth = position279, tokenIndex279, depth279
								if buffer[position] != rune('N') {
									goto l265
								}
								position++
							}
						l279:
							{
								position281, tokenIndex281, depth281 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l282
								}
								position++
								goto l281
							l282:
								position, tokenIndex, depth = position281, tokenIndex281, depth281
								if buffer[position] != rune('Z') {
									goto l265
								}
								position++
							}
						l281:
							break
						case 'A', 'a':
							{
								position283, tokenIndex283, depth283 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l284
								}
								position++
								goto l283
							l284:
								position, tokenIndex, depth = position283, tokenIndex283, depth283
								if buffer[position] != rune('A') {
									goto l265
								}
								position++
							}
						l283:
							break
						case 'Z':
							if buffer[position] != rune('Z') {
								goto l265
							}
							position++
							break
						case 'z':
							if buffer[position] != rune('z') {
								goto l265
							}
							position++
							break
						default:
							{
								position285, tokenIndex285, depth285 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l286
								}
								position++
								goto l285
							l286:
								position, tokenIndex, depth = position285, tokenIndex285, depth285
								if buffer[position] != rune('B') {
									goto l265
								}
								position++
							}
						l285:
							break
						}
					}

				}
			l267:
				if !_rules[ruleSpace]() {
					goto l265
				}
				depth--
				add(ruleCMP_OP, position266)
			}
			return true
		l265:
			position, tokenIndex, depth = position265, tokenIndex265, depth265
			return false
		},
		/* 31 DATA_TYPE <- <(((&('I' | 'i') (('i' / 'I') ('n' / 'N') ('t' / 'T'))) | (&('F' | 'f') (('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))) | (&('B' | 'b') (('b' / 'B') ('y' / 'Y') ('t' / 'T') ('e' / 'E'))) | (&('W' | 'w') (('w' / 'W') ('o' / 'O') ('r' / 'R') ('d' / 'D'))) | (&('D' | 'd') (('d' / 'D') ('w' / 'W') ('o' / 'O') ('r' / 'R') ('d' / 'D')))) Space)> */
		func() bool {
			position287, tokenIndex287, depth287 := position, tokenIndex, depth
			{
				position288 := position
				depth++
				{
					switch buffer[position] {
					case 'I', 'i':
						{
							position290, tokenIndex290, depth290 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l291
							}
							position++
							goto l290
						l291:
							position, tokenIndex, depth = position290, tokenIndex290, depth290
							if buffer[position] != rune('I') {
								goto l287
							}
							position++
						}
					l290:
						{
							position292, tokenIndex292, depth292 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l293
							}
							position++
							goto l292
						l293:
							position, tokenIndex, depth = position292, tokenIndex292, depth292
							if buffer[position] != rune('N') {
								goto l287
							}
							position++
						}
					l292:
						{
							position294, tokenIndex294, depth294 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l295
							}
							position++
							goto l294
						l295:
							position, tokenIndex, depth = position294, tokenIndex294, depth294
							if buffer[position] != rune('T') {
								goto l287
							}
							position++
						}
					l294:
						break
					case 'F', 'f':
						{
							position296, tokenIndex296, depth296 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l297
							}
							position++
							goto l296
						l297:
							position, tokenIndex, depth = position296, tokenIndex296, depth296
							if buffer[position] != rune('F') {
								goto l287
							}
							position++
						}
					l296:
						{
							position298, tokenIndex298, depth298 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l299
							}
							position++
							goto l298
						l299:
							position, tokenIndex, depth = position298, tokenIndex298, depth298
							if buffer[position] != rune('L') {
								goto l287
							}
							position++
						}
					l298:
						{
							position300, tokenIndex300, depth300 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l301
							}
							position++
							goto l300
						l301:
							position, tokenIndex, depth = position300, tokenIndex300, depth300
							if buffer[position] != rune('O') {
								goto l287
							}
							position++
						}
					l300:
						{
							position302, tokenIndex302, depth302 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l303
							}
							position++
							goto l302
						l303:
							position, tokenIndex, depth = position302, tokenIndex302, depth302
							if buffer[position] != rune('A') {
								goto l287
							}
							position++
						}
					l302:
						{
							position304, tokenIndex304, depth304 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l305
							}
							position++
							goto l304
						l305:
							position, tokenIndex, depth = position304, tokenIndex304, depth304
							if buffer[position] != rune('T') {
								goto l287
							}
							position++
						}
					l304:
						break
					case 'B', 'b':
						{
							position306, tokenIndex306, depth306 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l307
							}
							position++
							goto l306
						l307:
							position, tokenIndex, depth = position306, tokenIndex306, depth306
							if buffer[position] != rune('B') {
								goto l287
							}
							position++
						}
					l306:
						{
							position308, tokenIndex308, depth308 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l309
							}
							position++
							goto l308
						l309:
							position, tokenIndex, depth = position308, tokenIndex308, depth308
							if buffer[position] != rune('Y') {
								goto l287
							}
							position++
						}
					l308:
						{
							position310, tokenIndex310, depth310 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l311
							}
							position++
							goto l310
						l311:
							position, tokenIndex, depth = position310, tokenIndex310, depth310
							if buffer[position] != rune('T') {
								goto l287
							}
							position++
						}
					l310:
						{
							position312, tokenIndex312, depth312 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l313
							}
							position++
							goto l312
						l313:
							position, tokenIndex, depth = position312, tokenIndex312, depth312
							if buffer[position] != rune('E') {
								goto l287
							}
							position++
						}
					l312:
						break
					case 'W', 'w':
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l315
							}
							position++
							goto l314
						l315:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
							if buffer[position] != rune('W') {
								goto l287
							}
							position++
						}
					l314:
						{
							position316, tokenIndex316, depth316 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l317
							}
							position++
							goto l316
						l317:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							if buffer[position] != rune('O') {
								goto l287
							}
							position++
						}
					l316:
						{
							position318, tokenIndex318, depth318 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l319
							}
							position++
							goto l318
						l319:
							position, tokenIndex, depth = position318, tokenIndex318, depth318
							if buffer[position] != rune('R') {
								goto l287
							}
							position++
						}
					l318:
						{
							position320, tokenIndex320, depth320 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l321
							}
							position++
							goto l320
						l321:
							position, tokenIndex, depth = position320, tokenIndex320, depth320
							if buffer[position] != rune('D') {
								goto l287
							}
							position++
						}
					l320:
						break
					default:
						{
							position322, tokenIndex322, depth322 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l323
							}
							position++
							goto l322
						l323:
							position, tokenIndex, depth = position322, tokenIndex322, depth322
							if buffer[position] != rune('D') {
								goto l287
							}
							position++
						}
					l322:
						{
							position324, tokenIndex324, depth324 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l325
							}
							position++
							goto l324
						l325:
							position, tokenIndex, depth = position324, tokenIndex324, depth324
							if buffer[position] != rune('W') {
								goto l287
							}
							position++
						}
					l324:
						{
							position326, tokenIndex326, depth326 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l327
							}
							position++
							goto l326
						l327:
							position, tokenIndex, depth = position326, tokenIndex326, depth326
							if buffer[position] != rune('O') {
								goto l287
							}
							position++
						}
					l326:
						{
							position328, tokenIndex328, depth328 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l329
							}
							position++
							goto l328
						l329:
							position, tokenIndex, depth = position328, tokenIndex328, depth328
							if buffer[position] != rune('R') {
								goto l287
							}
							position++
						}
					l328:
						{
							position330, tokenIndex330, depth330 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l331
							}
							position++
							goto l330
						l331:
							position, tokenIndex, depth = position330, tokenIndex330, depth330
							if buffer[position] != rune('D') {
								goto l287
							}
							position++
						}
					l330:
						break
					}
				}

				if !_rules[ruleSpace]() {
					goto l287
				}
				depth--
				add(ruleDATA_TYPE, position288)
			}
			return true
		l287:
			position, tokenIndex, depth = position287, tokenIndex287, depth287
			return false
		},
		/* 32 LBRK <- <('[' Spacing)> */
		func() bool {
			position332, tokenIndex332, depth332 := position, tokenIndex, depth
			{
				position333 := position
				depth++
				if buffer[position] != rune('[') {
					goto l332
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l332
				}
				depth--
				add(ruleLBRK, position333)
			}
			return true
		l332:
			position, tokenIndex, depth = position332, tokenIndex332, depth332
			return false
		},
		/* 33 RBRK <- <(']' Spacing)> */
		func() bool {
			position334, tokenIndex334, depth334 := position, tokenIndex, depth
			{
				position335 := position
				depth++
				if buffer[position] != rune(']') {
					goto l334
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l334
				}
				depth--
				add(ruleRBRK, position335)
			}
			return true
		l334:
			position, tokenIndex, depth = position334, tokenIndex334, depth334
			return false
		},
		/* 34 COMMA <- <(',' Spacing)> */
		func() bool {
			position336, tokenIndex336, depth336 := position, tokenIndex, depth
			{
				position337 := position
				depth++
				if buffer[position] != rune(',') {
					goto l336
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l336
				}
				depth--
				add(ruleCOMMA, position337)
			}
			return true
		l336:
			position, tokenIndex, depth = position336, tokenIndex336, depth336
			return false
		},
		/* 35 SEMICOLON <- <(';' Spacing)> */
		func() bool {
			position338, tokenIndex338, depth338 := position, tokenIndex, depth
			{
				position339 := position
				depth++
				if buffer[position] != rune(';') {
					goto l338
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l338
				}
				depth--
				add(ruleSEMICOLON, position339)
			}
			return true
		l338:
			position, tokenIndex, depth = position338, tokenIndex338, depth338
			return false
		},
		/* 36 COLON <- <(':' Spacing)> */
		func() bool {
			position340, tokenIndex340, depth340 := position, tokenIndex, depth
			{
				position341 := position
				depth++
				if buffer[position] != rune(':') {
					goto l340
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l340
				}
				depth--
				add(ruleCOLON, position341)
			}
			return true
		l340:
			position, tokenIndex, depth = position340, tokenIndex340, depth340
			return false
		},
		/* 37 MINUS <- <('-' Spacing)> */
		func() bool {
			position342, tokenIndex342, depth342 := position, tokenIndex, depth
			{
				position343 := position
				depth++
				if buffer[position] != rune('-') {
					goto l342
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l342
				}
				depth--
				add(ruleMINUS, position343)
			}
			return true
		l342:
			position, tokenIndex, depth = position342, tokenIndex342, depth342
			return false
		},
		/* 38 NL <- <'\n'> */
		func() bool {
			position344, tokenIndex344, depth344 := position, tokenIndex, depth
			{
				position345 := position
				depth++
				if buffer[position] != rune('\n') {
					goto l344
				}
				position++
				depth--
				add(ruleNL, position345)
			}
			return true
		l344:
			position, tokenIndex, depth = position344, tokenIndex344, depth344
			return false
		},
		/* 39 EOT <- <!.> */
		func() bool {
			position346, tokenIndex346, depth346 := position, tokenIndex, depth
			{
				position347 := position
				depth++
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					if !matchDot() {
						goto l348
					}
					goto l346
				l348:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
				}
				depth--
				add(ruleEOT, position347)
			}
			return true
		l346:
			position, tokenIndex, depth = position346, tokenIndex346, depth346
			return false
		},
		/* 40 Literal <- <((FloatLiteral / ((&('"') StringLiteral) | (&('\'') CharLiteral) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') IntegerLiteral))) Spacing)> */
		func() bool {
			position349, tokenIndex349, depth349 := position, tokenIndex, depth
			{
				position350 := position
				depth++
				{
					position351, tokenIndex351, depth351 := position, tokenIndex, depth
					if !_rules[ruleFloatLiteral]() {
						goto l352
					}
					goto l351
				l352:
					position, tokenIndex, depth = position351, tokenIndex351, depth351
					{
						switch buffer[position] {
						case '"':
							if !_rules[ruleStringLiteral]() {
								goto l349
							}
							break
						case '\'':
							if !_rules[ruleCharLiteral]() {
								goto l349
							}
							break
						default:
							if !_rules[ruleIntegerLiteral]() {
								goto l349
							}
							break
						}
					}

				}
			l351:
				if !_rules[ruleSpacing]() {
					goto l349
				}
				depth--
				add(ruleLiteral, position350)
			}
			return true
		l349:
			position, tokenIndex, depth = position349, tokenIndex349, depth349
			return false
		},
		/* 41 IntegerLiteral <- <(<(MINUS? (HexNumeral / BinaryNumeral / OctalNumeral / DecimalNumeral))> Spacing Action34)> */
		func() bool {
			position354, tokenIndex354, depth354 := position, tokenIndex, depth
			{
				position355 := position
				depth++
				{
					position356 := position
					depth++
					{
						position357, tokenIndex357, depth357 := position, tokenIndex, depth
						if !_rules[ruleMINUS]() {
							goto l357
						}
						goto l358
					l357:
						position, tokenIndex, depth = position357, tokenIndex357, depth357
					}
				l358:
					{
						position359, tokenIndex359, depth359 := position, tokenIndex, depth
						if !_rules[ruleHexNumeral]() {
							goto l360
						}
						goto l359
					l360:
						position, tokenIndex, depth = position359, tokenIndex359, depth359
						if !_rules[ruleBinaryNumeral]() {
							goto l361
						}
						goto l359
					l361:
						position, tokenIndex, depth = position359, tokenIndex359, depth359
						if !_rules[ruleOctalNumeral]() {
							goto l362
						}
						goto l359
					l362:
						position, tokenIndex, depth = position359, tokenIndex359, depth359
						if !_rules[ruleDecimalNumeral]() {
							goto l354
						}
					}
				l359:
					depth--
					add(rulePegText, position356)
				}
				if !_rules[ruleSpacing]() {
					goto l354
				}
				if !_rules[ruleAction34]() {
					goto l354
				}
				depth--
				add(ruleIntegerLiteral, position355)
			}
			return true
		l354:
			position, tokenIndex, depth = position354, tokenIndex354, depth354
			return false
		},
		/* 42 DecimalNumeral <- <('0' / ([1-9] ('_'* [0-9])*))> */
		func() bool {
			position363, tokenIndex363, depth363 := position, tokenIndex, depth
			{
				position364 := position
				depth++
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if buffer[position] != rune('0') {
						goto l366
					}
					position++
					goto l365
				l366:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
					if c := buffer[position]; c < rune('1') || c > rune('9') {
						goto l363
					}
					position++
				l367:
					{
						position368, tokenIndex368, depth368 := position, tokenIndex, depth
					l369:
						{
							position370, tokenIndex370, depth370 := position, tokenIndex, depth
							if buffer[position] != rune('_') {
								goto l370
							}
							position++
							goto l369
						l370:
							position, tokenIndex, depth = position370, tokenIndex370, depth370
						}
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l368
						}
						position++
						goto l367
					l368:
						position, tokenIndex, depth = position368, tokenIndex368, depth368
					}
				}
			l365:
				depth--
				add(ruleDecimalNumeral, position364)
			}
			return true
		l363:
			position, tokenIndex, depth = position363, tokenIndex363, depth363
			return false
		},
		/* 43 HexNumeral <- <((('0' 'x') / ('0' 'X')) HexDigits)> */
		func() bool {
			position371, tokenIndex371, depth371 := position, tokenIndex, depth
			{
				position372 := position
				depth++
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if buffer[position] != rune('0') {
						goto l374
					}
					position++
					if buffer[position] != rune('x') {
						goto l374
					}
					position++
					goto l373
				l374:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
					if buffer[position] != rune('0') {
						goto l371
					}
					position++
					if buffer[position] != rune('X') {
						goto l371
					}
					position++
				}
			l373:
				if !_rules[ruleHexDigits]() {
					goto l371
				}
				depth--
				add(ruleHexNumeral, position372)
			}
			return true
		l371:
			position, tokenIndex, depth = position371, tokenIndex371, depth371
			return false
		},
		/* 44 BinaryNumeral <- <((('0' 'b') / ('0' 'B')) ('0' / '1') ('_'* ('0' / '1'))*)> */
		func() bool {
			position375, tokenIndex375, depth375 := position, tokenIndex, depth
			{
				position376 := position
				depth++
				{
					position377, tokenIndex377, depth377 := position, tokenIndex, depth
					if buffer[position] != rune('0') {
						goto l378
					}
					position++
					if buffer[position] != rune('b') {
						goto l378
					}
					position++
					goto l377
				l378:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
					if buffer[position] != rune('0') {
						goto l375
					}
					position++
					if buffer[position] != rune('B') {
						goto l375
					}
					position++
				}
			l377:
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if buffer[position] != rune('0') {
						goto l380
					}
					position++
					goto l379
				l380:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
					if buffer[position] != rune('1') {
						goto l375
					}
					position++
				}
			l379:
			l381:
				{
					position382, tokenIndex382, depth382 := position, tokenIndex, depth
				l383:
					{
						position384, tokenIndex384, depth384 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l384
						}
						position++
						goto l383
					l384:
						position, tokenIndex, depth = position384, tokenIndex384, depth384
					}
					{
						position385, tokenIndex385, depth385 := position, tokenIndex, depth
						if buffer[position] != rune('0') {
							goto l386
						}
						position++
						goto l385
					l386:
						position, tokenIndex, depth = position385, tokenIndex385, depth385
						if buffer[position] != rune('1') {
							goto l382
						}
						position++
					}
				l385:
					goto l381
				l382:
					position, tokenIndex, depth = position382, tokenIndex382, depth382
				}
				depth--
				add(ruleBinaryNumeral, position376)
			}
			return true
		l375:
			position, tokenIndex, depth = position375, tokenIndex375, depth375
			return false
		},
		/* 45 OctalNumeral <- <('0' ('_'* [0-7])+)> */
		func() bool {
			position387, tokenIndex387, depth387 := position, tokenIndex, depth
			{
				position388 := position
				depth++
				if buffer[position] != rune('0') {
					goto l387
				}
				position++
			l391:
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if buffer[position] != rune('_') {
						goto l392
					}
					position++
					goto l391
				l392:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
				}
				if c := buffer[position]; c < rune('0') || c > rune('7') {
					goto l387
				}
				position++
			l389:
				{
					position390, tokenIndex390, depth390 := position, tokenIndex, depth
				l393:
					{
						position394, tokenIndex394, depth394 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l394
						}
						position++
						goto l393
					l394:
						position, tokenIndex, depth = position394, tokenIndex394, depth394
					}
					if c := buffer[position]; c < rune('0') || c > rune('7') {
						goto l390
					}
					position++
					goto l389
				l390:
					position, tokenIndex, depth = position390, tokenIndex390, depth390
				}
				depth--
				add(ruleOctalNumeral, position388)
			}
			return true
		l387:
			position, tokenIndex, depth = position387, tokenIndex387, depth387
			return false
		},
		/* 46 FloatLiteral <- <(HexFloat / DecimalFloat)> */
		func() bool {
			position395, tokenIndex395, depth395 := position, tokenIndex, depth
			{
				position396 := position
				depth++
				{
					position397, tokenIndex397, depth397 := position, tokenIndex, depth
					if !_rules[ruleHexFloat]() {
						goto l398
					}
					goto l397
				l398:
					position, tokenIndex, depth = position397, tokenIndex397, depth397
					if !_rules[ruleDecimalFloat]() {
						goto l395
					}
				}
			l397:
				depth--
				add(ruleFloatLiteral, position396)
			}
			return true
		l395:
			position, tokenIndex, depth = position395, tokenIndex395, depth395
			return false
		},
		/* 47 DecimalFloat <- <((Digits '.' Digits? Exponent?) / ('.' Digits Exponent?) / (Digits Exponent?))> */
		func() bool {
			position399, tokenIndex399, depth399 := position, tokenIndex, depth
			{
				position400 := position
				depth++
				{
					position401, tokenIndex401, depth401 := position, tokenIndex, depth
					if !_rules[ruleDigits]() {
						goto l402
					}
					if buffer[position] != rune('.') {
						goto l402
					}
					position++
					{
						position403, tokenIndex403, depth403 := position, tokenIndex, depth
						if !_rules[ruleDigits]() {
							goto l403
						}
						goto l404
					l403:
						position, tokenIndex, depth = position403, tokenIndex403, depth403
					}
				l404:
					{
						position405, tokenIndex405, depth405 := position, tokenIndex, depth
						if !_rules[ruleExponent]() {
							goto l405
						}
						goto l406
					l405:
						position, tokenIndex, depth = position405, tokenIndex405, depth405
					}
				l406:
					goto l401
				l402:
					position, tokenIndex, depth = position401, tokenIndex401, depth401
					if buffer[position] != rune('.') {
						goto l407
					}
					position++
					if !_rules[ruleDigits]() {
						goto l407
					}
					{
						position408, tokenIndex408, depth408 := position, tokenIndex, depth
						if !_rules[ruleExponent]() {
							goto l408
						}
						goto l409
					l408:
						position, tokenIndex, depth = position408, tokenIndex408, depth408
					}
				l409:
					goto l401
				l407:
					position, tokenIndex, depth = position401, tokenIndex401, depth401
					if !_rules[ruleDigits]() {
						goto l399
					}
					{
						position410, tokenIndex410, depth410 := position, tokenIndex, depth
						if !_rules[ruleExponent]() {
							goto l410
						}
						goto l411
					l410:
						position, tokenIndex, depth = position410, tokenIndex410, depth410
					}
				l411:
				}
			l401:
				depth--
				add(ruleDecimalFloat, position400)
			}
			return true
		l399:
			position, tokenIndex, depth = position399, tokenIndex399, depth399
			return false
		},
		/* 48 Exponent <- <(('e' / 'E') ('+' / '-')? Digits)> */
		func() bool {
			position412, tokenIndex412, depth412 := position, tokenIndex, depth
			{
				position413 := position
				depth++
				{
					position414, tokenIndex414, depth414 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l415
					}
					position++
					goto l414
				l415:
					position, tokenIndex, depth = position414, tokenIndex414, depth414
					if buffer[position] != rune('E') {
						goto l412
					}
					position++
				}
			l414:
				{
					position416, tokenIndex416, depth416 := position, tokenIndex, depth
					{
						position418, tokenIndex418, depth418 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l419
						}
						position++
						goto l418
					l419:
						position, tokenIndex, depth = position418, tokenIndex418, depth418
						if buffer[position] != rune('-') {
							goto l416
						}
						position++
					}
				l418:
					goto l417
				l416:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
				}
			l417:
				if !_rules[ruleDigits]() {
					goto l412
				}
				depth--
				add(ruleExponent, position413)
			}
			return true
		l412:
			position, tokenIndex, depth = position412, tokenIndex412, depth412
			return false
		},
		/* 49 HexFloat <- <(HexSignificand BinaryExponent)> */
		func() bool {
			position420, tokenIndex420, depth420 := position, tokenIndex, depth
			{
				position421 := position
				depth++
				if !_rules[ruleHexSignificand]() {
					goto l420
				}
				if !_rules[ruleBinaryExponent]() {
					goto l420
				}
				depth--
				add(ruleHexFloat, position421)
			}
			return true
		l420:
			position, tokenIndex, depth = position420, tokenIndex420, depth420
			return false
		},
		/* 50 HexSignificand <- <(((('0' 'x') / ('0' 'X')) HexDigits? '.' HexDigits) / (HexNumeral '.'?))> */
		func() bool {
			position422, tokenIndex422, depth422 := position, tokenIndex, depth
			{
				position423 := position
				depth++
				{
					position424, tokenIndex424, depth424 := position, tokenIndex, depth
					{
						position426, tokenIndex426, depth426 := position, tokenIndex, depth
						if buffer[position] != rune('0') {
							goto l427
						}
						position++
						if buffer[position] != rune('x') {
							goto l427
						}
						position++
						goto l426
					l427:
						position, tokenIndex, depth = position426, tokenIndex426, depth426
						if buffer[position] != rune('0') {
							goto l425
						}
						position++
						if buffer[position] != rune('X') {
							goto l425
						}
						position++
					}
				l426:
					{
						position428, tokenIndex428, depth428 := position, tokenIndex, depth
						if !_rules[ruleHexDigits]() {
							goto l428
						}
						goto l429
					l428:
						position, tokenIndex, depth = position428, tokenIndex428, depth428
					}
				l429:
					if buffer[position] != rune('.') {
						goto l425
					}
					position++
					if !_rules[ruleHexDigits]() {
						goto l425
					}
					goto l424
				l425:
					position, tokenIndex, depth = position424, tokenIndex424, depth424
					if !_rules[ruleHexNumeral]() {
						goto l422
					}
					{
						position430, tokenIndex430, depth430 := position, tokenIndex, depth
						if buffer[position] != rune('.') {
							goto l430
						}
						position++
						goto l431
					l430:
						position, tokenIndex, depth = position430, tokenIndex430, depth430
					}
				l431:
				}
			l424:
				depth--
				add(ruleHexSignificand, position423)
			}
			return true
		l422:
			position, tokenIndex, depth = position422, tokenIndex422, depth422
			return false
		},
		/* 51 BinaryExponent <- <(('p' / 'P') ('+' / '-')? Digits)> */
		func() bool {
			position432, tokenIndex432, depth432 := position, tokenIndex, depth
			{
				position433 := position
				depth++
				{
					position434, tokenIndex434, depth434 := position, tokenIndex, depth
					if buffer[position] != rune('p') {
						goto l435
					}
					position++
					goto l434
				l435:
					position, tokenIndex, depth = position434, tokenIndex434, depth434
					if buffer[position] != rune('P') {
						goto l432
					}
					position++
				}
			l434:
				{
					position436, tokenIndex436, depth436 := position, tokenIndex, depth
					{
						position438, tokenIndex438, depth438 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l439
						}
						position++
						goto l438
					l439:
						position, tokenIndex, depth = position438, tokenIndex438, depth438
						if buffer[position] != rune('-') {
							goto l436
						}
						position++
					}
				l438:
					goto l437
				l436:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
				}
			l437:
				if !_rules[ruleDigits]() {
					goto l432
				}
				depth--
				add(ruleBinaryExponent, position433)
			}
			return true
		l432:
			position, tokenIndex, depth = position432, tokenIndex432, depth432
			return false
		},
		/* 52 Digits <- <([0-9] ('_'* [0-9])*)> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l440
				}
				position++
			l442:
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
				l444:
					{
						position445, tokenIndex445, depth445 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l445
						}
						position++
						goto l444
					l445:
						position, tokenIndex, depth = position445, tokenIndex445, depth445
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l443
					}
					position++
					goto l442
				l443:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
				}
				depth--
				add(ruleDigits, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 53 HexDigits <- <(HexDigit ('_'* HexDigit)*)> */
		func() bool {
			position446, tokenIndex446, depth446 := position, tokenIndex, depth
			{
				position447 := position
				depth++
				if !_rules[ruleHexDigit]() {
					goto l446
				}
			l448:
				{
					position449, tokenIndex449, depth449 := position, tokenIndex, depth
				l450:
					{
						position451, tokenIndex451, depth451 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l451
						}
						position++
						goto l450
					l451:
						position, tokenIndex, depth = position451, tokenIndex451, depth451
					}
					if !_rules[ruleHexDigit]() {
						goto l449
					}
					goto l448
				l449:
					position, tokenIndex, depth = position449, tokenIndex449, depth449
				}
				depth--
				add(ruleHexDigits, position447)
			}
			return true
		l446:
			position, tokenIndex, depth = position446, tokenIndex446, depth446
			return false
		},
		/* 54 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position452, tokenIndex452, depth452 := position, tokenIndex, depth
			{
				position453 := position
				depth++
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l452
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l452
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l452
						}
						position++
						break
					}
				}

				depth--
				add(ruleHexDigit, position453)
			}
			return true
		l452:
			position, tokenIndex, depth = position452, tokenIndex452, depth452
			return false
		},
		/* 55 CharLiteral <- <('\'' (Escape / (!('\'' / '\\') .)) '\'')> */
		func() bool {
			position455, tokenIndex455, depth455 := position, tokenIndex, depth
			{
				position456 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l455
				}
				position++
				{
					position457, tokenIndex457, depth457 := position, tokenIndex, depth
					if !_rules[ruleEscape]() {
						goto l458
					}
					goto l457
				l458:
					position, tokenIndex, depth = position457, tokenIndex457, depth457
					{
						position459, tokenIndex459, depth459 := position, tokenIndex, depth
						{
							position460, tokenIndex460, depth460 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l461
							}
							position++
							goto l460
						l461:
							position, tokenIndex, depth = position460, tokenIndex460, depth460
							if buffer[position] != rune('\\') {
								goto l459
							}
							position++
						}
					l460:
						goto l455
					l459:
						position, tokenIndex, depth = position459, tokenIndex459, depth459
					}
					if !matchDot() {
						goto l455
					}
				}
			l457:
				if buffer[position] != rune('\'') {
					goto l455
				}
				position++
				depth--
				add(ruleCharLiteral, position456)
			}
			return true
		l455:
			position, tokenIndex, depth = position455, tokenIndex455, depth455
			return false
		},
		/* 56 StringLiteral <- <(<('"' (Escape / (!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('"') '"')) .))* '"')> Action35)> */
		func() bool {
			position462, tokenIndex462, depth462 := position, tokenIndex, depth
			{
				position463 := position
				depth++
				{
					position464 := position
					depth++
					if buffer[position] != rune('"') {
						goto l462
					}
					position++
				l465:
					{
						position466, tokenIndex466, depth466 := position, tokenIndex, depth
						{
							position467, tokenIndex467, depth467 := position, tokenIndex, depth
							if !_rules[ruleEscape]() {
								goto l468
							}
							goto l467
						l468:
							position, tokenIndex, depth = position467, tokenIndex467, depth467
							{
								position469, tokenIndex469, depth469 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '\r':
										if buffer[position] != rune('\r') {
											goto l469
										}
										position++
										break
									case '\n':
										if buffer[position] != rune('\n') {
											goto l469
										}
										position++
										break
									case '\\':
										if buffer[position] != rune('\\') {
											goto l469
										}
										position++
										break
									default:
										if buffer[position] != rune('"') {
											goto l469
										}
										position++
										break
									}
								}

								goto l466
							l469:
								position, tokenIndex, depth = position469, tokenIndex469, depth469
							}
							if !matchDot() {
								goto l466
							}
						}
					l467:
						goto l465
					l466:
						position, tokenIndex, depth = position466, tokenIndex466, depth466
					}
					if buffer[position] != rune('"') {
						goto l462
					}
					position++
					depth--
					add(rulePegText, position464)
				}
				if !_rules[ruleAction35]() {
					goto l462
				}
				depth--
				add(ruleStringLiteral, position463)
			}
			return true
		l462:
			position, tokenIndex, depth = position462, tokenIndex462, depth462
			return false
		},
		/* 57 Escape <- <('\\' ((&('u') UnicodeEscape) | (&('\\') '\\') | (&('\'') '\'') | (&('"') '"') | (&('r') 'r') | (&('f') 'f') | (&('n') 'n') | (&('t') 't') | (&('b') 'b') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7') OctalEscape)))> */
		func() bool {
			position471, tokenIndex471, depth471 := position, tokenIndex, depth
			{
				position472 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l471
				}
				position++
				{
					switch buffer[position] {
					case 'u':
						if !_rules[ruleUnicodeEscape]() {
							goto l471
						}
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l471
						}
						position++
						break
					case '\'':
						if buffer[position] != rune('\'') {
							goto l471
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l471
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l471
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l471
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l471
						}
						position++
						break
					case 't':
						if buffer[position] != rune('t') {
							goto l471
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l471
						}
						position++
						break
					default:
						if !_rules[ruleOctalEscape]() {
							goto l471
						}
						break
					}
				}

				depth--
				add(ruleEscape, position472)
			}
			return true
		l471:
			position, tokenIndex, depth = position471, tokenIndex471, depth471
			return false
		},
		/* 58 OctalEscape <- <(([0-3] [0-7] [0-7]) / ([0-7] [0-7]) / [0-7])> */
		func() bool {
			position474, tokenIndex474, depth474 := position, tokenIndex, depth
			{
				position475 := position
				depth++
				{
					position476, tokenIndex476, depth476 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('3') {
						goto l477
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('7') {
						goto l477
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('7') {
						goto l477
					}
					position++
					goto l476
				l477:
					position, tokenIndex, depth = position476, tokenIndex476, depth476
					if c := buffer[position]; c < rune('0') || c > rune('7') {
						goto l478
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('7') {
						goto l478
					}
					position++
					goto l476
				l478:
					position, tokenIndex, depth = position476, tokenIndex476, depth476
					if c := buffer[position]; c < rune('0') || c > rune('7') {
						goto l474
					}
					position++
				}
			l476:
				depth--
				add(ruleOctalEscape, position475)
			}
			return true
		l474:
			position, tokenIndex, depth = position474, tokenIndex474, depth474
			return false
		},
		/* 59 UnicodeEscape <- <('u'+ HexDigit HexDigit HexDigit HexDigit)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				if buffer[position] != rune('u') {
					goto l479
				}
				position++
			l481:
				{
					position482, tokenIndex482, depth482 := position, tokenIndex, depth
					if buffer[position] != rune('u') {
						goto l482
					}
					position++
					goto l481
				l482:
					position, tokenIndex, depth = position482, tokenIndex482, depth482
				}
				if !_rules[ruleHexDigit]() {
					goto l479
				}
				if !_rules[ruleHexDigit]() {
					goto l479
				}
				if !_rules[ruleHexDigit]() {
					goto l479
				}
				if !_rules[ruleHexDigit]() {
					goto l479
				}
				depth--
				add(ruleUnicodeEscape, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 61 Action0 <- <{p.line++}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 62 Action1 <- <{p.AddAssembly()}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 63 Action2 <- <{p.AddAssembly();p.AddComment()}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		nil,
		/* 65 Action3 <- <{p.Push(&Comment{});p.Push(text)}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 66 Action4 <- <{p.Push(&Label{})}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 67 Action5 <- <{p.Push(lookup(bbasm.INT,text))}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 68 Action6 <- <{p.Push(lookup(bbasm.ADD,text))}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 69 Action7 <- <{p.Push(lookup(bbasm.INT,text))}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 70 Action8 <- <{p.Push(lookup(bbasm.INT,text))}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 71 Action9 <- <{p.Push(lookup(bbasm.A,text))}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 72 Action10 <- <{p.AddPseudoDataValue()}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 73 Action11 <- <{p.AddPseudoDataValue()}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 74 Action12 <- <{p.AddPseudoDataValue()}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 75 Action13 <- <{p.Push(text);p.AddPseudoDataValue()}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 76 Action14 <- <{p.AddOperand(true)}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 77 Action15 <- <{p.AddOperand(false)}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 78 Action16 <- <{p.AddOperand(true)}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 79 Action17 <- <{p.AddOperand(false)}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 80 Action18 <- <{p.Push(text)}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 81 Action19 <- <{p.PushInst(bbasm.EXIT)}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 82 Action20 <- <{p.PushInst(bbasm.RET)}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 83 Action21 <- <{p.PushInst(bbasm.NOP)}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 84 Action22 <- <{p.PushInst(bbasm.CALL)}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 85 Action23 <- <{p.PushInst(bbasm.PUSH)}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 86 Action24 <- <{p.PushInst(bbasm.POP)}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 87 Action25 <- <{p.PushInst(bbasm.JMP)}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 88 Action26 <- <{p.PushInst(bbasm.IN)}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 89 Action27 <- <{p.PushInst(bbasm.OUT)}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 90 Action28 <- <{p.PushInst(bbasm.CAL)}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 91 Action29 <- <{p.PushInst(bbasm.LD)}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 92 Action30 <- <{p.PushInst(bbasm.CMP)}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 93 Action31 <- <{p.PushInst(bbasm.JPC)}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 94 Action32 <- <{p.Push(&PseudoBlock{})}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 95 Action33 <- <{p.Push(&PseudoData{})}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 96 Action34 <- <{p.Push(text);p.AddInteger()}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 97 Action35 <- <{p.Push(text)}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
	}
	p.rules = _rules
}
