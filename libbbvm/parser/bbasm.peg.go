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
	ruleCAL_OP
	ruleCMP_OP
	ruleDATA_TYPE
	ruleLBRK
	ruleRBRK
	ruleCOMMA
	ruleSEMICOLON
	ruleCOLON
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
	"CAL_OP",
	"CMP_OP",
	"DATA_TYPE",
	"LBRK",
	"RBRK",
	"COMMA",
	"SEMICOLON",
	"COLON",
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
	rules  [88]func() bool
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
			/*strconv.Quote(*/
			e.p.Buffer[begin:end] /*)*/)
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
			p.PopAssembly()
		case ruleAction2:
			p.PopAssembly()
			p.AddComment()
		case ruleAction3:
			p.Push(&asm.Comment{})
			p.Push(text)
		case ruleAction4:
			p.Push(&asm.Label{})
			p.Push(text)
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
			p.Push(text)
		case ruleAction11:
			p.Push(text)
		case ruleAction12:
			p.Push(CreateOperand(text, false, true))
		case ruleAction13:
			p.Push(CreateOperand(text, false, false))
		case ruleAction14:
			p.Push(CreateOperand(text, true, true))
		case ruleAction15:
			p.Push(CreateOperand(text, true, false))
		case ruleAction16:
			p.PushInst(asm.OP_EXIT)
		case ruleAction17:
			p.PushInst(asm.OP_RET)
		case ruleAction18:
			p.PushInst(asm.OP_NOP)
		case ruleAction19:
			p.PushInst(asm.OP_CALL)
		case ruleAction20:
			p.PushInst(asm.OP_PUSH)
		case ruleAction21:
			p.PushInst(asm.OP_POP)
		case ruleAction22:
			p.PushInst(asm.OP_JMP)
		case ruleAction23:
			p.PushInst(asm.OP_IN)
		case ruleAction24:
			p.PushInst(asm.OP_OUT)
		case ruleAction25:
			p.PushInst(asm.OP_CAL)
		case ruleAction26:
			p.PushInst(asm.OP_LD)
		case ruleAction27:
			p.PushInst(asm.OP_CMP)
		case ruleAction28:
			p.PushInst(asm.OP_JPC)
		case ruleAction29:
			p.Push(&asm.PseudoBlock{})

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
										position10 := position
										depth++
										if !_rules[ruleIdentifier]() {
											goto l8
										}
										depth--
										add(rulePegText, position10)
									}
									{
										add(ruleAction4, position)
									}
									if !_rules[ruleSpacing]() {
										goto l8
									}
									{
										position12 := position
										depth++
										if buffer[position] != rune(':') {
											goto l8
										}
										position++
										if !_rules[ruleSpacing]() {
											goto l8
										}
										depth--
										add(ruleCOLON, position12)
									}
									depth--
									add(ruleLabel, position9)
								}
								goto l7
							l8:
								position, tokenIndex, depth = position7, tokenIndex7, depth7
								{
									switch buffer[position] {
									case '.':
										{
											position14 := position
											depth++
											{
												position15 := position
												depth++
												if buffer[position] != rune('.') {
													goto l4
												}
												position++
												{
													position16, tokenIndex16, depth16 := position, tokenIndex, depth
													if buffer[position] != rune('b') {
														goto l17
													}
													position++
													goto l16
												l17:
													position, tokenIndex, depth = position16, tokenIndex16, depth16
													if buffer[position] != rune('B') {
														goto l4
													}
													position++
												}
											l16:
												{
													position18, tokenIndex18, depth18 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l19
													}
													position++
													goto l18
												l19:
													position, tokenIndex, depth = position18, tokenIndex18, depth18
													if buffer[position] != rune('L') {
														goto l4
													}
													position++
												}
											l18:
												{
													position20, tokenIndex20, depth20 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l21
													}
													position++
													goto l20
												l21:
													position, tokenIndex, depth = position20, tokenIndex20, depth20
													if buffer[position] != rune('O') {
														goto l4
													}
													position++
												}
											l20:
												{
													position22, tokenIndex22, depth22 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l23
													}
													position++
													goto l22
												l23:
													position, tokenIndex, depth = position22, tokenIndex22, depth22
													if buffer[position] != rune('C') {
														goto l4
													}
													position++
												}
											l22:
												{
													position24, tokenIndex24, depth24 := position, tokenIndex, depth
													if buffer[position] != rune('k') {
														goto l25
													}
													position++
													goto l24
												l25:
													position, tokenIndex, depth = position24, tokenIndex24, depth24
													if buffer[position] != rune('K') {
														goto l4
													}
													position++
												}
											l24:
												if !_rules[ruleSpace]() {
													goto l4
												}
												{
													add(ruleAction29, position)
												}
												depth--
												add(ruleBLOCK, position15)
											}
											{
												position27 := position
												depth++
												if !_rules[ruleIntegerLiteral]() {
													goto l4
												}
												depth--
												add(rulePegText, position27)
											}
											{
												add(ruleAction10, position)
											}
											if !_rules[ruleSpace]() {
												goto l4
											}
											{
												position29 := position
												depth++
												if !_rules[ruleIntegerLiteral]() {
													goto l4
												}
												depth--
												add(rulePegText, position29)
											}
											{
												add(ruleAction11, position)
											}
											depth--
											add(rulePseudo, position14)
										}
										break
									case ';':
										if !_rules[ruleComment]() {
											goto l4
										}
										break
									default:
										{
											position31 := position
											depth++
											{
												position32, tokenIndex32, depth32 := position, tokenIndex, depth
												{
													position34, tokenIndex34, depth34 := position, tokenIndex, depth
													{
														position36 := position
														depth++
														{
															position37, tokenIndex37, depth37 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l38
															}
															position++
															goto l37
														l38:
															position, tokenIndex, depth = position37, tokenIndex37, depth37
															if buffer[position] != rune('P') {
																goto l35
															}
															position++
														}
													l37:
														{
															position39, tokenIndex39, depth39 := position, tokenIndex, depth
															if buffer[position] != rune('u') {
																goto l40
															}
															position++
															goto l39
														l40:
															position, tokenIndex, depth = position39, tokenIndex39, depth39
															if buffer[position] != rune('U') {
																goto l35
															}
															position++
														}
													l39:
														{
															position41, tokenIndex41, depth41 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l42
															}
															position++
															goto l41
														l42:
															position, tokenIndex, depth = position41, tokenIndex41, depth41
															if buffer[position] != rune('S') {
																goto l35
															}
															position++
														}
													l41:
														{
															position43, tokenIndex43, depth43 := position, tokenIndex, depth
															if buffer[position] != rune('h') {
																goto l44
															}
															position++
															goto l43
														l44:
															position, tokenIndex, depth = position43, tokenIndex43, depth43
															if buffer[position] != rune('H') {
																goto l35
															}
															position++
														}
													l43:
														if !_rules[ruleSpace]() {
															goto l35
														}
														{
															add(ruleAction20, position)
														}
														depth--
														add(rulePUSH, position36)
													}
													goto l34
												l35:
													position, tokenIndex, depth = position34, tokenIndex34, depth34
													{
														switch buffer[position] {
														case 'J', 'j':
															{
																position47 := position
																depth++
																{
																	position48, tokenIndex48, depth48 := position, tokenIndex, depth
																	if buffer[position] != rune('j') {
																		goto l49
																	}
																	position++
																	goto l48
																l49:
																	position, tokenIndex, depth = position48, tokenIndex48, depth48
																	if buffer[position] != rune('J') {
																		goto l33
																	}
																	position++
																}
															l48:
																{
																	position50, tokenIndex50, depth50 := position, tokenIndex, depth
																	if buffer[position] != rune('m') {
																		goto l51
																	}
																	position++
																	goto l50
																l51:
																	position, tokenIndex, depth = position50, tokenIndex50, depth50
																	if buffer[position] != rune('M') {
																		goto l33
																	}
																	position++
																}
															l50:
																{
																	position52, tokenIndex52, depth52 := position, tokenIndex, depth
																	if buffer[position] != rune('p') {
																		goto l53
																	}
																	position++
																	goto l52
																l53:
																	position, tokenIndex, depth = position52, tokenIndex52, depth52
																	if buffer[position] != rune('P') {
																		goto l33
																	}
																	position++
																}
															l52:
																if !_rules[ruleSpace]() {
																	goto l33
																}
																{
																	add(ruleAction22, position)
																}
																depth--
																add(ruleJMP, position47)
															}
															break
														case 'P', 'p':
															{
																position55 := position
																depth++
																{
																	position56, tokenIndex56, depth56 := position, tokenIndex, depth
																	if buffer[position] != rune('p') {
																		goto l57
																	}
																	position++
																	goto l56
																l57:
																	position, tokenIndex, depth = position56, tokenIndex56, depth56
																	if buffer[position] != rune('P') {
																		goto l33
																	}
																	position++
																}
															l56:
																{
																	position58, tokenIndex58, depth58 := position, tokenIndex, depth
																	if buffer[position] != rune('o') {
																		goto l59
																	}
																	position++
																	goto l58
																l59:
																	position, tokenIndex, depth = position58, tokenIndex58, depth58
																	if buffer[position] != rune('O') {
																		goto l33
																	}
																	position++
																}
															l58:
																{
																	position60, tokenIndex60, depth60 := position, tokenIndex, depth
																	if buffer[position] != rune('p') {
																		goto l61
																	}
																	position++
																	goto l60
																l61:
																	position, tokenIndex, depth = position60, tokenIndex60, depth60
																	if buffer[position] != rune('P') {
																		goto l33
																	}
																	position++
																}
															l60:
																if !_rules[ruleSpace]() {
																	goto l33
																}
																{
																	add(ruleAction21, position)
																}
																depth--
																add(rulePOP, position55)
															}
															break
														default:
															{
																position63 := position
																depth++
																{
																	position64, tokenIndex64, depth64 := position, tokenIndex, depth
																	if buffer[position] != rune('c') {
																		goto l65
																	}
																	position++
																	goto l64
																l65:
																	position, tokenIndex, depth = position64, tokenIndex64, depth64
																	if buffer[position] != rune('C') {
																		goto l33
																	}
																	position++
																}
															l64:
																{
																	position66, tokenIndex66, depth66 := position, tokenIndex, depth
																	if buffer[position] != rune('a') {
																		goto l67
																	}
																	position++
																	goto l66
																l67:
																	position, tokenIndex, depth = position66, tokenIndex66, depth66
																	if buffer[position] != rune('A') {
																		goto l33
																	}
																	position++
																}
															l66:
																{
																	position68, tokenIndex68, depth68 := position, tokenIndex, depth
																	if buffer[position] != rune('l') {
																		goto l69
																	}
																	position++
																	goto l68
																l69:
																	position, tokenIndex, depth = position68, tokenIndex68, depth68
																	if buffer[position] != rune('L') {
																		goto l33
																	}
																	position++
																}
															l68:
																{
																	position70, tokenIndex70, depth70 := position, tokenIndex, depth
																	if buffer[position] != rune('l') {
																		goto l71
																	}
																	position++
																	goto l70
																l71:
																	position, tokenIndex, depth = position70, tokenIndex70, depth70
																	if buffer[position] != rune('L') {
																		goto l33
																	}
																	position++
																}
															l70:
																if !_rules[ruleSpace]() {
																	goto l33
																}
																{
																	add(ruleAction19, position)
																}
																depth--
																add(ruleCALL, position63)
															}
															break
														}
													}

												}
											l34:
												if !_rules[ruleOperand]() {
													goto l33
												}
												goto l32
											l33:
												position, tokenIndex, depth = position32, tokenIndex32, depth32
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
															goto l73
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
															goto l73
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
															goto l73
														}
														position++
													}
												l79:
													if !_rules[ruleSpace]() {
														goto l73
													}
													{
														add(ruleAction25, position)
													}
													depth--
													add(ruleCAL, position74)
												}
												{
													position82 := position
													depth++
													if !_rules[ruleDATA_TYPE]() {
														goto l73
													}
													depth--
													add(rulePegText, position82)
												}
												{
													add(ruleAction5, position)
												}
												{
													position84 := position
													depth++
													{
														position85 := position
														depth++
														{
															position86, tokenIndex86, depth86 := position, tokenIndex, depth
															{
																position88, tokenIndex88, depth88 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l89
																}
																position++
																goto l88
															l89:
																position, tokenIndex, depth = position88, tokenIndex88, depth88
																if buffer[position] != rune('M') {
																	goto l87
																}
																position++
															}
														l88:
															{
																position90, tokenIndex90, depth90 := position, tokenIndex, depth
																if buffer[position] != rune('u') {
																	goto l91
																}
																position++
																goto l90
															l91:
																position, tokenIndex, depth = position90, tokenIndex90, depth90
																if buffer[position] != rune('U') {
																	goto l87
																}
																position++
															}
														l90:
															{
																position92, tokenIndex92, depth92 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l93
																}
																position++
																goto l92
															l93:
																position, tokenIndex, depth = position92, tokenIndex92, depth92
																if buffer[position] != rune('L') {
																	goto l87
																}
																position++
															}
														l92:
															goto l86
														l87:
															position, tokenIndex, depth = position86, tokenIndex86, depth86
															{
																switch buffer[position] {
																case 'M', 'm':
																	{
																		position95, tokenIndex95, depth95 := position, tokenIndex, depth
																		if buffer[position] != rune('m') {
																			goto l96
																		}
																		position++
																		goto l95
																	l96:
																		position, tokenIndex, depth = position95, tokenIndex95, depth95
																		if buffer[position] != rune('M') {
																			goto l73
																		}
																		position++
																	}
																l95:
																	{
																		position97, tokenIndex97, depth97 := position, tokenIndex, depth
																		if buffer[position] != rune('o') {
																			goto l98
																		}
																		position++
																		goto l97
																	l98:
																		position, tokenIndex, depth = position97, tokenIndex97, depth97
																		if buffer[position] != rune('O') {
																			goto l73
																		}
																		position++
																	}
																l97:
																	{
																		position99, tokenIndex99, depth99 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l100
																		}
																		position++
																		goto l99
																	l100:
																		position, tokenIndex, depth = position99, tokenIndex99, depth99
																		if buffer[position] != rune('D') {
																			goto l73
																		}
																		position++
																	}
																l99:
																	break
																case 'D', 'd':
																	{
																		position101, tokenIndex101, depth101 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l102
																		}
																		position++
																		goto l101
																	l102:
																		position, tokenIndex, depth = position101, tokenIndex101, depth101
																		if buffer[position] != rune('D') {
																			goto l73
																		}
																		position++
																	}
																l101:
																	{
																		position103, tokenIndex103, depth103 := position, tokenIndex, depth
																		if buffer[position] != rune('i') {
																			goto l104
																		}
																		position++
																		goto l103
																	l104:
																		position, tokenIndex, depth = position103, tokenIndex103, depth103
																		if buffer[position] != rune('I') {
																			goto l73
																		}
																		position++
																	}
																l103:
																	{
																		position105, tokenIndex105, depth105 := position, tokenIndex, depth
																		if buffer[position] != rune('v') {
																			goto l106
																		}
																		position++
																		goto l105
																	l106:
																		position, tokenIndex, depth = position105, tokenIndex105, depth105
																		if buffer[position] != rune('V') {
																			goto l73
																		}
																		position++
																	}
																l105:
																	break
																case 'S', 's':
																	{
																		position107, tokenIndex107, depth107 := position, tokenIndex, depth
																		if buffer[position] != rune('s') {
																			goto l108
																		}
																		position++
																		goto l107
																	l108:
																		position, tokenIndex, depth = position107, tokenIndex107, depth107
																		if buffer[position] != rune('S') {
																			goto l73
																		}
																		position++
																	}
																l107:
																	{
																		position109, tokenIndex109, depth109 := position, tokenIndex, depth
																		if buffer[position] != rune('u') {
																			goto l110
																		}
																		position++
																		goto l109
																	l110:
																		position, tokenIndex, depth = position109, tokenIndex109, depth109
																		if buffer[position] != rune('U') {
																			goto l73
																		}
																		position++
																	}
																l109:
																	{
																		position111, tokenIndex111, depth111 := position, tokenIndex, depth
																		if buffer[position] != rune('b') {
																			goto l112
																		}
																		position++
																		goto l111
																	l112:
																		position, tokenIndex, depth = position111, tokenIndex111, depth111
																		if buffer[position] != rune('B') {
																			goto l73
																		}
																		position++
																	}
																l111:
																	break
																default:
																	{
																		position113, tokenIndex113, depth113 := position, tokenIndex, depth
																		if buffer[position] != rune('a') {
																			goto l114
																		}
																		position++
																		goto l113
																	l114:
																		position, tokenIndex, depth = position113, tokenIndex113, depth113
																		if buffer[position] != rune('A') {
																			goto l73
																		}
																		position++
																	}
																l113:
																	{
																		position115, tokenIndex115, depth115 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l116
																		}
																		position++
																		goto l115
																	l116:
																		position, tokenIndex, depth = position115, tokenIndex115, depth115
																		if buffer[position] != rune('D') {
																			goto l73
																		}
																		position++
																	}
																l115:
																	{
																		position117, tokenIndex117, depth117 := position, tokenIndex, depth
																		if buffer[position] != rune('d') {
																			goto l118
																		}
																		position++
																		goto l117
																	l118:
																		position, tokenIndex, depth = position117, tokenIndex117, depth117
																		if buffer[position] != rune('D') {
																			goto l73
																		}
																		position++
																	}
																l117:
																	break
																}
															}

														}
													l86:
														if !_rules[ruleSpace]() {
															goto l73
														}
														depth--
														add(ruleCAL_OP, position85)
													}
													depth--
													add(rulePegText, position84)
												}
												{
													add(ruleAction6, position)
												}
												if !_rules[ruleOperand]() {
													goto l73
												}
												if !_rules[ruleCOMMA]() {
													goto l73
												}
												if !_rules[ruleOperand]() {
													goto l73
												}
												goto l32
											l73:
												position, tokenIndex, depth = position32, tokenIndex32, depth32
												{
													switch buffer[position] {
													case 'J', 'j':
														{
															position121 := position
															depth++
															{
																position122, tokenIndex122, depth122 := position, tokenIndex, depth
																if buffer[position] != rune('j') {
																	goto l123
																}
																position++
																goto l122
															l123:
																position, tokenIndex, depth = position122, tokenIndex122, depth122
																if buffer[position] != rune('J') {
																	goto l4
																}
																position++
															}
														l122:
															{
																position124, tokenIndex124, depth124 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l125
																}
																position++
																goto l124
															l125:
																position, tokenIndex, depth = position124, tokenIndex124, depth124
																if buffer[position] != rune('P') {
																	goto l4
																}
																position++
															}
														l124:
															{
																position126, tokenIndex126, depth126 := position, tokenIndex, depth
																if buffer[position] != rune('c') {
																	goto l127
																}
																position++
																goto l126
															l127:
																position, tokenIndex, depth = position126, tokenIndex126, depth126
																if buffer[position] != rune('C') {
																	goto l4
																}
																position++
															}
														l126:
															if !_rules[ruleSpace]() {
																goto l4
															}
															{
																add(ruleAction28, position)
															}
															depth--
															add(ruleJPC, position121)
														}
														{
															position129 := position
															depth++
															{
																position130 := position
																depth++
																{
																	position131, tokenIndex131, depth131 := position, tokenIndex, depth
																	{
																		position133, tokenIndex133, depth133 := position, tokenIndex, depth
																		if buffer[position] != rune('b') {
																			goto l134
																		}
																		position++
																		goto l133
																	l134:
																		position, tokenIndex, depth = position133, tokenIndex133, depth133
																		if buffer[position] != rune('B') {
																			goto l132
																		}
																		position++
																	}
																l133:
																	goto l131
																l132:
																	position, tokenIndex, depth = position131, tokenIndex131, depth131
																	{
																		position136, tokenIndex136, depth136 := position, tokenIndex, depth
																		if buffer[position] != rune('a') {
																			goto l137
																		}
																		position++
																		goto l136
																	l137:
																		position, tokenIndex, depth = position136, tokenIndex136, depth136
																		if buffer[position] != rune('A') {
																			goto l135
																		}
																		position++
																	}
																l136:
																	goto l131
																l135:
																	position, tokenIndex, depth = position131, tokenIndex131, depth131
																	{
																		switch buffer[position] {
																		case 'N', 'n':
																			{
																				position139, tokenIndex139, depth139 := position, tokenIndex, depth
																				if buffer[position] != rune('n') {
																					goto l140
																				}
																				position++
																				goto l139
																			l140:
																				position, tokenIndex, depth = position139, tokenIndex139, depth139
																				if buffer[position] != rune('N') {
																					goto l4
																				}
																				position++
																			}
																		l139:
																			{
																				position141, tokenIndex141, depth141 := position, tokenIndex, depth
																				if buffer[position] != rune('z') {
																					goto l142
																				}
																				position++
																				goto l141
																			l142:
																				position, tokenIndex, depth = position141, tokenIndex141, depth141
																				if buffer[position] != rune('Z') {
																					goto l4
																				}
																				position++
																			}
																		l141:
																			break
																		case 'A', 'a':
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
																					goto l4
																				}
																				position++
																			}
																		l143:
																			{
																				position145, tokenIndex145, depth145 := position, tokenIndex, depth
																				if buffer[position] != rune('e') {
																					goto l146
																				}
																				position++
																				goto l145
																			l146:
																				position, tokenIndex, depth = position145, tokenIndex145, depth145
																				if buffer[position] != rune('E') {
																					goto l4
																				}
																				position++
																			}
																		l145:
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
																				position147, tokenIndex147, depth147 := position, tokenIndex, depth
																				if buffer[position] != rune('b') {
																					goto l148
																				}
																				position++
																				goto l147
																			l148:
																				position, tokenIndex, depth = position147, tokenIndex147, depth147
																				if buffer[position] != rune('B') {
																					goto l4
																				}
																				position++
																			}
																		l147:
																			{
																				position149, tokenIndex149, depth149 := position, tokenIndex, depth
																				if buffer[position] != rune('e') {
																					goto l150
																				}
																				position++
																				goto l149
																			l150:
																				position, tokenIndex, depth = position149, tokenIndex149, depth149
																				if buffer[position] != rune('E') {
																					goto l4
																				}
																				position++
																			}
																		l149:
																			break
																		}
																	}

																}
															l131:
																if !_rules[ruleSpace]() {
																	goto l4
																}
																depth--
																add(ruleCMP_OP, position130)
															}
															depth--
															add(rulePegText, position129)
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
															position152 := position
															depth++
															{
																position153, tokenIndex153, depth153 := position, tokenIndex, depth
																if buffer[position] != rune('c') {
																	goto l154
																}
																position++
																goto l153
															l154:
																position, tokenIndex, depth = position153, tokenIndex153, depth153
																if buffer[position] != rune('C') {
																	goto l4
																}
																position++
															}
														l153:
															{
																position155, tokenIndex155, depth155 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l156
																}
																position++
																goto l155
															l156:
																position, tokenIndex, depth = position155, tokenIndex155, depth155
																if buffer[position] != rune('M') {
																	goto l4
																}
																position++
															}
														l155:
															{
																position157, tokenIndex157, depth157 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l158
																}
																position++
																goto l157
															l158:
																position, tokenIndex, depth = position157, tokenIndex157, depth157
																if buffer[position] != rune('P') {
																	goto l4
																}
																position++
															}
														l157:
															if !_rules[ruleSpace]() {
																goto l4
															}
															{
																add(ruleAction27, position)
															}
															depth--
															add(ruleCMP, position152)
														}
														{
															position160 := position
															depth++
															if !_rules[ruleDATA_TYPE]() {
																goto l4
															}
															depth--
															add(rulePegText, position160)
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
															position162 := position
															depth++
															{
																position163, tokenIndex163, depth163 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l164
																}
																position++
																goto l163
															l164:
																position, tokenIndex, depth = position163, tokenIndex163, depth163
																if buffer[position] != rune('L') {
																	goto l4
																}
																position++
															}
														l163:
															{
																position165, tokenIndex165, depth165 := position, tokenIndex, depth
																if buffer[position] != rune('d') {
																	goto l166
																}
																position++
																goto l165
															l166:
																position, tokenIndex, depth = position165, tokenIndex165, depth165
																if buffer[position] != rune('D') {
																	goto l4
																}
																position++
															}
														l165:
															if !_rules[ruleSpace]() {
																goto l4
															}
															{
																add(ruleAction26, position)
															}
															depth--
															add(ruleLD, position162)
														}
														{
															position168 := position
															depth++
															if !_rules[ruleDATA_TYPE]() {
																goto l4
															}
															depth--
															add(rulePegText, position168)
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
															position170 := position
															depth++
															{
																position171, tokenIndex171, depth171 := position, tokenIndex, depth
																if buffer[position] != rune('n') {
																	goto l172
																}
																position++
																goto l171
															l172:
																position, tokenIndex, depth = position171, tokenIndex171, depth171
																if buffer[position] != rune('N') {
																	goto l4
																}
																position++
															}
														l171:
															{
																position173, tokenIndex173, depth173 := position, tokenIndex, depth
																if buffer[position] != rune('o') {
																	goto l174
																}
																position++
																goto l173
															l174:
																position, tokenIndex, depth = position173, tokenIndex173, depth173
																if buffer[position] != rune('O') {
																	goto l4
																}
																position++
															}
														l173:
															{
																position175, tokenIndex175, depth175 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l176
																}
																position++
																goto l175
															l176:
																position, tokenIndex, depth = position175, tokenIndex175, depth175
																if buffer[position] != rune('P') {
																	goto l4
																}
																position++
															}
														l175:
															if !_rules[ruleSpacing]() {
																goto l4
															}
															{
																add(ruleAction18, position)
															}
															depth--
															add(ruleNOP, position170)
														}
														break
													case 'R', 'r':
														{
															position178 := position
															depth++
															{
																position179, tokenIndex179, depth179 := position, tokenIndex, depth
																if buffer[position] != rune('r') {
																	goto l180
																}
																position++
																goto l179
															l180:
																position, tokenIndex, depth = position179, tokenIndex179, depth179
																if buffer[position] != rune('R') {
																	goto l4
																}
																position++
															}
														l179:
															{
																position181, tokenIndex181, depth181 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l182
																}
																position++
																goto l181
															l182:
																position, tokenIndex, depth = position181, tokenIndex181, depth181
																if buffer[position] != rune('E') {
																	goto l4
																}
																position++
															}
														l181:
															{
																position183, tokenIndex183, depth183 := position, tokenIndex, depth
																if buffer[position] != rune('t') {
																	goto l184
																}
																position++
																goto l183
															l184:
																position, tokenIndex, depth = position183, tokenIndex183, depth183
																if buffer[position] != rune('T') {
																	goto l4
																}
																position++
															}
														l183:
															if !_rules[ruleSpacing]() {
																goto l4
															}
															{
																add(ruleAction17, position)
															}
															depth--
															add(ruleRET, position178)
														}
														break
													case 'E', 'e':
														{
															position186 := position
															depth++
															{
																position187, tokenIndex187, depth187 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l188
																}
																position++
																goto l187
															l188:
																position, tokenIndex, depth = position187, tokenIndex187, depth187
																if buffer[position] != rune('E') {
																	goto l4
																}
																position++
															}
														l187:
															{
																position189, tokenIndex189, depth189 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l190
																}
																position++
																goto l189
															l190:
																position, tokenIndex, depth = position189, tokenIndex189, depth189
																if buffer[position] != rune('X') {
																	goto l4
																}
																position++
															}
														l189:
															{
																position191, tokenIndex191, depth191 := position, tokenIndex, depth
																if buffer[position] != rune('i') {
																	goto l192
																}
																position++
																goto l191
															l192:
																position, tokenIndex, depth = position191, tokenIndex191, depth191
																if buffer[position] != rune('I') {
																	goto l4
																}
																position++
															}
														l191:
															{
																position193, tokenIndex193, depth193 := position, tokenIndex, depth
																if buffer[position] != rune('t') {
																	goto l194
																}
																position++
																goto l193
															l194:
																position, tokenIndex, depth = position193, tokenIndex193, depth193
																if buffer[position] != rune('T') {
																	goto l4
																}
																position++
															}
														l193:
															if !_rules[ruleSpacing]() {
																goto l4
															}
															{
																add(ruleAction16, position)
															}
															depth--
															add(ruleEXIT, position186)
														}
														break
													default:
														{
															position196, tokenIndex196, depth196 := position, tokenIndex, depth
															{
																position198 := position
																depth++
																{
																	position199, tokenIndex199, depth199 := position, tokenIndex, depth
																	if buffer[position] != rune('i') {
																		goto l200
																	}
																	position++
																	goto l199
																l200:
																	position, tokenIndex, depth = position199, tokenIndex199, depth199
																	if buffer[position] != rune('I') {
																		goto l197
																	}
																	position++
																}
															l199:
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
																		goto l197
																	}
																	position++
																}
															l201:
																if !_rules[ruleSpace]() {
																	goto l197
																}
																{
																	add(ruleAction23, position)
																}
																depth--
																add(ruleIN, position198)
															}
															goto l196
														l197:
															position, tokenIndex, depth = position196, tokenIndex196, depth196
															{
																position204 := position
																depth++
																{
																	position205, tokenIndex205, depth205 := position, tokenIndex, depth
																	if buffer[position] != rune('o') {
																		goto l206
																	}
																	position++
																	goto l205
																l206:
																	position, tokenIndex, depth = position205, tokenIndex205, depth205
																	if buffer[position] != rune('O') {
																		goto l4
																	}
																	position++
																}
															l205:
																{
																	position207, tokenIndex207, depth207 := position, tokenIndex, depth
																	if buffer[position] != rune('u') {
																		goto l208
																	}
																	position++
																	goto l207
																l208:
																	position, tokenIndex, depth = position207, tokenIndex207, depth207
																	if buffer[position] != rune('U') {
																		goto l4
																	}
																	position++
																}
															l207:
																{
																	position209, tokenIndex209, depth209 := position, tokenIndex, depth
																	if buffer[position] != rune('t') {
																		goto l210
																	}
																	position++
																	goto l209
																l210:
																	position, tokenIndex, depth = position209, tokenIndex209, depth209
																	if buffer[position] != rune('T') {
																		goto l4
																	}
																	position++
																}
															l209:
																if !_rules[ruleSpace]() {
																	goto l4
																}
																{
																	add(ruleAction24, position)
																}
																depth--
																add(ruleOUT, position204)
															}
														}
													l196:
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
										l32:
											depth--
											add(ruleInst, position31)
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
								position213, tokenIndex213, depth213 := position, tokenIndex, depth
								if !_rules[ruleComment]() {
									goto l213
								}
								{
									add(ruleAction2, position)
								}
								goto l214
							l213:
								position, tokenIndex, depth = position213, tokenIndex213, depth213
							}
						l214:
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
					position217 := position
					depth++
					{
						position218, tokenIndex218, depth218 := position, tokenIndex, depth
						if !matchDot() {
							goto l218
						}
						goto l0
					l218:
						position, tokenIndex, depth = position218, tokenIndex218, depth218
					}
					depth--
					add(ruleEOT, position217)
				}
				depth--
				add(ruleStart, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Line <- <((Label / ((&('.') Pseudo) | (&(';') Comment) | (&('C' | 'E' | 'I' | 'J' | 'L' | 'N' | 'O' | 'P' | 'R' | 'c' | 'e' | 'i' | 'j' | 'l' | 'n' | 'o' | 'p' | 'r') Inst))) Action1 (Comment Action2)?)> */
		nil,
		/* 2 Comment <- <(SEMICOLON <(!NL .)*> Action3)> */
		func() bool {
			position220, tokenIndex220, depth220 := position, tokenIndex, depth
			{
				position221 := position
				depth++
				{
					position222 := position
					depth++
					if buffer[position] != rune(';') {
						goto l220
					}
					position++
					if !_rules[ruleSpacing]() {
						goto l220
					}
					depth--
					add(ruleSEMICOLON, position222)
				}
				{
					position223 := position
					depth++
				l224:
					{
						position225, tokenIndex225, depth225 := position, tokenIndex, depth
						{
							position226, tokenIndex226, depth226 := position, tokenIndex, depth
							if !_rules[ruleNL]() {
								goto l226
							}
							goto l225
						l226:
							position, tokenIndex, depth = position226, tokenIndex226, depth226
						}
						if !matchDot() {
							goto l225
						}
						goto l224
					l225:
						position, tokenIndex, depth = position225, tokenIndex225, depth225
					}
					depth--
					add(rulePegText, position223)
				}
				{
					add(ruleAction3, position)
				}
				depth--
				add(ruleComment, position221)
			}
			return true
		l220:
			position, tokenIndex, depth = position220, tokenIndex220, depth220
			return false
		},
		/* 3 Label <- <(<Identifier> Action4 Spacing COLON)> */
		nil,
		/* 4 Inst <- <(((PUSH / ((&('J' | 'j') JMP) | (&('P' | 'p') POP) | (&('C' | 'c') CALL))) Operand) / (CAL <DATA_TYPE> Action5 <CAL_OP> Action6 Operand COMMA Operand) / ((&('J' | 'j') (JPC <CMP_OP> Action9 Operand)) | (&('C' | 'c') (CMP <DATA_TYPE> Action8 Operand COMMA Operand)) | (&('L' | 'l') (LD <DATA_TYPE> Action7 Operand COMMA Operand)) | (&('N' | 'n') NOP) | (&('R' | 'r') RET) | (&('E' | 'e') EXIT) | (&('I' | 'O' | 'i' | 'o') ((IN / OUT) Operand COMMA Operand))))> */
		nil,
		/* 5 Pseudo <- <(BLOCK <IntegerLiteral> Action10 Space <IntegerLiteral> Action11)> */
		nil,
		/* 6 Operand <- <(((LBRK <Identifier> RBRK Action13) / (<IntegerLiteral> Action14) / ((&('"' | '\'' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Literal) | (&('[') (LBRK <IntegerLiteral> RBRK Action15)) | (&('$' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<Identifier> Action12)))) Spacing)> */
		func() bool {
			position231, tokenIndex231, depth231 := position, tokenIndex, depth
			{
				position232 := position
				depth++
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					if !_rules[ruleLBRK]() {
						goto l234
					}
					{
						position235 := position
						depth++
						if !_rules[ruleIdentifier]() {
							goto l234
						}
						depth--
						add(rulePegText, position235)
					}
					if !_rules[ruleRBRK]() {
						goto l234
					}
					{
						add(ruleAction13, position)
					}
					goto l233
				l234:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
					{
						position238 := position
						depth++
						if !_rules[ruleIntegerLiteral]() {
							goto l237
						}
						depth--
						add(rulePegText, position238)
					}
					{
						add(ruleAction14, position)
					}
					goto l233
				l237:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
					{
						switch buffer[position] {
						case '"', '\'', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							{
								position241 := position
								depth++
								{
									position242, tokenIndex242, depth242 := position, tokenIndex, depth
									{
										position244 := position
										depth++
										{
											position245, tokenIndex245, depth245 := position, tokenIndex, depth
											{
												position247 := position
												depth++
												{
													position248 := position
													depth++
													{
														position249, tokenIndex249, depth249 := position, tokenIndex, depth
														{
															position251, tokenIndex251, depth251 := position, tokenIndex, depth
															if buffer[position] != rune('0') {
																goto l252
															}
															position++
															if buffer[position] != rune('x') {
																goto l252
															}
															position++
															goto l251
														l252:
															position, tokenIndex, depth = position251, tokenIndex251, depth251
															if buffer[position] != rune('0') {
																goto l250
															}
															position++
															if buffer[position] != rune('X') {
																goto l250
															}
															position++
														}
													l251:
														{
															position253, tokenIndex253, depth253 := position, tokenIndex, depth
															if !_rules[ruleHexDigits]() {
																goto l253
															}
															goto l254
														l253:
															position, tokenIndex, depth = position253, tokenIndex253, depth253
														}
													l254:
														if buffer[position] != rune('.') {
															goto l250
														}
														position++
														if !_rules[ruleHexDigits]() {
															goto l250
														}
														goto l249
													l250:
														position, tokenIndex, depth = position249, tokenIndex249, depth249
														if !_rules[ruleHexNumeral]() {
															goto l246
														}
														{
															position255, tokenIndex255, depth255 := position, tokenIndex, depth
															if buffer[position] != rune('.') {
																goto l255
															}
															position++
															goto l256
														l255:
															position, tokenIndex, depth = position255, tokenIndex255, depth255
														}
													l256:
													}
												l249:
													depth--
													add(ruleHexSignificand, position248)
												}
												{
													position257 := position
													depth++
													{
														position258, tokenIndex258, depth258 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l259
														}
														position++
														goto l258
													l259:
														position, tokenIndex, depth = position258, tokenIndex258, depth258
														if buffer[position] != rune('P') {
															goto l246
														}
														position++
													}
												l258:
													{
														position260, tokenIndex260, depth260 := position, tokenIndex, depth
														{
															position262, tokenIndex262, depth262 := position, tokenIndex, depth
															if buffer[position] != rune('+') {
																goto l263
															}
															position++
															goto l262
														l263:
															position, tokenIndex, depth = position262, tokenIndex262, depth262
															if buffer[position] != rune('-') {
																goto l260
															}
															position++
														}
													l262:
														goto l261
													l260:
														position, tokenIndex, depth = position260, tokenIndex260, depth260
													}
												l261:
													if !_rules[ruleDigits]() {
														goto l246
													}
													depth--
													add(ruleBinaryExponent, position257)
												}
												{
													position264, tokenIndex264, depth264 := position, tokenIndex, depth
													{
														switch buffer[position] {
														case 'D':
															if buffer[position] != rune('D') {
																goto l264
															}
															position++
															break
														case 'd':
															if buffer[position] != rune('d') {
																goto l264
															}
															position++
															break
														case 'F':
															if buffer[position] != rune('F') {
																goto l264
															}
															position++
															break
														default:
															if buffer[position] != rune('f') {
																goto l264
															}
															position++
															break
														}
													}

													goto l265
												l264:
													position, tokenIndex, depth = position264, tokenIndex264, depth264
												}
											l265:
												depth--
												add(ruleHexFloat, position247)
											}
											goto l245
										l246:
											position, tokenIndex, depth = position245, tokenIndex245, depth245
											{
												position267 := position
												depth++
												{
													position268, tokenIndex268, depth268 := position, tokenIndex, depth
													if !_rules[ruleDigits]() {
														goto l269
													}
													if buffer[position] != rune('.') {
														goto l269
													}
													position++
													{
														position270, tokenIndex270, depth270 := position, tokenIndex, depth
														if !_rules[ruleDigits]() {
															goto l270
														}
														goto l271
													l270:
														position, tokenIndex, depth = position270, tokenIndex270, depth270
													}
												l271:
													{
														position272, tokenIndex272, depth272 := position, tokenIndex, depth
														if !_rules[ruleExponent]() {
															goto l272
														}
														goto l273
													l272:
														position, tokenIndex, depth = position272, tokenIndex272, depth272
													}
												l273:
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
													goto l268
												l269:
													position, tokenIndex, depth = position268, tokenIndex268, depth268
													if buffer[position] != rune('.') {
														goto l277
													}
													position++
													if !_rules[ruleDigits]() {
														goto l277
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
														position280, tokenIndex280, depth280 := position, tokenIndex, depth
														{
															switch buffer[position] {
															case 'D':
																if buffer[position] != rune('D') {
																	goto l280
																}
																position++
																break
															case 'd':
																if buffer[position] != rune('d') {
																	goto l280
																}
																position++
																break
															case 'F':
																if buffer[position] != rune('F') {
																	goto l280
																}
																position++
																break
															default:
																if buffer[position] != rune('f') {
																	goto l280
																}
																position++
																break
															}
														}

														goto l281
													l280:
														position, tokenIndex, depth = position280, tokenIndex280, depth280
													}
												l281:
													goto l268
												l277:
													position, tokenIndex, depth = position268, tokenIndex268, depth268
													if !_rules[ruleDigits]() {
														goto l283
													}
													if !_rules[ruleExponent]() {
														goto l283
													}
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
													goto l268
												l283:
													position, tokenIndex, depth = position268, tokenIndex268, depth268
													if !_rules[ruleDigits]() {
														goto l243
													}
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
														switch buffer[position] {
														case 'D':
															if buffer[position] != rune('D') {
																goto l243
															}
															position++
															break
														case 'd':
															if buffer[position] != rune('d') {
																goto l243
															}
															position++
															break
														case 'F':
															if buffer[position] != rune('F') {
																goto l243
															}
															position++
															break
														default:
															if buffer[position] != rune('f') {
																goto l243
															}
															position++
															break
														}
													}

												}
											l268:
												depth--
												add(ruleDecimalFloat, position267)
											}
										}
									l245:
										depth--
										add(ruleFloatLiteral, position244)
									}
									goto l242
								l243:
									position, tokenIndex, depth = position242, tokenIndex242, depth242
									{
										switch buffer[position] {
										case '"':
											{
												position291 := position
												depth++
												if buffer[position] != rune('"') {
													goto l231
												}
												position++
											l292:
												{
													position293, tokenIndex293, depth293 := position, tokenIndex, depth
													{
														position294, tokenIndex294, depth294 := position, tokenIndex, depth
														if !_rules[ruleEscape]() {
															goto l295
														}
														goto l294
													l295:
														position, tokenIndex, depth = position294, tokenIndex294, depth294
														{
															position296, tokenIndex296, depth296 := position, tokenIndex, depth
															{
																switch buffer[position] {
																case '\r':
																	if buffer[position] != rune('\r') {
																		goto l296
																	}
																	position++
																	break
																case '\n':
																	if buffer[position] != rune('\n') {
																		goto l296
																	}
																	position++
																	break
																case '\\':
																	if buffer[position] != rune('\\') {
																		goto l296
																	}
																	position++
																	break
																default:
																	if buffer[position] != rune('"') {
																		goto l296
																	}
																	position++
																	break
																}
															}

															goto l293
														l296:
															position, tokenIndex, depth = position296, tokenIndex296, depth296
														}
														if !matchDot() {
															goto l293
														}
													}
												l294:
													goto l292
												l293:
													position, tokenIndex, depth = position293, tokenIndex293, depth293
												}
												if buffer[position] != rune('"') {
													goto l231
												}
												position++
												depth--
												add(ruleStringLiteral, position291)
											}
											break
										case '\'':
											{
												position298 := position
												depth++
												if buffer[position] != rune('\'') {
													goto l231
												}
												position++
												{
													position299, tokenIndex299, depth299 := position, tokenIndex, depth
													if !_rules[ruleEscape]() {
														goto l300
													}
													goto l299
												l300:
													position, tokenIndex, depth = position299, tokenIndex299, depth299
													{
														position301, tokenIndex301, depth301 := position, tokenIndex, depth
														{
															position302, tokenIndex302, depth302 := position, tokenIndex, depth
															if buffer[position] != rune('\'') {
																goto l303
															}
															position++
															goto l302
														l303:
															position, tokenIndex, depth = position302, tokenIndex302, depth302
															if buffer[position] != rune('\\') {
																goto l301
															}
															position++
														}
													l302:
														goto l231
													l301:
														position, tokenIndex, depth = position301, tokenIndex301, depth301
													}
													if !matchDot() {
														goto l231
													}
												}
											l299:
												if buffer[position] != rune('\'') {
													goto l231
												}
												position++
												depth--
												add(ruleCharLiteral, position298)
											}
											break
										default:
											if !_rules[ruleIntegerLiteral]() {
												goto l231
											}
											break
										}
									}

								}
							l242:
								if !_rules[ruleSpacing]() {
									goto l231
								}
								depth--
								add(ruleLiteral, position241)
							}
							break
						case '[':
							if !_rules[ruleLBRK]() {
								goto l231
							}
							{
								position304 := position
								depth++
								if !_rules[ruleIntegerLiteral]() {
									goto l231
								}
								depth--
								add(rulePegText, position304)
							}
							if !_rules[ruleRBRK]() {
								goto l231
							}
							{
								add(ruleAction15, position)
							}
							break
						default:
							{
								position306 := position
								depth++
								if !_rules[ruleIdentifier]() {
									goto l231
								}
								depth--
								add(rulePegText, position306)
							}
							{
								add(ruleAction12, position)
							}
							break
						}
					}

				}
			l233:
				if !_rules[ruleSpacing]() {
					goto l231
				}
				depth--
				add(ruleOperand, position232)
			}
			return true
		l231:
			position, tokenIndex, depth = position231, tokenIndex231, depth231
			return false
		},
		/* 7 Spacing <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position309 := position
				depth++
			l310:
				{
					position311, tokenIndex311, depth311 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l311
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l311
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l311
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l311
							}
							position++
							break
						}
					}

					goto l310
				l311:
					position, tokenIndex, depth = position311, tokenIndex311, depth311
				}
				depth--
				add(ruleSpacing, position309)
			}
			return true
		},
		/* 8 Space <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))+> */
		func() bool {
			position313, tokenIndex313, depth313 := position, tokenIndex, depth
			{
				position314 := position
				depth++
				{
					switch buffer[position] {
					case '\f':
						if buffer[position] != rune('\f') {
							goto l313
						}
						position++
						break
					case '\r':
						if buffer[position] != rune('\r') {
							goto l313
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l313
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l313
						}
						position++
						break
					}
				}

			l315:
				{
					position316, tokenIndex316, depth316 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l316
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l316
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l316
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l316
							}
							position++
							break
						}
					}

					goto l315
				l316:
					position, tokenIndex, depth = position316, tokenIndex316, depth316
				}
				depth--
				add(ruleSpace, position314)
			}
			return true
		l313:
			position, tokenIndex, depth = position313, tokenIndex313, depth313
			return false
		},
		/* 9 Identifier <- <(Letter LetterOrDigit* Spacing)> */
		func() bool {
			position319, tokenIndex319, depth319 := position, tokenIndex, depth
			{
				position320 := position
				depth++
				{
					position321 := position
					depth++
					{
						switch buffer[position] {
						case '$', '_':
							{
								position323, tokenIndex323, depth323 := position, tokenIndex, depth
								if buffer[position] != rune('_') {
									goto l324
								}
								position++
								goto l323
							l324:
								position, tokenIndex, depth = position323, tokenIndex323, depth323
								if buffer[position] != rune('$') {
									goto l319
								}
								position++
							}
						l323:
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l319
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l319
							}
							position++
							break
						}
					}

					depth--
					add(ruleLetter, position321)
				}
			l325:
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					{
						position327 := position
						depth++
						{
							switch buffer[position] {
							case '$', '_':
								{
									position329, tokenIndex329, depth329 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l330
									}
									position++
									goto l329
								l330:
									position, tokenIndex, depth = position329, tokenIndex329, depth329
									if buffer[position] != rune('$') {
										goto l326
									}
									position++
								}
							l329:
								break
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l326
								}
								position++
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
						add(ruleLetterOrDigit, position327)
					}
					goto l325
				l326:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
				}
				if !_rules[ruleSpacing]() {
					goto l319
				}
				depth--
				add(ruleIdentifier, position320)
			}
			return true
		l319:
			position, tokenIndex, depth = position319, tokenIndex319, depth319
			return false
		},
		/* 10 Letter <- <((&('$' | '_') ('_' / '$')) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 11 LetterOrDigit <- <((&('$' | '_') ('_' / '$')) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 12 EXIT <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('t' / 'T') Spacing Action16)> */
		nil,
		/* 13 RET <- <(('r' / 'R') ('e' / 'E') ('t' / 'T') Spacing Action17)> */
		nil,
		/* 14 NOP <- <(('n' / 'N') ('o' / 'O') ('p' / 'P') Spacing Action18)> */
		nil,
		/* 15 CALL <- <(('c' / 'C') ('a' / 'A') ('l' / 'L') ('l' / 'L') Space Action19)> */
		nil,
		/* 16 PUSH <- <(('p' / 'P') ('u' / 'U') ('s' / 'S') ('h' / 'H') Space Action20)> */
		nil,
		/* 17 POP <- <(('p' / 'P') ('o' / 'O') ('p' / 'P') Space Action21)> */
		nil,
		/* 18 JMP <- <(('j' / 'J') ('m' / 'M') ('p' / 'P') Space Action22)> */
		nil,
		/* 19 IN <- <(('i' / 'I') ('n' / 'N') Space Action23)> */
		nil,
		/* 20 OUT <- <(('o' / 'O') ('u' / 'U') ('t' / 'T') Space Action24)> */
		nil,
		/* 21 CAL <- <(('c' / 'C') ('a' / 'A') ('l' / 'L') Space Action25)> */
		nil,
		/* 22 LD <- <(('l' / 'L') ('d' / 'D') Space Action26)> */
		nil,
		/* 23 CMP <- <(('c' / 'C') ('m' / 'M') ('p' / 'P') Space Action27)> */
		nil,
		/* 24 JPC <- <(('j' / 'J') ('p' / 'P') ('c' / 'C') Space Action28)> */
		nil,
		/* 25 BLOCK <- <('.' ('b' / 'B') ('l' / 'L') ('o' / 'O') ('c' / 'C') ('k' / 'K') Space Action29)> */
		nil,
		/* 26 CAL_OP <- <(((('m' / 'M') ('u' / 'U') ('l' / 'L')) / ((&('M' | 'm') (('m' / 'M') ('o' / 'O') ('d' / 'D'))) | (&('D' | 'd') (('d' / 'D') ('i' / 'I') ('v' / 'V'))) | (&('S' | 's') (('s' / 'S') ('u' / 'U') ('b' / 'B'))) | (&('A' | 'a') (('a' / 'A') ('d' / 'D') ('d' / 'D'))))) Space)> */
		nil,
		/* 27 CMP_OP <- <((('b' / 'B') / ('a' / 'A') / ((&('N' | 'n') (('n' / 'N') ('z' / 'Z'))) | (&('A' | 'a') (('a' / 'A') ('e' / 'E'))) | (&('Z') 'Z') | (&('z') 'z') | (&('B' | 'b') (('b' / 'B') ('e' / 'E'))))) Space)> */
		nil,
		/* 28 DATA_TYPE <- <(((&('I' | 'i') (('i' / 'I') ('n' / 'N') ('t' / 'T'))) | (&('F' | 'f') (('f' / 'F') ('l' / 'L') ('o' / 'O') ('a' / 'A') ('t' / 'T'))) | (&('B' | 'b') (('b' / 'B') ('y' / 'Y') ('t' / 'T') ('e' / 'E'))) | (&('W' | 'w') (('w' / 'W') ('o' / 'O') ('r' / 'R') ('d' / 'D'))) | (&('D' | 'd') (('d' / 'D') ('w' / 'W') ('o' / 'O') ('r' / 'R') ('d' / 'D')))) Space)> */
		func() bool {
			position349, tokenIndex349, depth349 := position, tokenIndex, depth
			{
				position350 := position
				depth++
				{
					switch buffer[position] {
					case 'I', 'i':
						{
							position352, tokenIndex352, depth352 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l353
							}
							position++
							goto l352
						l353:
							position, tokenIndex, depth = position352, tokenIndex352, depth352
							if buffer[position] != rune('I') {
								goto l349
							}
							position++
						}
					l352:
						{
							position354, tokenIndex354, depth354 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l355
							}
							position++
							goto l354
						l355:
							position, tokenIndex, depth = position354, tokenIndex354, depth354
							if buffer[position] != rune('N') {
								goto l349
							}
							position++
						}
					l354:
						{
							position356, tokenIndex356, depth356 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l357
							}
							position++
							goto l356
						l357:
							position, tokenIndex, depth = position356, tokenIndex356, depth356
							if buffer[position] != rune('T') {
								goto l349
							}
							position++
						}
					l356:
						break
					case 'F', 'f':
						{
							position358, tokenIndex358, depth358 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l359
							}
							position++
							goto l358
						l359:
							position, tokenIndex, depth = position358, tokenIndex358, depth358
							if buffer[position] != rune('F') {
								goto l349
							}
							position++
						}
					l358:
						{
							position360, tokenIndex360, depth360 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l361
							}
							position++
							goto l360
						l361:
							position, tokenIndex, depth = position360, tokenIndex360, depth360
							if buffer[position] != rune('L') {
								goto l349
							}
							position++
						}
					l360:
						{
							position362, tokenIndex362, depth362 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l363
							}
							position++
							goto l362
						l363:
							position, tokenIndex, depth = position362, tokenIndex362, depth362
							if buffer[position] != rune('O') {
								goto l349
							}
							position++
						}
					l362:
						{
							position364, tokenIndex364, depth364 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l365
							}
							position++
							goto l364
						l365:
							position, tokenIndex, depth = position364, tokenIndex364, depth364
							if buffer[position] != rune('A') {
								goto l349
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
								goto l349
							}
							position++
						}
					l366:
						break
					case 'B', 'b':
						{
							position368, tokenIndex368, depth368 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l369
							}
							position++
							goto l368
						l369:
							position, tokenIndex, depth = position368, tokenIndex368, depth368
							if buffer[position] != rune('B') {
								goto l349
							}
							position++
						}
					l368:
						{
							position370, tokenIndex370, depth370 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l371
							}
							position++
							goto l370
						l371:
							position, tokenIndex, depth = position370, tokenIndex370, depth370
							if buffer[position] != rune('Y') {
								goto l349
							}
							position++
						}
					l370:
						{
							position372, tokenIndex372, depth372 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l373
							}
							position++
							goto l372
						l373:
							position, tokenIndex, depth = position372, tokenIndex372, depth372
							if buffer[position] != rune('T') {
								goto l349
							}
							position++
						}
					l372:
						{
							position374, tokenIndex374, depth374 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l375
							}
							position++
							goto l374
						l375:
							position, tokenIndex, depth = position374, tokenIndex374, depth374
							if buffer[position] != rune('E') {
								goto l349
							}
							position++
						}
					l374:
						break
					case 'W', 'w':
						{
							position376, tokenIndex376, depth376 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l377
							}
							position++
							goto l376
						l377:
							position, tokenIndex, depth = position376, tokenIndex376, depth376
							if buffer[position] != rune('W') {
								goto l349
							}
							position++
						}
					l376:
						{
							position378, tokenIndex378, depth378 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l379
							}
							position++
							goto l378
						l379:
							position, tokenIndex, depth = position378, tokenIndex378, depth378
							if buffer[position] != rune('O') {
								goto l349
							}
							position++
						}
					l378:
						{
							position380, tokenIndex380, depth380 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l381
							}
							position++
							goto l380
						l381:
							position, tokenIndex, depth = position380, tokenIndex380, depth380
							if buffer[position] != rune('R') {
								goto l349
							}
							position++
						}
					l380:
						{
							position382, tokenIndex382, depth382 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l383
							}
							position++
							goto l382
						l383:
							position, tokenIndex, depth = position382, tokenIndex382, depth382
							if buffer[position] != rune('D') {
								goto l349
							}
							position++
						}
					l382:
						break
					default:
						{
							position384, tokenIndex384, depth384 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l385
							}
							position++
							goto l384
						l385:
							position, tokenIndex, depth = position384, tokenIndex384, depth384
							if buffer[position] != rune('D') {
								goto l349
							}
							position++
						}
					l384:
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
								goto l349
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
								goto l349
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
								goto l349
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
								goto l349
							}
							position++
						}
					l392:
						break
					}
				}

				if !_rules[ruleSpace]() {
					goto l349
				}
				depth--
				add(ruleDATA_TYPE, position350)
			}
			return true
		l349:
			position, tokenIndex, depth = position349, tokenIndex349, depth349
			return false
		},
		/* 29 LBRK <- <('[' Spacing)> */
		func() bool {
			position394, tokenIndex394, depth394 := position, tokenIndex, depth
			{
				position395 := position
				depth++
				if buffer[position] != rune('[') {
					goto l394
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l394
				}
				depth--
				add(ruleLBRK, position395)
			}
			return true
		l394:
			position, tokenIndex, depth = position394, tokenIndex394, depth394
			return false
		},
		/* 30 RBRK <- <(']' Spacing)> */
		func() bool {
			position396, tokenIndex396, depth396 := position, tokenIndex, depth
			{
				position397 := position
				depth++
				if buffer[position] != rune(']') {
					goto l396
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l396
				}
				depth--
				add(ruleRBRK, position397)
			}
			return true
		l396:
			position, tokenIndex, depth = position396, tokenIndex396, depth396
			return false
		},
		/* 31 COMMA <- <(',' Spacing)> */
		func() bool {
			position398, tokenIndex398, depth398 := position, tokenIndex, depth
			{
				position399 := position
				depth++
				if buffer[position] != rune(',') {
					goto l398
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l398
				}
				depth--
				add(ruleCOMMA, position399)
			}
			return true
		l398:
			position, tokenIndex, depth = position398, tokenIndex398, depth398
			return false
		},
		/* 32 SEMICOLON <- <(';' Spacing)> */
		nil,
		/* 33 COLON <- <(':' Spacing)> */
		nil,
		/* 34 NL <- <'\n'> */
		func() bool {
			position402, tokenIndex402, depth402 := position, tokenIndex, depth
			{
				position403 := position
				depth++
				if buffer[position] != rune('\n') {
					goto l402
				}
				position++
				depth--
				add(ruleNL, position403)
			}
			return true
		l402:
			position, tokenIndex, depth = position402, tokenIndex402, depth402
			return false
		},
		/* 35 EOT <- <!.> */
		nil,
		/* 36 Literal <- <((FloatLiteral / ((&('"') StringLiteral) | (&('\'') CharLiteral) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') IntegerLiteral))) Spacing)> */
		nil,
		/* 37 IntegerLiteral <- <(HexNumeral / BinaryNumeral / OctalNumeral / DecimalNumeral)> */
		func() bool {
			position406, tokenIndex406, depth406 := position, tokenIndex, depth
			{
				position407 := position
				depth++
				{
					position408, tokenIndex408, depth408 := position, tokenIndex, depth
					if !_rules[ruleHexNumeral]() {
						goto l409
					}
					goto l408
				l409:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					{
						position411 := position
						depth++
						{
							position412, tokenIndex412, depth412 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l413
							}
							position++
							if buffer[position] != rune('b') {
								goto l413
							}
							position++
							goto l412
						l413:
							position, tokenIndex, depth = position412, tokenIndex412, depth412
							if buffer[position] != rune('0') {
								goto l410
							}
							position++
							if buffer[position] != rune('B') {
								goto l410
							}
							position++
						}
					l412:
						{
							position414, tokenIndex414, depth414 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l415
							}
							position++
							goto l414
						l415:
							position, tokenIndex, depth = position414, tokenIndex414, depth414
							if buffer[position] != rune('1') {
								goto l410
							}
							position++
						}
					l414:
					l416:
						{
							position417, tokenIndex417, depth417 := position, tokenIndex, depth
						l418:
							{
								position419, tokenIndex419, depth419 := position, tokenIndex, depth
								if buffer[position] != rune('_') {
									goto l419
								}
								position++
								goto l418
							l419:
								position, tokenIndex, depth = position419, tokenIndex419, depth419
							}
							{
								position420, tokenIndex420, depth420 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l421
								}
								position++
								goto l420
							l421:
								position, tokenIndex, depth = position420, tokenIndex420, depth420
								if buffer[position] != rune('1') {
									goto l417
								}
								position++
							}
						l420:
							goto l416
						l417:
							position, tokenIndex, depth = position417, tokenIndex417, depth417
						}
						depth--
						add(ruleBinaryNumeral, position411)
					}
					goto l408
				l410:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					{
						position423 := position
						depth++
						if buffer[position] != rune('0') {
							goto l422
						}
						position++
					l426:
						{
							position427, tokenIndex427, depth427 := position, tokenIndex, depth
							if buffer[position] != rune('_') {
								goto l427
							}
							position++
							goto l426
						l427:
							position, tokenIndex, depth = position427, tokenIndex427, depth427
						}
						if c := buffer[position]; c < rune('0') || c > rune('7') {
							goto l422
						}
						position++
					l424:
						{
							position425, tokenIndex425, depth425 := position, tokenIndex, depth
						l428:
							{
								position429, tokenIndex429, depth429 := position, tokenIndex, depth
								if buffer[position] != rune('_') {
									goto l429
								}
								position++
								goto l428
							l429:
								position, tokenIndex, depth = position429, tokenIndex429, depth429
							}
							if c := buffer[position]; c < rune('0') || c > rune('7') {
								goto l425
							}
							position++
							goto l424
						l425:
							position, tokenIndex, depth = position425, tokenIndex425, depth425
						}
						depth--
						add(ruleOctalNumeral, position423)
					}
					goto l408
				l422:
					position, tokenIndex, depth = position408, tokenIndex408, depth408
					{
						position430 := position
						depth++
						{
							position431, tokenIndex431, depth431 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l432
							}
							position++
							goto l431
						l432:
							position, tokenIndex, depth = position431, tokenIndex431, depth431
							if c := buffer[position]; c < rune('1') || c > rune('9') {
								goto l406
							}
							position++
						l433:
							{
								position434, tokenIndex434, depth434 := position, tokenIndex, depth
							l435:
								{
									position436, tokenIndex436, depth436 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l436
									}
									position++
									goto l435
								l436:
									position, tokenIndex, depth = position436, tokenIndex436, depth436
								}
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l434
								}
								position++
								goto l433
							l434:
								position, tokenIndex, depth = position434, tokenIndex434, depth434
							}
						}
					l431:
						depth--
						add(ruleDecimalNumeral, position430)
					}
				}
			l408:
				depth--
				add(ruleIntegerLiteral, position407)
			}
			return true
		l406:
			position, tokenIndex, depth = position406, tokenIndex406, depth406
			return false
		},
		/* 38 DecimalNumeral <- <('0' / ([1-9] ('_'* [0-9])*))> */
		nil,
		/* 39 HexNumeral <- <((('0' 'x') / ('0' 'X')) HexDigits)> */
		func() bool {
			position438, tokenIndex438, depth438 := position, tokenIndex, depth
			{
				position439 := position
				depth++
				{
					position440, tokenIndex440, depth440 := position, tokenIndex, depth
					if buffer[position] != rune('0') {
						goto l441
					}
					position++
					if buffer[position] != rune('x') {
						goto l441
					}
					position++
					goto l440
				l441:
					position, tokenIndex, depth = position440, tokenIndex440, depth440
					if buffer[position] != rune('0') {
						goto l438
					}
					position++
					if buffer[position] != rune('X') {
						goto l438
					}
					position++
				}
			l440:
				if !_rules[ruleHexDigits]() {
					goto l438
				}
				depth--
				add(ruleHexNumeral, position439)
			}
			return true
		l438:
			position, tokenIndex, depth = position438, tokenIndex438, depth438
			return false
		},
		/* 40 BinaryNumeral <- <((('0' 'b') / ('0' 'B')) ('0' / '1') ('_'* ('0' / '1'))*)> */
		nil,
		/* 41 OctalNumeral <- <('0' ('_'* [0-7])+)> */
		nil,
		/* 42 FloatLiteral <- <(HexFloat / DecimalFloat)> */
		nil,
		/* 43 DecimalFloat <- <((Digits '.' Digits? Exponent? ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?) / ('.' Digits Exponent? ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?) / (Digits Exponent ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?) / (Digits Exponent? ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))))> */
		nil,
		/* 44 Exponent <- <(('e' / 'E') ('+' / '-')? Digits)> */
		func() bool {
			position446, tokenIndex446, depth446 := position, tokenIndex, depth
			{
				position447 := position
				depth++
				{
					position448, tokenIndex448, depth448 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l449
					}
					position++
					goto l448
				l449:
					position, tokenIndex, depth = position448, tokenIndex448, depth448
					if buffer[position] != rune('E') {
						goto l446
					}
					position++
				}
			l448:
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					{
						position452, tokenIndex452, depth452 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l453
						}
						position++
						goto l452
					l453:
						position, tokenIndex, depth = position452, tokenIndex452, depth452
						if buffer[position] != rune('-') {
							goto l450
						}
						position++
					}
				l452:
					goto l451
				l450:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
				}
			l451:
				if !_rules[ruleDigits]() {
					goto l446
				}
				depth--
				add(ruleExponent, position447)
			}
			return true
		l446:
			position, tokenIndex, depth = position446, tokenIndex446, depth446
			return false
		},
		/* 45 HexFloat <- <(HexSignificand BinaryExponent ((&('D') 'D') | (&('d') 'd') | (&('F') 'F') | (&('f') 'f'))?)> */
		nil,
		/* 46 HexSignificand <- <(((('0' 'x') / ('0' 'X')) HexDigits? '.' HexDigits) / (HexNumeral '.'?))> */
		nil,
		/* 47 BinaryExponent <- <(('p' / 'P') ('+' / '-')? Digits)> */
		nil,
		/* 48 Digits <- <([0-9] ('_'* [0-9])*)> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l457
				}
				position++
			l459:
				{
					position460, tokenIndex460, depth460 := position, tokenIndex, depth
				l461:
					{
						position462, tokenIndex462, depth462 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l462
						}
						position++
						goto l461
					l462:
						position, tokenIndex, depth = position462, tokenIndex462, depth462
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l460
					}
					position++
					goto l459
				l460:
					position, tokenIndex, depth = position460, tokenIndex460, depth460
				}
				depth--
				add(ruleDigits, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 49 HexDigits <- <(HexDigit ('_'* HexDigit)*)> */
		func() bool {
			position463, tokenIndex463, depth463 := position, tokenIndex, depth
			{
				position464 := position
				depth++
				if !_rules[ruleHexDigit]() {
					goto l463
				}
			l465:
				{
					position466, tokenIndex466, depth466 := position, tokenIndex, depth
				l467:
					{
						position468, tokenIndex468, depth468 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l468
						}
						position++
						goto l467
					l468:
						position, tokenIndex, depth = position468, tokenIndex468, depth468
					}
					if !_rules[ruleHexDigit]() {
						goto l466
					}
					goto l465
				l466:
					position, tokenIndex, depth = position466, tokenIndex466, depth466
				}
				depth--
				add(ruleHexDigits, position464)
			}
			return true
		l463:
			position, tokenIndex, depth = position463, tokenIndex463, depth463
			return false
		},
		/* 50 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l469
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l469
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l469
						}
						position++
						break
					}
				}

				depth--
				add(ruleHexDigit, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 51 CharLiteral <- <('\'' (Escape / (!('\'' / '\\') .)) '\'')> */
		nil,
		/* 52 StringLiteral <- <('"' (Escape / (!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('"') '"')) .))* '"')> */
		nil,
		/* 53 Escape <- <('\\' ((&('u') UnicodeEscape) | (&('\\') '\\') | (&('\'') '\'') | (&('"') '"') | (&('r') 'r') | (&('f') 'f') | (&('n') 'n') | (&('t') 't') | (&('b') 'b') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7') OctalEscape)))> */
		func() bool {
			position474, tokenIndex474, depth474 := position, tokenIndex, depth
			{
				position475 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l474
				}
				position++
				{
					switch buffer[position] {
					case 'u':
						{
							position477 := position
							depth++
							if buffer[position] != rune('u') {
								goto l474
							}
							position++
						l478:
							{
								position479, tokenIndex479, depth479 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l479
								}
								position++
								goto l478
							l479:
								position, tokenIndex, depth = position479, tokenIndex479, depth479
							}
							if !_rules[ruleHexDigit]() {
								goto l474
							}
							if !_rules[ruleHexDigit]() {
								goto l474
							}
							if !_rules[ruleHexDigit]() {
								goto l474
							}
							if !_rules[ruleHexDigit]() {
								goto l474
							}
							depth--
							add(ruleUnicodeEscape, position477)
						}
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l474
						}
						position++
						break
					case '\'':
						if buffer[position] != rune('\'') {
							goto l474
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l474
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l474
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l474
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l474
						}
						position++
						break
					case 't':
						if buffer[position] != rune('t') {
							goto l474
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l474
						}
						position++
						break
					default:
						{
							position480 := position
							depth++
							{
								position481, tokenIndex481, depth481 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('3') {
									goto l482
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l482
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l482
								}
								position++
								goto l481
							l482:
								position, tokenIndex, depth = position481, tokenIndex481, depth481
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l483
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l483
								}
								position++
								goto l481
							l483:
								position, tokenIndex, depth = position481, tokenIndex481, depth481
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l474
								}
								position++
							}
						l481:
							depth--
							add(ruleOctalEscape, position480)
						}
						break
					}
				}

				depth--
				add(ruleEscape, position475)
			}
			return true
		l474:
			position, tokenIndex, depth = position474, tokenIndex474, depth474
			return false
		},
		/* 54 OctalEscape <- <(([0-3] [0-7] [0-7]) / ([0-7] [0-7]) / [0-7])> */
		nil,
		/* 55 UnicodeEscape <- <('u'+ HexDigit HexDigit HexDigit HexDigit)> */
		nil,
		/* 57 Action0 <- <{p.line++}> */
		nil,
		/* 58 Action1 <- <{p.PopAssembly()}> */
		nil,
		/* 59 Action2 <- <{p.PopAssembly();p.AddComment()}> */
		nil,
		nil,
		/* 61 Action3 <- <{p.Push(&asm.Comment{});p.Push(text)}> */
		nil,
		/* 62 Action4 <- <{p.Push(&asm.Label{});p.Push(text)}> */
		nil,
		/* 63 Action5 <- <{p.Push(lookup(asm.T_INT,text))}> */
		nil,
		/* 64 Action6 <- <{p.Push(lookup(asm.CAL_ADD,text))}> */
		nil,
		/* 65 Action7 <- <{p.Push(lookup(asm.T_INT,text))}> */
		nil,
		/* 66 Action8 <- <{p.Push(lookup(asm.T_INT,text))}> */
		nil,
		/* 67 Action9 <- <{p.Push(lookup(asm.CMP_A,text))}> */
		nil,
		/* 68 Action10 <- <{p.Push(text)}> */
		nil,
		/* 69 Action11 <- <{p.Push(text)}> */
		nil,
		/* 70 Action12 <- <{p.Push(CreateOperand(text,false,true))}> */
		nil,
		/* 71 Action13 <- <{p.Push(CreateOperand(text,false,false))}> */
		nil,
		/* 72 Action14 <- <{p.Push(CreateOperand(text,true,true))}> */
		nil,
		/* 73 Action15 <- <{p.Push(CreateOperand(text,true,false))}> */
		nil,
		/* 74 Action16 <- <{p.PushInst(asm.OP_EXIT)}> */
		nil,
		/* 75 Action17 <- <{p.PushInst(asm.OP_RET)}> */
		nil,
		/* 76 Action18 <- <{p.PushInst(asm.OP_NOP)}> */
		nil,
		/* 77 Action19 <- <{p.PushInst(asm.OP_CALL)}> */
		nil,
		/* 78 Action20 <- <{p.PushInst(asm.OP_PUSH)}> */
		nil,
		/* 79 Action21 <- <{p.PushInst(asm.OP_POP)}> */
		nil,
		/* 80 Action22 <- <{p.PushInst(asm.OP_JMP)}> */
		nil,
		/* 81 Action23 <- <{p.PushInst(asm.OP_IN)}> */
		nil,
		/* 82 Action24 <- <{p.PushInst(asm.OP_OUT)}> */
		nil,
		/* 83 Action25 <- <{p.PushInst(asm.OP_CAL)}> */
		nil,
		/* 84 Action26 <- <{p.PushInst(asm.OP_LD)}> */
		nil,
		/* 85 Action27 <- <{p.PushInst(asm.OP_CMP)}> */
		nil,
		/* 86 Action28 <- <{p.PushInst(asm.OP_JPC)}> */
		nil,
		/* 87 Action29 <- <{p.Push(&asm.PseudoBlock{})}> */
		nil,
	}
	p.rules = _rules
}
