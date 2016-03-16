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
	rules  [98]func() bool
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
			p.PushInst(asm.OP_EXIT)
		case ruleAction20:
			p.PushInst(asm.OP_RET)
		case ruleAction21:
			p.PushInst(asm.OP_NOP)
		case ruleAction22:
			p.PushInst(asm.OP_CALL)
		case ruleAction23:
			p.PushInst(asm.OP_PUSH)
		case ruleAction24:
			p.PushInst(asm.OP_POP)
		case ruleAction25:
			p.PushInst(asm.OP_JMP)
		case ruleAction26:
			p.PushInst(asm.OP_IN)
		case ruleAction27:
			p.PushInst(asm.OP_OUT)
		case ruleAction28:
			p.PushInst(asm.OP_CAL)
		case ruleAction29:
			p.PushInst(asm.OP_LD)
		case ruleAction30:
			p.PushInst(asm.OP_CMP)
		case ruleAction31:
			p.PushInst(asm.OP_JPC)
		case ruleAction32:
			p.Push(&asm.PseudoBlock{})
		case ruleAction33:
			p.Push(&asm.PseudoData{})
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
														add(ruleAction32, position)
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
														add(ruleAction33, position)
													}
													depth--
													add(ruleDATA, position28)
												}
												if !_rules[ruleIdentifier]() {
													goto l4
												}
												{
													position38, tokenIndex38, depth38 := position, tokenIndex, depth
													{
														position40 := position
														depth++
														{
															position41, tokenIndex41, depth41 := position, tokenIndex, depth
															if !_rules[ruleDATA_TYPE]() {
																goto l42
															}
															goto l41
														l42:
															position, tokenIndex, depth = position41, tokenIndex41, depth41
															{
																position43, tokenIndex43, depth43 := position, tokenIndex, depth
																{
																	position45, tokenIndex45, depth45 := position, tokenIndex, depth
																	if buffer[position] != rune('c') {
																		goto l46
																	}
																	position++
																	goto l45
																l46:
																	position, tokenIndex, depth = position45, tokenIndex45, depth45
																	if buffer[position] != rune('C') {
																		goto l44
																	}
																	position++
																}
															l45:
																{
																	position47, tokenIndex47, depth47 := position, tokenIndex, depth
																	if buffer[position] != rune('h') {
																		goto l48
																	}
																	position++
																	goto l47
																l48:
																	position, tokenIndex, depth = position47, tokenIndex47, depth47
																	if buffer[position] != rune('H') {
																		goto l44
																	}
																	position++
																}
															l47:
																{
																	position49, tokenIndex49, depth49 := position, tokenIndex, depth
																	if buffer[position] != rune('a') {
																		goto l50
																	}
																	position++
																	goto l49
																l50:
																	position, tokenIndex, depth = position49, tokenIndex49, depth49
																	if buffer[position] != rune('A') {
																		goto l44
																	}
																	position++
																}
															l49:
																{
																	position51, tokenIndex51, depth51 := position, tokenIndex, depth
																	if buffer[position] != rune('r') {
																		goto l52
																	}
																	position++
																	goto l51
																l52:
																	position, tokenIndex, depth = position51, tokenIndex51, depth51
																	if buffer[position] != rune('R') {
																		goto l44
																	}
																	position++
																}
															l51:
																goto l43
															l44:
																position, tokenIndex, depth = position43, tokenIndex43, depth43
																{
																	position53, tokenIndex53, depth53 := position, tokenIndex, depth
																	if buffer[position] != rune('b') {
																		goto l54
																	}
																	position++
																	goto l53
																l54:
																	position, tokenIndex, depth = position53, tokenIndex53, depth53
																	if buffer[position] != rune('B') {
																		goto l38
																	}
																	position++
																}
															l53:
																{
																	position55, tokenIndex55, depth55 := position, tokenIndex, depth
																	if buffer[position] != rune('i') {
																		goto l56
																	}
																	position++
																	goto l55
																l56:
																	position, tokenIndex, depth = position55, tokenIndex55, depth55
																	if buffer[position] != rune('I') {
																		goto l38
																	}
																	position++
																}
															l55:
																{
																	position57, tokenIndex57, depth57 := position, tokenIndex, depth
																	if buffer[position] != rune('n') {
																		goto l58
																	}
																	position++
																	goto l57
																l58:
																	position, tokenIndex, depth = position57, tokenIndex57, depth57
																	if buffer[position] != rune('N') {
																		goto l38
																	}
																	position++
																}
															l57:
															}
														l43:
															if !_rules[ruleSpace]() {
																goto l38
															}
														}
													l41:
														depth--
														add(rulePSEUDO_DATA_TYPE, position40)
													}
													goto l39
												l38:
													position, tokenIndex, depth = position38, tokenIndex38, depth38
												}
											l39:
												if !_rules[rulePseudoDataValue]() {
													goto l4
												}
											l59:
												{
													position60, tokenIndex60, depth60 := position, tokenIndex, depth
													if !_rules[ruleCOMMA]() {
														goto l60
													}
													if !_rules[rulePseudoDataValue]() {
														goto l60
													}
													goto l59
												l60:
													position, tokenIndex, depth = position60, tokenIndex60, depth60
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
											position61 := position
											depth++
											{
												position62, tokenIndex62, depth62 := position, tokenIndex, depth
												{
													position64, tokenIndex64, depth64 := position, tokenIndex, depth
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
																goto l65
															}
															position++
														}
													l67:
														{
															position69, tokenIndex69, depth69 := position, tokenIndex, depth
															if buffer[position] != rune('u') {
																goto l70
															}
															position++
															goto l69
														l70:
															position, tokenIndex, depth = position69, tokenIndex69, depth69
															if buffer[position] != rune('U') {
																goto l65
															}
															position++
														}
													l69:
														{
															position71, tokenIndex71, depth71 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l72
															}
															position++
															goto l71
														l72:
															position, tokenIndex, depth = position71, tokenIndex71, depth71
															if buffer[position] != rune('S') {
																goto l65
															}
															position++
														}
													l71:
														{
															position73, tokenIndex73, depth73 := position, tokenIndex, depth
															if buffer[position] != rune('h') {
																goto l74
															}
															position++
															goto l73
														l74:
															position, tokenIndex, depth = position73, tokenIndex73, depth73
															if buffer[position] != rune('H') {
																goto l65
															}
															position++
														}
													l73:
														if !_rules[ruleSpace]() {
															goto l65
														}
														{
															add(ruleAction23, position)
														}
														depth--
														add(rulePUSH, position66)
													}
													goto l64
												l65:
													position, tokenIndex, depth = position64, tokenIndex64, depth64
													{
														switch buffer[position] {
														case 'J', 'j':
															{
																position77 := position
																depth++
																{
																	position78, tokenIndex78, depth78 := position, tokenIndex, depth
																	if buffer[position] != rune('j') {
																		goto l79
																	}
																	position++
																	goto l78
																l79:
																	position, tokenIndex, depth = position78, tokenIndex78, depth78
																	if buffer[position] != rune('J') {
																		goto l63
																	}
																	position++
																}
															l78:
																{
																	position80, tokenIndex80, depth80 := position, tokenIndex, depth
																	if buffer[position] != rune('m') {
																		goto l81
																	}
																	position++
																	goto l80
																l81:
																	position, tokenIndex, depth = position80, tokenIndex80, depth80
																	if buffer[position] != rune('M') {
																		goto l63
																	}
																	position++
																}
															l80:
																{
																	position82, tokenIndex82, depth82 := position, tokenIndex, depth
																	if buffer[position] != rune('p') {
																		goto l83
																	}
																	position++
																	goto l82
																l83:
																	position, tokenIndex, depth = position82, tokenIndex82, depth82
																	if buffer[position] != rune('P') {
																		goto l63
																	}
																	position++
																}
															l82:
																if !_rules[ruleSpace]() {
																	goto l63
																}
																{
																	add(ruleAction25, position)
																}
																depth--
																add(ruleJMP, position77)
															}
															break
														case 'P', 'p':
															{
																position85 := position
																depth++
																{
																	position86, tokenIndex86, depth86 := position, tokenIndex, depth
																	if buffer[position] != rune('p') {
																		goto l87
																	}
																	position++
																	goto l86
																l87:
																	position, tokenIndex, depth = position86, tokenIndex86, depth86
																	if buffer[position] != rune('P') {
																		goto l63
																	}
																	position++
																}
															l86:
																{
																	position88, tokenIndex88, depth88 := position, tokenIndex, depth
																	if buffer[position] != rune('o') {
																		goto l89
																	}
																	position++
																	goto l88
																l89:
																	position, tokenIndex, depth = position88, tokenIndex88, depth88
																	if buffer[position] != rune('O') {
																		goto l63
																	}
																	position++
																}
															l88:
																{
																	position90, tokenIndex90, depth90 := position, tokenIndex, depth
																	if buffer[position] != rune('p') {
																		goto l91
																	}
																	position++
																	goto l90
																l91:
																	position, tokenIndex, depth = position90, tokenIndex90, depth90
																	if buffer[position] != rune('P') {
																		goto l63
																	}
																	position++
																}
															l90:
																if !_rules[ruleSpace]() {
																	goto l63
																}
																{
																	add(ruleAction24, position)
																}
																depth--
																add(rulePOP, position85)
															}
															break
														default:
															{
																position93 := position
																depth++
																{
																	position94, tokenIndex94, depth94 := position, tokenIndex, depth
																	if buffer[position] != rune('c') {
																		goto l95
																	}
																	position++
																	goto l94
																l95:
																	position, tokenIndex, depth = position94, tokenIndex94, depth94
																	if buffer[position] != rune('C') {
																		goto l63
																	}
																	position++
																}
															l94:
																{
																	position96, tokenIndex96, depth96 := position, tokenIndex, depth
																	if buffer[position] != rune('a') {
																		goto l97
																	}
																	position++
																	goto l96
																l97:
																	position, tokenIndex, depth = position96, tokenIndex96, depth96
																	if buffer[position] != rune('A') {
																		goto l63
																	}
																	position++
																}
															l96:
																{
																	position98, tokenIndex98, depth98 := position, tokenIndex, depth
																	if buffer[position] != rune('l') {
																		goto l99
																	}
																	position++
																	goto l98
																l99:
																	position, tokenIndex, depth = position98, tokenIndex98, depth98
																	if buffer[position] != rune('L') {
																		goto l63
																	}
																	position++
																}
															l98:
																{
																	position100, tokenIndex100, depth100 := position, tokenIndex, depth
																	if buffer[position] != rune('l') {
																		goto l101
																	}
																	position++
																	goto l100
																l101:
																	position, tokenIndex, depth = position100, tokenIndex100, depth100
																	if buffer[position] != rune('L') {
																		goto l63
																	}
																	position++
																}
															l100:
																if !_rules[ruleSpace]() {
																	goto l63
																}
																{
																	add(ruleAction22, position)
																}
																depth--
																add(ruleCALL, position93)
															}
															break
														}
													}

												}
											l64:
												if !_rules[ruleOperand]() {
													goto l63
												}
												goto l62
											l63:
												position, tokenIndex, depth = position62, tokenIndex62, depth62
												{
													position104 := position
													depth++
													{
														position105, tokenIndex105, depth105 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l106
														}
														position++
														goto l105
													l106:
														position, tokenIndex, depth = position105, tokenIndex105, depth105
														if buffer[position] != rune('C') {
															goto l103
														}
														position++
													}
												l105:
													{
														position107, tokenIndex107, depth107 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l108
														}
														position++
														goto l107
													l108:
														position, tokenIndex, depth = position107, tokenIndex107, depth107
														if buffer[position] != rune('A') {
															goto l103
														}
														position++
													}
												l107:
													{
														position109, tokenIndex109, depth109 := position, tokenIndex, depth
														if buffer[position] != rune('l') {
															goto l110
														}
														position++
														goto l109
													l110:
														position, tokenIndex, depth = position109, tokenIndex109, depth109
														if buffer[position] != rune('L') {
															goto l103
														}
														position++
													}
												l109:
													if !_rules[ruleSpace]() {
														goto l103
													}
													{
														add(ruleAction28, position)
													}
													depth--
													add(ruleCAL, position104)
												}
												{
													position112 := position
													depth++
													if !_rules[ruleDATA_TYPE]() {
														goto l103
													}
													depth--
													add(rulePegText, position112)
												}
												{
													add(ruleAction5, position)
												}
												{
													position114 := position
													depth++
													{
														position115 := position
														depth++
														{
															position116, tokenIndex116, depth116 := position, tokenIndex, depth
															{
																position118, tokenIndex118, depth118 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l119
																}
																position++
																goto l118
															l119:
																position, tokenIndex, depth = position118, tokenIndex118, depth118
																if buffer[position] != rune('M') {
																	goto l117
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
																	goto l117
																}
																position++
															}
														l120:
															{
																position122, tokenIndex122, depth122 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l123
																}
																position++
																goto l122
															l123:
																position, tokenIndex, depth = position122, tokenIndex122, depth122
																if buffer[position] != rune('L') {
																	goto l117
																}
																position++
															}
														l122:
															goto l116
														l117:
															position, tokenIndex, depth = position116, tokenIndex116, depth116
															{
																switch buffer[position] {
																case 'M', 'm':
																	{
																		position125, tokenIndex125, depth125 := position, tokenIndex, depth
																		if buffer[position] != rune('m') {
																			goto l126
																		}
																		position++
																		goto l125
																	l126:
																		position, tokenIndex, depth = position125, tokenIndex125, depth125
																		if buffer[position] != rune('M') {
																			goto l103
																		}
																		position++
																	}
																l125:
																	{
																		position127, tokenIndex127, depth127 := position, tokenIndex, depth
																		if buffer[position] != rune('o') {
																			goto l128
																		}
																		position++
																		goto l127
																	l128:
																		position, tokenIndex, depth = position127, tokenIndex127, depth127
																		if buffer[position] != rune('O') {
																			goto l103
																		}
																		position++
																	}
																l127:
																	{
																		position129, tokenIndex129, depth129 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l130
																		}
																		position++
																		goto l129
																	l130:
																		position, tokenIndex, depth = position129, tokenIndex129, depth129
																		if buffer[position] != rune('D') {
																			goto l103
																		}
																		position++
																	}
																l129:
																	break
																case 'D', 'd':
																	{
																		position131, tokenIndex131, depth131 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l132
																		}
																		position++
																		goto l131
																	l132:
																		position, tokenIndex, depth = position131, tokenIndex131, depth131
																		if buffer[position] != rune('D') {
																			goto l103
																		}
																		position++
																	}
																l131:
																	{
																		position133, tokenIndex133, depth133 := position, tokenIndex, depth
																		if buffer[position] != rune('i') {
																			goto l134
																		}
																		position++
																		goto l133
																	l134:
																		position, tokenIndex, depth = position133, tokenIndex133, depth133
																		if buffer[position] != rune('I') {
																			goto l103
																		}
																		position++
																	}
																l133:
																	{
																		position135, tokenIndex135, depth135 := position, tokenIndex, depth
																		if buffer[position] != rune('v') {
																			goto l136
																		}
																		position++
																		goto l135
																	l136:
																		position, tokenIndex, depth = position135, tokenIndex135, depth135
																		if buffer[position] != rune('V') {
																			goto l103
																		}
																		position++
																	}
																l135:
																	break
																case 'S', 's':
																	{
																		position137, tokenIndex137, depth137 := position, tokenIndex, depth
																		if buffer[position] != rune('s') {
																			goto l138
																		}
																		position++
																		goto l137
																	l138:
																		position, tokenIndex, depth = position137, tokenIndex137, depth137
																		if buffer[position] != rune('S') {
																			goto l103
																		}
																		position++
																	}
																l137:
																	{
																		position139, tokenIndex139, depth139 := position, tokenIndex, depth
																		if buffer[position] != rune('u') {
																			goto l140
																		}
																		position++
																		goto l139
																	l140:
																		position, tokenIndex, depth = position139, tokenIndex139, depth139
																		if buffer[position] != rune('U') {
																			goto l103
																		}
																		position++
																	}
																l139:
																	{
																		position141, tokenIndex141, depth141 := position, tokenIndex, depth
																		if buffer[position] != rune('b') {
																			goto l142
																		}
																		position++
																		goto l141
																	l142:
																		position, tokenIndex, depth = position141, tokenIndex141, depth141
																		if buffer[position] != rune('B') {
																			goto l103
																		}
																		position++
																	}
																l141:
																	break
																default:
																	{
																		position143, tokenIndex143, depth143 := position, tokenIndex, depth
																		if buffer[position] != rune('a') {
																			goto l144
																		}
																		position++
																		goto l143
																	l144:
																		position, tokenIndex, depth = position143, tokenIndex143, depth143
																		if buffer[position] != rune('A') {
																			goto l103
																		}
																		position++
																	}
																l143:
																	{
																		position145, tokenIndex145, depth145 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l146
																		}
																		position++
																		goto l145
																	l146:
																		position, tokenIndex, depth = position145, tokenIndex145, depth145
																		if buffer[position] != rune('D') {
																			goto l103
																		}
																		position++
																	}
																l145:
																	{
																		position147, tokenIndex147, depth147 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l148
																		}
																		position++
																		goto l147
																	l148:
																		position, tokenIndex, depth = position147, tokenIndex147, depth147
																		if buffer[position] != rune('D') {
																			goto l103
																		}
																		position++
																	}
																l147:
																	break
																}
															}

														}
													l116:
														if !_rules[ruleSpace]() {
															goto l103
														}
														depth--
														add(ruleCAL_OP, position115)
													}
													depth--
													add(rulePegText, position114)
												}
												{
													add(ruleAction6, position)
												}
												if !_rules[ruleOperand]() {
													goto l103
												}
												if !_rules[ruleCOMMA]() {
													goto l103
												}
												if !_rules[ruleOperand]() {
													goto l103
												}
												goto l62
											l103:
												position, tokenIndex, depth = position62, tokenIndex62, depth62
												{
													switch buffer[position] {
													case 'J', 'j':
														{
															position151 := position
															depth++
															{
																position152, tokenIndex152, depth152 := position, tokenIndex, depth
																if buffer[position] != rune('j') {
																	goto l153
																}
																position++
																goto l152
															l153:
																position, tokenIndex, depth = position152, tokenIndex152, depth152
																if buffer[position] != rune('J') {
																	goto l4
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
																	goto l4
																}
																position++
															}
														l154:
															{
																position156, tokenIndex156, depth156 := position, tokenIndex, depth
																if buffer[position] != rune('c') {
																	goto l157
																}
																position++
																goto l156
															l157:
																position, tokenIndex, depth = position156, tokenIndex156, depth156
																if buffer[position] != rune('C') {
																	goto l4
																}
																position++
															}
														l156:
															if !_rules[ruleSpace]() {
																goto l4
															}
															{
																add(ruleAction31, position)
															}
															depth--
															add(ruleJPC, position151)
														}
														{
															position159 := position
															depth++
															{
																position160 := position
																depth++
																{
																	position161, tokenIndex161, depth161 := position, tokenIndex, depth
																	{
																		position163, tokenIndex163, depth163 := position, tokenIndex, depth
																		if buffer[position] != rune('b') {
																			goto l164
																		}
																		position++
																		goto l163
																	l164:
																		position, tokenIndex, depth = position163, tokenIndex163, depth163
																		if buffer[position] != rune('B') {
																			goto l162
																		}
																		position++
																	}
																l163:
																	{
																		position165, tokenIndex165, depth165 := position, tokenIndex, depth
																		if buffer[position] != rune('e') {
																			goto l166
																		}
																		position++
																		goto l165
																	l166:
																		position, tokenIndex, depth = position165, tokenIndex165, depth165
																		if buffer[position] != rune('E') {
																			goto l162
																		}
																		position++
																	}
																l165:
																	goto l161
																l162:
																	position, tokenIndex, depth = position161, tokenIndex161, depth161
																	{
																		position168, tokenIndex168, depth168 := position, tokenIndex, depth
																		if buffer[position] != rune('a') {
																			goto l169
																		}
																		position++
																		goto l168
																	l169:
																		position, tokenIndex, depth = position168, tokenIndex168, depth168
																		if buffer[position] != rune('A') {
																			goto l167
																		}
																		position++
																	}
																l168:
																	{
																		position170, tokenIndex170, depth170 := position, tokenIndex, depth
																		if buffer[position] != rune('e') {
																			goto l171
																		}
																		position++
																		goto l170
																	l171:
																		position, tokenIndex, depth = position170, tokenIndex170, depth170
																		if buffer[position] != rune('E') {
																			goto l167
																		}
																		position++
																	}
																l170:
																	goto l161
																l167:
																	position, tokenIndex, depth = position161, tokenIndex161, depth161
																	{
																		switch buffer[position] {
																		case 'N', 'n':
																			{
																				position173, tokenIndex173, depth173 := position, tokenIndex, depth
																				if buffer[position] != rune('n') {
																					goto l174
																				}
																				position++
																				goto l173
																			l174:
																				position, tokenIndex, depth = position173, tokenIndex173, depth173
																				if buffer[position] != rune('N') {
																					goto l4
																				}
																				position++
																			}
																		l173:
																			{
																				position175, tokenIndex175, depth175 := position, tokenIndex, depth
																				if buffer[position] != rune('z') {
																					goto l176
																				}
																				position++
																				goto l175
																			l176:
																				position, tokenIndex, depth = position175, tokenIndex175, depth175
																				if buffer[position] != rune('Z') {
																					goto l4
																				}
																				position++
																			}
																		l175:
																			break
																		case 'A', 'a':
																			{
																				position177, tokenIndex177, depth177 := position, tokenIndex, depth
																				if buffer[position] != rune('a') {
																					goto l178
																				}
																				position++
																				goto l177
																			l178:
																				position, tokenIndex, depth = position177, tokenIndex177, depth177
																				if buffer[position] != rune('A') {
																					goto l4
																				}
																				position++
																			}
																		l177:
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
																				position179, tokenIndex179, depth179 := position, tokenIndex, depth
																				if buffer[position] != rune('b') {
																					goto l180
																				}
																				position++
																				goto l179
																			l180:
																				position, tokenIndex, depth = position179, tokenIndex179, depth179
																				if buffer[position] != rune('B') {
																					goto l4
																				}
																				position++
																			}
																		l179:
																			break
																		}
																	}

																}
															l161:
																if !_rules[ruleSpace]() {
																	goto l4
																}
																depth--
																add(ruleCMP_OP, position160)
															}
															depth--
															add(rulePegText, position159)
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
															position182 := position
															depth++
															{
																position183, tokenIndex183, depth183 := position, tokenIndex, depth
																if buffer[position] != rune('c') {
																	goto l184
																}
																position++
																goto l183
															l184:
																position, tokenIndex, depth = position183, tokenIndex183, depth183
																if buffer[position] != rune('C') {
																	goto l4
																}
																position++
															}
														l183:
															{
																position185, tokenIndex185, depth185 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l186
																}
																position++
																goto l185
															l186:
																position, tokenIndex, depth = position185, tokenIndex185, depth185
																if buffer[position] != rune('M') {
																	goto l4
																}
																position++
															}
														l185:
															{
																position187, tokenIndex187, depth187 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l188
																}
																position++
																goto l187
															l188:
																position, tokenIndex, depth = position187, tokenIndex187, depth187
																if buffer[position] != rune('P') {
																	goto l4
																}
																position++
															}
														l187:
															if !_rules[ruleSpace]() {
																goto l4
															}
															{
																add(ruleAction30, position)
															}
															depth--
															add(ruleCMP, position182)
														}
														{
															position190 := position
															depth++
															if !_rules[ruleDATA_TYPE]() {
																goto l4
															}
															depth--
															add(rulePegText, position190)
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
															position192 := position
															depth++
															{
																position193, tokenIndex193, depth193 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l194
																}
																position++
																goto l193
															l194:
																position, tokenIndex, depth = position193, tokenIndex193, depth193
																if buffer[position] != rune('L') {
																	goto l4
																}
																position++
															}
														l193:
															{
																position195, tokenIndex195, depth195 := position, tokenIndex, depth
																if buffer[position] != rune('d') {
																	goto l196
																}
																position++
																goto l195
															l196:
																position, tokenIndex, depth = position195, tokenIndex195, depth195
																if buffer[position] != rune('D') {
																	goto l4
																}
																position++
															}
														l195:
															if !_rules[ruleSpace]() {
																goto l4
															}
															{
																add(ruleAction29, position)
															}
															depth--
															add(ruleLD, position192)
														}
														{
															position198 := position
															depth++
															if !_rules[ruleDATA_TYPE]() {
																goto l4
															}
															depth--
															add(rulePegText, position198)
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
															position200 := position
															depth++
															{
																position201, tokenIndex201, depth201 := position, tokenIndex, depth
																if buffer[position] != rune('n') {
																	goto l202
																}
																position++
																goto l201
															l202:
																position, tokenIndex, depth = position201, tokenIndex201, depth201
																if buffer[position] != rune('N') {
																	goto l4
																}
																position++
															}
														l201:
															{
																position203, tokenIndex203, depth203 := position, tokenIndex, depth
																if buffer[position] != rune('o') {
																	goto l204
																}
																position++
																goto l203
															l204:
																position, tokenIndex, depth = position203, tokenIndex203, depth203
																if buffer[position] != rune('O') {
																	goto l4
																}
																position++
															}
														l203:
															{
																position205, tokenIndex205, depth205 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l206
																}
																position++
																goto l205
															l206:
																position, tokenIndex, depth = position205, tokenIndex205, depth205
																if buffer[position] != rune('P') {
																	goto l4
																}
																position++
															}
														l205:
															if !_rules[ruleSpacing]() {
																goto l4
															}
															{
																add(ruleAction21, position)
															}
															depth--
															add(ruleNOP, position200)
														}
														break
													case 'R', 'r':
														{
															position208 := position
															depth++
															{
																position209, tokenIndex209, depth209 := position, tokenIndex, depth
																if buffer[position] != rune('r') {
																	goto l210
																}
																position++
																goto l209
															l210:
																position, tokenIndex, depth = position209, tokenIndex209, depth209
																if buffer[position] != rune('R') {
																	goto l4
																}
																position++
															}
														l209:
															{
																position211, tokenIndex211, depth211 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l212
																}
																position++
																goto l211
															l212:
																position, tokenIndex, depth = position211, tokenIndex211, depth211
																if buffer[position] != rune('E') {
																	goto l4
																}
																position++
															}
														l211:
															{
																position213, tokenIndex213, depth213 := position, tokenIndex, depth
																if buffer[position] != rune('t') {
																	goto l214
																}
																position++
																goto l213
															l214:
																position, tokenIndex, depth = position213, tokenIndex213, depth213
																if buffer[position] != rune('T') {
																	goto l4
																}
																position++
															}
														l213:
															if !_rules[ruleSpacing]() {
																goto l4
															}
															{
																add(ruleAction20, position)
															}
															depth--
															add(ruleRET, position208)
														}
														break
													case 'E', 'e':
														{
															position216 := position
															depth++
															{
																position217, tokenIndex217, depth217 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l218
																}
																position++
																goto l217
															l218:
																position, tokenIndex, depth = position217, tokenIndex217, depth217
																if buffer[position] != rune('E') {
																	goto l4
																}
																position++
															}
														l217:
															{
																position219, tokenIndex219, depth219 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l220
																}
																position++
																goto l219
															l220:
																position, tokenIndex, depth = position219, tokenIndex219, depth219
																if buffer[position] != rune('X') {
																	goto l4
																}
																position++
															}
														l219:
															{
																position221, tokenIndex221, depth221 := position, tokenIndex, depth
																if buffer[position] != rune('i') {
																	goto l222
																}
																position++
																goto l221
															l222:
																position, tokenIndex, depth = position221, tokenIndex221, depth221
																if buffer[position] != rune('I') {
																	goto l4
																}
																position++
															}
														l221:
															{
																position223, tokenIndex223, depth223 := position, tokenIndex, depth
																if buffer[position] != rune('t') {
																	goto l224
																}
																position++
																goto l223
															l224:
																position, tokenIndex, depth = position223, tokenIndex223, depth223
																if buffer[position] != rune('T') {
																	goto l4
																}
																position++
															}
														l223:
															if !_rules[ruleSpacing]() {
																goto l4
															}
															{
																add(ruleAction19, position)
															}
															depth--
															add(ruleEXIT, position216)
														}
														break
													default:
														{
															position226, tokenIndex226, depth226 := position, tokenIndex, depth
															{
																position228 := position
																depth++
																{
																	position229, tokenIndex229, depth229 := position, tokenIndex, depth
																	if buffer[position] != rune('i') {
																		goto l230
																	}
																	position++
																	goto l229
																l230:
																	position, tokenIndex, depth = position229, tokenIndex229, depth229
																	if buffer[position] != rune('I') {
																		goto l227
																	}
																	position++
																}
															l229:
																{
																	position231, tokenIndex231, depth231 := position, tokenIndex, depth
																	if buffer[position] != rune('n') {
																		goto l232
																	}
																	position++
																	goto l231
																l232:
																	position, tokenIndex, depth = position231, tokenIndex231, depth231
																	if buffer[position] != rune('N') {
																		goto l227
																	}
																	position++
																}
															l231:
																if !_rules[ruleSpace]() {
																	goto l227
																}
																{
																	add(ruleAction26, position)
																}
																depth--
																add(ruleIN, position228)
															}
															goto l226
														l227:
															position, tokenIndex, depth = position226, tokenIndex226, depth226
															{
																position234 := position
																depth++
																{
																	position235, tokenIndex235, depth235 := position, tokenIndex, depth
																	if buffer[position] != rune('o') {
																		goto l236
																	}
																	position++
																	goto l235
																l236:
																	position, tokenIndex, depth = position235, tokenIndex235, depth235
																	if buffer[position] != rune('O') {
																		goto l4
																	}
																	position++
																}
															l235:
																{
																	position237, tokenIndex237, depth237 := position, tokenIndex, depth
																	if buffer[position] != rune('u') {
																		goto l238
																	}
																	position++
																	goto l237
																l238:
																	position, tokenIndex, depth = position237, tokenIndex237, depth237
																	if buffer[position] != rune('U') {
																		goto l4
																	}
																	position++
																}
															l237:
																{
																	position239, tokenIndex239, depth239 := position, tokenIndex, depth
																	if buffer[position] != rune('t') {
																		goto l240
																	}
																	position++
																	goto l239
																l240:
																	position, tokenIndex, depth = position239, tokenIndex239, depth239
																	if buffer[position] != rune('T') {
																		goto l4
																	}
																	position++
																}
															l239:
																if !_rules[ruleSpace]() {
																	goto l4
																}
																{
																	add(ruleAction27, position)
																}
																depth--
																add(ruleOUT, position234)
															}
														}
													l226:
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
										l62:
											depth--
											add(ruleInst, position61)
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
								position243, tokenIndex243, depth243 := position, tokenIndex, depth
								if !_rules[ruleComment]() {
									goto l243
								}
								{
									add(ruleAction2, position)
								}
								goto l244
							l243:
								position, tokenIndex, depth = position243, tokenIndex243, depth243
							}
						l244:
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
					position247 := position
					depth++
					{
						position248, tokenIndex248, depth248 := position, tokenIndex, depth
						if !matchDot() {
							goto l248
						}
						goto l0
					l248:
						position, tokenIndex, depth = position248, tokenIndex248, depth248
					}
					depth--
					add(ruleEOT, position247)
				}
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					{
						position251 := position
						depth++
						{
							position252, tokenIndex252, depth252 := position, tokenIndex, depth
							{
								position254 := position
								depth++
								{
									position255, tokenIndex255, depth255 := position, tokenIndex, depth
									{
										position257 := position
										depth++
										{
											position258 := position
											depth++
											{
												position259, tokenIndex259, depth259 := position, tokenIndex, depth
												{
													position261, tokenIndex261, depth261 := position, tokenIndex, depth
													if buffer[position] != rune('0') {
														goto l262
													}
													position++
													if buffer[position] != rune('x') {
														goto l262
													}
													position++
													goto l261
												l262:
													position, tokenIndex, depth = position261, tokenIndex261, depth261
													if buffer[position] != rune('0') {
														goto l260
													}
													position++
													if buffer[position] != rune('X') {
														goto l260
													}
													position++
												}
											l261:
												{
													position263, tokenIndex263, depth263 := position, tokenIndex, depth
													if !_rules[ruleHexDigits]() {
														goto l263
													}
													goto l264
												l263:
													position, tokenIndex, depth = position263, tokenIndex263, depth263
												}
											l264:
												if buffer[position] != rune('.') {
													goto l260
												}
												position++
												if !_rules[ruleHexDigits]() {
													goto l260
												}
												goto l259
											l260:
												position, tokenIndex, depth = position259, tokenIndex259, depth259
												if !_rules[ruleHexNumeral]() {
													goto l256
												}
												{
													position265, tokenIndex265, depth265 := position, tokenIndex, depth
													if buffer[position] != rune('.') {
														goto l265
													}
													position++
													goto l266
												l265:
													position, tokenIndex, depth = position265, tokenIndex265, depth265
												}
											l266:
											}
										l259:
											depth--
											add(ruleHexSignificand, position258)
										}
										{
											position267 := position
											depth++
											{
												position268, tokenIndex268, depth268 := position, tokenIndex, depth
												if buffer[position] != rune('p') {
													goto l269
												}
												position++
												goto l268
											l269:
												position, tokenIndex, depth = position268, tokenIndex268, depth268
												if buffer[position] != rune('P') {
													goto l256
												}
												position++
											}
										l268:
											{
												position270, tokenIndex270, depth270 := position, tokenIndex, depth
												{
													position272, tokenIndex272, depth272 := position, tokenIndex, depth
													if buffer[position] != rune('+') {
														goto l273
													}
													position++
													goto l272
												l273:
													position, tokenIndex, depth = position272, tokenIndex272, depth272
													if buffer[position] != rune('-') {
														goto l270
													}
													position++
												}
											l272:
												goto l271
											l270:
												position, tokenIndex, depth = position270, tokenIndex270, depth270
											}
										l271:
											if !_rules[ruleDigits]() {
												goto l256
											}
											depth--
											add(ruleBinaryExponent, position267)
										}
										{
											position274, tokenIndex274, depth274 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case 'D':
													if buffer[position] != rune('D') {
														goto l274
													}
													position++
													break
												case 'd':
													if buffer[position] != rune('d') {
														goto l274
													}
													position++
													break
												case 'F':
													if buffer[position] != rune('F') {
														goto l274
													}
													position++
													break
												default:
													if buffer[position] != rune('f') {
														goto l274
													}
													position++
													break
												}
											}

											goto l275
										l274:
											position, tokenIndex, depth = position274, tokenIndex274, depth274
										}
									l275:
										depth--
										add(ruleHexFloat, position257)
									}
									goto l255
								l256:
									position, tokenIndex, depth = position255, tokenIndex255, depth255
									{
										position277 := position
										depth++
										{
											position278, tokenIndex278, depth278 := position, tokenIndex, depth
											if !_rules[ruleDigits]() {
												goto l279
											}
											if buffer[position] != rune('.') {
												goto l279
											}
											position++
											{
												position280, tokenIndex280, depth280 := position, tokenIndex, depth
												if !_rules[ruleDigits]() {
													goto l280
												}
												goto l281
											l280:
												position, tokenIndex, depth = position280, tokenIndex280, depth280
											}
										l281:
											{
												position282, tokenIndex282, depth282 := position, tokenIndex, depth
												if !_rules[ruleExponent]() {
													goto l282
												}
												goto l283
											l282:
												position, tokenIndex, depth = position282, tokenIndex282, depth282
											}
										l283:
											{
												position284, tokenIndex284, depth284 := position, tokenIndex, depth
												{
													switch buffer[position] {
													case 'D':
														if buffer[position] != rune('D') {
															goto l284
														}
														position++
														break
													case 'd':
														if buffer[position] != rune('d') {
															goto l284
														}
														position++
														break
													case 'F':
														if buffer[position] != rune('F') {
															goto l284
														}
														position++
														break
													default:
														if buffer[position] != rune('f') {
															goto l284
														}
														position++
														break
													}
												}

												goto l285
											l284:
												position, tokenIndex, depth = position284, tokenIndex284, depth284
											}
										l285:
											goto l278
										l279:
											position, tokenIndex, depth = position278, tokenIndex278, depth278
											if buffer[position] != rune('.') {
												goto l287
											}
											position++
											if !_rules[ruleDigits]() {
												goto l287
											}
											{
												position288, tokenIndex288, depth288 := position, tokenIndex, depth
												if !_rules[ruleExponent]() {
													goto l288
												}
												goto l289
											l288:
												position, tokenIndex, depth = position288, tokenIndex288, depth288
											}
										l289:
											{
												position290, tokenIndex290, depth290 := position, tokenIndex, depth
												{
													switch buffer[position] {
													case 'D':
														if buffer[position] != rune('D') {
															goto l290
														}
														position++
														break
													case 'd':
														if buffer[position] != rune('d') {
															goto l290
														}
														position++
														break
													case 'F':
														if buffer[position] != rune('F') {
															goto l290
														}
														position++
														break
													default:
														if buffer[position] != rune('f') {
															goto l290
														}
														position++
														break
													}
												}

												goto l291
											l290:
												position, tokenIndex, depth = position290, tokenIndex290, depth290
											}
										l291:
											goto l278
										l287:
											position, tokenIndex, depth = position278, tokenIndex278, depth278
											if !_rules[ruleDigits]() {
												goto l293
											}
											if !_rules[ruleExponent]() {
												goto l293
											}
											{
												position294, tokenIndex294, depth294 := position, tokenIndex, depth
												{
													switch buffer[position] {
													case 'D':
														if buffer[position] != rune('D') {
															goto l294
														}
														position++
														break
													case 'd':
														if buffer[position] != rune('d') {
															goto l294
														}
														position++
														break
													case 'F':
														if buffer[position] != rune('F') {
															goto l294
														}
														position++
														break
													default:
														if buffer[position] != rune('f') {
															goto l294
														}
														position++
														break
													}
												}

												goto l295
											l294:
												position, tokenIndex, depth = position294, tokenIndex294, depth294
											}
										l295:
											goto l278
										l293:
											position, tokenIndex, depth = position278, tokenIndex278, depth278
											if !_rules[ruleDigits]() {
												goto l253
											}
											{
												position297, tokenIndex297, depth297 := position, tokenIndex, depth
												if !_rules[ruleExponent]() {
													goto l297
												}
												goto l298
											l297:
												position, tokenIndex, depth = position297, tokenIndex297, depth297
											}
										l298:
											{
												switch buffer[position] {
												case 'D':
													if buffer[position] != rune('D') {
														goto l253
													}
													position++
													break
												case 'd':
													if buffer[position] != rune('d') {
														goto l253
													}
													position++
													break
												case 'F':
													if buffer[position] != rune('F') {
														goto l253
													}
													position++
													break
												default:
													if buffer[position] != rune('f') {
														goto l253
													}
													position++
													break
												}
											}

										}
									l278:
										depth--
										add(ruleDecimalFloat, position277)
									}
								}
							l255:
								depth--
								add(ruleFloatLiteral, position254)
							}
							goto l252
						l253:
							position, tokenIndex, depth = position252, tokenIndex252, depth252
							{
								switch buffer[position] {
								case '"':
									if !_rules[ruleStringLiteral]() {
										goto l249
									}
									break
								case '\'':
									{
										position301 := position
										depth++
										if buffer[position] != rune('\'') {
											goto l249
										}
										position++
										{
											position302, tokenIndex302, depth302 := position, tokenIndex, depth
											if !_rules[ruleEscape]() {
												goto l303
											}
											goto l302
										l303:
											position, tokenIndex, depth = position302, tokenIndex302, depth302
											{
												position304, tokenIndex304, depth304 := position, tokenIndex, depth
												{
													position305, tokenIndex305, depth305 := position, tokenIndex, depth
													if buffer[position] != rune('\'') {
														goto l306
													}
													position++
													goto l305
												l306:
													position, tokenIndex, depth = position305, tokenIndex305, depth305
													if buffer[position] != rune('\\') {
														goto l304
													}
													position++
												}
											l305:
												goto l249
											l304:
												position, tokenIndex, depth = position304, tokenIndex304, depth304
											}
											if !matchDot() {
												goto l249
											}
										}
									l302:
										if buffer[position] != rune('\'') {
											goto l249
										}
										position++
										depth--
										add(ruleCharLiteral, position301)
									}
									break
								default:
									if !_rules[ruleIntegerLiteral]() {
										goto l249
									}
									break
								}
							}

						}
					l252:
						if !_rules[ruleSpacing]() {
							goto l249
						}
						depth--
						add(ruleLiteral, position251)
					}
					goto l250
				l249:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
				}
			l250:
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
			position308, tokenIndex308, depth308 := position, tokenIndex, depth
			{
				position309 := position
				depth++
				{
					position310 := position
					depth++
					if buffer[position] != rune(';') {
						goto l308
					}
					position++
					if !_rules[ruleSpacing]() {
						goto l308
					}
					depth--
					add(ruleSEMICOLON, position310)
				}
				{
					position311 := position
					depth++
				l312:
					{
						position313, tokenIndex313, depth313 := position, tokenIndex, depth
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							if !_rules[ruleNL]() {
								goto l314
							}
							goto l313
						l314:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
						}
						if !matchDot() {
							goto l313
						}
						goto l312
					l313:
						position, tokenIndex, depth = position313, tokenIndex313, depth313
					}
					depth--
					add(rulePegText, position311)
				}
				{
					add(ruleAction3, position)
				}
				depth--
				add(ruleComment, position309)
			}
			return true
		l308:
			position, tokenIndex, depth = position308, tokenIndex308, depth308
			return false
		},
		/* 3 Label <- <(Action4 Identifier Spacing COLON)> */
		nil,
		/* 4 Inst <- <(((PUSH / ((&('J' | 'j') JMP) | (&('P' | 'p') POP) | (&('C' | 'c') CALL))) Operand) / (CAL <DATA_TYPE> Action5 <CAL_OP> Action6 Operand COMMA Operand) / ((&('J' | 'j') (JPC <CMP_OP> Action9 Operand)) | (&('C' | 'c') (CMP <DATA_TYPE> Action8 Operand COMMA Operand)) | (&('L' | 'l') (LD <DATA_TYPE> Action7 Operand COMMA Operand)) | (&('N' | 'n') NOP) | (&('R' | 'r') RET) | (&('E' | 'e') EXIT) | (&('I' | 'O' | 'i' | 'o') ((IN / OUT) Operand COMMA Operand))))> */
		nil,
		/* 5 Pseudo <- <((BLOCK IntegerLiteral IntegerLiteral) / (DATA Identifier PSEUDO_DATA_TYPE? PseudoDataValue (COMMA PseudoDataValue)*))> */
		nil,
		/* 6 PseudoDataValue <- <((&('%') (<('%' HexDigits '%')> Spacing Action13)) | (&('"') (StringLiteral Action12)) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (IntegerLiteral Action10)) | (&('$' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (Identifier Action11)))> */
		func() bool {
			position319, tokenIndex319, depth319 := position, tokenIndex, depth
			{
				position320 := position
				depth++
				{
					switch buffer[position] {
					case '%':
						{
							position322 := position
							depth++
							if buffer[position] != rune('%') {
								goto l319
							}
							position++
							if !_rules[ruleHexDigits]() {
								goto l319
							}
							if buffer[position] != rune('%') {
								goto l319
							}
							position++
							depth--
							add(rulePegText, position322)
						}
						if !_rules[ruleSpacing]() {
							goto l319
						}
						{
							add(ruleAction13, position)
						}
						break
					case '"':
						if !_rules[ruleStringLiteral]() {
							goto l319
						}
						{
							add(ruleAction12, position)
						}
						break
					case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if !_rules[ruleIntegerLiteral]() {
							goto l319
						}
						{
							add(ruleAction10, position)
						}
						break
					default:
						if !_rules[ruleIdentifier]() {
							goto l319
						}
						{
							add(ruleAction11, position)
						}
						break
					}
				}

				depth--
				add(rulePseudoDataValue, position320)
			}
			return true
		l319:
			position, tokenIndex, depth = position319, tokenIndex319, depth319
			return false
		},
		/* 7 PSEUDO_DATA_TYPE <- <(DATA_TYPE / (((('c' / 'C') ('h' / 'H') ('a' / 'A') ('r' / 'R')) / (('b' / 'B') ('i' / 'I') ('n' / 'N'))) Space))> */
		nil,
		/* 8 Operand <- <(((LBRK Identifier RBRK Action15) / ((&('[') (LBRK IntegerLiteral RBRK Action17)) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') (IntegerLiteral Action16)) | (&('$' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (Identifier Action14)))) Spacing)> */
		func() bool {
			position328, tokenIndex328, depth328 := position, tokenIndex, depth
			{
				position329 := position
				depth++
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					if !_rules[ruleLBRK]() {
						goto l331
					}
					if !_rules[ruleIdentifier]() {
						goto l331
					}
					if !_rules[ruleRBRK]() {
						goto l331
					}
					{
						add(ruleAction15, position)
					}
					goto l330
				l331:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
					{
						switch buffer[position] {
						case '[':
							if !_rules[ruleLBRK]() {
								goto l328
							}
							if !_rules[ruleIntegerLiteral]() {
								goto l328
							}
							if !_rules[ruleRBRK]() {
								goto l328
							}
							{
								add(ruleAction17, position)
							}
							break
						case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if !_rules[ruleIntegerLiteral]() {
								goto l328
							}
							{
								add(ruleAction16, position)
							}
							break
						default:
							if !_rules[ruleIdentifier]() {
								goto l328
							}
							{
								add(ruleAction14, position)
							}
							break
						}
					}

				}
			l330:
				if !_rules[ruleSpacing]() {
					goto l328
				}
				depth--
				add(ruleOperand, position329)
			}
			return true
		l328:
			position, tokenIndex, depth = position328, tokenIndex328, depth328
			return false
		},
		/* 9 Spacing <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position338 := position
				depth++
			l339:
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l340
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l340
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l340
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l340
							}
							position++
							break
						}
					}

					goto l339
				l340:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
				}
				depth--
				add(ruleSpacing, position338)
			}
			return true
		},
		/* 10 Space <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))+> */
		func() bool {
			position342, tokenIndex342, depth342 := position, tokenIndex, depth
			{
				position343 := position
				depth++
				{
					switch buffer[position] {
					case '\f':
						if buffer[position] != rune('\f') {
							goto l342
						}
						position++
						break
					case '\r':
						if buffer[position] != rune('\r') {
							goto l342
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l342
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l342
						}
						position++
						break
					}
				}

			l344:
				{
					position345, tokenIndex345, depth345 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l345
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l345
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l345
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l345
							}
							position++
							break
						}
					}

					goto l344
				l345:
					position, tokenIndex, depth = position345, tokenIndex345, depth345
				}
				depth--
				add(ruleSpace, position343)
			}
			return true
		l342:
			position, tokenIndex, depth = position342, tokenIndex342, depth342
			return false
		},
		/* 11 Identifier <- <(<(Letter LetterOrDigit*)> Spacing Action18)> */
		func() bool {
			position348, tokenIndex348, depth348 := position, tokenIndex, depth
			{
				position349 := position
				depth++
				{
					position350 := position
					depth++
					{
						position351 := position
						depth++
						{
							switch buffer[position] {
							case '$', '_':
								{
									position353, tokenIndex353, depth353 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l354
									}
									position++
									goto l353
								l354:
									position, tokenIndex, depth = position353, tokenIndex353, depth353
									if buffer[position] != rune('$') {
										goto l348
									}
									position++
								}
							l353:
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l348
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l348
								}
								position++
								break
							}
						}

						depth--
						add(ruleLetter, position351)
					}
				l355:
					{
						position356, tokenIndex356, depth356 := position, tokenIndex, depth
						{
							position357 := position
							depth++
							{
								switch buffer[position] {
								case '$', '_':
									{
										position359, tokenIndex359, depth359 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l360
										}
										position++
										goto l359
									l360:
										position, tokenIndex, depth = position359, tokenIndex359, depth359
										if buffer[position] != rune('$') {
											goto l356
										}
										position++
									}
								l359:
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l356
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l356
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l356
									}
									position++
									break
								}
							}

							depth--
							add(ruleLetterOrDigit, position357)
						}
						goto l355
					l356:
						position, tokenIndex, depth = position356, tokenIndex356, depth356
					}
					depth--
					add(rulePegText, position350)
				}
				if !_rules[ruleSpacing]() {
					goto l348
				}
				{
					add(ruleAction18, position)
				}
				depth--
				add(ruleIdentifier, position349)
			}
			return true
		l348:
			position, tokenIndex, depth = position348, tokenIndex348, depth348
			return false
		},
		/* 12 Letter <- <((&('$' | '_') ('_' / '$')) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 13 LetterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 14 EXIT <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('t' / 'T') Spacing Action19)> */
		nil,
		/* 15 RET <- <(('r' / 'R') ('e' / 'E') ('t' / 'T') Spacing Action20)> */
		nil,
		/* 16 NOP <- <(('n' / 'N') ('o' / 'O') ('p' / 'P') Spacing Action21)> */
		nil,
		/* 17 CALL <- <(('c' / 'C') ('a' / 'A') ('l' / 'L') ('l' / 'L') Space Action22)> */
		nil,
		/* 18 PUSH <- <(('p' / 'P') ('u' / 'U') ('s' / 'S') ('h' / 'H') Space Action23)> */
		nil,
		/* 19 POP <- <(('p' / 'P') ('o' / 'O') ('p' / 'P') Space Action24)> */
		nil,
		/* 20 JMP <- <(('j' / 'J') ('m' / 'M') ('p' / 'P') Space Action25)> */
		nil,
		/* 21 IN <- <(('i' / 'I') ('n' / 'N') Space Action26)> */
		nil,
		/* 22 OUT <- <(('o' / 'O') ('u' / 'U') ('t' / 'T') Space Action27)> */
		nil,
		/* 23 CAL <- <(('c' / 'C') ('a' / 'A') ('l' / 'L') Space Action28)> */
		nil,
		/* 24 LD <- <(('l' / 'L') ('d' / 'D') Space Action29)> */
		nil,
		/* 25 CMP <- <(('c' / 'C') ('m' / 'M') ('p' / 'P') Space Action30)> */
		nil,
		/* 26 JPC <- <(('j' / 'J') ('p' / 'P') ('c' / 'C') Space Action31)> */
		nil,
		/* 27 BLOCK <- <('.' ('b' / 'B') ('l' / 'L') ('o' / 'O') ('c' / 'C') ('k' / 'K') Space Action32)> */
		nil,
		/* 28 DATA <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') Space Action33)> */
		nil,
		/* 29 CAL_OP <- <(((('m' / 'M') ('u' / 'U') ('l' / 'L')) / ((&('M' | 'm') (('m' / 'M') ('o' / 'O') ('d' / 'D'))) | (&('D' | 'd') (('d' / 'D') ('i' / 'I') ('v' / 'V'))) | (&('S' | 's') (('s' / 'S') ('u' / 'U') ('b' / 'B'))) | (&('A' | 'a') (('a' / 'A') ('d' / 'D') ('d' / 'D'))))) Space)> */
		nil,
		/* 30 CMP_OP <- <(((('b' / 'B') ('e' / 'E')) / (('a' / 'A') ('e' / 'E')) / ((&('N' | 'n') (('n' / 'N') ('z' / 'Z'))) | (&('A' | 'a') ('a' / 'A')) | (&('Z') 'Z') | (&('z') 'z') | (&('B' | 'b') ('b' / 'B')))) Space)> */
		nil,
		/* 31 DATA_TYPE <- <(((&('I' | 'i') (('i' / 'I') ('n' / 'N') ('t' / 'T'))) | (&('F' | 'f') (('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))) | (&('B' | 'b') (('b' / 'B') ('y' / 'Y') ('t' / 'T') ('e' / 'E'))) | (&('W' | 'w') (('w' / 'W') ('o' / 'O') ('r' / 'R') ('d' / 'D'))) | (&('D' | 'd') (('d' / 'D') ('w' / 'W') ('o' / 'O') ('r' / 'R') ('d' / 'D')))) Space)> */
		func() bool {
			position381, tokenIndex381, depth381 := position, tokenIndex, depth
			{
				position382 := position
				depth++
				{
					switch buffer[position] {
					case 'I', 'i':
						{
							position384, tokenIndex384, depth384 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l385
							}
							position++
							goto l384
						l385:
							position, tokenIndex, depth = position384, tokenIndex384, depth384
							if buffer[position] != rune('I') {
								goto l381
							}
							position++
						}
					l384:
						{
							position386, tokenIndex386, depth386 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l387
							}
							position++
							goto l386
						l387:
							position, tokenIndex, depth = position386, tokenIndex386, depth386
							if buffer[position] != rune('N') {
								goto l381
							}
							position++
						}
					l386:
						{
							position388, tokenIndex388, depth388 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l389
							}
							position++
							goto l388
						l389:
							position, tokenIndex, depth = position388, tokenIndex388, depth388
							if buffer[position] != rune('T') {
								goto l381
							}
							position++
						}
					l388:
						break
					case 'F', 'f':
						{
							position390, tokenIndex390, depth390 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l391
							}
							position++
							goto l390
						l391:
							position, tokenIndex, depth = position390, tokenIndex390, depth390
							if buffer[position] != rune('F') {
								goto l381
							}
							position++
						}
					l390:
						{
							position392, tokenIndex392, depth392 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l393
							}
							position++
							goto l392
						l393:
							position, tokenIndex, depth = position392, tokenIndex392, depth392
							if buffer[position] != rune('L') {
								goto l381
							}
							position++
						}
					l392:
						{
							position394, tokenIndex394, depth394 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l395
							}
							position++
							goto l394
						l395:
							position, tokenIndex, depth = position394, tokenIndex394, depth394
							if buffer[position] != rune('O') {
								goto l381
							}
							position++
						}
					l394:
						{
							position396, tokenIndex396, depth396 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l397
							}
							position++
							goto l396
						l397:
							position, tokenIndex, depth = position396, tokenIndex396, depth396
							if buffer[position] != rune('A') {
								goto l381
							}
							position++
						}
					l396:
						{
							position398, tokenIndex398, depth398 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l399
							}
							position++
							goto l398
						l399:
							position, tokenIndex, depth = position398, tokenIndex398, depth398
							if buffer[position] != rune('T') {
								goto l381
							}
							position++
						}
					l398:
						break
					case 'B', 'b':
						{
							position400, tokenIndex400, depth400 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l401
							}
							position++
							goto l400
						l401:
							position, tokenIndex, depth = position400, tokenIndex400, depth400
							if buffer[position] != rune('B') {
								goto l381
							}
							position++
						}
					l400:
						{
							position402, tokenIndex402, depth402 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l403
							}
							position++
							goto l402
						l403:
							position, tokenIndex, depth = position402, tokenIndex402, depth402
							if buffer[position] != rune('Y') {
								goto l381
							}
							position++
						}
					l402:
						{
							position404, tokenIndex404, depth404 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l405
							}
							position++
							goto l404
						l405:
							position, tokenIndex, depth = position404, tokenIndex404, depth404
							if buffer[position] != rune('T') {
								goto l381
							}
							position++
						}
					l404:
						{
							position406, tokenIndex406, depth406 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l407
							}
							position++
							goto l406
						l407:
							position, tokenIndex, depth = position406, tokenIndex406, depth406
							if buffer[position] != rune('E') {
								goto l381
							}
							position++
						}
					l406:
						break
					case 'W', 'w':
						{
							position408, tokenIndex408, depth408 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l409
							}
							position++
							goto l408
						l409:
							position, tokenIndex, depth = position408, tokenIndex408, depth408
							if buffer[position] != rune('W') {
								goto l381
							}
							position++
						}
					l408:
						{
							position410, tokenIndex410, depth410 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l411
							}
							position++
							goto l410
						l411:
							position, tokenIndex, depth = position410, tokenIndex410, depth410
							if buffer[position] != rune('O') {
								goto l381
							}
							position++
						}
					l410:
						{
							position412, tokenIndex412, depth412 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l413
							}
							position++
							goto l412
						l413:
							position, tokenIndex, depth = position412, tokenIndex412, depth412
							if buffer[position] != rune('R') {
								goto l381
							}
							position++
						}
					l412:
						{
							position414, tokenIndex414, depth414 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l415
							}
							position++
							goto l414
						l415:
							position, tokenIndex, depth = position414, tokenIndex414, depth414
							if buffer[position] != rune('D') {
								goto l381
							}
							position++
						}
					l414:
						break
					default:
						{
							position416, tokenIndex416, depth416 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l417
							}
							position++
							goto l416
						l417:
							position, tokenIndex, depth = position416, tokenIndex416, depth416
							if buffer[position] != rune('D') {
								goto l381
							}
							position++
						}
					l416:
						{
							position418, tokenIndex418, depth418 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l419
							}
							position++
							goto l418
						l419:
							position, tokenIndex, depth = position418, tokenIndex418, depth418
							if buffer[position] != rune('W') {
								goto l381
							}
							position++
						}
					l418:
						{
							position420, tokenIndex420, depth420 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l421
							}
							position++
							goto l420
						l421:
							position, tokenIndex, depth = position420, tokenIndex420, depth420
							if buffer[position] != rune('O') {
								goto l381
							}
							position++
						}
					l420:
						{
							position422, tokenIndex422, depth422 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l423
							}
							position++
							goto l422
						l423:
							position, tokenIndex, depth = position422, tokenIndex422, depth422
							if buffer[position] != rune('R') {
								goto l381
							}
							position++
						}
					l422:
						{
							position424, tokenIndex424, depth424 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l425
							}
							position++
							goto l424
						l425:
							position, tokenIndex, depth = position424, tokenIndex424, depth424
							if buffer[position] != rune('D') {
								goto l381
							}
							position++
						}
					l424:
						break
					}
				}

				if !_rules[ruleSpace]() {
					goto l381
				}
				depth--
				add(ruleDATA_TYPE, position382)
			}
			return true
		l381:
			position, tokenIndex, depth = position381, tokenIndex381, depth381
			return false
		},
		/* 32 LBRK <- <('[' Spacing)> */
		func() bool {
			position426, tokenIndex426, depth426 := position, tokenIndex, depth
			{
				position427 := position
				depth++
				if buffer[position] != rune('[') {
					goto l426
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l426
				}
				depth--
				add(ruleLBRK, position427)
			}
			return true
		l426:
			position, tokenIndex, depth = position426, tokenIndex426, depth426
			return false
		},
		/* 33 RBRK <- <(']' Spacing)> */
		func() bool {
			position428, tokenIndex428, depth428 := position, tokenIndex, depth
			{
				position429 := position
				depth++
				if buffer[position] != rune(']') {
					goto l428
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l428
				}
				depth--
				add(ruleRBRK, position429)
			}
			return true
		l428:
			position, tokenIndex, depth = position428, tokenIndex428, depth428
			return false
		},
		/* 34 COMMA <- <(',' Spacing)> */
		func() bool {
			position430, tokenIndex430, depth430 := position, tokenIndex, depth
			{
				position431 := position
				depth++
				if buffer[position] != rune(',') {
					goto l430
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l430
				}
				depth--
				add(ruleCOMMA, position431)
			}
			return true
		l430:
			position, tokenIndex, depth = position430, tokenIndex430, depth430
			return false
		},
		/* 35 SEMICOLON <- <(';' Spacing)> */
		nil,
		/* 36 COLON <- <(':' Spacing)> */
		nil,
		/* 37 MINUS <- <('-' Spacing)> */
		nil,
		/* 38 NL <- <'\n'> */
		func() bool {
			position435, tokenIndex435, depth435 := position, tokenIndex, depth
			{
				position436 := position
				depth++
				if buffer[position] != rune('\n') {
					goto l435
				}
				position++
				depth--
				add(ruleNL, position436)
			}
			return true
		l435:
			position, tokenIndex, depth = position435, tokenIndex435, depth435
			return false
		},
		/* 39 EOT <- <!.> */
		nil,
		/* 40 Literal <- <((FloatLiteral / ((&('"') StringLiteral) | (&('\'') CharLiteral) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') IntegerLiteral))) Spacing)> */
		nil,
		/* 41 IntegerLiteral <- <(<(MINUS? (HexNumeral / BinaryNumeral / OctalNumeral / DecimalNumeral))> Spacing Action34)> */
		func() bool {
			position439, tokenIndex439, depth439 := position, tokenIndex, depth
			{
				position440 := position
				depth++
				{
					position441 := position
					depth++
					{
						position442, tokenIndex442, depth442 := position, tokenIndex, depth
						{
							position444 := position
							depth++
							if buffer[position] != rune('-') {
								goto l442
							}
							position++
							if !_rules[ruleSpacing]() {
								goto l442
							}
							depth--
							add(ruleMINUS, position444)
						}
						goto l443
					l442:
						position, tokenIndex, depth = position442, tokenIndex442, depth442
					}
				l443:
					{
						position445, tokenIndex445, depth445 := position, tokenIndex, depth
						if !_rules[ruleHexNumeral]() {
							goto l446
						}
						goto l445
					l446:
						position, tokenIndex, depth = position445, tokenIndex445, depth445
						{
							position448 := position
							depth++
							{
								position449, tokenIndex449, depth449 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l450
								}
								position++
								if buffer[position] != rune('b') {
									goto l450
								}
								position++
								goto l449
							l450:
								position, tokenIndex, depth = position449, tokenIndex449, depth449
								if buffer[position] != rune('0') {
									goto l447
								}
								position++
								if buffer[position] != rune('B') {
									goto l447
								}
								position++
							}
						l449:
							{
								position451, tokenIndex451, depth451 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l452
								}
								position++
								goto l451
							l452:
								position, tokenIndex, depth = position451, tokenIndex451, depth451
								if buffer[position] != rune('1') {
									goto l447
								}
								position++
							}
						l451:
						l453:
							{
								position454, tokenIndex454, depth454 := position, tokenIndex, depth
							l455:
								{
									position456, tokenIndex456, depth456 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l456
									}
									position++
									goto l455
								l456:
									position, tokenIndex, depth = position456, tokenIndex456, depth456
								}
								{
									position457, tokenIndex457, depth457 := position, tokenIndex, depth
									if buffer[position] != rune('0') {
										goto l458
									}
									position++
									goto l457
								l458:
									position, tokenIndex, depth = position457, tokenIndex457, depth457
									if buffer[position] != rune('1') {
										goto l454
									}
									position++
								}
							l457:
								goto l453
							l454:
								position, tokenIndex, depth = position454, tokenIndex454, depth454
							}
							depth--
							add(ruleBinaryNumeral, position448)
						}
						goto l445
					l447:
						position, tokenIndex, depth = position445, tokenIndex445, depth445
						{
							position460 := position
							depth++
							if buffer[position] != rune('0') {
								goto l459
							}
							position++
						l463:
							{
								position464, tokenIndex464, depth464 := position, tokenIndex, depth
								if buffer[position] != rune('_') {
									goto l464
								}
								position++
								goto l463
							l464:
								position, tokenIndex, depth = position464, tokenIndex464, depth464
							}
							if c := buffer[position]; c < rune('0') || c > rune('7') {
								goto l459
							}
							position++
						l461:
							{
								position462, tokenIndex462, depth462 := position, tokenIndex, depth
							l465:
								{
									position466, tokenIndex466, depth466 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l466
									}
									position++
									goto l465
								l466:
									position, tokenIndex, depth = position466, tokenIndex466, depth466
								}
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l462
								}
								position++
								goto l461
							l462:
								position, tokenIndex, depth = position462, tokenIndex462, depth462
							}
							depth--
							add(ruleOctalNumeral, position460)
						}
						goto l445
					l459:
						position, tokenIndex, depth = position445, tokenIndex445, depth445
						{
							position467 := position
							depth++
							{
								position468, tokenIndex468, depth468 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l469
								}
								position++
								goto l468
							l469:
								position, tokenIndex, depth = position468, tokenIndex468, depth468
								if c := buffer[position]; c < rune('1') || c > rune('9') {
									goto l439
								}
								position++
							l470:
								{
									position471, tokenIndex471, depth471 := position, tokenIndex, depth
								l472:
									{
										position473, tokenIndex473, depth473 := position, tokenIndex, depth
										if buffer[position] != rune('_') {
											goto l473
										}
										position++
										goto l472
									l473:
										position, tokenIndex, depth = position473, tokenIndex473, depth473
									}
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l471
									}
									position++
									goto l470
								l471:
									position, tokenIndex, depth = position471, tokenIndex471, depth471
								}
							}
						l468:
							depth--
							add(ruleDecimalNumeral, position467)
						}
					}
				l445:
					depth--
					add(rulePegText, position441)
				}
				if !_rules[ruleSpacing]() {
					goto l439
				}
				{
					add(ruleAction34, position)
				}
				depth--
				add(ruleIntegerLiteral, position440)
			}
			return true
		l439:
			position, tokenIndex, depth = position439, tokenIndex439, depth439
			return false
		},
		/* 42 DecimalNumeral <- <('0' / ([1-9] ('_'* [0-9])*))> */
		nil,
		/* 43 HexNumeral <- <((('0' 'x') / ('0' 'X')) HexDigits)> */
		func() bool {
			position476, tokenIndex476, depth476 := position, tokenIndex, depth
			{
				position477 := position
				depth++
				{
					position478, tokenIndex478, depth478 := position, tokenIndex, depth
					if buffer[position] != rune('0') {
						goto l479
					}
					position++
					if buffer[position] != rune('x') {
						goto l479
					}
					position++
					goto l478
				l479:
					position, tokenIndex, depth = position478, tokenIndex478, depth478
					if buffer[position] != rune('0') {
						goto l476
					}
					position++
					if buffer[position] != rune('X') {
						goto l476
					}
					position++
				}
			l478:
				if !_rules[ruleHexDigits]() {
					goto l476
				}
				depth--
				add(ruleHexNumeral, position477)
			}
			return true
		l476:
			position, tokenIndex, depth = position476, tokenIndex476, depth476
			return false
		},
		/* 44 BinaryNumeral <- <((('0' 'b') / ('0' 'B')) ('0' / '1') ('_'* ('0' / '1'))*)> */
		nil,
		/* 45 OctalNumeral <- <('0' ('_'* [0-7])+)> */
		nil,
		/* 46 FloatLiteral <- <(HexFloat / DecimalFloat)> */
		nil,
		/* 47 DecimalFloat <- <((Digits '.' Digits? Exponent? ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?) / ('.' Digits Exponent? ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?) / (Digits Exponent ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?) / (Digits Exponent? ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))))> */
		nil,
		/* 48 Exponent <- <(('e' / 'E') ('+' / '-')? Digits)> */
		func() bool {
			position484, tokenIndex484, depth484 := position, tokenIndex, depth
			{
				position485 := position
				depth++
				{
					position486, tokenIndex486, depth486 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l487
					}
					position++
					goto l486
				l487:
					position, tokenIndex, depth = position486, tokenIndex486, depth486
					if buffer[position] != rune('E') {
						goto l484
					}
					position++
				}
			l486:
				{
					position488, tokenIndex488, depth488 := position, tokenIndex, depth
					{
						position490, tokenIndex490, depth490 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l491
						}
						position++
						goto l490
					l491:
						position, tokenIndex, depth = position490, tokenIndex490, depth490
						if buffer[position] != rune('-') {
							goto l488
						}
						position++
					}
				l490:
					goto l489
				l488:
					position, tokenIndex, depth = position488, tokenIndex488, depth488
				}
			l489:
				if !_rules[ruleDigits]() {
					goto l484
				}
				depth--
				add(ruleExponent, position485)
			}
			return true
		l484:
			position, tokenIndex, depth = position484, tokenIndex484, depth484
			return false
		},
		/* 49 HexFloat <- <(HexSignificand BinaryExponent ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?)> */
		nil,
		/* 50 HexSignificand <- <(((('0' 'x') / ('0' 'X')) HexDigits? '.' HexDigits) / (HexNumeral '.'?))> */
		nil,
		/* 51 BinaryExponent <- <(('p' / 'P') ('+' / '-')? Digits)> */
		nil,
		/* 52 Digits <- <([0-9] ('_'* [0-9])*)> */
		func() bool {
			position495, tokenIndex495, depth495 := position, tokenIndex, depth
			{
				position496 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l495
				}
				position++
			l497:
				{
					position498, tokenIndex498, depth498 := position, tokenIndex, depth
				l499:
					{
						position500, tokenIndex500, depth500 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l500
						}
						position++
						goto l499
					l500:
						position, tokenIndex, depth = position500, tokenIndex500, depth500
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l498
					}
					position++
					goto l497
				l498:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
				}
				depth--
				add(ruleDigits, position496)
			}
			return true
		l495:
			position, tokenIndex, depth = position495, tokenIndex495, depth495
			return false
		},
		/* 53 HexDigits <- <(HexDigit ('_'* HexDigit)*)> */
		func() bool {
			position501, tokenIndex501, depth501 := position, tokenIndex, depth
			{
				position502 := position
				depth++
				if !_rules[ruleHexDigit]() {
					goto l501
				}
			l503:
				{
					position504, tokenIndex504, depth504 := position, tokenIndex, depth
				l505:
					{
						position506, tokenIndex506, depth506 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l506
						}
						position++
						goto l505
					l506:
						position, tokenIndex, depth = position506, tokenIndex506, depth506
					}
					if !_rules[ruleHexDigit]() {
						goto l504
					}
					goto l503
				l504:
					position, tokenIndex, depth = position504, tokenIndex504, depth504
				}
				depth--
				add(ruleHexDigits, position502)
			}
			return true
		l501:
			position, tokenIndex, depth = position501, tokenIndex501, depth501
			return false
		},
		/* 54 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position507, tokenIndex507, depth507 := position, tokenIndex, depth
			{
				position508 := position
				depth++
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l507
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l507
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l507
						}
						position++
						break
					}
				}

				depth--
				add(ruleHexDigit, position508)
			}
			return true
		l507:
			position, tokenIndex, depth = position507, tokenIndex507, depth507
			return false
		},
		/* 55 CharLiteral <- <('\'' (Escape / (!('\'' / '\\') .)) '\'')> */
		nil,
		/* 56 StringLiteral <- <(<('"' (Escape / (!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('"') '"')) .))* '"')> Action35)> */
		func() bool {
			position511, tokenIndex511, depth511 := position, tokenIndex, depth
			{
				position512 := position
				depth++
				{
					position513 := position
					depth++
					if buffer[position] != rune('"') {
						goto l511
					}
					position++
				l514:
					{
						position515, tokenIndex515, depth515 := position, tokenIndex, depth
						{
							position516, tokenIndex516, depth516 := position, tokenIndex, depth
							if !_rules[ruleEscape]() {
								goto l517
							}
							goto l516
						l517:
							position, tokenIndex, depth = position516, tokenIndex516, depth516
							{
								position518, tokenIndex518, depth518 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '\r':
										if buffer[position] != rune('\r') {
											goto l518
										}
										position++
										break
									case '\n':
										if buffer[position] != rune('\n') {
											goto l518
										}
										position++
										break
									case '\\':
										if buffer[position] != rune('\\') {
											goto l518
										}
										position++
										break
									default:
										if buffer[position] != rune('"') {
											goto l518
										}
										position++
										break
									}
								}

								goto l515
							l518:
								position, tokenIndex, depth = position518, tokenIndex518, depth518
							}
							if !matchDot() {
								goto l515
							}
						}
					l516:
						goto l514
					l515:
						position, tokenIndex, depth = position515, tokenIndex515, depth515
					}
					if buffer[position] != rune('"') {
						goto l511
					}
					position++
					depth--
					add(rulePegText, position513)
				}
				{
					add(ruleAction35, position)
				}
				depth--
				add(ruleStringLiteral, position512)
			}
			return true
		l511:
			position, tokenIndex, depth = position511, tokenIndex511, depth511
			return false
		},
		/* 57 Escape <- <('\\' ((&('u') UnicodeEscape) | (&('\\') '\\') | (&('\'') '\'') | (&('"') '"') | (&('r') 'r') | (&('f') 'f') | (&('n') 'n') | (&('t') 't') | (&('b') 'b') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7') OctalEscape)))> */
		func() bool {
			position521, tokenIndex521, depth521 := position, tokenIndex, depth
			{
				position522 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l521
				}
				position++
				{
					switch buffer[position] {
					case 'u':
						{
							position524 := position
							depth++
							if buffer[position] != rune('u') {
								goto l521
							}
							position++
						l525:
							{
								position526, tokenIndex526, depth526 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l526
								}
								position++
								goto l525
							l526:
								position, tokenIndex, depth = position526, tokenIndex526, depth526
							}
							if !_rules[ruleHexDigit]() {
								goto l521
							}
							if !_rules[ruleHexDigit]() {
								goto l521
							}
							if !_rules[ruleHexDigit]() {
								goto l521
							}
							if !_rules[ruleHexDigit]() {
								goto l521
							}
							depth--
							add(ruleUnicodeEscape, position524)
						}
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l521
						}
						position++
						break
					case '\'':
						if buffer[position] != rune('\'') {
							goto l521
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l521
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l521
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l521
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l521
						}
						position++
						break
					case 't':
						if buffer[position] != rune('t') {
							goto l521
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l521
						}
						position++
						break
					default:
						{
							position527 := position
							depth++
							{
								position528, tokenIndex528, depth528 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('3') {
									goto l529
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l529
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l529
								}
								position++
								goto l528
							l529:
								position, tokenIndex, depth = position528, tokenIndex528, depth528
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l530
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l530
								}
								position++
								goto l528
							l530:
								position, tokenIndex, depth = position528, tokenIndex528, depth528
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l521
								}
								position++
							}
						l528:
							depth--
							add(ruleOctalEscape, position527)
						}
						break
					}
				}

				depth--
				add(ruleEscape, position522)
			}
			return true
		l521:
			position, tokenIndex, depth = position521, tokenIndex521, depth521
			return false
		},
		/* 58 OctalEscape <- <(([0-3] [0-7] [0-7]) / ([0-7] [0-7]) / [0-7])> */
		nil,
		/* 59 UnicodeEscape <- <('u'+ HexDigit HexDigit HexDigit HexDigit)> */
		nil,
		/* 61 Action0 <- <{p.line++}> */
		nil,
		/* 62 Action1 <- <{p.AddAssembly()}> */
		nil,
		/* 63 Action2 <- <{p.AddAssembly();p.AddComment()}> */
		nil,
		nil,
		/* 65 Action3 <- <{p.Push(&asm.Comment{});p.Push(text)}> */
		nil,
		/* 66 Action4 <- <{p.Push(&asm.Label{})}> */
		nil,
		/* 67 Action5 <- <{p.Push(lookup(asm.T_INT,text))}> */
		nil,
		/* 68 Action6 <- <{p.Push(lookup(asm.CAL_ADD,text))}> */
		nil,
		/* 69 Action7 <- <{p.Push(lookup(asm.T_INT,text))}> */
		nil,
		/* 70 Action8 <- <{p.Push(lookup(asm.T_INT,text))}> */
		nil,
		/* 71 Action9 <- <{p.Push(lookup(asm.CMP_A,text))}> */
		nil,
		/* 72 Action10 <- <{p.AddPseudoDataValue()}> */
		nil,
		/* 73 Action11 <- <{p.AddPseudoDataValue()}> */
		nil,
		/* 74 Action12 <- <{p.AddPseudoDataValue()}> */
		nil,
		/* 75 Action13 <- <{p.Push(text);p.AddPseudoDataValue()}> */
		nil,
		/* 76 Action14 <- <{p.AddOperand(true)}> */
		nil,
		/* 77 Action15 <- <{p.AddOperand(false)}> */
		nil,
		/* 78 Action16 <- <{p.AddOperand(true)}> */
		nil,
		/* 79 Action17 <- <{p.AddOperand(false)}> */
		nil,
		/* 80 Action18 <- <{p.Push(text)}> */
		nil,
		/* 81 Action19 <- <{p.PushInst(asm.OP_EXIT)}> */
		nil,
		/* 82 Action20 <- <{p.PushInst(asm.OP_RET)}> */
		nil,
		/* 83 Action21 <- <{p.PushInst(asm.OP_NOP)}> */
		nil,
		/* 84 Action22 <- <{p.PushInst(asm.OP_CALL)}> */
		nil,
		/* 85 Action23 <- <{p.PushInst(asm.OP_PUSH)}> */
		nil,
		/* 86 Action24 <- <{p.PushInst(asm.OP_POP)}> */
		nil,
		/* 87 Action25 <- <{p.PushInst(asm.OP_JMP)}> */
		nil,
		/* 88 Action26 <- <{p.PushInst(asm.OP_IN)}> */
		nil,
		/* 89 Action27 <- <{p.PushInst(asm.OP_OUT)}> */
		nil,
		/* 90 Action28 <- <{p.PushInst(asm.OP_CAL)}> */
		nil,
		/* 91 Action29 <- <{p.PushInst(asm.OP_LD)}> */
		nil,
		/* 92 Action30 <- <{p.PushInst(asm.OP_CMP)}> */
		nil,
		/* 93 Action31 <- <{p.PushInst(asm.OP_JPC)}> */
		nil,
		/* 94 Action32 <- <{p.Push(&asm.PseudoBlock{})}> */
		nil,
		/* 95 Action33 <- <{p.Push(&asm.PseudoData{})}> */
		nil,
		/* 96 Action34 <- <{p.Push(text);p.AddInteger()}> */
		nil,
		/* 97 Action35 <- <{p.Push(text)}> */
		nil,
	}
	p.rules = _rules
}
