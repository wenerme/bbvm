package parser

import (
	"fmt"
	"github.com/wenerme/bbvm/libbbvm/asm"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleStart
	ruleLine
	ruleComment
	ruleLabel
	ruleInst
	rulePseudo
	rulePseudoDataValue
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

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"Start",
	"Line",
	"Comment",
	"Label",
	"Inst",
	"Pseudo",
	"PseudoDataValue",
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

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next uint32, depth int)
	Expand(index int) tokenTree
	Tokens() <-chan token32
	AST() *node32
	Error() []token32
	trim(length int)
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

func (ast *node32) Print(buffer string) {
	ast.print(0, buffer)
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
		for i, _ := range states {
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
							write(token32{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
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
					write(token32{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
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
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

/*func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2 * len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}*/

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type BBAsm struct {
	parser

	Buffer string
	buffer []rune
	rules  [96]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	tokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer string, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range []rune(buffer) {
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
	p *BBAsm
}

func (e *parseError) Error() string {
	tokens, error := e.p.tokenTree.Error(), "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n",
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			/*strconv.Quote(*/ e.p.Buffer[begin:end] /*)*/)
	}

	return error
}

func (p *BBAsm) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *BBAsm) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *BBAsm) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.tokenTree.Tokens() {
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
			p.Push(&asm.Comment{})
			p.Push(text)
		case ruleAction4:
			p.Push(&asm.Label{})
		case ruleAction5:
			p.Push(lookup(asm.T_INT, text))
		case ruleAction6:
			p.Push(lookup(asm.CAL_ADD, text))
		case ruleAction7:
			p.Push(lookup(asm.T_INT, text))
		case ruleAction8:
			p.Push(lookup(asm.T_INT, text))
		case ruleAction9:
			p.Push(lookup(asm.CMP_A, text))
		case ruleAction10:
			p.AddPseudoDataValue()
		case ruleAction11:
			p.AddPseudoDataValue()
		case ruleAction12:
			p.AddPseudoDataValue()
		case ruleAction13:
			p.AddOperand(true)
		case ruleAction14:
			p.AddOperand(false)
		case ruleAction15:
			p.AddOperand(true)
		case ruleAction16:
			p.AddOperand(false)
		case ruleAction17:
			p.Push(text)
		case ruleAction18:
			p.PushInst(asm.OP_EXIT)
		case ruleAction19:
			p.PushInst(asm.OP_RET)
		case ruleAction20:
			p.PushInst(asm.OP_NOP)
		case ruleAction21:
			p.PushInst(asm.OP_CALL)
		case ruleAction22:
			p.PushInst(asm.OP_PUSH)
		case ruleAction23:
			p.PushInst(asm.OP_POP)
		case ruleAction24:
			p.PushInst(asm.OP_JMP)
		case ruleAction25:
			p.PushInst(asm.OP_IN)
		case ruleAction26:
			p.PushInst(asm.OP_OUT)
		case ruleAction27:
			p.PushInst(asm.OP_CAL)
		case ruleAction28:
			p.PushInst(asm.OP_LD)
		case ruleAction29:
			p.PushInst(asm.OP_CMP)
		case ruleAction30:
			p.PushInst(asm.OP_JPC)
		case ruleAction31:
			p.Push(&asm.PseudoBlock{})
		case ruleAction32:
			p.Push(&asm.PseudoData{})
		case ruleAction33:
			p.Push(text)
			p.AddInteger()
		case ruleAction34:
			p.Push(text)

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *BBAsm) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens32{tree: make([]token32, math.MaxInt16)}
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokenTree = tree
		if matches {
			p.tokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
	}

	matchDot := func() bool {
		if buffer[position] != end_symbol {
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
		/* 0 Start <- <((Spacing Line? NL Action0)* EOT Literal?)> */
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
						{
							position6 := position
							depth++
							{
								position7, tokenIndex7, depth7 := position, tokenIndex, depth
								{
									position9 := position
									depth++
									{
										add(ruleAction4, position)
									}
									if !_rules[ruleIdentifier]() {
										goto l8
									}
									if !_rules[ruleSpacing]() {
										goto l8
									}
									{
										position11 := position
										depth++
										if buffer[position] != rune(':') {
											goto l8
										}
										position++
										if !_rules[ruleSpacing]() {
											goto l8
										}
										depth--
										add(ruleCOLON, position11)
									}
									depth--
									add(ruleLabel, position9)
								}
								goto l7
							l8:
								position, tokenIndex, depth = position7, tokenIndex7, depth7
								{
									switch buffer[position] {
									case '.', 'D', 'd':
										{
											position13 := position
											depth++
											{
												position14, tokenIndex14, depth14 := position, tokenIndex, depth
												{
													position16 := position
													depth++
													if buffer[position] != rune('.') {
														goto l15
													}
													position++
													{
														position17, tokenIndex17, depth17 := position, tokenIndex, depth
														if buffer[position] != rune('b') {
															goto l18
														}
														position++
														goto l17
													l18:
														position, tokenIndex, depth = position17, tokenIndex17, depth17
														if buffer[position] != rune('B') {
															goto l15
														}
														position++
													}
												l17:
													{
														position19, tokenIndex19, depth19 := position, tokenIndex, depth
														if buffer[position] != rune('l') {
															goto l20
														}
														position++
														goto l19
													l20:
														position, tokenIndex, depth = position19, tokenIndex19, depth19
														if buffer[position] != rune('L') {
															goto l15
														}
														position++
													}
												l19:
													{
														position21, tokenIndex21, depth21 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l22
														}
														position++
														goto l21
													l22:
														position, tokenIndex, depth = position21, tokenIndex21, depth21
														if buffer[position] != rune('O') {
															goto l15
														}
														position++
													}
												l21:
													{
														position23, tokenIndex23, depth23 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l24
														}
														position++
														goto l23
													l24:
														position, tokenIndex, depth = position23, tokenIndex23, depth23
														if buffer[position] != rune('C') {
															goto l15
														}
														position++
													}
												l23:
													{
														position25, tokenIndex25, depth25 := position, tokenIndex, depth
														if buffer[position] != rune('k') {
															goto l26
														}
														position++
														goto l25
													l26:
														position, tokenIndex, depth = position25, tokenIndex25, depth25
														if buffer[position] != rune('K') {
															goto l15
														}
														position++
													}
												l25:
													if !_rules[ruleSpace]() {
														goto l15
													}
													{
														add(ruleAction31, position)
													}
													depth--
													add(ruleBLOCK, position16)
												}
												if !_rules[ruleIntegerLiteral]() {
													goto l15
												}
												if !_rules[ruleIntegerLiteral]() {
													goto l15
												}
												goto l14
											l15:
												position, tokenIndex, depth = position14, tokenIndex14, depth14
												{
													position28 := position
													depth++
													{
														position29, tokenIndex29, depth29 := position, tokenIndex, depth
														if buffer[position] != rune('d') {
															goto l30
														}
														position++
														goto l29
													l30:
														position, tokenIndex, depth = position29, tokenIndex29, depth29
														if buffer[position] != rune('D') {
															goto l4
														}
														position++
													}
												l29:
													{
														position31, tokenIndex31, depth31 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l32
														}
														position++
														goto l31
													l32:
														position, tokenIndex, depth = position31, tokenIndex31, depth31
														if buffer[position] != rune('A') {
															goto l4
														}
														position++
													}
												l31:
													{
														position33, tokenIndex33, depth33 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l34
														}
														position++
														goto l33
													l34:
														position, tokenIndex, depth = position33, tokenIndex33, depth33
														if buffer[position] != rune('T') {
															goto l4
														}
														position++
													}
												l33:
													{
														position35, tokenIndex35, depth35 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l36
														}
														position++
														goto l35
													l36:
														position, tokenIndex, depth = position35, tokenIndex35, depth35
														if buffer[position] != rune('A') {
															goto l4
														}
														position++
													}
												l35:
													if !_rules[ruleSpace]() {
														goto l4
													}
													{
														add(ruleAction32, position)
													}
													depth--
													add(ruleDATA, position28)
												}
												if !_rules[ruleIdentifier]() {
													goto l4
												}
												{
													position38, tokenIndex38, depth38 := position, tokenIndex, depth
													if !_rules[ruleDATA_TYPE]() {
														goto l38
													}
													goto l39
												l38:
													position, tokenIndex, depth = position38, tokenIndex38, depth38
												}
											l39:
												if !_rules[rulePseudoDataValue]() {
													goto l4
												}
											l40:
												{
													position41, tokenIndex41, depth41 := position, tokenIndex, depth
													if !_rules[ruleCOMMA]() {
														goto l41
													}
													if !_rules[rulePseudoDataValue]() {
														goto l41
													}
													goto l40
												l41:
													position, tokenIndex, depth = position41, tokenIndex41, depth41
												}
											}
										l14:
											depth--
											add(rulePseudo, position13)
										}
										break
									case ';':
										if !_rules[ruleComment]() {
											goto l4
										}
										break
									default:
										{
											position42 := position
											depth++
											{
												position43, tokenIndex43, depth43 := position, tokenIndex, depth
												{
													position45, tokenIndex45, depth45 := position, tokenIndex, depth
													{
														position47 := position
														depth++
														{
															position48, tokenIndex48, depth48 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l49
															}
															position++
															goto l48
														l49:
															position, tokenIndex, depth = position48, tokenIndex48, depth48
															if buffer[position] != rune('P') {
																goto l46
															}
															position++
														}
													l48:
														{
															position50, tokenIndex50, depth50 := position, tokenIndex, depth
															if buffer[position] != rune('u') {
																goto l51
															}
															position++
															goto l50
														l51:
															position, tokenIndex, depth = position50, tokenIndex50, depth50
															if buffer[position] != rune('U') {
																goto l46
															}
															position++
														}
													l50:
														{
															position52, tokenIndex52, depth52 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l53
															}
															position++
															goto l52
														l53:
															position, tokenIndex, depth = position52, tokenIndex52, depth52
															if buffer[position] != rune('S') {
																goto l46
															}
															position++
														}
													l52:
														{
															position54, tokenIndex54, depth54 := position, tokenIndex, depth
															if buffer[position] != rune('h') {
																goto l55
															}
															position++
															goto l54
														l55:
															position, tokenIndex, depth = position54, tokenIndex54, depth54
															if buffer[position] != rune('H') {
																goto l46
															}
															position++
														}
													l54:
														if !_rules[ruleSpace]() {
															goto l46
														}
														{
															add(ruleAction22, position)
														}
														depth--
														add(rulePUSH, position47)
													}
													goto l45
												l46:
													position, tokenIndex, depth = position45, tokenIndex45, depth45
													{
														switch buffer[position] {
														case 'J', 'j':
															{
																position58 := position
																depth++
																{
																	position59, tokenIndex59, depth59 := position, tokenIndex, depth
																	if buffer[position] != rune('j') {
																		goto l60
																	}
																	position++
																	goto l59
																l60:
																	position, tokenIndex, depth = position59, tokenIndex59, depth59
																	if buffer[position] != rune('J') {
																		goto l44
																	}
																	position++
																}
															l59:
																{
																	position61, tokenIndex61, depth61 := position, tokenIndex, depth
																	if buffer[position] != rune('m') {
																		goto l62
																	}
																	position++
																	goto l61
																l62:
																	position, tokenIndex, depth = position61, tokenIndex61, depth61
																	if buffer[position] != rune('M') {
																		goto l44
																	}
																	position++
																}
															l61:
																{
																	position63, tokenIndex63, depth63 := position, tokenIndex, depth
																	if buffer[position] != rune('p') {
																		goto l64
																	}
																	position++
																	goto l63
																l64:
																	position, tokenIndex, depth = position63, tokenIndex63, depth63
																	if buffer[position] != rune('P') {
																		goto l44
																	}
																	position++
																}
															l63:
																if !_rules[ruleSpace]() {
																	goto l44
																}
																{
																	add(ruleAction24, position)
																}
																depth--
																add(ruleJMP, position58)
															}
															break
														case 'P', 'p':
															{
																position66 := position
																depth++
																{
																	position67, tokenIndex67, depth67 := position, tokenIndex, depth
																	if buffer[position] != rune('p') {
																		goto l68
																	}
																	position++
																	goto l67
																l68:
																	position, tokenIndex, depth = position67, tokenIndex67, depth67
																	if buffer[position] != rune('P') {
																		goto l44
																	}
																	position++
																}
															l67:
																{
																	position69, tokenIndex69, depth69 := position, tokenIndex, depth
																	if buffer[position] != rune('o') {
																		goto l70
																	}
																	position++
																	goto l69
																l70:
																	position, tokenIndex, depth = position69, tokenIndex69, depth69
																	if buffer[position] != rune('O') {
																		goto l44
																	}
																	position++
																}
															l69:
																{
																	position71, tokenIndex71, depth71 := position, tokenIndex, depth
																	if buffer[position] != rune('p') {
																		goto l72
																	}
																	position++
																	goto l71
																l72:
																	position, tokenIndex, depth = position71, tokenIndex71, depth71
																	if buffer[position] != rune('P') {
																		goto l44
																	}
																	position++
																}
															l71:
																if !_rules[ruleSpace]() {
																	goto l44
																}
																{
																	add(ruleAction23, position)
																}
																depth--
																add(rulePOP, position66)
															}
															break
														default:
															{
																position74 := position
																depth++
																{
																	position75, tokenIndex75, depth75 := position, tokenIndex, depth
																	if buffer[position] != rune('c') {
																		goto l76
																	}
																	position++
																	goto l75
																l76:
																	position, tokenIndex, depth = position75, tokenIndex75, depth75
																	if buffer[position] != rune('C') {
																		goto l44
																	}
																	position++
																}
															l75:
																{
																	position77, tokenIndex77, depth77 := position, tokenIndex, depth
																	if buffer[position] != rune('a') {
																		goto l78
																	}
																	position++
																	goto l77
																l78:
																	position, tokenIndex, depth = position77, tokenIndex77, depth77
																	if buffer[position] != rune('A') {
																		goto l44
																	}
																	position++
																}
															l77:
																{
																	position79, tokenIndex79, depth79 := position, tokenIndex, depth
																	if buffer[position] != rune('l') {
																		goto l80
																	}
																	position++
																	goto l79
																l80:
																	position, tokenIndex, depth = position79, tokenIndex79, depth79
																	if buffer[position] != rune('L') {
																		goto l44
																	}
																	position++
																}
															l79:
																{
																	position81, tokenIndex81, depth81 := position, tokenIndex, depth
																	if buffer[position] != rune('l') {
																		goto l82
																	}
																	position++
																	goto l81
																l82:
																	position, tokenIndex, depth = position81, tokenIndex81, depth81
																	if buffer[position] != rune('L') {
																		goto l44
																	}
																	position++
																}
															l81:
																if !_rules[ruleSpace]() {
																	goto l44
																}
																{
																	add(ruleAction21, position)
																}
																depth--
																add(ruleCALL, position74)
															}
															break
														}
													}

												}
											l45:
												if !_rules[ruleOperand]() {
													goto l44
												}
												goto l43
											l44:
												position, tokenIndex, depth = position43, tokenIndex43, depth43
												{
													position85 := position
													depth++
													{
														position86, tokenIndex86, depth86 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l87
														}
														position++
														goto l86
													l87:
														position, tokenIndex, depth = position86, tokenIndex86, depth86
														if buffer[position] != rune('C') {
															goto l84
														}
														position++
													}
												l86:
													{
														position88, tokenIndex88, depth88 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l89
														}
														position++
														goto l88
													l89:
														position, tokenIndex, depth = position88, tokenIndex88, depth88
														if buffer[position] != rune('A') {
															goto l84
														}
														position++
													}
												l88:
													{
														position90, tokenIndex90, depth90 := position, tokenIndex, depth
														if buffer[position] != rune('l') {
															goto l91
														}
														position++
														goto l90
													l91:
														position, tokenIndex, depth = position90, tokenIndex90, depth90
														if buffer[position] != rune('L') {
															goto l84
														}
														position++
													}
												l90:
													if !_rules[ruleSpace]() {
														goto l84
													}
													{
														add(ruleAction27, position)
													}
													depth--
													add(ruleCAL, position85)
												}
												{
													position93 := position
													depth++
													if !_rules[ruleDATA_TYPE]() {
														goto l84
													}
													depth--
													add(rulePegText, position93)
												}
												{
													add(ruleAction5, position)
												}
												{
													position95 := position
													depth++
													{
														position96 := position
														depth++
														{
															position97, tokenIndex97, depth97 := position, tokenIndex, depth
															{
																position99, tokenIndex99, depth99 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l100
																}
																position++
																goto l99
															l100:
																position, tokenIndex, depth = position99, tokenIndex99, depth99
																if buffer[position] != rune('M') {
																	goto l98
																}
																position++
															}
														l99:
															{
																position101, tokenIndex101, depth101 := position, tokenIndex, depth
																if buffer[position] != rune('u') {
																	goto l102
																}
																position++
																goto l101
															l102:
																position, tokenIndex, depth = position101, tokenIndex101, depth101
																if buffer[position] != rune('U') {
																	goto l98
																}
																position++
															}
														l101:
															{
																position103, tokenIndex103, depth103 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l104
																}
																position++
																goto l103
															l104:
																position, tokenIndex, depth = position103, tokenIndex103, depth103
																if buffer[position] != rune('L') {
																	goto l98
																}
																position++
															}
														l103:
															goto l97
														l98:
															position, tokenIndex, depth = position97, tokenIndex97, depth97
															{
																switch buffer[position] {
																case 'M', 'm':
																	{
																		position106, tokenIndex106, depth106 := position, tokenIndex, depth
																		if buffer[position] != rune('m') {
																			goto l107
																		}
																		position++
																		goto l106
																	l107:
																		position, tokenIndex, depth = position106, tokenIndex106, depth106
																		if buffer[position] != rune('M') {
																			goto l84
																		}
																		position++
																	}
																l106:
																	{
																		position108, tokenIndex108, depth108 := position, tokenIndex, depth
																		if buffer[position] != rune('o') {
																			goto l109
																		}
																		position++
																		goto l108
																	l109:
																		position, tokenIndex, depth = position108, tokenIndex108, depth108
																		if buffer[position] != rune('O') {
																			goto l84
																		}
																		position++
																	}
																l108:
																	{
																		position110, tokenIndex110, depth110 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l111
																		}
																		position++
																		goto l110
																	l111:
																		position, tokenIndex, depth = position110, tokenIndex110, depth110
																		if buffer[position] != rune('D') {
																			goto l84
																		}
																		position++
																	}
																l110:
																	break
																case 'D', 'd':
																	{
																		position112, tokenIndex112, depth112 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l113
																		}
																		position++
																		goto l112
																	l113:
																		position, tokenIndex, depth = position112, tokenIndex112, depth112
																		if buffer[position] != rune('D') {
																			goto l84
																		}
																		position++
																	}
																l112:
																	{
																		position114, tokenIndex114, depth114 := position, tokenIndex, depth
																		if buffer[position] != rune('i') {
																			goto l115
																		}
																		position++
																		goto l114
																	l115:
																		position, tokenIndex, depth = position114, tokenIndex114, depth114
																		if buffer[position] != rune('I') {
																			goto l84
																		}
																		position++
																	}
																l114:
																	{
																		position116, tokenIndex116, depth116 := position, tokenIndex, depth
																		if buffer[position] != rune('v') {
																			goto l117
																		}
																		position++
																		goto l116
																	l117:
																		position, tokenIndex, depth = position116, tokenIndex116, depth116
																		if buffer[position] != rune('V') {
																			goto l84
																		}
																		position++
																	}
																l116:
																	break
																case 'S', 's':
																	{
																		position118, tokenIndex118, depth118 := position, tokenIndex, depth
																		if buffer[position] != rune('s') {
																			goto l119
																		}
																		position++
																		goto l118
																	l119:
																		position, tokenIndex, depth = position118, tokenIndex118, depth118
																		if buffer[position] != rune('S') {
																			goto l84
																		}
																		position++
																	}
																l118:
																	{
																		position120, tokenIndex120, depth120 := position, tokenIndex, depth
																		if buffer[position] != rune('u') {
																			goto l121
																		}
																		position++
																		goto l120
																	l121:
																		position, tokenIndex, depth = position120, tokenIndex120, depth120
																		if buffer[position] != rune('U') {
																			goto l84
																		}
																		position++
																	}
																l120:
																	{
																		position122, tokenIndex122, depth122 := position, tokenIndex, depth
																		if buffer[position] != rune('b') {
																			goto l123
																		}
																		position++
																		goto l122
																	l123:
																		position, tokenIndex, depth = position122, tokenIndex122, depth122
																		if buffer[position] != rune('B') {
																			goto l84
																		}
																		position++
																	}
																l122:
																	break
																default:
																	{
																		position124, tokenIndex124, depth124 := position, tokenIndex, depth
																		if buffer[position] != rune('a') {
																			goto l125
																		}
																		position++
																		goto l124
																	l125:
																		position, tokenIndex, depth = position124, tokenIndex124, depth124
																		if buffer[position] != rune('A') {
																			goto l84
																		}
																		position++
																	}
																l124:
																	{
																		position126, tokenIndex126, depth126 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l127
																		}
																		position++
																		goto l126
																	l127:
																		position, tokenIndex, depth = position126, tokenIndex126, depth126
																		if buffer[position] != rune('D') {
																			goto l84
																		}
																		position++
																	}
																l126:
																	{
																		position128, tokenIndex128, depth128 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l129
																		}
																		position++
																		goto l128
																	l129:
																		position, tokenIndex, depth = position128, tokenIndex128, depth128
																		if buffer[position] != rune('D') {
																			goto l84
																		}
																		position++
																	}
																l128:
																	break
																}
															}

														}
													l97:
														if !_rules[ruleSpace]() {
															goto l84
														}
														depth--
														add(ruleCAL_OP, position96)
													}
													depth--
													add(rulePegText, position95)
												}
												{
													add(ruleAction6, position)
												}
												if !_rules[ruleOperand]() {
													goto l84
												}
												if !_rules[ruleCOMMA]() {
													goto l84
												}
												if !_rules[ruleOperand]() {
													goto l84
												}
												goto l43
											l84:
												position, tokenIndex, depth = position43, tokenIndex43, depth43
												{
													switch buffer[position] {
													case 'J', 'j':
														{
															position132 := position
															depth++
															{
																position133, tokenIndex133, depth133 := position, tokenIndex, depth
																if buffer[position] != rune('j') {
																	goto l134
																}
																position++
																goto l133
															l134:
																position, tokenIndex, depth = position133, tokenIndex133, depth133
																if buffer[position] != rune('J') {
																	goto l4
																}
																position++
															}
														l133:
															{
																position135, tokenIndex135, depth135 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l136
																}
																position++
																goto l135
															l136:
																position, tokenIndex, depth = position135, tokenIndex135, depth135
																if buffer[position] != rune('P') {
																	goto l4
																}
																position++
															}
														l135:
															{
																position137, tokenIndex137, depth137 := position, tokenIndex, depth
																if buffer[position] != rune('c') {
																	goto l138
																}
																position++
																goto l137
															l138:
																position, tokenIndex, depth = position137, tokenIndex137, depth137
																if buffer[position] != rune('C') {
																	goto l4
																}
																position++
															}
														l137:
															if !_rules[ruleSpace]() {
																goto l4
															}
															{
																add(ruleAction30, position)
															}
															depth--
															add(ruleJPC, position132)
														}
														{
															position140 := position
															depth++
															{
																position141 := position
																depth++
																{
																	position142, tokenIndex142, depth142 := position, tokenIndex, depth
																	{
																		position144, tokenIndex144, depth144 := position, tokenIndex, depth
																		if buffer[position] != rune('b') {
																			goto l145
																		}
																		position++
																		goto l144
																	l145:
																		position, tokenIndex, depth = position144, tokenIndex144, depth144
																		if buffer[position] != rune('B') {
																			goto l143
																		}
																		position++
																	}
																l144:
																	{
																		position146, tokenIndex146, depth146 := position, tokenIndex, depth
																		if buffer[position] != rune('e') {
																			goto l147
																		}
																		position++
																		goto l146
																	l147:
																		position, tokenIndex, depth = position146, tokenIndex146, depth146
																		if buffer[position] != rune('E') {
																			goto l143
																		}
																		position++
																	}
																l146:
																	goto l142
																l143:
																	position, tokenIndex, depth = position142, tokenIndex142, depth142
																	{
																		position149, tokenIndex149, depth149 := position, tokenIndex, depth
																		if buffer[position] != rune('a') {
																			goto l150
																		}
																		position++
																		goto l149
																	l150:
																		position, tokenIndex, depth = position149, tokenIndex149, depth149
																		if buffer[position] != rune('A') {
																			goto l148
																		}
																		position++
																	}
																l149:
																	{
																		position151, tokenIndex151, depth151 := position, tokenIndex, depth
																		if buffer[position] != rune('e') {
																			goto l152
																		}
																		position++
																		goto l151
																	l152:
																		position, tokenIndex, depth = position151, tokenIndex151, depth151
																		if buffer[position] != rune('E') {
																			goto l148
																		}
																		position++
																	}
																l151:
																	goto l142
																l148:
																	position, tokenIndex, depth = position142, tokenIndex142, depth142
																	{
																		switch buffer[position] {
																		case 'N', 'n':
																			{
																				position154, tokenIndex154, depth154 := position, tokenIndex, depth
																				if buffer[position] != rune('n') {
																					goto l155
																				}
																				position++
																				goto l154
																			l155:
																				position, tokenIndex, depth = position154, tokenIndex154, depth154
																				if buffer[position] != rune('N') {
																					goto l4
																				}
																				position++
																			}
																		l154:
																			{
																				position156, tokenIndex156, depth156 := position, tokenIndex, depth
																				if buffer[position] != rune('z') {
																					goto l157
																				}
																				position++
																				goto l156
																			l157:
																				position, tokenIndex, depth = position156, tokenIndex156, depth156
																				if buffer[position] != rune('Z') {
																					goto l4
																				}
																				position++
																			}
																		l156:
																			break
																		case 'A', 'a':
																			{
																				position158, tokenIndex158, depth158 := position, tokenIndex, depth
																				if buffer[position] != rune('a') {
																					goto l159
																				}
																				position++
																				goto l158
																			l159:
																				position, tokenIndex, depth = position158, tokenIndex158, depth158
																				if buffer[position] != rune('A') {
																					goto l4
																				}
																				position++
																			}
																		l158:
																			break
																		case 'Z':
																			if buffer[position] != rune('Z') {
																				goto l4
																			}
																			position++
																			break
																		case 'z':
																			if buffer[position] != rune('z') {
																				goto l4
																			}
																			position++
																			break
																		default:
																			{
																				position160, tokenIndex160, depth160 := position, tokenIndex, depth
																				if buffer[position] != rune('b') {
																					goto l161
																				}
																				position++
																				goto l160
																			l161:
																				position, tokenIndex, depth = position160, tokenIndex160, depth160
																				if buffer[position] != rune('B') {
																					goto l4
																				}
																				position++
																			}
																		l160:
																			break
																		}
																	}

																}
															l142:
																if !_rules[ruleSpace]() {
																	goto l4
																}
																depth--
																add(ruleCMP_OP, position141)
															}
															depth--
															add(rulePegText, position140)
														}
														{
															add(ruleAction9, position)
														}
														if !_rules[ruleOperand]() {
															goto l4
														}
														break
													case 'C', 'c':
														{
															position163 := position
															depth++
															{
																position164, tokenIndex164, depth164 := position, tokenIndex, depth
																if buffer[position] != rune('c') {
																	goto l165
																}
																position++
																goto l164
															l165:
																position, tokenIndex, depth = position164, tokenIndex164, depth164
																if buffer[position] != rune('C') {
																	goto l4
																}
																position++
															}
														l164:
															{
																position166, tokenIndex166, depth166 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l167
																}
																position++
																goto l166
															l167:
																position, tokenIndex, depth = position166, tokenIndex166, depth166
																if buffer[position] != rune('M') {
																	goto l4
																}
																position++
															}
														l166:
															{
																position168, tokenIndex168, depth168 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l169
																}
																position++
																goto l168
															l169:
																position, tokenIndex, depth = position168, tokenIndex168, depth168
																if buffer[position] != rune('P') {
																	goto l4
																}
																position++
															}
														l168:
															if !_rules[ruleSpace]() {
																goto l4
															}
															{
																add(ruleAction29, position)
															}
															depth--
															add(ruleCMP, position163)
														}
														{
															position171 := position
															depth++
															if !_rules[ruleDATA_TYPE]() {
																goto l4
															}
															depth--
															add(rulePegText, position171)
														}
														{
															add(ruleAction8, position)
														}
														if !_rules[ruleOperand]() {
															goto l4
														}
														if !_rules[ruleCOMMA]() {
															goto l4
														}
														if !_rules[ruleOperand]() {
															goto l4
														}
														break
													case 'L', 'l':
														{
															position173 := position
															depth++
															{
																position174, tokenIndex174, depth174 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l175
																}
																position++
																goto l174
															l175:
																position, tokenIndex, depth = position174, tokenIndex174, depth174
																if buffer[position] != rune('L') {
																	goto l4
																}
																position++
															}
														l174:
															{
																position176, tokenIndex176, depth176 := position, tokenIndex, depth
																if buffer[position] != rune('d') {
																	goto l177
																}
																position++
																goto l176
															l177:
																position, tokenIndex, depth = position176, tokenIndex176, depth176
																if buffer[position] != rune('D') {
																	goto l4
																}
																position++
															}
														l176:
															if !_rules[ruleSpace]() {
																goto l4
															}
															{
																add(ruleAction28, position)
															}
															depth--
															add(ruleLD, position173)
														}
														{
															position179 := position
															depth++
															if !_rules[ruleDATA_TYPE]() {
																goto l4
															}
															depth--
															add(rulePegText, position179)
														}
														{
															add(ruleAction7, position)
														}
														if !_rules[ruleOperand]() {
															goto l4
														}
														if !_rules[ruleCOMMA]() {
															goto l4
														}
														if !_rules[ruleOperand]() {
															goto l4
														}
														break
													case 'N', 'n':
														{
															position181 := position
															depth++
															{
																position182, tokenIndex182, depth182 := position, tokenIndex, depth
																if buffer[position] != rune('n') {
																	goto l183
																}
																position++
																goto l182
															l183:
																position, tokenIndex, depth = position182, tokenIndex182, depth182
																if buffer[position] != rune('N') {
																	goto l4
																}
																position++
															}
														l182:
															{
																position184, tokenIndex184, depth184 := position, tokenIndex, depth
																if buffer[position] != rune('o') {
																	goto l185
																}
																position++
																goto l184
															l185:
																position, tokenIndex, depth = position184, tokenIndex184, depth184
																if buffer[position] != rune('O') {
																	goto l4
																}
																position++
															}
														l184:
															{
																position186, tokenIndex186, depth186 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l187
																}
																position++
																goto l186
															l187:
																position, tokenIndex, depth = position186, tokenIndex186, depth186
																if buffer[position] != rune('P') {
																	goto l4
																}
																position++
															}
														l186:
															if !_rules[ruleSpacing]() {
																goto l4
															}
															{
																add(ruleAction20, position)
															}
															depth--
															add(ruleNOP, position181)
														}
														break
													case 'R', 'r':
														{
															position189 := position
															depth++
															{
																position190, tokenIndex190, depth190 := position, tokenIndex, depth
																if buffer[position] != rune('r') {
																	goto l191
																}
																position++
																goto l190
															l191:
																position, tokenIndex, depth = position190, tokenIndex190, depth190
																if buffer[position] != rune('R') {
																	goto l4
																}
																position++
															}
														l190:
															{
																position192, tokenIndex192, depth192 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l193
																}
																position++
																goto l192
															l193:
																position, tokenIndex, depth = position192, tokenIndex192, depth192
																if buffer[position] != rune('E') {
																	goto l4
																}
																position++
															}
														l192:
															{
																position194, tokenIndex194, depth194 := position, tokenIndex, depth
																if buffer[position] != rune('t') {
																	goto l195
																}
																position++
																goto l194
															l195:
																position, tokenIndex, depth = position194, tokenIndex194, depth194
																if buffer[position] != rune('T') {
																	goto l4
																}
																position++
															}
														l194:
															if !_rules[ruleSpacing]() {
																goto l4
															}
															{
																add(ruleAction19, position)
															}
															depth--
															add(ruleRET, position189)
														}
														break
													case 'E', 'e':
														{
															position197 := position
															depth++
															{
																position198, tokenIndex198, depth198 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l199
																}
																position++
																goto l198
															l199:
																position, tokenIndex, depth = position198, tokenIndex198, depth198
																if buffer[position] != rune('E') {
																	goto l4
																}
																position++
															}
														l198:
															{
																position200, tokenIndex200, depth200 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l201
																}
																position++
																goto l200
															l201:
																position, tokenIndex, depth = position200, tokenIndex200, depth200
																if buffer[position] != rune('X') {
																	goto l4
																}
																position++
															}
														l200:
															{
																position202, tokenIndex202, depth202 := position, tokenIndex, depth
																if buffer[position] != rune('i') {
																	goto l203
																}
																position++
																goto l202
															l203:
																position, tokenIndex, depth = position202, tokenIndex202, depth202
																if buffer[position] != rune('I') {
																	goto l4
																}
																position++
															}
														l202:
															{
																position204, tokenIndex204, depth204 := position, tokenIndex, depth
																if buffer[position] != rune('t') {
																	goto l205
																}
																position++
																goto l204
															l205:
																position, tokenIndex, depth = position204, tokenIndex204, depth204
																if buffer[position] != rune('T') {
																	goto l4
																}
																position++
															}
														l204:
															if !_rules[ruleSpacing]() {
																goto l4
															}
															{
																add(ruleAction18, position)
															}
															depth--
															add(ruleEXIT, position197)
														}
														break
													default:
														{
															position207, tokenIndex207, depth207 := position, tokenIndex, depth
															{
																position209 := position
																depth++
																{
																	position210, tokenIndex210, depth210 := position, tokenIndex, depth
																	if buffer[position] != rune('i') {
																		goto l211
																	}
																	position++
																	goto l210
																l211:
																	position, tokenIndex, depth = position210, tokenIndex210, depth210
																	if buffer[position] != rune('I') {
																		goto l208
																	}
																	position++
																}
															l210:
																{
																	position212, tokenIndex212, depth212 := position, tokenIndex, depth
																	if buffer[position] != rune('n') {
																		goto l213
																	}
																	position++
																	goto l212
																l213:
																	position, tokenIndex, depth = position212, tokenIndex212, depth212
																	if buffer[position] != rune('N') {
																		goto l208
																	}
																	position++
																}
															l212:
																if !_rules[ruleSpace]() {
																	goto l208
																}
																{
																	add(ruleAction25, position)
																}
																depth--
																add(ruleIN, position209)
															}
															goto l207
														l208:
															position, tokenIndex, depth = position207, tokenIndex207, depth207
															{
																position215 := position
																depth++
																{
																	position216, tokenIndex216, depth216 := position, tokenIndex, depth
																	if buffer[position] != rune('o') {
																		goto l217
																	}
																	position++
																	goto l216
																l217:
																	position, tokenIndex, depth = position216, tokenIndex216, depth216
																	if buffer[position] != rune('O') {
																		goto l4
																	}
																	position++
																}
															l216:
																{
																	position218, tokenIndex218, depth218 := position, tokenIndex, depth
																	if buffer[position] != rune('u') {
																		goto l219
																	}
																	position++
																	goto l218
																l219:
																	position, tokenIndex, depth = position218, tokenIndex218, depth218
																	if buffer[position] != rune('U') {
																		goto l4
																	}
																	position++
																}
															l218:
																{
																	position220, tokenIndex220, depth220 := position, tokenIndex, depth
																	if buffer[position] != rune('t') {
																		goto l221
																	}
																	position++
																	goto l220
																l221:
																	position, tokenIndex, depth = position220, tokenIndex220, depth220
																	if buffer[position] != rune('T') {
																		goto l4
																	}
																	position++
																}
															l220:
																if !_rules[ruleSpace]() {
																	goto l4
																}
																{
																	add(ruleAction26, position)
																}
																depth--
																add(ruleOUT, position215)
															}
														}
													l207:
														if !_rules[ruleOperand]() {
															goto l4
														}
														if !_rules[ruleCOMMA]() {
															goto l4
														}
														if !_rules[ruleOperand]() {
															goto l4
														}
														break
													}
												}

											}
										l43:
											depth--
											add(ruleInst, position42)
										}
										break
									}
								}

							}
						l7:
							{
								add(ruleAction1, position)
							}
							{
								position224, tokenIndex224, depth224 := position, tokenIndex, depth
								if !_rules[ruleComment]() {
									goto l224
								}
								{
									add(ruleAction2, position)
								}
								goto l225
							l224:
								position, tokenIndex, depth = position224, tokenIndex224, depth224
							}
						l225:
							depth--
							add(ruleLine, position6)
						}
						goto l5
					l4:
						position, tokenIndex, depth = position4, tokenIndex4, depth4
					}
				l5:
					if !_rules[ruleNL]() {
						goto l3
					}
					{
						add(ruleAction0, position)
					}
					goto l2
				l3:
					position, tokenIndex, depth = position3, tokenIndex3, depth3
				}
				{
					position228 := position
					depth++
					{
						position229, tokenIndex229, depth229 := position, tokenIndex, depth
						if !matchDot() {
							goto l229
						}
						goto l0
					l229:
						position, tokenIndex, depth = position229, tokenIndex229, depth229
					}
					depth--
					add(ruleEOT, position228)
				}
				{
					position230, tokenIndex230, depth230 := position, tokenIndex, depth
					{
						position232 := position
						depth++
						{
							position233, tokenIndex233, depth233 := position, tokenIndex, depth
							{
								position235 := position
								depth++
								{
									position236, tokenIndex236, depth236 := position, tokenIndex, depth
									{
										position238 := position
										depth++
										{
											position239 := position
											depth++
											{
												position240, tokenIndex240, depth240 := position, tokenIndex, depth
												{
													position242, tokenIndex242, depth242 := position, tokenIndex, depth
													if buffer[position] != rune('0') {
														goto l243
													}
													position++
													if buffer[position] != rune('x') {
														goto l243
													}
													position++
													goto l242
												l243:
													position, tokenIndex, depth = position242, tokenIndex242, depth242
													if buffer[position] != rune('0') {
														goto l241
													}
													position++
													if buffer[position] != rune('X') {
														goto l241
													}
													position++
												}
											l242:
												{
													position244, tokenIndex244, depth244 := position, tokenIndex, depth
													if !_rules[ruleHexDigits]() {
														goto l244
													}
													goto l245
												l244:
													position, tokenIndex, depth = position244, tokenIndex244, depth244
												}
											l245:
												if buffer[position] != rune('.') {
													goto l241
												}
												position++
												if !_rules[ruleHexDigits]() {
													goto l241
												}
												goto l240
											l241:
												position, tokenIndex, depth = position240, tokenIndex240, depth240
												if !_rules[ruleHexNumeral]() {
													goto l237
												}
												{
													position246, tokenIndex246, depth246 := position, tokenIndex, depth
													if buffer[position] != rune('.') {
														goto l246
													}
													position++
													goto l247
												l246:
													position, tokenIndex, depth = position246, tokenIndex246, depth246
												}
											l247:
											}
										l240:
											depth--
											add(ruleHexSignificand, position239)
										}
										{
											position248 := position
											depth++
											{
												position249, tokenIndex249, depth249 := position, tokenIndex, depth
												if buffer[position] != rune('p') {
													goto l250
												}
												position++
												goto l249
											l250:
												position, tokenIndex, depth = position249, tokenIndex249, depth249
												if buffer[position] != rune('P') {
													goto l237
												}
												position++
											}
										l249:
											{
												position251, tokenIndex251, depth251 := position, tokenIndex, depth
												{
													position253, tokenIndex253, depth253 := position, tokenIndex, depth
													if buffer[position] != rune('+') {
														goto l254
													}
													position++
													goto l253
												l254:
													position, tokenIndex, depth = position253, tokenIndex253, depth253
													if buffer[position] != rune('-') {
														goto l251
													}
													position++
												}
											l253:
												goto l252
											l251:
												position, tokenIndex, depth = position251, tokenIndex251, depth251
											}
										l252:
											if !_rules[ruleDigits]() {
												goto l237
											}
											depth--
											add(ruleBinaryExponent, position248)
										}
										{
											position255, tokenIndex255, depth255 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 'D':
													if buffer[position] != rune('D') {
														goto l255
													}
													position++
													break
												case 'd':
													if buffer[position] != rune('d') {
														goto l255
													}
													position++
													break
												case 'F':
													if buffer[position] != rune('F') {
														goto l255
													}
													position++
													break
												default:
													if buffer[position] != rune('f') {
														goto l255
													}
													position++
													break
												}
											}

											goto l256
										l255:
											position, tokenIndex, depth = position255, tokenIndex255, depth255
										}
									l256:
										depth--
										add(ruleHexFloat, position238)
									}
									goto l236
								l237:
									position, tokenIndex, depth = position236, tokenIndex236, depth236
									{
										position258 := position
										depth++
										{
											position259, tokenIndex259, depth259 := position, tokenIndex, depth
											if !_rules[ruleDigits]() {
												goto l260
											}
											if buffer[position] != rune('.') {
												goto l260
											}
											position++
											{
												position261, tokenIndex261, depth261 := position, tokenIndex, depth
												if !_rules[ruleDigits]() {
													goto l261
												}
												goto l262
											l261:
												position, tokenIndex, depth = position261, tokenIndex261, depth261
											}
										l262:
											{
												position263, tokenIndex263, depth263 := position, tokenIndex, depth
												if !_rules[ruleExponent]() {
													goto l263
												}
												goto l264
											l263:
												position, tokenIndex, depth = position263, tokenIndex263, depth263
											}
										l264:
											{
												position265, tokenIndex265, depth265 := position, tokenIndex, depth
												{
													switch buffer[position] {
													case 'D':
														if buffer[position] != rune('D') {
															goto l265
														}
														position++
														break
													case 'd':
														if buffer[position] != rune('d') {
															goto l265
														}
														position++
														break
													case 'F':
														if buffer[position] != rune('F') {
															goto l265
														}
														position++
														break
													default:
														if buffer[position] != rune('f') {
															goto l265
														}
														position++
														break
													}
												}

												goto l266
											l265:
												position, tokenIndex, depth = position265, tokenIndex265, depth265
											}
										l266:
											goto l259
										l260:
											position, tokenIndex, depth = position259, tokenIndex259, depth259
											if buffer[position] != rune('.') {
												goto l268
											}
											position++
											if !_rules[ruleDigits]() {
												goto l268
											}
											{
												position269, tokenIndex269, depth269 := position, tokenIndex, depth
												if !_rules[ruleExponent]() {
													goto l269
												}
												goto l270
											l269:
												position, tokenIndex, depth = position269, tokenIndex269, depth269
											}
										l270:
											{
												position271, tokenIndex271, depth271 := position, tokenIndex, depth
												{
													switch buffer[position] {
													case 'D':
														if buffer[position] != rune('D') {
															goto l271
														}
														position++
														break
													case 'd':
														if buffer[position] != rune('d') {
															goto l271
														}
														position++
														break
													case 'F':
														if buffer[position] != rune('F') {
															goto l271
														}
														position++
														break
													default:
														if buffer[position] != rune('f') {
															goto l271
														}
														position++
														break
													}
												}

												goto l272
											l271:
												position, tokenIndex, depth = position271, tokenIndex271, depth271
											}
										l272:
											goto l259
										l268:
											position, tokenIndex, depth = position259, tokenIndex259, depth259
											if !_rules[ruleDigits]() {
												goto l274
											}
											if !_rules[ruleExponent]() {
												goto l274
											}
											{
												position275, tokenIndex275, depth275 := position, tokenIndex, depth
												{
													switch buffer[position] {
													case 'D':
														if buffer[position] != rune('D') {
															goto l275
														}
														position++
														break
													case 'd':
														if buffer[position] != rune('d') {
															goto l275
														}
														position++
														break
													case 'F':
														if buffer[position] != rune('F') {
															goto l275
														}
														position++
														break
													default:
														if buffer[position] != rune('f') {
															goto l275
														}
														position++
														break
													}
												}

												goto l276
											l275:
												position, tokenIndex, depth = position275, tokenIndex275, depth275
											}
										l276:
											goto l259
										l274:
											position, tokenIndex, depth = position259, tokenIndex259, depth259
											if !_rules[ruleDigits]() {
												goto l234
											}
											{
												position278, tokenIndex278, depth278 := position, tokenIndex, depth
												if !_rules[ruleExponent]() {
													goto l278
												}
												goto l279
											l278:
												position, tokenIndex, depth = position278, tokenIndex278, depth278
											}
										l279:
											{
												switch buffer[position] {
												case 'D':
													if buffer[position] != rune('D') {
														goto l234
													}
													position++
													break
												case 'd':
													if buffer[position] != rune('d') {
														goto l234
													}
													position++
													break
												case 'F':
													if buffer[position] != rune('F') {
														goto l234
													}
													position++
													break
												default:
													if buffer[position] != rune('f') {
														goto l234
													}
													position++
													break
												}
											}

										}
									l259:
										depth--
										add(ruleDecimalFloat, position258)
									}
								}
							l236:
								depth--
								add(ruleFloatLiteral, position235)
							}
							goto l233
						l234:
							position, tokenIndex, depth = position233, tokenIndex233, depth233
							{
								switch buffer[position] {
								case '"':
									if !_rules[ruleStringLiteral]() {
										goto l230
									}
									break
								case '\'':
									{
										position282 := position
										depth++
										if buffer[position] != rune('\'') {
											goto l230
										}
										position++
										{
											position283, tokenIndex283, depth283 := position, tokenIndex, depth
											if !_rules[ruleEscape]() {
												goto l284
											}
											goto l283
										l284:
											position, tokenIndex, depth = position283, tokenIndex283, depth283
											{
												position285, tokenIndex285, depth285 := position, tokenIndex, depth
												{
													position286, tokenIndex286, depth286 := position, tokenIndex, depth
													if buffer[position] != rune('\'') {
														goto l287
													}
													position++
													goto l286
												l287:
													position, tokenIndex, depth = position286, tokenIndex286, depth286
													if buffer[position] != rune('\\') {
														goto l285
													}
													position++
												}
											l286:
												goto l230
											l285:
												position, tokenIndex, depth = position285, tokenIndex285, depth285
											}
											if !matchDot() {
												goto l230
											}
										}
									l283:
										if buffer[position] != rune('\'') {
											goto l230
										}
										position++
										depth--
										add(ruleCharLiteral, position282)
									}
									break
								default:
									if !_rules[ruleIntegerLiteral]() {
										goto l230
									}
									break
								}
							}

						}
					l233:
						if !_rules[ruleSpacing]() {
							goto l230
						}
						depth--
						add(ruleLiteral, position232)
					}
					goto l231
				l230:
					position, tokenIndex, depth = position230, tokenIndex230, depth230
				}
			l231:
				depth--
				add(ruleStart, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Line <- <((Label / ((&('.' | 'D' | 'd') Pseudo) | (&(';') Comment) | (&('C' | 'E' | 'I' | 'J' | 'L' | 'N' | 'O' | 'P' | 'R' | 'c' | 'e' | 'i' | 'j' | 'l' | 'n' | 'o' | 'p' | 'r') Inst))) Action1 (Comment Action2)?)> */
		nil,
		/* 2 Comment <- <(SEMICOLON <(!NL .)*> Action3)> */
		func() bool {
			position289, tokenIndex289, depth289 := position, tokenIndex, depth
			{
				position290 := position
				depth++
				{
					position291 := position
					depth++
					if buffer[position] != rune(';') {
						goto l289
					}
					position++
					if !_rules[ruleSpacing]() {
						goto l289
					}
					depth--
					add(ruleSEMICOLON, position291)
				}
				{
					position292 := position
					depth++
				l293:
					{
						position294, tokenIndex294, depth294 := position, tokenIndex, depth
						{
							position295, tokenIndex295, depth295 := position, tokenIndex, depth
							if !_rules[ruleNL]() {
								goto l295
							}
							goto l294
						l295:
							position, tokenIndex, depth = position295, tokenIndex295, depth295
						}
						if !matchDot() {
							goto l294
						}
						goto l293
					l294:
						position, tokenIndex, depth = position294, tokenIndex294, depth294
					}
					depth--
					add(rulePegText, position292)
				}
				{
					add(ruleAction3, position)
				}
				depth--
				add(ruleComment, position290)
			}
			return true
		l289:
			position, tokenIndex, depth = position289, tokenIndex289, depth289
			return false
		},
		/* 3 Label <- <(Action4 Identifier Spacing COLON)> */
		nil,
		/* 4 Inst <- <(((PUSH / ((&('J' | 'j') JMP) | (&('P' | 'p') POP) | (&('C' | 'c') CALL))) Operand) / (CAL <DATA_TYPE> Action5 <CAL_OP> Action6 Operand COMMA Operand) / ((&('J' | 'j') (JPC <CMP_OP> Action9 Operand)) | (&('C' | 'c') (CMP <DATA_TYPE> Action8 Operand COMMA Operand)) | (&('L' | 'l') (LD <DATA_TYPE> Action7 Operand COMMA Operand)) | (&('N' | 'n') NOP) | (&('R' | 'r') RET) | (&('E' | 'e') EXIT) | (&('I' | 'O' | 'i' | 'o') ((IN / OUT) Operand COMMA Operand))))> */
		nil,
		/* 5 Pseudo <- <((BLOCK IntegerLiteral IntegerLiteral) / (DATA Identifier DATA_TYPE? PseudoDataValue (COMMA PseudoDataValue)*))> */
		nil,
		/* 6 PseudoDataValue <- <((&('"') (StringLiteral Action12)) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (IntegerLiteral Action10)) | (&('$' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (Identifier Action11)))> */
		func() bool {
			position300, tokenIndex300, depth300 := position, tokenIndex, depth
			{
				position301 := position
				depth++
				{
					switch buffer[position] {
					case '"':
						if !_rules[ruleStringLiteral]() {
							goto l300
						}
						{
							add(ruleAction12, position)
						}
						break
					case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if !_rules[ruleIntegerLiteral]() {
							goto l300
						}
						{
							add(ruleAction10, position)
						}
						break
					default:
						if !_rules[ruleIdentifier]() {
							goto l300
						}
						{
							add(ruleAction11, position)
						}
						break
					}
				}

				depth--
				add(rulePseudoDataValue, position301)
			}
			return true
		l300:
			position, tokenIndex, depth = position300, tokenIndex300, depth300
			return false
		},
		/* 7 Operand <- <(((LBRK Identifier RBRK Action14) / ((&('[') (LBRK IntegerLiteral RBRK Action16)) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (IntegerLiteral Action15)) | (&('$' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (Identifier Action13)))) Spacing)> */
		func() bool {
			position306, tokenIndex306, depth306 := position, tokenIndex, depth
			{
				position307 := position
				depth++
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					if !_rules[ruleLBRK]() {
						goto l309
					}
					if !_rules[ruleIdentifier]() {
						goto l309
					}
					if !_rules[ruleRBRK]() {
						goto l309
					}
					{
						add(ruleAction14, position)
					}
					goto l308
				l309:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
					{
						switch buffer[position] {
						case '[':
							if !_rules[ruleLBRK]() {
								goto l306
							}
							if !_rules[ruleIntegerLiteral]() {
								goto l306
							}
							if !_rules[ruleRBRK]() {
								goto l306
							}
							{
								add(ruleAction16, position)
							}
							break
						case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if !_rules[ruleIntegerLiteral]() {
								goto l306
							}
							{
								add(ruleAction15, position)
							}
							break
						default:
							if !_rules[ruleIdentifier]() {
								goto l306
							}
							{
								add(ruleAction13, position)
							}
							break
						}
					}

				}
			l308:
				if !_rules[ruleSpacing]() {
					goto l306
				}
				depth--
				add(ruleOperand, position307)
			}
			return true
		l306:
			position, tokenIndex, depth = position306, tokenIndex306, depth306
			return false
		},
		/* 8 Spacing <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position316 := position
				depth++
			l317:
				{
					position318, tokenIndex318, depth318 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l318
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l318
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l318
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l318
							}
							position++
							break
						}
					}

					goto l317
				l318:
					position, tokenIndex, depth = position318, tokenIndex318, depth318
				}
				depth--
				add(ruleSpacing, position316)
			}
			return true
		},
		/* 9 Space <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))+> */
		func() bool {
			position320, tokenIndex320, depth320 := position, tokenIndex, depth
			{
				position321 := position
				depth++
				{
					switch buffer[position] {
					case '\f':
						if buffer[position] != rune('\f') {
							goto l320
						}
						position++
						break
					case '\r':
						if buffer[position] != rune('\r') {
							goto l320
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l320
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l320
						}
						position++
						break
					}
				}

			l322:
				{
					position323, tokenIndex323, depth323 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l323
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l323
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l323
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l323
							}
							position++
							break
						}
					}

					goto l322
				l323:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
				}
				depth--
				add(ruleSpace, position321)
			}
			return true
		l320:
			position, tokenIndex, depth = position320, tokenIndex320, depth320
			return false
		},
		/* 10 Identifier <- <(<(Letter LetterOrDigit*)> Spacing Action17)> */
		func() bool {
			position326, tokenIndex326, depth326 := position, tokenIndex, depth
			{
				position327 := position
				depth++
				{
					position328 := position
					depth++
					{
						position329 := position
						depth++
						{
							switch buffer[position] {
							case '$', '_':
								{
									position331, tokenIndex331, depth331 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l332
									}
									position++
									goto l331
								l332:
									position, tokenIndex, depth = position331, tokenIndex331, depth331
									if buffer[position] != rune('$') {
										goto l326
									}
									position++
								}
							l331:
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l326
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l326
								}
								position++
								break
							}
						}

						depth--
						add(ruleLetter, position329)
					}
				l333:
					{
						position334, tokenIndex334, depth334 := position, tokenIndex, depth
						{
							position335 := position
							depth++
							{
								switch buffer[position] {
								case '$', '_':
									{
										position337, tokenIndex337, depth337 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l338
										}
										position++
										goto l337
									l338:
										position, tokenIndex, depth = position337, tokenIndex337, depth337
										if buffer[position] != rune('$') {
											goto l334
										}
										position++
									}
								l337:
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l334
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l334
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l334
									}
									position++
									break
								}
							}

							depth--
							add(ruleLetterOrDigit, position335)
						}
						goto l333
					l334:
						position, tokenIndex, depth = position334, tokenIndex334, depth334
					}
					depth--
					add(rulePegText, position328)
				}
				if !_rules[ruleSpacing]() {
					goto l326
				}
				{
					add(ruleAction17, position)
				}
				depth--
				add(ruleIdentifier, position327)
			}
			return true
		l326:
			position, tokenIndex, depth = position326, tokenIndex326, depth326
			return false
		},
		/* 11 Letter <- <((&('$' | '_') ('_' / '$')) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 12 LetterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 13 EXIT <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('t' / 'T') Spacing Action18)> */
		nil,
		/* 14 RET <- <(('r' / 'R') ('e' / 'E') ('t' / 'T') Spacing Action19)> */
		nil,
		/* 15 NOP <- <(('n' / 'N') ('o' / 'O') ('p' / 'P') Spacing Action20)> */
		nil,
		/* 16 CALL <- <(('c' / 'C') ('a' / 'A') ('l' / 'L') ('l' / 'L') Space Action21)> */
		nil,
		/* 17 PUSH <- <(('p' / 'P') ('u' / 'U') ('s' / 'S') ('h' / 'H') Space Action22)> */
		nil,
		/* 18 POP <- <(('p' / 'P') ('o' / 'O') ('p' / 'P') Space Action23)> */
		nil,
		/* 19 JMP <- <(('j' / 'J') ('m' / 'M') ('p' / 'P') Space Action24)> */
		nil,
		/* 20 IN <- <(('i' / 'I') ('n' / 'N') Space Action25)> */
		nil,
		/* 21 OUT <- <(('o' / 'O') ('u' / 'U') ('t' / 'T') Space Action26)> */
		nil,
		/* 22 CAL <- <(('c' / 'C') ('a' / 'A') ('l' / 'L') Space Action27)> */
		nil,
		/* 23 LD <- <(('l' / 'L') ('d' / 'D') Space Action28)> */
		nil,
		/* 24 CMP <- <(('c' / 'C') ('m' / 'M') ('p' / 'P') Space Action29)> */
		nil,
		/* 25 JPC <- <(('j' / 'J') ('p' / 'P') ('c' / 'C') Space Action30)> */
		nil,
		/* 26 BLOCK <- <('.' ('b' / 'B') ('l' / 'L') ('o' / 'O') ('c' / 'C') ('k' / 'K') Space Action31)> */
		nil,
		/* 27 DATA <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') Space Action32)> */
		nil,
		/* 28 CAL_OP <- <(((('m' / 'M') ('u' / 'U') ('l' / 'L')) / ((&('M' | 'm') (('m' / 'M') ('o' / 'O') ('d' / 'D'))) | (&('D' | 'd') (('d' / 'D') ('i' / 'I') ('v' / 'V'))) | (&('S' | 's') (('s' / 'S') ('u' / 'U') ('b' / 'B'))) | (&('A' | 'a') (('a' / 'A') ('d' / 'D') ('d' / 'D'))))) Space)> */
		nil,
		/* 29 CMP_OP <- <(((('b' / 'B') ('e' / 'E')) / (('a' / 'A') ('e' / 'E')) / ((&('N' | 'n') (('n' / 'N') ('z' / 'Z'))) | (&('A' | 'a') ('a' / 'A')) | (&('Z') 'Z') | (&('z') 'z') | (&('B' | 'b') ('b' / 'B')))) Space)> */
		nil,
		/* 30 DATA_TYPE <- <(((&('I' | 'i') (('i' / 'I') ('n' / 'N') ('t' / 'T'))) | (&('F' | 'f') (('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))) | (&('B' | 'b') (('b' / 'B') ('y' / 'Y') ('t' / 'T') ('e' / 'E'))) | (&('W' | 'w') (('w' / 'W') ('o' / 'O') ('r' / 'R') ('d' / 'D'))) | (&('D' | 'd') (('d' / 'D') ('w' / 'W') ('o' / 'O') ('r' / 'R') ('d' / 'D')))) Space)> */
		func() bool {
			position359, tokenIndex359, depth359 := position, tokenIndex, depth
			{
				position360 := position
				depth++
				{
					switch buffer[position] {
					case 'I', 'i':
						{
							position362, tokenIndex362, depth362 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l363
							}
							position++
							goto l362
						l363:
							position, tokenIndex, depth = position362, tokenIndex362, depth362
							if buffer[position] != rune('I') {
								goto l359
							}
							position++
						}
					l362:
						{
							position364, tokenIndex364, depth364 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l365
							}
							position++
							goto l364
						l365:
							position, tokenIndex, depth = position364, tokenIndex364, depth364
							if buffer[position] != rune('N') {
								goto l359
							}
							position++
						}
					l364:
						{
							position366, tokenIndex366, depth366 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l367
							}
							position++
							goto l366
						l367:
							position, tokenIndex, depth = position366, tokenIndex366, depth366
							if buffer[position] != rune('T') {
								goto l359
							}
							position++
						}
					l366:
						break
					case 'F', 'f':
						{
							position368, tokenIndex368, depth368 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l369
							}
							position++
							goto l368
						l369:
							position, tokenIndex, depth = position368, tokenIndex368, depth368
							if buffer[position] != rune('F') {
								goto l359
							}
							position++
						}
					l368:
						{
							position370, tokenIndex370, depth370 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l371
							}
							position++
							goto l370
						l371:
							position, tokenIndex, depth = position370, tokenIndex370, depth370
							if buffer[position] != rune('L') {
								goto l359
							}
							position++
						}
					l370:
						{
							position372, tokenIndex372, depth372 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l373
							}
							position++
							goto l372
						l373:
							position, tokenIndex, depth = position372, tokenIndex372, depth372
							if buffer[position] != rune('O') {
								goto l359
							}
							position++
						}
					l372:
						{
							position374, tokenIndex374, depth374 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l375
							}
							position++
							goto l374
						l375:
							position, tokenIndex, depth = position374, tokenIndex374, depth374
							if buffer[position] != rune('A') {
								goto l359
							}
							position++
						}
					l374:
						{
							position376, tokenIndex376, depth376 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l377
							}
							position++
							goto l376
						l377:
							position, tokenIndex, depth = position376, tokenIndex376, depth376
							if buffer[position] != rune('T') {
								goto l359
							}
							position++
						}
					l376:
						break
					case 'B', 'b':
						{
							position378, tokenIndex378, depth378 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l379
							}
							position++
							goto l378
						l379:
							position, tokenIndex, depth = position378, tokenIndex378, depth378
							if buffer[position] != rune('B') {
								goto l359
							}
							position++
						}
					l378:
						{
							position380, tokenIndex380, depth380 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l381
							}
							position++
							goto l380
						l381:
							position, tokenIndex, depth = position380, tokenIndex380, depth380
							if buffer[position] != rune('Y') {
								goto l359
							}
							position++
						}
					l380:
						{
							position382, tokenIndex382, depth382 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l383
							}
							position++
							goto l382
						l383:
							position, tokenIndex, depth = position382, tokenIndex382, depth382
							if buffer[position] != rune('T') {
								goto l359
							}
							position++
						}
					l382:
						{
							position384, tokenIndex384, depth384 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l385
							}
							position++
							goto l384
						l385:
							position, tokenIndex, depth = position384, tokenIndex384, depth384
							if buffer[position] != rune('E') {
								goto l359
							}
							position++
						}
					l384:
						break
					case 'W', 'w':
						{
							position386, tokenIndex386, depth386 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l387
							}
							position++
							goto l386
						l387:
							position, tokenIndex, depth = position386, tokenIndex386, depth386
							if buffer[position] != rune('W') {
								goto l359
							}
							position++
						}
					l386:
						{
							position388, tokenIndex388, depth388 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l389
							}
							position++
							goto l388
						l389:
							position, tokenIndex, depth = position388, tokenIndex388, depth388
							if buffer[position] != rune('O') {
								goto l359
							}
							position++
						}
					l388:
						{
							position390, tokenIndex390, depth390 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l391
							}
							position++
							goto l390
						l391:
							position, tokenIndex, depth = position390, tokenIndex390, depth390
							if buffer[position] != rune('R') {
								goto l359
							}
							position++
						}
					l390:
						{
							position392, tokenIndex392, depth392 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l393
							}
							position++
							goto l392
						l393:
							position, tokenIndex, depth = position392, tokenIndex392, depth392
							if buffer[position] != rune('D') {
								goto l359
							}
							position++
						}
					l392:
						break
					default:
						{
							position394, tokenIndex394, depth394 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l395
							}
							position++
							goto l394
						l395:
							position, tokenIndex, depth = position394, tokenIndex394, depth394
							if buffer[position] != rune('D') {
								goto l359
							}
							position++
						}
					l394:
						{
							position396, tokenIndex396, depth396 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l397
							}
							position++
							goto l396
						l397:
							position, tokenIndex, depth = position396, tokenIndex396, depth396
							if buffer[position] != rune('W') {
								goto l359
							}
							position++
						}
					l396:
						{
							position398, tokenIndex398, depth398 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l399
							}
							position++
							goto l398
						l399:
							position, tokenIndex, depth = position398, tokenIndex398, depth398
							if buffer[position] != rune('O') {
								goto l359
							}
							position++
						}
					l398:
						{
							position400, tokenIndex400, depth400 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l401
							}
							position++
							goto l400
						l401:
							position, tokenIndex, depth = position400, tokenIndex400, depth400
							if buffer[position] != rune('R') {
								goto l359
							}
							position++
						}
					l400:
						{
							position402, tokenIndex402, depth402 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l403
							}
							position++
							goto l402
						l403:
							position, tokenIndex, depth = position402, tokenIndex402, depth402
							if buffer[position] != rune('D') {
								goto l359
							}
							position++
						}
					l402:
						break
					}
				}

				if !_rules[ruleSpace]() {
					goto l359
				}
				depth--
				add(ruleDATA_TYPE, position360)
			}
			return true
		l359:
			position, tokenIndex, depth = position359, tokenIndex359, depth359
			return false
		},
		/* 31 LBRK <- <('[' Spacing)> */
		func() bool {
			position404, tokenIndex404, depth404 := position, tokenIndex, depth
			{
				position405 := position
				depth++
				if buffer[position] != rune('[') {
					goto l404
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l404
				}
				depth--
				add(ruleLBRK, position405)
			}
			return true
		l404:
			position, tokenIndex, depth = position404, tokenIndex404, depth404
			return false
		},
		/* 32 RBRK <- <(']' Spacing)> */
		func() bool {
			position406, tokenIndex406, depth406 := position, tokenIndex, depth
			{
				position407 := position
				depth++
				if buffer[position] != rune(']') {
					goto l406
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l406
				}
				depth--
				add(ruleRBRK, position407)
			}
			return true
		l406:
			position, tokenIndex, depth = position406, tokenIndex406, depth406
			return false
		},
		/* 33 COMMA <- <(',' Spacing)> */
		func() bool {
			position408, tokenIndex408, depth408 := position, tokenIndex, depth
			{
				position409 := position
				depth++
				if buffer[position] != rune(',') {
					goto l408
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l408
				}
				depth--
				add(ruleCOMMA, position409)
			}
			return true
		l408:
			position, tokenIndex, depth = position408, tokenIndex408, depth408
			return false
		},
		/* 34 SEMICOLON <- <(';' Spacing)> */
		nil,
		/* 35 COLON <- <(':' Spacing)> */
		nil,
		/* 36 MINUS <- <('-' Spacing)> */
		nil,
		/* 37 NL <- <'\n'> */
		func() bool {
			position413, tokenIndex413, depth413 := position, tokenIndex, depth
			{
				position414 := position
				depth++
				if buffer[position] != rune('\n') {
					goto l413
				}
				position++
				depth--
				add(ruleNL, position414)
			}
			return true
		l413:
			position, tokenIndex, depth = position413, tokenIndex413, depth413
			return false
		},
		/* 38 EOT <- <!.> */
		nil,
		/* 39 Literal <- <((FloatLiteral / ((&('"') StringLiteral) | (&('\'') CharLiteral) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') IntegerLiteral))) Spacing)> */
		nil,
		/* 40 IntegerLiteral <- <(<(MINUS? (HexNumeral / BinaryNumeral / OctalNumeral / DecimalNumeral))> Spacing Action33)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				{
					position419 := position
					depth++
					{
						position420, tokenIndex420, depth420 := position, tokenIndex, depth
						{
							position422 := position
							depth++
							if buffer[position] != rune('-') {
								goto l420
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l420
							}
							depth--
							add(ruleMINUS, position422)
						}
						goto l421
					l420:
						position, tokenIndex, depth = position420, tokenIndex420, depth420
					}
				l421:
					{
						position423, tokenIndex423, depth423 := position, tokenIndex, depth
						if !_rules[ruleHexNumeral]() {
							goto l424
						}
						goto l423
					l424:
						position, tokenIndex, depth = position423, tokenIndex423, depth423
						{
							position426 := position
							depth++
							{
								position427, tokenIndex427, depth427 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l428
								}
								position++
								if buffer[position] != rune('b') {
									goto l428
								}
								position++
								goto l427
							l428:
								position, tokenIndex, depth = position427, tokenIndex427, depth427
								if buffer[position] != rune('0') {
									goto l425
								}
								position++
								if buffer[position] != rune('B') {
									goto l425
								}
								position++
							}
						l427:
							{
								position429, tokenIndex429, depth429 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l430
								}
								position++
								goto l429
							l430:
								position, tokenIndex, depth = position429, tokenIndex429, depth429
								if buffer[position] != rune('1') {
									goto l425
								}
								position++
							}
						l429:
						l431:
							{
								position432, tokenIndex432, depth432 := position, tokenIndex, depth
							l433:
								{
									position434, tokenIndex434, depth434 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l434
									}
									position++
									goto l433
								l434:
									position, tokenIndex, depth = position434, tokenIndex434, depth434
								}
								{
									position435, tokenIndex435, depth435 := position, tokenIndex, depth
									if buffer[position] != rune('0') {
										goto l436
									}
									position++
									goto l435
								l436:
									position, tokenIndex, depth = position435, tokenIndex435, depth435
									if buffer[position] != rune('1') {
										goto l432
									}
									position++
								}
							l435:
								goto l431
							l432:
								position, tokenIndex, depth = position432, tokenIndex432, depth432
							}
							depth--
							add(ruleBinaryNumeral, position426)
						}
						goto l423
					l425:
						position, tokenIndex, depth = position423, tokenIndex423, depth423
						{
							position438 := position
							depth++
							if buffer[position] != rune('0') {
								goto l437
							}
							position++
						l441:
							{
								position442, tokenIndex442, depth442 := position, tokenIndex, depth
								if buffer[position] != rune('_') {
									goto l442
								}
								position++
								goto l441
							l442:
								position, tokenIndex, depth = position442, tokenIndex442, depth442
							}
							if c := buffer[position]; c < rune('0') || c > rune('7') {
								goto l437
							}
							position++
						l439:
							{
								position440, tokenIndex440, depth440 := position, tokenIndex, depth
							l443:
								{
									position444, tokenIndex444, depth444 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l444
									}
									position++
									goto l443
								l444:
									position, tokenIndex, depth = position444, tokenIndex444, depth444
								}
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l440
								}
								position++
								goto l439
							l440:
								position, tokenIndex, depth = position440, tokenIndex440, depth440
							}
							depth--
							add(ruleOctalNumeral, position438)
						}
						goto l423
					l437:
						position, tokenIndex, depth = position423, tokenIndex423, depth423
						{
							position445 := position
							depth++
							{
								position446, tokenIndex446, depth446 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l447
								}
								position++
								goto l446
							l447:
								position, tokenIndex, depth = position446, tokenIndex446, depth446
								if c := buffer[position]; c < rune('1') || c > rune('9') {
									goto l417
								}
								position++
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
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l449
									}
									position++
									goto l448
								l449:
									position, tokenIndex, depth = position449, tokenIndex449, depth449
								}
							}
						l446:
							depth--
							add(ruleDecimalNumeral, position445)
						}
					}
				l423:
					depth--
					add(rulePegText, position419)
				}
				if !_rules[ruleSpacing]() {
					goto l417
				}
				{
					add(ruleAction33, position)
				}
				depth--
				add(ruleIntegerLiteral, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 41 DecimalNumeral <- <('0' / ([1-9] ('_'* [0-9])*))> */
		nil,
		/* 42 HexNumeral <- <((('0' 'x') / ('0' 'X')) HexDigits)> */
		func() bool {
			position454, tokenIndex454, depth454 := position, tokenIndex, depth
			{
				position455 := position
				depth++
				{
					position456, tokenIndex456, depth456 := position, tokenIndex, depth
					if buffer[position] != rune('0') {
						goto l457
					}
					position++
					if buffer[position] != rune('x') {
						goto l457
					}
					position++
					goto l456
				l457:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
					if buffer[position] != rune('0') {
						goto l454
					}
					position++
					if buffer[position] != rune('X') {
						goto l454
					}
					position++
				}
			l456:
				if !_rules[ruleHexDigits]() {
					goto l454
				}
				depth--
				add(ruleHexNumeral, position455)
			}
			return true
		l454:
			position, tokenIndex, depth = position454, tokenIndex454, depth454
			return false
		},
		/* 43 BinaryNumeral <- <((('0' 'b') / ('0' 'B')) ('0' / '1') ('_'* ('0' / '1'))*)> */
		nil,
		/* 44 OctalNumeral <- <('0' ('_'* [0-7])+)> */
		nil,
		/* 45 FloatLiteral <- <(HexFloat / DecimalFloat)> */
		nil,
		/* 46 DecimalFloat <- <((Digits '.' Digits? Exponent? ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?) / ('.' Digits Exponent? ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?) / (Digits Exponent ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?) / (Digits Exponent? ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))))> */
		nil,
		/* 47 Exponent <- <(('e' / 'E') ('+' / '-')? Digits)> */
		func() bool {
			position462, tokenIndex462, depth462 := position, tokenIndex, depth
			{
				position463 := position
				depth++
				{
					position464, tokenIndex464, depth464 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l465
					}
					position++
					goto l464
				l465:
					position, tokenIndex, depth = position464, tokenIndex464, depth464
					if buffer[position] != rune('E') {
						goto l462
					}
					position++
				}
			l464:
				{
					position466, tokenIndex466, depth466 := position, tokenIndex, depth
					{
						position468, tokenIndex468, depth468 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l469
						}
						position++
						goto l468
					l469:
						position, tokenIndex, depth = position468, tokenIndex468, depth468
						if buffer[position] != rune('-') {
							goto l466
						}
						position++
					}
				l468:
					goto l467
				l466:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
				}
			l467:
				if !_rules[ruleDigits]() {
					goto l462
				}
				depth--
				add(ruleExponent, position463)
			}
			return true
		l462:
			position, tokenIndex, depth = position462, tokenIndex462, depth462
			return false
		},
		/* 48 HexFloat <- <(HexSignificand BinaryExponent ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?)> */
		nil,
		/* 49 HexSignificand <- <(((('0' 'x') / ('0' 'X')) HexDigits? '.' HexDigits) / (HexNumeral '.'?))> */
		nil,
		/* 50 BinaryExponent <- <(('p' / 'P') ('+' / '-')? Digits)> */
		nil,
		/* 51 Digits <- <([0-9] ('_'* [0-9])*)> */
		func() bool {
			position473, tokenIndex473, depth473 := position, tokenIndex, depth
			{
				position474 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l473
				}
				position++
			l475:
				{
					position476, tokenIndex476, depth476 := position, tokenIndex, depth
				l477:
					{
						position478, tokenIndex478, depth478 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l478
						}
						position++
						goto l477
					l478:
						position, tokenIndex, depth = position478, tokenIndex478, depth478
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l476
					}
					position++
					goto l475
				l476:
					position, tokenIndex, depth = position476, tokenIndex476, depth476
				}
				depth--
				add(ruleDigits, position474)
			}
			return true
		l473:
			position, tokenIndex, depth = position473, tokenIndex473, depth473
			return false
		},
		/* 52 HexDigits <- <(HexDigit ('_'* HexDigit)*)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				if !_rules[ruleHexDigit]() {
					goto l479
				}
			l481:
				{
					position482, tokenIndex482, depth482 := position, tokenIndex, depth
				l483:
					{
						position484, tokenIndex484, depth484 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l484
						}
						position++
						goto l483
					l484:
						position, tokenIndex, depth = position484, tokenIndex484, depth484
					}
					if !_rules[ruleHexDigit]() {
						goto l482
					}
					goto l481
				l482:
					position, tokenIndex, depth = position482, tokenIndex482, depth482
				}
				depth--
				add(ruleHexDigits, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 53 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l485
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l485
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l485
						}
						position++
						break
					}
				}

				depth--
				add(ruleHexDigit, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 54 CharLiteral <- <('\'' (Escape / (!('\'' / '\\') .)) '\'')> */
		nil,
		/* 55 StringLiteral <- <(<('"' (Escape / (!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('"') '"')) .))* '"')> Action34)> */
		func() bool {
			position489, tokenIndex489, depth489 := position, tokenIndex, depth
			{
				position490 := position
				depth++
				{
					position491 := position
					depth++
					if buffer[position] != rune('"') {
						goto l489
					}
					position++
				l492:
					{
						position493, tokenIndex493, depth493 := position, tokenIndex, depth
						{
							position494, tokenIndex494, depth494 := position, tokenIndex, depth
							if !_rules[ruleEscape]() {
								goto l495
							}
							goto l494
						l495:
							position, tokenIndex, depth = position494, tokenIndex494, depth494
							{
								position496, tokenIndex496, depth496 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '\r':
										if buffer[position] != rune('\r') {
											goto l496
										}
										position++
										break
									case '\n':
										if buffer[position] != rune('\n') {
											goto l496
										}
										position++
										break
									case '\\':
										if buffer[position] != rune('\\') {
											goto l496
										}
										position++
										break
									default:
										if buffer[position] != rune('"') {
											goto l496
										}
										position++
										break
									}
								}

								goto l493
							l496:
								position, tokenIndex, depth = position496, tokenIndex496, depth496
							}
							if !matchDot() {
								goto l493
							}
						}
					l494:
						goto l492
					l493:
						position, tokenIndex, depth = position493, tokenIndex493, depth493
					}
					if buffer[position] != rune('"') {
						goto l489
					}
					position++
					depth--
					add(rulePegText, position491)
				}
				{
					add(ruleAction34, position)
				}
				depth--
				add(ruleStringLiteral, position490)
			}
			return true
		l489:
			position, tokenIndex, depth = position489, tokenIndex489, depth489
			return false
		},
		/* 56 Escape <- <('\\' ((&('u') UnicodeEscape) | (&('\\') '\\') | (&('\'') '\'') | (&('"') '"') | (&('r') 'r') | (&('f') 'f') | (&('n') 'n') | (&('t') 't') | (&('b') 'b') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7') OctalEscape)))> */
		func() bool {
			position499, tokenIndex499, depth499 := position, tokenIndex, depth
			{
				position500 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l499
				}
				position++
				{
					switch buffer[position] {
					case 'u':
						{
							position502 := position
							depth++
							if buffer[position] != rune('u') {
								goto l499
							}
							position++
						l503:
							{
								position504, tokenIndex504, depth504 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l504
								}
								position++
								goto l503
							l504:
								position, tokenIndex, depth = position504, tokenIndex504, depth504
							}
							if !_rules[ruleHexDigit]() {
								goto l499
							}
							if !_rules[ruleHexDigit]() {
								goto l499
							}
							if !_rules[ruleHexDigit]() {
								goto l499
							}
							if !_rules[ruleHexDigit]() {
								goto l499
							}
							depth--
							add(ruleUnicodeEscape, position502)
						}
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l499
						}
						position++
						break
					case '\'':
						if buffer[position] != rune('\'') {
							goto l499
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l499
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l499
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l499
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l499
						}
						position++
						break
					case 't':
						if buffer[position] != rune('t') {
							goto l499
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l499
						}
						position++
						break
					default:
						{
							position505 := position
							depth++
							{
								position506, tokenIndex506, depth506 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('3') {
									goto l507
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l507
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l507
								}
								position++
								goto l506
							l507:
								position, tokenIndex, depth = position506, tokenIndex506, depth506
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l508
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l508
								}
								position++
								goto l506
							l508:
								position, tokenIndex, depth = position506, tokenIndex506, depth506
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l499
								}
								position++
							}
						l506:
							depth--
							add(ruleOctalEscape, position505)
						}
						break
					}
				}

				depth--
				add(ruleEscape, position500)
			}
			return true
		l499:
			position, tokenIndex, depth = position499, tokenIndex499, depth499
			return false
		},
		/* 57 OctalEscape <- <(([0-3] [0-7] [0-7]) / ([0-7] [0-7]) / [0-7])> */
		nil,
		/* 58 UnicodeEscape <- <('u'+ HexDigit HexDigit HexDigit HexDigit)> */
		nil,
		/* 60 Action0 <- <{p.line++}> */
		nil,
		/* 61 Action1 <- <{p.AddAssembly()}> */
		nil,
		/* 62 Action2 <- <{p.AddAssembly();p.AddComment()}> */
		nil,
		nil,
		/* 64 Action3 <- <{p.Push(&asm.Comment{});p.Push(text)}> */
		nil,
		/* 65 Action4 <- <{p.Push(&asm.Label{})}> */
		nil,
		/* 66 Action5 <- <{p.Push(lookup(asm.T_INT,text))}> */
		nil,
		/* 67 Action6 <- <{p.Push(lookup(asm.CAL_ADD,text))}> */
		nil,
		/* 68 Action7 <- <{p.Push(lookup(asm.T_INT,text))}> */
		nil,
		/* 69 Action8 <- <{p.Push(lookup(asm.T_INT,text))}> */
		nil,
		/* 70 Action9 <- <{p.Push(lookup(asm.CMP_A,text))}> */
		nil,
		/* 71 Action10 <- <{p.AddPseudoDataValue()}> */
		nil,
		/* 72 Action11 <- <{p.AddPseudoDataValue()}> */
		nil,
		/* 73 Action12 <- <{p.AddPseudoDataValue()}> */
		nil,
		/* 74 Action13 <- <{p.AddOperand(true)}> */
		nil,
		/* 75 Action14 <- <{p.AddOperand(false)}> */
		nil,
		/* 76 Action15 <- <{p.AddOperand(true)}> */
		nil,
		/* 77 Action16 <- <{p.AddOperand(false)}> */
		nil,
		/* 78 Action17 <- <{p.Push(text)}> */
		nil,
		/* 79 Action18 <- <{p.PushInst(asm.OP_EXIT)}> */
		nil,
		/* 80 Action19 <- <{p.PushInst(asm.OP_RET)}> */
		nil,
		/* 81 Action20 <- <{p.PushInst(asm.OP_NOP)}> */
		nil,
		/* 82 Action21 <- <{p.PushInst(asm.OP_CALL)}> */
		nil,
		/* 83 Action22 <- <{p.PushInst(asm.OP_PUSH)}> */
		nil,
		/* 84 Action23 <- <{p.PushInst(asm.OP_POP)}> */
		nil,
		/* 85 Action24 <- <{p.PushInst(asm.OP_JMP)}> */
		nil,
		/* 86 Action25 <- <{p.PushInst(asm.OP_IN)}> */
		nil,
		/* 87 Action26 <- <{p.PushInst(asm.OP_OUT)}> */
		nil,
		/* 88 Action27 <- <{p.PushInst(asm.OP_CAL)}> */
		nil,
		/* 89 Action28 <- <{p.PushInst(asm.OP_LD)}> */
		nil,
		/* 90 Action29 <- <{p.PushInst(asm.OP_CMP)}> */
		nil,
		/* 91 Action30 <- <{p.PushInst(asm.OP_JPC)}> */
		nil,
		/* 92 Action31 <- <{p.Push(&asm.PseudoBlock{})}> */
		nil,
		/* 93 Action32 <- <{p.Push(&asm.PseudoData{})}> */
		nil,
		/* 94 Action33 <- <{p.Push(text);p.AddInteger()}> */
		nil,
		/* 95 Action34 <- <{p.Push(text)}> */
		nil,
	}
	p.rules = _rules
}
