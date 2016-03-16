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
		/* 0 Start <- <((Spacing Line? NL Action0)* EOT)> */
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
																	goto l142
																l143:
																	position, tokenIndex, depth = position142, tokenIndex142, depth142
																	{
																		position147, tokenIndex147, depth147 := position, tokenIndex, depth
																		if buffer[position] != rune('a') {
																			goto l148
																		}
																		position++
																		goto l147
																	l148:
																		position, tokenIndex, depth = position147, tokenIndex147, depth147
																		if buffer[position] != rune('A') {
																			goto l146
																		}
																		position++
																	}
																l147:
																	goto l142
																l146:
																	position, tokenIndex, depth = position142, tokenIndex142, depth142
																	{
																		switch buffer[position] {
																		case 'N', 'n':
																			{
																				position150, tokenIndex150, depth150 := position, tokenIndex, depth
																				if buffer[position] != rune('n') {
																					goto l151
																				}
																				position++
																				goto l150
																			l151:
																				position, tokenIndex, depth = position150, tokenIndex150, depth150
																				if buffer[position] != rune('N') {
																					goto l4
																				}
																				position++
																			}
																		l150:
																			{
																				position152, tokenIndex152, depth152 := position, tokenIndex, depth
																				if buffer[position] != rune('z') {
																					goto l153
																				}
																				position++
																				goto l152
																			l153:
																				position, tokenIndex, depth = position152, tokenIndex152, depth152
																				if buffer[position] != rune('Z') {
																					goto l4
																				}
																				position++
																			}
																		l152:
																			break
																		case 'A', 'a':
																			{
																				position154, tokenIndex154, depth154 := position, tokenIndex, depth
																				if buffer[position] != rune('a') {
																					goto l155
																				}
																				position++
																				goto l154
																			l155:
																				position, tokenIndex, depth = position154, tokenIndex154, depth154
																				if buffer[position] != rune('A') {
																					goto l4
																				}
																				position++
																			}
																		l154:
																			{
																				position156, tokenIndex156, depth156 := position, tokenIndex, depth
																				if buffer[position] != rune('e') {
																					goto l157
																				}
																				position++
																				goto l156
																			l157:
																				position, tokenIndex, depth = position156, tokenIndex156, depth156
																				if buffer[position] != rune('E') {
																					goto l4
																				}
																				position++
																			}
																		l156:
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
																				position158, tokenIndex158, depth158 := position, tokenIndex, depth
																				if buffer[position] != rune('b') {
																					goto l159
																				}
																				position++
																				goto l158
																			l159:
																				position, tokenIndex, depth = position158, tokenIndex158, depth158
																				if buffer[position] != rune('B') {
																					goto l4
																				}
																				position++
																			}
																		l158:
																			{
																				position160, tokenIndex160, depth160 := position, tokenIndex, depth
																				if buffer[position] != rune('e') {
																					goto l161
																				}
																				position++
																				goto l160
																			l161:
																				position, tokenIndex, depth = position160, tokenIndex160, depth160
																				if buffer[position] != rune('E') {
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
			position231, tokenIndex231, depth231 := position, tokenIndex, depth
			{
				position232 := position
				depth++
				{
					position233 := position
					depth++
					if buffer[position] != rune(';') {
						goto l231
					}
					position++
					if !_rules[ruleSpacing]() {
						goto l231
					}
					depth--
					add(ruleSEMICOLON, position233)
				}
				{
					position234 := position
					depth++
				l235:
					{
						position236, tokenIndex236, depth236 := position, tokenIndex, depth
						{
							position237, tokenIndex237, depth237 := position, tokenIndex, depth
							if !_rules[ruleNL]() {
								goto l237
							}
							goto l236
						l237:
							position, tokenIndex, depth = position237, tokenIndex237, depth237
						}
						if !matchDot() {
							goto l236
						}
						goto l235
					l236:
						position, tokenIndex, depth = position236, tokenIndex236, depth236
					}
					depth--
					add(rulePegText, position234)
				}
				{
					add(ruleAction3, position)
				}
				depth--
				add(ruleComment, position232)
			}
			return true
		l231:
			position, tokenIndex, depth = position231, tokenIndex231, depth231
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
			position242, tokenIndex242, depth242 := position, tokenIndex, depth
			{
				position243 := position
				depth++
				{
					switch buffer[position] {
					case '"':
						if !_rules[ruleStringLiteral]() {
							goto l242
						}
						{
							add(ruleAction12, position)
						}
						break
					case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if !_rules[ruleIntegerLiteral]() {
							goto l242
						}
						{
							add(ruleAction10, position)
						}
						break
					default:
						if !_rules[ruleIdentifier]() {
							goto l242
						}
						{
							add(ruleAction11, position)
						}
						break
					}
				}

				depth--
				add(rulePseudoDataValue, position243)
			}
			return true
		l242:
			position, tokenIndex, depth = position242, tokenIndex242, depth242
			return false
		},
		/* 7 Operand <- <(((LBRK Identifier RBRK Action14) / (IntegerLiteral Action15) / ((&('"' | '\'' | '-' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Literal) | (&('[') (LBRK IntegerLiteral RBRK Action16)) | (&('$' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (Identifier Action13)))) Spacing)> */
		func() bool {
			position248, tokenIndex248, depth248 := position, tokenIndex, depth
			{
				position249 := position
				depth++
				{
					position250, tokenIndex250, depth250 := position, tokenIndex, depth
					if !_rules[ruleLBRK]() {
						goto l251
					}
					if !_rules[ruleIdentifier]() {
						goto l251
					}
					if !_rules[ruleRBRK]() {
						goto l251
					}
					{
						add(ruleAction14, position)
					}
					goto l250
				l251:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					if !_rules[ruleIntegerLiteral]() {
						goto l253
					}
					{
						add(ruleAction15, position)
					}
					goto l250
				l253:
					position, tokenIndex, depth = position250, tokenIndex250, depth250
					{
						switch buffer[position] {
						case '"', '\'', '-', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							{
								position256 := position
								depth++
								{
									position257, tokenIndex257, depth257 := position, tokenIndex, depth
									{
										position259 := position
										depth++
										{
											position260, tokenIndex260, depth260 := position, tokenIndex, depth
											{
												position262 := position
												depth++
												{
													position263 := position
													depth++
													{
														position264, tokenIndex264, depth264 := position, tokenIndex, depth
														{
															position266, tokenIndex266, depth266 := position, tokenIndex, depth
															if buffer[position] != rune('0') {
																goto l267
															}
															position++
															if buffer[position] != rune('x') {
																goto l267
															}
															position++
															goto l266
														l267:
															position, tokenIndex, depth = position266, tokenIndex266, depth266
															if buffer[position] != rune('0') {
																goto l265
															}
															position++
															if buffer[position] != rune('X') {
																goto l265
															}
															position++
														}
													l266:
														{
															position268, tokenIndex268, depth268 := position, tokenIndex, depth
															if !_rules[ruleHexDigits]() {
																goto l268
															}
															goto l269
														l268:
															position, tokenIndex, depth = position268, tokenIndex268, depth268
														}
													l269:
														if buffer[position] != rune('.') {
															goto l265
														}
														position++
														if !_rules[ruleHexDigits]() {
															goto l265
														}
														goto l264
													l265:
														position, tokenIndex, depth = position264, tokenIndex264, depth264
														if !_rules[ruleHexNumeral]() {
															goto l261
														}
														{
															position270, tokenIndex270, depth270 := position, tokenIndex, depth
															if buffer[position] != rune('.') {
																goto l270
															}
															position++
															goto l271
														l270:
															position, tokenIndex, depth = position270, tokenIndex270, depth270
														}
													l271:
													}
												l264:
													depth--
													add(ruleHexSignificand, position263)
												}
												{
													position272 := position
													depth++
													{
														position273, tokenIndex273, depth273 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l274
														}
														position++
														goto l273
													l274:
														position, tokenIndex, depth = position273, tokenIndex273, depth273
														if buffer[position] != rune('P') {
															goto l261
														}
														position++
													}
												l273:
													{
														position275, tokenIndex275, depth275 := position, tokenIndex, depth
														{
															position277, tokenIndex277, depth277 := position, tokenIndex, depth
															if buffer[position] != rune('+') {
																goto l278
															}
															position++
															goto l277
														l278:
															position, tokenIndex, depth = position277, tokenIndex277, depth277
															if buffer[position] != rune('-') {
																goto l275
															}
															position++
														}
													l277:
														goto l276
													l275:
														position, tokenIndex, depth = position275, tokenIndex275, depth275
													}
												l276:
													if !_rules[ruleDigits]() {
														goto l261
													}
													depth--
													add(ruleBinaryExponent, position272)
												}
												{
													position279, tokenIndex279, depth279 := position, tokenIndex, depth
													{
														switch buffer[position] {
														case 'D':
															if buffer[position] != rune('D') {
																goto l279
															}
															position++
															break
														case 'd':
															if buffer[position] != rune('d') {
																goto l279
															}
															position++
															break
														case 'F':
															if buffer[position] != rune('F') {
																goto l279
															}
															position++
															break
														default:
															if buffer[position] != rune('f') {
																goto l279
															}
															position++
															break
														}
													}

													goto l280
												l279:
													position, tokenIndex, depth = position279, tokenIndex279, depth279
												}
											l280:
												depth--
												add(ruleHexFloat, position262)
											}
											goto l260
										l261:
											position, tokenIndex, depth = position260, tokenIndex260, depth260
											{
												position282 := position
												depth++
												{
													position283, tokenIndex283, depth283 := position, tokenIndex, depth
													if !_rules[ruleDigits]() {
														goto l284
													}
													if buffer[position] != rune('.') {
														goto l284
													}
													position++
													{
														position285, tokenIndex285, depth285 := position, tokenIndex, depth
														if !_rules[ruleDigits]() {
															goto l285
														}
														goto l286
													l285:
														position, tokenIndex, depth = position285, tokenIndex285, depth285
													}
												l286:
													{
														position287, tokenIndex287, depth287 := position, tokenIndex, depth
														if !_rules[ruleExponent]() {
															goto l287
														}
														goto l288
													l287:
														position, tokenIndex, depth = position287, tokenIndex287, depth287
													}
												l288:
													{
														position289, tokenIndex289, depth289 := position, tokenIndex, depth
														{
															switch buffer[position] {
															case 'D':
																if buffer[position] != rune('D') {
																	goto l289
																}
																position++
																break
															case 'd':
																if buffer[position] != rune('d') {
																	goto l289
																}
																position++
																break
															case 'F':
																if buffer[position] != rune('F') {
																	goto l289
																}
																position++
																break
															default:
																if buffer[position] != rune('f') {
																	goto l289
																}
																position++
																break
															}
														}

														goto l290
													l289:
														position, tokenIndex, depth = position289, tokenIndex289, depth289
													}
												l290:
													goto l283
												l284:
													position, tokenIndex, depth = position283, tokenIndex283, depth283
													if buffer[position] != rune('.') {
														goto l292
													}
													position++
													if !_rules[ruleDigits]() {
														goto l292
													}
													{
														position293, tokenIndex293, depth293 := position, tokenIndex, depth
														if !_rules[ruleExponent]() {
															goto l293
														}
														goto l294
													l293:
														position, tokenIndex, depth = position293, tokenIndex293, depth293
													}
												l294:
													{
														position295, tokenIndex295, depth295 := position, tokenIndex, depth
														{
															switch buffer[position] {
															case 'D':
																if buffer[position] != rune('D') {
																	goto l295
																}
																position++
																break
															case 'd':
																if buffer[position] != rune('d') {
																	goto l295
																}
																position++
																break
															case 'F':
																if buffer[position] != rune('F') {
																	goto l295
																}
																position++
																break
															default:
																if buffer[position] != rune('f') {
																	goto l295
																}
																position++
																break
															}
														}

														goto l296
													l295:
														position, tokenIndex, depth = position295, tokenIndex295, depth295
													}
												l296:
													goto l283
												l292:
													position, tokenIndex, depth = position283, tokenIndex283, depth283
													if !_rules[ruleDigits]() {
														goto l298
													}
													if !_rules[ruleExponent]() {
														goto l298
													}
													{
														position299, tokenIndex299, depth299 := position, tokenIndex, depth
														{
															switch buffer[position] {
															case 'D':
																if buffer[position] != rune('D') {
																	goto l299
																}
																position++
																break
															case 'd':
																if buffer[position] != rune('d') {
																	goto l299
																}
																position++
																break
															case 'F':
																if buffer[position] != rune('F') {
																	goto l299
																}
																position++
																break
															default:
																if buffer[position] != rune('f') {
																	goto l299
																}
																position++
																break
															}
														}

														goto l300
													l299:
														position, tokenIndex, depth = position299, tokenIndex299, depth299
													}
												l300:
													goto l283
												l298:
													position, tokenIndex, depth = position283, tokenIndex283, depth283
													if !_rules[ruleDigits]() {
														goto l258
													}
													{
														position302, tokenIndex302, depth302 := position, tokenIndex, depth
														if !_rules[ruleExponent]() {
															goto l302
														}
														goto l303
													l302:
														position, tokenIndex, depth = position302, tokenIndex302, depth302
													}
												l303:
													{
														switch buffer[position] {
														case 'D':
															if buffer[position] != rune('D') {
																goto l258
															}
															position++
															break
														case 'd':
															if buffer[position] != rune('d') {
																goto l258
															}
															position++
															break
														case 'F':
															if buffer[position] != rune('F') {
																goto l258
															}
															position++
															break
														default:
															if buffer[position] != rune('f') {
																goto l258
															}
															position++
															break
														}
													}

												}
											l283:
												depth--
												add(ruleDecimalFloat, position282)
											}
										}
									l260:
										depth--
										add(ruleFloatLiteral, position259)
									}
									goto l257
								l258:
									position, tokenIndex, depth = position257, tokenIndex257, depth257
									{
										switch buffer[position] {
										case '"':
											if !_rules[ruleStringLiteral]() {
												goto l248
											}
											break
										case '\'':
											{
												position306 := position
												depth++
												if buffer[position] != rune('\'') {
													goto l248
												}
												position++
												{
													position307, tokenIndex307, depth307 := position, tokenIndex, depth
													if !_rules[ruleEscape]() {
														goto l308
													}
													goto l307
												l308:
													position, tokenIndex, depth = position307, tokenIndex307, depth307
													{
														position309, tokenIndex309, depth309 := position, tokenIndex, depth
														{
															position310, tokenIndex310, depth310 := position, tokenIndex, depth
															if buffer[position] != rune('\'') {
																goto l311
															}
															position++
															goto l310
														l311:
															position, tokenIndex, depth = position310, tokenIndex310, depth310
															if buffer[position] != rune('\\') {
																goto l309
															}
															position++
														}
													l310:
														goto l248
													l309:
														position, tokenIndex, depth = position309, tokenIndex309, depth309
													}
													if !matchDot() {
														goto l248
													}
												}
											l307:
												if buffer[position] != rune('\'') {
													goto l248
												}
												position++
												depth--
												add(ruleCharLiteral, position306)
											}
											break
										default:
											if !_rules[ruleIntegerLiteral]() {
												goto l248
											}
											break
										}
									}

								}
							l257:
								if !_rules[ruleSpacing]() {
									goto l248
								}
								depth--
								add(ruleLiteral, position256)
							}
							break
						case '[':
							if !_rules[ruleLBRK]() {
								goto l248
							}
							if !_rules[ruleIntegerLiteral]() {
								goto l248
							}
							if !_rules[ruleRBRK]() {
								goto l248
							}
							{
								add(ruleAction16, position)
							}
							break
						default:
							if !_rules[ruleIdentifier]() {
								goto l248
							}
							{
								add(ruleAction13, position)
							}
							break
						}
					}

				}
			l250:
				if !_rules[ruleSpacing]() {
					goto l248
				}
				depth--
				add(ruleOperand, position249)
			}
			return true
		l248:
			position, tokenIndex, depth = position248, tokenIndex248, depth248
			return false
		},
		/* 8 Spacing <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position315 := position
				depth++
			l316:
				{
					position317, tokenIndex317, depth317 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l317
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l317
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l317
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l317
							}
							position++
							break
						}
					}

					goto l316
				l317:
					position, tokenIndex, depth = position317, tokenIndex317, depth317
				}
				depth--
				add(ruleSpacing, position315)
			}
			return true
		},
		/* 9 Space <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))+> */
		func() bool {
			position319, tokenIndex319, depth319 := position, tokenIndex, depth
			{
				position320 := position
				depth++
				{
					switch buffer[position] {
					case '\f':
						if buffer[position] != rune('\f') {
							goto l319
						}
						position++
						break
					case '\r':
						if buffer[position] != rune('\r') {
							goto l319
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l319
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l319
						}
						position++
						break
					}
				}

			l321:
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l322
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l322
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l322
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l322
							}
							position++
							break
						}
					}

					goto l321
				l322:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
				}
				depth--
				add(ruleSpace, position320)
			}
			return true
		l319:
			position, tokenIndex, depth = position319, tokenIndex319, depth319
			return false
		},
		/* 10 Identifier <- <(<(Letter LetterOrDigit*)> Spacing Action17)> */
		func() bool {
			position325, tokenIndex325, depth325 := position, tokenIndex, depth
			{
				position326 := position
				depth++
				{
					position327 := position
					depth++
					{
						position328 := position
						depth++
						{
							switch buffer[position] {
							case '$', '_':
								{
									position330, tokenIndex330, depth330 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l331
									}
									position++
									goto l330
								l331:
									position, tokenIndex, depth = position330, tokenIndex330, depth330
									if buffer[position] != rune('$') {
										goto l325
									}
									position++
								}
							l330:
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l325
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l325
								}
								position++
								break
							}
						}

						depth--
						add(ruleLetter, position328)
					}
				l332:
					{
						position333, tokenIndex333, depth333 := position, tokenIndex, depth
						{
							position334 := position
							depth++
							{
								switch buffer[position] {
								case '$', '_':
									{
										position336, tokenIndex336, depth336 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l337
										}
										position++
										goto l336
									l337:
										position, tokenIndex, depth = position336, tokenIndex336, depth336
										if buffer[position] != rune('$') {
											goto l333
										}
										position++
									}
								l336:
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l333
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l333
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l333
									}
									position++
									break
								}
							}

							depth--
							add(ruleLetterOrDigit, position334)
						}
						goto l332
					l333:
						position, tokenIndex, depth = position333, tokenIndex333, depth333
					}
					depth--
					add(rulePegText, position327)
				}
				if !_rules[ruleSpacing]() {
					goto l325
				}
				{
					add(ruleAction17, position)
				}
				depth--
				add(ruleIdentifier, position326)
			}
			return true
		l325:
			position, tokenIndex, depth = position325, tokenIndex325, depth325
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
		/* 29 CMP_OP <- <((('b' / 'B') / ('a' / 'A') / ((&('N' | 'n') (('n' / 'N') ('z' / 'Z'))) | (&('A' | 'a') (('a' / 'A') ('e' / 'E'))) | (&('Z') 'Z') | (&('z') 'z') | (&('B' | 'b') (('b' / 'B') ('e' / 'E'))))) Space)> */
		nil,
		/* 30 DATA_TYPE <- <(((&('I' | 'i') (('i' / 'I') ('n' / 'N') ('t' / 'T'))) | (&('F' | 'f') (('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))) | (&('B' | 'b') (('b' / 'B') ('y' / 'Y') ('t' / 'T') ('e' / 'E'))) | (&('W' | 'w') (('w' / 'W') ('o' / 'O') ('r' / 'R') ('d' / 'D'))) | (&('D' | 'd') (('d' / 'D') ('w' / 'W') ('o' / 'O') ('r' / 'R') ('d' / 'D')))) Space)> */
		func() bool {
			position358, tokenIndex358, depth358 := position, tokenIndex, depth
			{
				position359 := position
				depth++
				{
					switch buffer[position] {
					case 'I', 'i':
						{
							position361, tokenIndex361, depth361 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l362
							}
							position++
							goto l361
						l362:
							position, tokenIndex, depth = position361, tokenIndex361, depth361
							if buffer[position] != rune('I') {
								goto l358
							}
							position++
						}
					l361:
						{
							position363, tokenIndex363, depth363 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l364
							}
							position++
							goto l363
						l364:
							position, tokenIndex, depth = position363, tokenIndex363, depth363
							if buffer[position] != rune('N') {
								goto l358
							}
							position++
						}
					l363:
						{
							position365, tokenIndex365, depth365 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l366
							}
							position++
							goto l365
						l366:
							position, tokenIndex, depth = position365, tokenIndex365, depth365
							if buffer[position] != rune('T') {
								goto l358
							}
							position++
						}
					l365:
						break
					case 'F', 'f':
						{
							position367, tokenIndex367, depth367 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l368
							}
							position++
							goto l367
						l368:
							position, tokenIndex, depth = position367, tokenIndex367, depth367
							if buffer[position] != rune('F') {
								goto l358
							}
							position++
						}
					l367:
						{
							position369, tokenIndex369, depth369 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l370
							}
							position++
							goto l369
						l370:
							position, tokenIndex, depth = position369, tokenIndex369, depth369
							if buffer[position] != rune('L') {
								goto l358
							}
							position++
						}
					l369:
						{
							position371, tokenIndex371, depth371 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l372
							}
							position++
							goto l371
						l372:
							position, tokenIndex, depth = position371, tokenIndex371, depth371
							if buffer[position] != rune('O') {
								goto l358
							}
							position++
						}
					l371:
						{
							position373, tokenIndex373, depth373 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l374
							}
							position++
							goto l373
						l374:
							position, tokenIndex, depth = position373, tokenIndex373, depth373
							if buffer[position] != rune('A') {
								goto l358
							}
							position++
						}
					l373:
						{
							position375, tokenIndex375, depth375 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l376
							}
							position++
							goto l375
						l376:
							position, tokenIndex, depth = position375, tokenIndex375, depth375
							if buffer[position] != rune('T') {
								goto l358
							}
							position++
						}
					l375:
						break
					case 'B', 'b':
						{
							position377, tokenIndex377, depth377 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l378
							}
							position++
							goto l377
						l378:
							position, tokenIndex, depth = position377, tokenIndex377, depth377
							if buffer[position] != rune('B') {
								goto l358
							}
							position++
						}
					l377:
						{
							position379, tokenIndex379, depth379 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l380
							}
							position++
							goto l379
						l380:
							position, tokenIndex, depth = position379, tokenIndex379, depth379
							if buffer[position] != rune('Y') {
								goto l358
							}
							position++
						}
					l379:
						{
							position381, tokenIndex381, depth381 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l382
							}
							position++
							goto l381
						l382:
							position, tokenIndex, depth = position381, tokenIndex381, depth381
							if buffer[position] != rune('T') {
								goto l358
							}
							position++
						}
					l381:
						{
							position383, tokenIndex383, depth383 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l384
							}
							position++
							goto l383
						l384:
							position, tokenIndex, depth = position383, tokenIndex383, depth383
							if buffer[position] != rune('E') {
								goto l358
							}
							position++
						}
					l383:
						break
					case 'W', 'w':
						{
							position385, tokenIndex385, depth385 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l386
							}
							position++
							goto l385
						l386:
							position, tokenIndex, depth = position385, tokenIndex385, depth385
							if buffer[position] != rune('W') {
								goto l358
							}
							position++
						}
					l385:
						{
							position387, tokenIndex387, depth387 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l388
							}
							position++
							goto l387
						l388:
							position, tokenIndex, depth = position387, tokenIndex387, depth387
							if buffer[position] != rune('O') {
								goto l358
							}
							position++
						}
					l387:
						{
							position389, tokenIndex389, depth389 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l390
							}
							position++
							goto l389
						l390:
							position, tokenIndex, depth = position389, tokenIndex389, depth389
							if buffer[position] != rune('R') {
								goto l358
							}
							position++
						}
					l389:
						{
							position391, tokenIndex391, depth391 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l392
							}
							position++
							goto l391
						l392:
							position, tokenIndex, depth = position391, tokenIndex391, depth391
							if buffer[position] != rune('D') {
								goto l358
							}
							position++
						}
					l391:
						break
					default:
						{
							position393, tokenIndex393, depth393 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l394
							}
							position++
							goto l393
						l394:
							position, tokenIndex, depth = position393, tokenIndex393, depth393
							if buffer[position] != rune('D') {
								goto l358
							}
							position++
						}
					l393:
						{
							position395, tokenIndex395, depth395 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l396
							}
							position++
							goto l395
						l396:
							position, tokenIndex, depth = position395, tokenIndex395, depth395
							if buffer[position] != rune('W') {
								goto l358
							}
							position++
						}
					l395:
						{
							position397, tokenIndex397, depth397 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l398
							}
							position++
							goto l397
						l398:
							position, tokenIndex, depth = position397, tokenIndex397, depth397
							if buffer[position] != rune('O') {
								goto l358
							}
							position++
						}
					l397:
						{
							position399, tokenIndex399, depth399 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l400
							}
							position++
							goto l399
						l400:
							position, tokenIndex, depth = position399, tokenIndex399, depth399
							if buffer[position] != rune('R') {
								goto l358
							}
							position++
						}
					l399:
						{
							position401, tokenIndex401, depth401 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l402
							}
							position++
							goto l401
						l402:
							position, tokenIndex, depth = position401, tokenIndex401, depth401
							if buffer[position] != rune('D') {
								goto l358
							}
							position++
						}
					l401:
						break
					}
				}

				if !_rules[ruleSpace]() {
					goto l358
				}
				depth--
				add(ruleDATA_TYPE, position359)
			}
			return true
		l358:
			position, tokenIndex, depth = position358, tokenIndex358, depth358
			return false
		},
		/* 31 LBRK <- <('[' Spacing)> */
		func() bool {
			position403, tokenIndex403, depth403 := position, tokenIndex, depth
			{
				position404 := position
				depth++
				if buffer[position] != rune('[') {
					goto l403
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l403
				}
				depth--
				add(ruleLBRK, position404)
			}
			return true
		l403:
			position, tokenIndex, depth = position403, tokenIndex403, depth403
			return false
		},
		/* 32 RBRK <- <(']' Spacing)> */
		func() bool {
			position405, tokenIndex405, depth405 := position, tokenIndex, depth
			{
				position406 := position
				depth++
				if buffer[position] != rune(']') {
					goto l405
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l405
				}
				depth--
				add(ruleRBRK, position406)
			}
			return true
		l405:
			position, tokenIndex, depth = position405, tokenIndex405, depth405
			return false
		},
		/* 33 COMMA <- <(',' Spacing)> */
		func() bool {
			position407, tokenIndex407, depth407 := position, tokenIndex, depth
			{
				position408 := position
				depth++
				if buffer[position] != rune(',') {
					goto l407
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l407
				}
				depth--
				add(ruleCOMMA, position408)
			}
			return true
		l407:
			position, tokenIndex, depth = position407, tokenIndex407, depth407
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
			position412, tokenIndex412, depth412 := position, tokenIndex, depth
			{
				position413 := position
				depth++
				if buffer[position] != rune('\n') {
					goto l412
				}
				position++
				depth--
				add(ruleNL, position413)
			}
			return true
		l412:
			position, tokenIndex, depth = position412, tokenIndex412, depth412
			return false
		},
		/* 38 EOT <- <!.> */
		nil,
		/* 39 Literal <- <((FloatLiteral / ((&('"') StringLiteral) | (&('\'') CharLiteral) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') IntegerLiteral))) Spacing)> */
		nil,
		/* 40 IntegerLiteral <- <(<(MINUS? (HexNumeral / BinaryNumeral / OctalNumeral / DecimalNumeral))> Spacing Action33)> */
		func() bool {
			position416, tokenIndex416, depth416 := position, tokenIndex, depth
			{
				position417 := position
				depth++
				{
					position418 := position
					depth++
					{
						position419, tokenIndex419, depth419 := position, tokenIndex, depth
						{
							position421 := position
							depth++
							if buffer[position] != rune('-') {
								goto l419
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l419
							}
							depth--
							add(ruleMINUS, position421)
						}
						goto l420
					l419:
						position, tokenIndex, depth = position419, tokenIndex419, depth419
					}
				l420:
					{
						position422, tokenIndex422, depth422 := position, tokenIndex, depth
						if !_rules[ruleHexNumeral]() {
							goto l423
						}
						goto l422
					l423:
						position, tokenIndex, depth = position422, tokenIndex422, depth422
						{
							position425 := position
							depth++
							{
								position426, tokenIndex426, depth426 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l427
								}
								position++
								if buffer[position] != rune('b') {
									goto l427
								}
								position++
								goto l426
							l427:
								position, tokenIndex, depth = position426, tokenIndex426, depth426
								if buffer[position] != rune('0') {
									goto l424
								}
								position++
								if buffer[position] != rune('B') {
									goto l424
								}
								position++
							}
						l426:
							{
								position428, tokenIndex428, depth428 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l429
								}
								position++
								goto l428
							l429:
								position, tokenIndex, depth = position428, tokenIndex428, depth428
								if buffer[position] != rune('1') {
									goto l424
								}
								position++
							}
						l428:
						l430:
							{
								position431, tokenIndex431, depth431 := position, tokenIndex, depth
							l432:
								{
									position433, tokenIndex433, depth433 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l433
									}
									position++
									goto l432
								l433:
									position, tokenIndex, depth = position433, tokenIndex433, depth433
								}
								{
									position434, tokenIndex434, depth434 := position, tokenIndex, depth
									if buffer[position] != rune('0') {
										goto l435
									}
									position++
									goto l434
								l435:
									position, tokenIndex, depth = position434, tokenIndex434, depth434
									if buffer[position] != rune('1') {
										goto l431
									}
									position++
								}
							l434:
								goto l430
							l431:
								position, tokenIndex, depth = position431, tokenIndex431, depth431
							}
							depth--
							add(ruleBinaryNumeral, position425)
						}
						goto l422
					l424:
						position, tokenIndex, depth = position422, tokenIndex422, depth422
						{
							position437 := position
							depth++
							if buffer[position] != rune('0') {
								goto l436
							}
							position++
						l440:
							{
								position441, tokenIndex441, depth441 := position, tokenIndex, depth
								if buffer[position] != rune('_') {
									goto l441
								}
								position++
								goto l440
							l441:
								position, tokenIndex, depth = position441, tokenIndex441, depth441
							}
							if c := buffer[position]; c < rune('0') || c > rune('7') {
								goto l436
							}
							position++
						l438:
							{
								position439, tokenIndex439, depth439 := position, tokenIndex, depth
							l442:
								{
									position443, tokenIndex443, depth443 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l443
									}
									position++
									goto l442
								l443:
									position, tokenIndex, depth = position443, tokenIndex443, depth443
								}
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l439
								}
								position++
								goto l438
							l439:
								position, tokenIndex, depth = position439, tokenIndex439, depth439
							}
							depth--
							add(ruleOctalNumeral, position437)
						}
						goto l422
					l436:
						position, tokenIndex, depth = position422, tokenIndex422, depth422
						{
							position444 := position
							depth++
							{
								position445, tokenIndex445, depth445 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l446
								}
								position++
								goto l445
							l446:
								position, tokenIndex, depth = position445, tokenIndex445, depth445
								if c := buffer[position]; c < rune('1') || c > rune('9') {
									goto l416
								}
								position++
							l447:
								{
									position448, tokenIndex448, depth448 := position, tokenIndex, depth
								l449:
									{
										position450, tokenIndex450, depth450 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l450
										}
										position++
										goto l449
									l450:
										position, tokenIndex, depth = position450, tokenIndex450, depth450
									}
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l448
									}
									position++
									goto l447
								l448:
									position, tokenIndex, depth = position448, tokenIndex448, depth448
								}
							}
						l445:
							depth--
							add(ruleDecimalNumeral, position444)
						}
					}
				l422:
					depth--
					add(rulePegText, position418)
				}
				if !_rules[ruleSpacing]() {
					goto l416
				}
				{
					add(ruleAction33, position)
				}
				depth--
				add(ruleIntegerLiteral, position417)
			}
			return true
		l416:
			position, tokenIndex, depth = position416, tokenIndex416, depth416
			return false
		},
		/* 41 DecimalNumeral <- <('0' / ([1-9] ('_'* [0-9])*))> */
		nil,
		/* 42 HexNumeral <- <((('0' 'x') / ('0' 'X')) HexDigits)> */
		func() bool {
			position453, tokenIndex453, depth453 := position, tokenIndex, depth
			{
				position454 := position
				depth++
				{
					position455, tokenIndex455, depth455 := position, tokenIndex, depth
					if buffer[position] != rune('0') {
						goto l456
					}
					position++
					if buffer[position] != rune('x') {
						goto l456
					}
					position++
					goto l455
				l456:
					position, tokenIndex, depth = position455, tokenIndex455, depth455
					if buffer[position] != rune('0') {
						goto l453
					}
					position++
					if buffer[position] != rune('X') {
						goto l453
					}
					position++
				}
			l455:
				if !_rules[ruleHexDigits]() {
					goto l453
				}
				depth--
				add(ruleHexNumeral, position454)
			}
			return true
		l453:
			position, tokenIndex, depth = position453, tokenIndex453, depth453
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
			position461, tokenIndex461, depth461 := position, tokenIndex, depth
			{
				position462 := position
				depth++
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l464
					}
					position++
					goto l463
				l464:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
					if buffer[position] != rune('E') {
						goto l461
					}
					position++
				}
			l463:
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					{
						position467, tokenIndex467, depth467 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l468
						}
						position++
						goto l467
					l468:
						position, tokenIndex, depth = position467, tokenIndex467, depth467
						if buffer[position] != rune('-') {
							goto l465
						}
						position++
					}
				l467:
					goto l466
				l465:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
				}
			l466:
				if !_rules[ruleDigits]() {
					goto l461
				}
				depth--
				add(ruleExponent, position462)
			}
			return true
		l461:
			position, tokenIndex, depth = position461, tokenIndex461, depth461
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
			position472, tokenIndex472, depth472 := position, tokenIndex, depth
			{
				position473 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l472
				}
				position++
			l474:
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
				l476:
					{
						position477, tokenIndex477, depth477 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l477
						}
						position++
						goto l476
					l477:
						position, tokenIndex, depth = position477, tokenIndex477, depth477
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l475
					}
					position++
					goto l474
				l475:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
				}
				depth--
				add(ruleDigits, position473)
			}
			return true
		l472:
			position, tokenIndex, depth = position472, tokenIndex472, depth472
			return false
		},
		/* 52 HexDigits <- <(HexDigit ('_'* HexDigit)*)> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				if !_rules[ruleHexDigit]() {
					goto l478
				}
			l480:
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
				l482:
					{
						position483, tokenIndex483, depth483 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l483
						}
						position++
						goto l482
					l483:
						position, tokenIndex, depth = position483, tokenIndex483, depth483
					}
					if !_rules[ruleHexDigit]() {
						goto l481
					}
					goto l480
				l481:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
				}
				depth--
				add(ruleHexDigits, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 53 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position484, tokenIndex484, depth484 := position, tokenIndex, depth
			{
				position485 := position
				depth++
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l484
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l484
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l484
						}
						position++
						break
					}
				}

				depth--
				add(ruleHexDigit, position485)
			}
			return true
		l484:
			position, tokenIndex, depth = position484, tokenIndex484, depth484
			return false
		},
		/* 54 CharLiteral <- <('\'' (Escape / (!('\'' / '\\') .)) '\'')> */
		nil,
		/* 55 StringLiteral <- <(<('"' (Escape / (!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('"') '"')) .))* '"')> Action34)> */
		func() bool {
			position488, tokenIndex488, depth488 := position, tokenIndex, depth
			{
				position489 := position
				depth++
				{
					position490 := position
					depth++
					if buffer[position] != rune('"') {
						goto l488
					}
					position++
				l491:
					{
						position492, tokenIndex492, depth492 := position, tokenIndex, depth
						{
							position493, tokenIndex493, depth493 := position, tokenIndex, depth
							if !_rules[ruleEscape]() {
								goto l494
							}
							goto l493
						l494:
							position, tokenIndex, depth = position493, tokenIndex493, depth493
							{
								position495, tokenIndex495, depth495 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '\r':
										if buffer[position] != rune('\r') {
											goto l495
										}
										position++
										break
									case '\n':
										if buffer[position] != rune('\n') {
											goto l495
										}
										position++
										break
									case '\\':
										if buffer[position] != rune('\\') {
											goto l495
										}
										position++
										break
									default:
										if buffer[position] != rune('"') {
											goto l495
										}
										position++
										break
									}
								}

								goto l492
							l495:
								position, tokenIndex, depth = position495, tokenIndex495, depth495
							}
							if !matchDot() {
								goto l492
							}
						}
					l493:
						goto l491
					l492:
						position, tokenIndex, depth = position492, tokenIndex492, depth492
					}
					if buffer[position] != rune('"') {
						goto l488
					}
					position++
					depth--
					add(rulePegText, position490)
				}
				{
					add(ruleAction34, position)
				}
				depth--
				add(ruleStringLiteral, position489)
			}
			return true
		l488:
			position, tokenIndex, depth = position488, tokenIndex488, depth488
			return false
		},
		/* 56 Escape <- <('\\' ((&('u') UnicodeEscape) | (&('\\') '\\') | (&('\'') '\'') | (&('"') '"') | (&('r') 'r') | (&('f') 'f') | (&('n') 'n') | (&('t') 't') | (&('b') 'b') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7') OctalEscape)))> */
		func() bool {
			position498, tokenIndex498, depth498 := position, tokenIndex, depth
			{
				position499 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l498
				}
				position++
				{
					switch buffer[position] {
					case 'u':
						{
							position501 := position
							depth++
							if buffer[position] != rune('u') {
								goto l498
							}
							position++
						l502:
							{
								position503, tokenIndex503, depth503 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l503
								}
								position++
								goto l502
							l503:
								position, tokenIndex, depth = position503, tokenIndex503, depth503
							}
							if !_rules[ruleHexDigit]() {
								goto l498
							}
							if !_rules[ruleHexDigit]() {
								goto l498
							}
							if !_rules[ruleHexDigit]() {
								goto l498
							}
							if !_rules[ruleHexDigit]() {
								goto l498
							}
							depth--
							add(ruleUnicodeEscape, position501)
						}
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l498
						}
						position++
						break
					case '\'':
						if buffer[position] != rune('\'') {
							goto l498
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l498
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l498
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l498
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l498
						}
						position++
						break
					case 't':
						if buffer[position] != rune('t') {
							goto l498
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l498
						}
						position++
						break
					default:
						{
							position504 := position
							depth++
							{
								position505, tokenIndex505, depth505 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('3') {
									goto l506
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l506
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l506
								}
								position++
								goto l505
							l506:
								position, tokenIndex, depth = position505, tokenIndex505, depth505
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l507
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l507
								}
								position++
								goto l505
							l507:
								position, tokenIndex, depth = position505, tokenIndex505, depth505
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l498
								}
								position++
							}
						l505:
							depth--
							add(ruleOctalEscape, position504)
						}
						break
					}
				}

				depth--
				add(ruleEscape, position499)
			}
			return true
		l498:
			position, tokenIndex, depth = position498, tokenIndex498, depth498
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
