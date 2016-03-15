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
								switch buffer[position] {
								case '.':
									{
										position8 := position
										depth++
										{
											position9 := position
											depth++
											if buffer[position] != rune('.') {
												goto l4
											}
											position++
											{
												position10, tokenIndex10, depth10 := position, tokenIndex, depth
												if buffer[position] != rune('b') {
													goto l11
												}
												position++
												goto l10
											l11:
												position, tokenIndex, depth = position10, tokenIndex10, depth10
												if buffer[position] != rune('B') {
													goto l4
												}
												position++
											}
										l10:
											{
												position12, tokenIndex12, depth12 := position, tokenIndex, depth
												if buffer[position] != rune('l') {
													goto l13
												}
												position++
												goto l12
											l13:
												position, tokenIndex, depth = position12, tokenIndex12, depth12
												if buffer[position] != rune('L') {
													goto l4
												}
												position++
											}
										l12:
											{
												position14, tokenIndex14, depth14 := position, tokenIndex, depth
												if buffer[position] != rune('o') {
													goto l15
												}
												position++
												goto l14
											l15:
												position, tokenIndex, depth = position14, tokenIndex14, depth14
												if buffer[position] != rune('O') {
													goto l4
												}
												position++
											}
										l14:
											{
												position16, tokenIndex16, depth16 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l17
												}
												position++
												goto l16
											l17:
												position, tokenIndex, depth = position16, tokenIndex16, depth16
												if buffer[position] != rune('C') {
													goto l4
												}
												position++
											}
										l16:
											{
												position18, tokenIndex18, depth18 := position, tokenIndex, depth
												if buffer[position] != rune('k') {
													goto l19
												}
												position++
												goto l18
											l19:
												position, tokenIndex, depth = position18, tokenIndex18, depth18
												if buffer[position] != rune('K') {
													goto l4
												}
												position++
											}
										l18:
											if !_rules[ruleSpace]() {
												goto l4
											}
											{
												add(ruleAction29, position)
											}
											depth--
											add(ruleBLOCK, position9)
										}
										{
											position21 := position
											depth++
											if !_rules[ruleIntegerLiteral]() {
												goto l4
											}
											depth--
											add(rulePegText, position21)
										}
										{
											add(ruleAction10, position)
										}
										if !_rules[ruleSpace]() {
											goto l4
										}
										{
											position23 := position
											depth++
											if !_rules[ruleIntegerLiteral]() {
												goto l4
											}
											depth--
											add(rulePegText, position23)
										}
										{
											add(ruleAction11, position)
										}
										depth--
										add(rulePseudo, position8)
									}
									break
								case ':':
									{
										position25 := position
										depth++
										{
											position26 := position
											depth++
											if buffer[position] != rune(':') {
												goto l4
											}
											position++
											if !_rules[ruleSpacing]() {
												goto l4
											}
											depth--
											add(ruleCOLON, position26)
										}
										{
											position27 := position
											depth++
											if !_rules[ruleIdentifier]() {
												goto l4
											}
											depth--
											add(rulePegText, position27)
										}
										{
											add(ruleAction4, position)
										}
										depth--
										add(ruleLabel, position25)
									}
									break
								case ';':
									if !_rules[ruleComment]() {
										goto l4
									}
									break
								default:
									{
										position29 := position
										depth++
										{
											position30, tokenIndex30, depth30 := position, tokenIndex, depth
											{
												position32, tokenIndex32, depth32 := position, tokenIndex, depth
												{
													position34 := position
													depth++
													{
														position35, tokenIndex35, depth35 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l36
														}
														position++
														goto l35
													l36:
														position, tokenIndex, depth = position35, tokenIndex35, depth35
														if buffer[position] != rune('P') {
															goto l33
														}
														position++
													}
												l35:
													{
														position37, tokenIndex37, depth37 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l38
														}
														position++
														goto l37
													l38:
														position, tokenIndex, depth = position37, tokenIndex37, depth37
														if buffer[position] != rune('U') {
															goto l33
														}
														position++
													}
												l37:
													{
														position39, tokenIndex39, depth39 := position, tokenIndex, depth
														if buffer[position] != rune('s') {
															goto l40
														}
														position++
														goto l39
													l40:
														position, tokenIndex, depth = position39, tokenIndex39, depth39
														if buffer[position] != rune('S') {
															goto l33
														}
														position++
													}
												l39:
													{
														position41, tokenIndex41, depth41 := position, tokenIndex, depth
														if buffer[position] != rune('h') {
															goto l42
														}
														position++
														goto l41
													l42:
														position, tokenIndex, depth = position41, tokenIndex41, depth41
														if buffer[position] != rune('H') {
															goto l33
														}
														position++
													}
												l41:
													if !_rules[ruleSpace]() {
														goto l33
													}
													{
														add(ruleAction20, position)
													}
													depth--
													add(rulePUSH, position34)
												}
												goto l32
											l33:
												position, tokenIndex, depth = position32, tokenIndex32, depth32
												{
													switch buffer[position] {
													case 'J', 'j':
														{
															position45 := position
															depth++
															{
																position46, tokenIndex46, depth46 := position, tokenIndex, depth
																if buffer[position] != rune('j') {
																	goto l47
																}
																position++
																goto l46
															l47:
																position, tokenIndex, depth = position46, tokenIndex46, depth46
																if buffer[position] != rune('J') {
																	goto l31
																}
																position++
															}
														l46:
															{
																position48, tokenIndex48, depth48 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l49
																}
																position++
																goto l48
															l49:
																position, tokenIndex, depth = position48, tokenIndex48, depth48
																if buffer[position] != rune('M') {
																	goto l31
																}
																position++
															}
														l48:
															{
																position50, tokenIndex50, depth50 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l51
																}
																position++
																goto l50
															l51:
																position, tokenIndex, depth = position50, tokenIndex50, depth50
																if buffer[position] != rune('P') {
																	goto l31
																}
																position++
															}
														l50:
															if !_rules[ruleSpace]() {
																goto l31
															}
															{
																add(ruleAction22, position)
															}
															depth--
															add(ruleJMP, position45)
														}
														break
													case 'P', 'p':
														{
															position53 := position
															depth++
															{
																position54, tokenIndex54, depth54 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l55
																}
																position++
																goto l54
															l55:
																position, tokenIndex, depth = position54, tokenIndex54, depth54
																if buffer[position] != rune('P') {
																	goto l31
																}
																position++
															}
														l54:
															{
																position56, tokenIndex56, depth56 := position, tokenIndex, depth
																if buffer[position] != rune('o') {
																	goto l57
																}
																position++
																goto l56
															l57:
																position, tokenIndex, depth = position56, tokenIndex56, depth56
																if buffer[position] != rune('O') {
																	goto l31
																}
																position++
															}
														l56:
															{
																position58, tokenIndex58, depth58 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l59
																}
																position++
																goto l58
															l59:
																position, tokenIndex, depth = position58, tokenIndex58, depth58
																if buffer[position] != rune('P') {
																	goto l31
																}
																position++
															}
														l58:
															if !_rules[ruleSpace]() {
																goto l31
															}
															{
																add(ruleAction21, position)
															}
															depth--
															add(rulePOP, position53)
														}
														break
													default:
														{
															position61 := position
															depth++
															{
																position62, tokenIndex62, depth62 := position, tokenIndex, depth
																if buffer[position] != rune('c') {
																	goto l63
																}
																position++
																goto l62
															l63:
																position, tokenIndex, depth = position62, tokenIndex62, depth62
																if buffer[position] != rune('C') {
																	goto l31
																}
																position++
															}
														l62:
															{
																position64, tokenIndex64, depth64 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l65
																}
																position++
																goto l64
															l65:
																position, tokenIndex, depth = position64, tokenIndex64, depth64
																if buffer[position] != rune('A') {
																	goto l31
																}
																position++
															}
														l64:
															{
																position66, tokenIndex66, depth66 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l67
																}
																position++
																goto l66
															l67:
																position, tokenIndex, depth = position66, tokenIndex66, depth66
																if buffer[position] != rune('L') {
																	goto l31
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
																	goto l31
																}
																position++
															}
														l68:
															if !_rules[ruleSpace]() {
																goto l31
															}
															{
																add(ruleAction19, position)
															}
															depth--
															add(ruleCALL, position61)
														}
														break
													}
												}

											}
										l32:
											if !_rules[ruleOperand]() {
												goto l31
											}
											goto l30
										l31:
											position, tokenIndex, depth = position30, tokenIndex30, depth30
											{
												position72 := position
												depth++
												{
													position73, tokenIndex73, depth73 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l74
													}
													position++
													goto l73
												l74:
													position, tokenIndex, depth = position73, tokenIndex73, depth73
													if buffer[position] != rune('C') {
														goto l71
													}
													position++
												}
											l73:
												{
													position75, tokenIndex75, depth75 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l76
													}
													position++
													goto l75
												l76:
													position, tokenIndex, depth = position75, tokenIndex75, depth75
													if buffer[position] != rune('A') {
														goto l71
													}
													position++
												}
											l75:
												{
													position77, tokenIndex77, depth77 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l78
													}
													position++
													goto l77
												l78:
													position, tokenIndex, depth = position77, tokenIndex77, depth77
													if buffer[position] != rune('L') {
														goto l71
													}
													position++
												}
											l77:
												if !_rules[ruleSpace]() {
													goto l71
												}
												{
													add(ruleAction25, position)
												}
												depth--
												add(ruleCAL, position72)
											}
											{
												position80 := position
												depth++
												if !_rules[ruleDATA_TYPE]() {
													goto l71
												}
												depth--
												add(rulePegText, position80)
											}
											{
												add(ruleAction5, position)
											}
											{
												position82 := position
												depth++
												{
													position83 := position
													depth++
													{
														position84, tokenIndex84, depth84 := position, tokenIndex, depth
														{
															position86, tokenIndex86, depth86 := position, tokenIndex, depth
															if buffer[position] != rune('m') {
																goto l87
															}
															position++
															goto l86
														l87:
															position, tokenIndex, depth = position86, tokenIndex86, depth86
															if buffer[position] != rune('M') {
																goto l85
															}
															position++
														}
													l86:
														{
															position88, tokenIndex88, depth88 := position, tokenIndex, depth
															if buffer[position] != rune('u') {
																goto l89
															}
															position++
															goto l88
														l89:
															position, tokenIndex, depth = position88, tokenIndex88, depth88
															if buffer[position] != rune('U') {
																goto l85
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
																goto l85
															}
															position++
														}
													l90:
														goto l84
													l85:
														position, tokenIndex, depth = position84, tokenIndex84, depth84
														{
															switch buffer[position] {
															case 'M', 'm':
																{
																	position93, tokenIndex93, depth93 := position, tokenIndex, depth
																	if buffer[position] != rune('m') {
																		goto l94
																	}
																	position++
																	goto l93
																l94:
																	position, tokenIndex, depth = position93, tokenIndex93, depth93
																	if buffer[position] != rune('M') {
																		goto l71
																	}
																	position++
																}
															l93:
																{
																	position95, tokenIndex95, depth95 := position, tokenIndex, depth
																	if buffer[position] != rune('o') {
																		goto l96
																	}
																	position++
																	goto l95
																l96:
																	position, tokenIndex, depth = position95, tokenIndex95, depth95
																	if buffer[position] != rune('O') {
																		goto l71
																	}
																	position++
																}
															l95:
																{
																	position97, tokenIndex97, depth97 := position, tokenIndex, depth
																	if buffer[position] != rune('d') {
																		goto l98
																	}
																	position++
																	goto l97
																l98:
																	position, tokenIndex, depth = position97, tokenIndex97, depth97
																	if buffer[position] != rune('D') {
																		goto l71
																	}
																	position++
																}
															l97:
																break
															case 'D', 'd':
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
																		goto l71
																	}
																	position++
																}
															l99:
																{
																	position101, tokenIndex101, depth101 := position, tokenIndex, depth
																	if buffer[position] != rune('i') {
																		goto l102
																	}
																	position++
																	goto l101
																l102:
																	position, tokenIndex, depth = position101, tokenIndex101, depth101
																	if buffer[position] != rune('I') {
																		goto l71
																	}
																	position++
																}
															l101:
																{
																	position103, tokenIndex103, depth103 := position, tokenIndex, depth
																	if buffer[position] != rune('v') {
																		goto l104
																	}
																	position++
																	goto l103
																l104:
																	position, tokenIndex, depth = position103, tokenIndex103, depth103
																	if buffer[position] != rune('V') {
																		goto l71
																	}
																	position++
																}
															l103:
																break
															case 'S', 's':
																{
																	position105, tokenIndex105, depth105 := position, tokenIndex, depth
																	if buffer[position] != rune('s') {
																		goto l106
																	}
																	position++
																	goto l105
																l106:
																	position, tokenIndex, depth = position105, tokenIndex105, depth105
																	if buffer[position] != rune('S') {
																		goto l71
																	}
																	position++
																}
															l105:
																{
																	position107, tokenIndex107, depth107 := position, tokenIndex, depth
																	if buffer[position] != rune('u') {
																		goto l108
																	}
																	position++
																	goto l107
																l108:
																	position, tokenIndex, depth = position107, tokenIndex107, depth107
																	if buffer[position] != rune('U') {
																		goto l71
																	}
																	position++
																}
															l107:
																{
																	position109, tokenIndex109, depth109 := position, tokenIndex, depth
																	if buffer[position] != rune('b') {
																		goto l110
																	}
																	position++
																	goto l109
																l110:
																	position, tokenIndex, depth = position109, tokenIndex109, depth109
																	if buffer[position] != rune('B') {
																		goto l71
																	}
																	position++
																}
															l109:
																break
															default:
																{
																	position111, tokenIndex111, depth111 := position, tokenIndex, depth
																	if buffer[position] != rune('a') {
																		goto l112
																	}
																	position++
																	goto l111
																l112:
																	position, tokenIndex, depth = position111, tokenIndex111, depth111
																	if buffer[position] != rune('A') {
																		goto l71
																	}
																	position++
																}
															l111:
																{
																	position113, tokenIndex113, depth113 := position, tokenIndex, depth
																	if buffer[position] != rune('d') {
																		goto l114
																	}
																	position++
																	goto l113
																l114:
																	position, tokenIndex, depth = position113, tokenIndex113, depth113
																	if buffer[position] != rune('D') {
																		goto l71
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
																		goto l71
																	}
																	position++
																}
															l115:
																break
															}
														}

													}
												l84:
													if !_rules[ruleSpace]() {
														goto l71
													}
													depth--
													add(ruleCAL_OP, position83)
												}
												depth--
												add(rulePegText, position82)
											}
											{
												add(ruleAction6, position)
											}
											if !_rules[ruleOperand]() {
												goto l71
											}
											if !_rules[ruleCOMMA]() {
												goto l71
											}
											if !_rules[ruleOperand]() {
												goto l71
											}
											goto l30
										l71:
											position, tokenIndex, depth = position30, tokenIndex30, depth30
											{
												switch buffer[position] {
												case 'J', 'j':
													{
														position119 := position
														depth++
														{
															position120, tokenIndex120, depth120 := position, tokenIndex, depth
															if buffer[position] != rune('j') {
																goto l121
															}
															position++
															goto l120
														l121:
															position, tokenIndex, depth = position120, tokenIndex120, depth120
															if buffer[position] != rune('J') {
																goto l4
															}
															position++
														}
													l120:
														{
															position122, tokenIndex122, depth122 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l123
															}
															position++
															goto l122
														l123:
															position, tokenIndex, depth = position122, tokenIndex122, depth122
															if buffer[position] != rune('P') {
																goto l4
															}
															position++
														}
													l122:
														{
															position124, tokenIndex124, depth124 := position, tokenIndex, depth
															if buffer[position] != rune('c') {
																goto l125
															}
															position++
															goto l124
														l125:
															position, tokenIndex, depth = position124, tokenIndex124, depth124
															if buffer[position] != rune('C') {
																goto l4
															}
															position++
														}
													l124:
														if !_rules[ruleSpace]() {
															goto l4
														}
														{
															add(ruleAction28, position)
														}
														depth--
														add(ruleJPC, position119)
													}
													{
														position127 := position
														depth++
														{
															position128 := position
															depth++
															{
																position129, tokenIndex129, depth129 := position, tokenIndex, depth
																{
																	position131, tokenIndex131, depth131 := position, tokenIndex, depth
																	if buffer[position] != rune('b') {
																		goto l132
																	}
																	position++
																	goto l131
																l132:
																	position, tokenIndex, depth = position131, tokenIndex131, depth131
																	if buffer[position] != rune('B') {
																		goto l130
																	}
																	position++
																}
															l131:
																goto l129
															l130:
																position, tokenIndex, depth = position129, tokenIndex129, depth129
																{
																	position134, tokenIndex134, depth134 := position, tokenIndex, depth
																	if buffer[position] != rune('a') {
																		goto l135
																	}
																	position++
																	goto l134
																l135:
																	position, tokenIndex, depth = position134, tokenIndex134, depth134
																	if buffer[position] != rune('A') {
																		goto l133
																	}
																	position++
																}
															l134:
																goto l129
															l133:
																position, tokenIndex, depth = position129, tokenIndex129, depth129
																{
																	switch buffer[position] {
																	case 'N', 'n':
																		{
																			position137, tokenIndex137, depth137 := position, tokenIndex, depth
																			if buffer[position] != rune('n') {
																				goto l138
																			}
																			position++
																			goto l137
																		l138:
																			position, tokenIndex, depth = position137, tokenIndex137, depth137
																			if buffer[position] != rune('N') {
																				goto l4
																			}
																			position++
																		}
																	l137:
																		{
																			position139, tokenIndex139, depth139 := position, tokenIndex, depth
																			if buffer[position] != rune('z') {
																				goto l140
																			}
																			position++
																			goto l139
																		l140:
																			position, tokenIndex, depth = position139, tokenIndex139, depth139
																			if buffer[position] != rune('Z') {
																				goto l4
																			}
																			position++
																		}
																	l139:
																		break
																	case 'A', 'a':
																		{
																			position141, tokenIndex141, depth141 := position, tokenIndex, depth
																			if buffer[position] != rune('a') {
																				goto l142
																			}
																			position++
																			goto l141
																		l142:
																			position, tokenIndex, depth = position141, tokenIndex141, depth141
																			if buffer[position] != rune('A') {
																				goto l4
																			}
																			position++
																		}
																	l141:
																		{
																			position143, tokenIndex143, depth143 := position, tokenIndex, depth
																			if buffer[position] != rune('e') {
																				goto l144
																			}
																			position++
																			goto l143
																		l144:
																			position, tokenIndex, depth = position143, tokenIndex143, depth143
																			if buffer[position] != rune('E') {
																				goto l4
																			}
																			position++
																		}
																	l143:
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
																			position145, tokenIndex145, depth145 := position, tokenIndex, depth
																			if buffer[position] != rune('b') {
																				goto l146
																			}
																			position++
																			goto l145
																		l146:
																			position, tokenIndex, depth = position145, tokenIndex145, depth145
																			if buffer[position] != rune('B') {
																				goto l4
																			}
																			position++
																		}
																	l145:
																		{
																			position147, tokenIndex147, depth147 := position, tokenIndex, depth
																			if buffer[position] != rune('e') {
																				goto l148
																			}
																			position++
																			goto l147
																		l148:
																			position, tokenIndex, depth = position147, tokenIndex147, depth147
																			if buffer[position] != rune('E') {
																				goto l4
																			}
																			position++
																		}
																	l147:
																		break
																	}
																}

															}
														l129:
															if !_rules[ruleSpace]() {
																goto l4
															}
															depth--
															add(ruleCMP_OP, position128)
														}
														depth--
														add(rulePegText, position127)
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
														position150 := position
														depth++
														{
															position151, tokenIndex151, depth151 := position, tokenIndex, depth
															if buffer[position] != rune('c') {
																goto l152
															}
															position++
															goto l151
														l152:
															position, tokenIndex, depth = position151, tokenIndex151, depth151
															if buffer[position] != rune('C') {
																goto l4
															}
															position++
														}
													l151:
														{
															position153, tokenIndex153, depth153 := position, tokenIndex, depth
															if buffer[position] != rune('m') {
																goto l154
															}
															position++
															goto l153
														l154:
															position, tokenIndex, depth = position153, tokenIndex153, depth153
															if buffer[position] != rune('M') {
																goto l4
															}
															position++
														}
													l153:
														{
															position155, tokenIndex155, depth155 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l156
															}
															position++
															goto l155
														l156:
															position, tokenIndex, depth = position155, tokenIndex155, depth155
															if buffer[position] != rune('P') {
																goto l4
															}
															position++
														}
													l155:
														if !_rules[ruleSpace]() {
															goto l4
														}
														{
															add(ruleAction27, position)
														}
														depth--
														add(ruleCMP, position150)
													}
													{
														position158 := position
														depth++
														if !_rules[ruleDATA_TYPE]() {
															goto l4
														}
														depth--
														add(rulePegText, position158)
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
														position160 := position
														depth++
														{
															position161, tokenIndex161, depth161 := position, tokenIndex, depth
															if buffer[position] != rune('l') {
																goto l162
															}
															position++
															goto l161
														l162:
															position, tokenIndex, depth = position161, tokenIndex161, depth161
															if buffer[position] != rune('L') {
																goto l4
															}
															position++
														}
													l161:
														{
															position163, tokenIndex163, depth163 := position, tokenIndex, depth
															if buffer[position] != rune('d') {
																goto l164
															}
															position++
															goto l163
														l164:
															position, tokenIndex, depth = position163, tokenIndex163, depth163
															if buffer[position] != rune('D') {
																goto l4
															}
															position++
														}
													l163:
														if !_rules[ruleSpace]() {
															goto l4
														}
														{
															add(ruleAction26, position)
														}
														depth--
														add(ruleLD, position160)
													}
													{
														position166 := position
														depth++
														if !_rules[ruleDATA_TYPE]() {
															goto l4
														}
														depth--
														add(rulePegText, position166)
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
														position168 := position
														depth++
														{
															position169, tokenIndex169, depth169 := position, tokenIndex, depth
															if buffer[position] != rune('n') {
																goto l170
															}
															position++
															goto l169
														l170:
															position, tokenIndex, depth = position169, tokenIndex169, depth169
															if buffer[position] != rune('N') {
																goto l4
															}
															position++
														}
													l169:
														{
															position171, tokenIndex171, depth171 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l172
															}
															position++
															goto l171
														l172:
															position, tokenIndex, depth = position171, tokenIndex171, depth171
															if buffer[position] != rune('O') {
																goto l4
															}
															position++
														}
													l171:
														{
															position173, tokenIndex173, depth173 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l174
															}
															position++
															goto l173
														l174:
															position, tokenIndex, depth = position173, tokenIndex173, depth173
															if buffer[position] != rune('P') {
																goto l4
															}
															position++
														}
													l173:
														if !_rules[ruleSpacing]() {
															goto l4
														}
														{
															add(ruleAction18, position)
														}
														depth--
														add(ruleNOP, position168)
													}
													break
												case 'R', 'r':
													{
														position176 := position
														depth++
														{
															position177, tokenIndex177, depth177 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l178
															}
															position++
															goto l177
														l178:
															position, tokenIndex, depth = position177, tokenIndex177, depth177
															if buffer[position] != rune('R') {
																goto l4
															}
															position++
														}
													l177:
														{
															position179, tokenIndex179, depth179 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l180
															}
															position++
															goto l179
														l180:
															position, tokenIndex, depth = position179, tokenIndex179, depth179
															if buffer[position] != rune('E') {
																goto l4
															}
															position++
														}
													l179:
														{
															position181, tokenIndex181, depth181 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l182
															}
															position++
															goto l181
														l182:
															position, tokenIndex, depth = position181, tokenIndex181, depth181
															if buffer[position] != rune('T') {
																goto l4
															}
															position++
														}
													l181:
														if !_rules[ruleSpacing]() {
															goto l4
														}
														{
															add(ruleAction17, position)
														}
														depth--
														add(ruleRET, position176)
													}
													break
												case 'E', 'e':
													{
														position184 := position
														depth++
														{
															position185, tokenIndex185, depth185 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l186
															}
															position++
															goto l185
														l186:
															position, tokenIndex, depth = position185, tokenIndex185, depth185
															if buffer[position] != rune('E') {
																goto l4
															}
															position++
														}
													l185:
														{
															position187, tokenIndex187, depth187 := position, tokenIndex, depth
															if buffer[position] != rune('x') {
																goto l188
															}
															position++
															goto l187
														l188:
															position, tokenIndex, depth = position187, tokenIndex187, depth187
															if buffer[position] != rune('X') {
																goto l4
															}
															position++
														}
													l187:
														{
															position189, tokenIndex189, depth189 := position, tokenIndex, depth
															if buffer[position] != rune('i') {
																goto l190
															}
															position++
															goto l189
														l190:
															position, tokenIndex, depth = position189, tokenIndex189, depth189
															if buffer[position] != rune('I') {
																goto l4
															}
															position++
														}
													l189:
														{
															position191, tokenIndex191, depth191 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l192
															}
															position++
															goto l191
														l192:
															position, tokenIndex, depth = position191, tokenIndex191, depth191
															if buffer[position] != rune('T') {
																goto l4
															}
															position++
														}
													l191:
														if !_rules[ruleSpacing]() {
															goto l4
														}
														{
															add(ruleAction16, position)
														}
														depth--
														add(ruleEXIT, position184)
													}
													break
												default:
													{
														position194, tokenIndex194, depth194 := position, tokenIndex, depth
														{
															position196 := position
															depth++
															{
																position197, tokenIndex197, depth197 := position, tokenIndex, depth
																if buffer[position] != rune('i') {
																	goto l198
																}
																position++
																goto l197
															l198:
																position, tokenIndex, depth = position197, tokenIndex197, depth197
																if buffer[position] != rune('I') {
																	goto l195
																}
																position++
															}
														l197:
															{
																position199, tokenIndex199, depth199 := position, tokenIndex, depth
																if buffer[position] != rune('n') {
																	goto l200
																}
																position++
																goto l199
															l200:
																position, tokenIndex, depth = position199, tokenIndex199, depth199
																if buffer[position] != rune('N') {
																	goto l195
																}
																position++
															}
														l199:
															if !_rules[ruleSpace]() {
																goto l195
															}
															{
																add(ruleAction23, position)
															}
															depth--
															add(ruleIN, position196)
														}
														goto l194
													l195:
														position, tokenIndex, depth = position194, tokenIndex194, depth194
														{
															position202 := position
															depth++
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
																if buffer[position] != rune('u') {
																	goto l206
																}
																position++
																goto l205
															l206:
																position, tokenIndex, depth = position205, tokenIndex205, depth205
																if buffer[position] != rune('U') {
																	goto l4
																}
																position++
															}
														l205:
															{
																position207, tokenIndex207, depth207 := position, tokenIndex, depth
																if buffer[position] != rune('t') {
																	goto l208
																}
																position++
																goto l207
															l208:
																position, tokenIndex, depth = position207, tokenIndex207, depth207
																if buffer[position] != rune('T') {
																	goto l4
																}
																position++
															}
														l207:
															if !_rules[ruleSpace]() {
																goto l4
															}
															{
																add(ruleAction24, position)
															}
															depth--
															add(ruleOUT, position202)
														}
													}
												l194:
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
									l30:
										depth--
										add(ruleInst, position29)
									}
									break
								}
							}

							{
								add(ruleAction1, position)
							}
							{
								position211, tokenIndex211, depth211 := position, tokenIndex, depth
								if !_rules[ruleComment]() {
									goto l211
								}
								{
									add(ruleAction2, position)
								}
								goto l212
							l211:
								position, tokenIndex, depth = position211, tokenIndex211, depth211
							}
						l212:
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
					position215 := position
					depth++
					{
						position216, tokenIndex216, depth216 := position, tokenIndex, depth
						if !matchDot() {
							goto l216
						}
						goto l0
					l216:
						position, tokenIndex, depth = position216, tokenIndex216, depth216
					}
					depth--
					add(ruleEOT, position215)
				}
				depth--
				add(ruleStart, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Line <- <(((&('.') Pseudo) | (&(':') Label) | (&(';') Comment) | (&('C' | 'E' | 'I' | 'J' | 'L' | 'N' | 'O' | 'P' | 'R' | 'c' | 'e' | 'i' | 'j' | 'l' | 'n' | 'o' | 'p' | 'r') Inst)) Action1 (Comment Action2)?)> */
		nil,
		/* 2 Comment <- <(SEMICOLON <(!NL .)*> Action3)> */
		func() bool {
			position218, tokenIndex218, depth218 := position, tokenIndex, depth
			{
				position219 := position
				depth++
				{
					position220 := position
					depth++
					if buffer[position] != rune(';') {
						goto l218
					}
					position++
					depth--
					add(ruleSEMICOLON, position220)
				}
				{
					position221 := position
					depth++
				l222:
					{
						position223, tokenIndex223, depth223 := position, tokenIndex, depth
						{
							position224, tokenIndex224, depth224 := position, tokenIndex, depth
							if !_rules[ruleNL]() {
								goto l224
							}
							goto l223
						l224:
							position, tokenIndex, depth = position224, tokenIndex224, depth224
						}
						if !matchDot() {
							goto l223
						}
						goto l222
					l223:
						position, tokenIndex, depth = position223, tokenIndex223, depth223
					}
					depth--
					add(rulePegText, position221)
				}
				{
					add(ruleAction3, position)
				}
				depth--
				add(ruleComment, position219)
			}
			return true
		l218:
			position, tokenIndex, depth = position218, tokenIndex218, depth218
			return false
		},
		/* 3 Label <- <(COLON <Identifier> Action4)> */
		nil,
		/* 4 Inst <- <(((PUSH / ((&('J' | 'j') JMP) | (&('P' | 'p') POP) | (&('C' | 'c') CALL))) Operand) / (CAL <DATA_TYPE> Action5 <CAL_OP> Action6 Operand COMMA Operand) / ((&('J' | 'j') (JPC <CMP_OP> Action9 Operand)) | (&('C' | 'c') (CMP <DATA_TYPE> Action8 Operand COMMA Operand)) | (&('L' | 'l') (LD <DATA_TYPE> Action7 Operand COMMA Operand)) | (&('N' | 'n') NOP) | (&('R' | 'r') RET) | (&('E' | 'e') EXIT) | (&('I' | 'O' | 'i' | 'o') ((IN / OUT) Operand COMMA Operand))))> */
		nil,
		/* 5 Pseudo <- <(BLOCK <IntegerLiteral> Action10 Space <IntegerLiteral> Action11)> */
		nil,
		/* 6 Operand <- <(((LBRK <Identifier> RBRK Action13) / (<IntegerLiteral> Action14) / ((&('"' | '\'' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Literal) | (&('[') (LBRK <IntegerLiteral> RBRK Action15)) | (&('$' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') (<Identifier> Action12)))) Spacing)> */
		func() bool {
			position229, tokenIndex229, depth229 := position, tokenIndex, depth
			{
				position230 := position
				depth++
				{
					position231, tokenIndex231, depth231 := position, tokenIndex, depth
					if !_rules[ruleLBRK]() {
						goto l232
					}
					{
						position233 := position
						depth++
						if !_rules[ruleIdentifier]() {
							goto l232
						}
						depth--
						add(rulePegText, position233)
					}
					if !_rules[ruleRBRK]() {
						goto l232
					}
					{
						add(ruleAction13, position)
					}
					goto l231
				l232:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
					{
						position236 := position
						depth++
						if !_rules[ruleIntegerLiteral]() {
							goto l235
						}
						depth--
						add(rulePegText, position236)
					}
					{
						add(ruleAction14, position)
					}
					goto l231
				l235:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
					{
						switch buffer[position] {
						case '"', '\'', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							{
								position239 := position
								depth++
								{
									position240, tokenIndex240, depth240 := position, tokenIndex, depth
									{
										position242 := position
										depth++
										{
											position243, tokenIndex243, depth243 := position, tokenIndex, depth
											{
												position245 := position
												depth++
												{
													position246 := position
													depth++
													{
														position247, tokenIndex247, depth247 := position, tokenIndex, depth
														{
															position249, tokenIndex249, depth249 := position, tokenIndex, depth
															if buffer[position] != rune('0') {
																goto l250
															}
															position++
															if buffer[position] != rune('x') {
																goto l250
															}
															position++
															goto l249
														l250:
															position, tokenIndex, depth = position249, tokenIndex249, depth249
															if buffer[position] != rune('0') {
																goto l248
															}
															position++
															if buffer[position] != rune('X') {
																goto l248
															}
															position++
														}
													l249:
														{
															position251, tokenIndex251, depth251 := position, tokenIndex, depth
															if !_rules[ruleHexDigits]() {
																goto l251
															}
															goto l252
														l251:
															position, tokenIndex, depth = position251, tokenIndex251, depth251
														}
													l252:
														if buffer[position] != rune('.') {
															goto l248
														}
														position++
														if !_rules[ruleHexDigits]() {
															goto l248
														}
														goto l247
													l248:
														position, tokenIndex, depth = position247, tokenIndex247, depth247
														if !_rules[ruleHexNumeral]() {
															goto l244
														}
														{
															position253, tokenIndex253, depth253 := position, tokenIndex, depth
															if buffer[position] != rune('.') {
																goto l253
															}
															position++
															goto l254
														l253:
															position, tokenIndex, depth = position253, tokenIndex253, depth253
														}
													l254:
													}
												l247:
													depth--
													add(ruleHexSignificand, position246)
												}
												{
													position255 := position
													depth++
													{
														position256, tokenIndex256, depth256 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l257
														}
														position++
														goto l256
													l257:
														position, tokenIndex, depth = position256, tokenIndex256, depth256
														if buffer[position] != rune('P') {
															goto l244
														}
														position++
													}
												l256:
													{
														position258, tokenIndex258, depth258 := position, tokenIndex, depth
														{
															position260, tokenIndex260, depth260 := position, tokenIndex, depth
															if buffer[position] != rune('+') {
																goto l261
															}
															position++
															goto l260
														l261:
															position, tokenIndex, depth = position260, tokenIndex260, depth260
															if buffer[position] != rune('-') {
																goto l258
															}
															position++
														}
													l260:
														goto l259
													l258:
														position, tokenIndex, depth = position258, tokenIndex258, depth258
													}
												l259:
													if !_rules[ruleDigits]() {
														goto l244
													}
													depth--
													add(ruleBinaryExponent, position255)
												}
												{
													position262, tokenIndex262, depth262 := position, tokenIndex, depth
													{
														switch buffer[position] {
														case 'D':
															if buffer[position] != rune('D') {
																goto l262
															}
															position++
															break
														case 'd':
															if buffer[position] != rune('d') {
																goto l262
															}
															position++
															break
														case 'F':
															if buffer[position] != rune('F') {
																goto l262
															}
															position++
															break
														default:
															if buffer[position] != rune('f') {
																goto l262
															}
															position++
															break
														}
													}

													goto l263
												l262:
													position, tokenIndex, depth = position262, tokenIndex262, depth262
												}
											l263:
												depth--
												add(ruleHexFloat, position245)
											}
											goto l243
										l244:
											position, tokenIndex, depth = position243, tokenIndex243, depth243
											{
												position265 := position
												depth++
												{
													position266, tokenIndex266, depth266 := position, tokenIndex, depth
													if !_rules[ruleDigits]() {
														goto l267
													}
													if buffer[position] != rune('.') {
														goto l267
													}
													position++
													{
														position268, tokenIndex268, depth268 := position, tokenIndex, depth
														if !_rules[ruleDigits]() {
															goto l268
														}
														goto l269
													l268:
														position, tokenIndex, depth = position268, tokenIndex268, depth268
													}
												l269:
													{
														position270, tokenIndex270, depth270 := position, tokenIndex, depth
														if !_rules[ruleExponent]() {
															goto l270
														}
														goto l271
													l270:
														position, tokenIndex, depth = position270, tokenIndex270, depth270
													}
												l271:
													{
														position272, tokenIndex272, depth272 := position, tokenIndex, depth
														{
															switch buffer[position] {
															case 'D':
																if buffer[position] != rune('D') {
																	goto l272
																}
																position++
																break
															case 'd':
																if buffer[position] != rune('d') {
																	goto l272
																}
																position++
																break
															case 'F':
																if buffer[position] != rune('F') {
																	goto l272
																}
																position++
																break
															default:
																if buffer[position] != rune('f') {
																	goto l272
																}
																position++
																break
															}
														}

														goto l273
													l272:
														position, tokenIndex, depth = position272, tokenIndex272, depth272
													}
												l273:
													goto l266
												l267:
													position, tokenIndex, depth = position266, tokenIndex266, depth266
													if buffer[position] != rune('.') {
														goto l275
													}
													position++
													if !_rules[ruleDigits]() {
														goto l275
													}
													{
														position276, tokenIndex276, depth276 := position, tokenIndex, depth
														if !_rules[ruleExponent]() {
															goto l276
														}
														goto l277
													l276:
														position, tokenIndex, depth = position276, tokenIndex276, depth276
													}
												l277:
													{
														position278, tokenIndex278, depth278 := position, tokenIndex, depth
														{
															switch buffer[position] {
															case 'D':
																if buffer[position] != rune('D') {
																	goto l278
																}
																position++
																break
															case 'd':
																if buffer[position] != rune('d') {
																	goto l278
																}
																position++
																break
															case 'F':
																if buffer[position] != rune('F') {
																	goto l278
																}
																position++
																break
															default:
																if buffer[position] != rune('f') {
																	goto l278
																}
																position++
																break
															}
														}

														goto l279
													l278:
														position, tokenIndex, depth = position278, tokenIndex278, depth278
													}
												l279:
													goto l266
												l275:
													position, tokenIndex, depth = position266, tokenIndex266, depth266
													if !_rules[ruleDigits]() {
														goto l281
													}
													if !_rules[ruleExponent]() {
														goto l281
													}
													{
														position282, tokenIndex282, depth282 := position, tokenIndex, depth
														{
															switch buffer[position] {
															case 'D':
																if buffer[position] != rune('D') {
																	goto l282
																}
																position++
																break
															case 'd':
																if buffer[position] != rune('d') {
																	goto l282
																}
																position++
																break
															case 'F':
																if buffer[position] != rune('F') {
																	goto l282
																}
																position++
																break
															default:
																if buffer[position] != rune('f') {
																	goto l282
																}
																position++
																break
															}
														}

														goto l283
													l282:
														position, tokenIndex, depth = position282, tokenIndex282, depth282
													}
												l283:
													goto l266
												l281:
													position, tokenIndex, depth = position266, tokenIndex266, depth266
													if !_rules[ruleDigits]() {
														goto l241
													}
													{
														position285, tokenIndex285, depth285 := position, tokenIndex, depth
														if !_rules[ruleExponent]() {
															goto l285
														}
														goto l286
													l285:
														position, tokenIndex, depth = position285, tokenIndex285, depth285
													}
												l286:
													{
														switch buffer[position] {
														case 'D':
															if buffer[position] != rune('D') {
																goto l241
															}
															position++
															break
														case 'd':
															if buffer[position] != rune('d') {
																goto l241
															}
															position++
															break
														case 'F':
															if buffer[position] != rune('F') {
																goto l241
															}
															position++
															break
														default:
															if buffer[position] != rune('f') {
																goto l241
															}
															position++
															break
														}
													}

												}
											l266:
												depth--
												add(ruleDecimalFloat, position265)
											}
										}
									l243:
										depth--
										add(ruleFloatLiteral, position242)
									}
									goto l240
								l241:
									position, tokenIndex, depth = position240, tokenIndex240, depth240
									{
										switch buffer[position] {
										case '"':
											{
												position289 := position
												depth++
												if buffer[position] != rune('"') {
													goto l229
												}
												position++
											l290:
												{
													position291, tokenIndex291, depth291 := position, tokenIndex, depth
													{
														position292, tokenIndex292, depth292 := position, tokenIndex, depth
														if !_rules[ruleEscape]() {
															goto l293
														}
														goto l292
													l293:
														position, tokenIndex, depth = position292, tokenIndex292, depth292
														{
															position294, tokenIndex294, depth294 := position, tokenIndex, depth
															{
																switch buffer[position] {
																case '\r':
																	if buffer[position] != rune('\r') {
																		goto l294
																	}
																	position++
																	break
																case '\n':
																	if buffer[position] != rune('\n') {
																		goto l294
																	}
																	position++
																	break
																case '\\':
																	if buffer[position] != rune('\\') {
																		goto l294
																	}
																	position++
																	break
																default:
																	if buffer[position] != rune('"') {
																		goto l294
																	}
																	position++
																	break
																}
															}

															goto l291
														l294:
															position, tokenIndex, depth = position294, tokenIndex294, depth294
														}
														if !matchDot() {
															goto l291
														}
													}
												l292:
													goto l290
												l291:
													position, tokenIndex, depth = position291, tokenIndex291, depth291
												}
												if buffer[position] != rune('"') {
													goto l229
												}
												position++
												depth--
												add(ruleStringLiteral, position289)
											}
											break
										case '\'':
											{
												position296 := position
												depth++
												if buffer[position] != rune('\'') {
													goto l229
												}
												position++
												{
													position297, tokenIndex297, depth297 := position, tokenIndex, depth
													if !_rules[ruleEscape]() {
														goto l298
													}
													goto l297
												l298:
													position, tokenIndex, depth = position297, tokenIndex297, depth297
													{
														position299, tokenIndex299, depth299 := position, tokenIndex, depth
														{
															position300, tokenIndex300, depth300 := position, tokenIndex, depth
															if buffer[position] != rune('\'') {
																goto l301
															}
															position++
															goto l300
														l301:
															position, tokenIndex, depth = position300, tokenIndex300, depth300
															if buffer[position] != rune('\\') {
																goto l299
															}
															position++
														}
													l300:
														goto l229
													l299:
														position, tokenIndex, depth = position299, tokenIndex299, depth299
													}
													if !matchDot() {
														goto l229
													}
												}
											l297:
												if buffer[position] != rune('\'') {
													goto l229
												}
												position++
												depth--
												add(ruleCharLiteral, position296)
											}
											break
										default:
											if !_rules[ruleIntegerLiteral]() {
												goto l229
											}
											break
										}
									}

								}
							l240:
								if !_rules[ruleSpacing]() {
									goto l229
								}
								depth--
								add(ruleLiteral, position239)
							}
							break
						case '[':
							if !_rules[ruleLBRK]() {
								goto l229
							}
							{
								position302 := position
								depth++
								if !_rules[ruleIntegerLiteral]() {
									goto l229
								}
								depth--
								add(rulePegText, position302)
							}
							if !_rules[ruleRBRK]() {
								goto l229
							}
							{
								add(ruleAction15, position)
							}
							break
						default:
							{
								position304 := position
								depth++
								if !_rules[ruleIdentifier]() {
									goto l229
								}
								depth--
								add(rulePegText, position304)
							}
							{
								add(ruleAction12, position)
							}
							break
						}
					}

				}
			l231:
				if !_rules[ruleSpacing]() {
					goto l229
				}
				depth--
				add(ruleOperand, position230)
			}
			return true
		l229:
			position, tokenIndex, depth = position229, tokenIndex229, depth229
			return false
		},
		/* 7 Spacing <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position307 := position
				depth++
			l308:
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l309
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l309
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l309
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l309
							}
							position++
							break
						}
					}

					goto l308
				l309:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
				}
				depth--
				add(ruleSpacing, position307)
			}
			return true
		},
		/* 8 Space <- <((&('\f') '\f') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))+> */
		func() bool {
			position311, tokenIndex311, depth311 := position, tokenIndex, depth
			{
				position312 := position
				depth++
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

			l313:
				{
					position314, tokenIndex314, depth314 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l314
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l314
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l314
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l314
							}
							position++
							break
						}
					}

					goto l313
				l314:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
				}
				depth--
				add(ruleSpace, position312)
			}
			return true
		l311:
			position, tokenIndex, depth = position311, tokenIndex311, depth311
			return false
		},
		/* 9 Identifier <- <(Letter LetterOrDigit* Spacing)> */
		func() bool {
			position317, tokenIndex317, depth317 := position, tokenIndex, depth
			{
				position318 := position
				depth++
				{
					position319 := position
					depth++
					{
						switch buffer[position] {
						case '$', '_':
							{
								position321, tokenIndex321, depth321 := position, tokenIndex, depth
								if buffer[position] != rune('_') {
									goto l322
								}
								position++
								goto l321
							l322:
								position, tokenIndex, depth = position321, tokenIndex321, depth321
								if buffer[position] != rune('$') {
									goto l317
								}
								position++
							}
						l321:
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l317
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l317
							}
							position++
							break
						}
					}

					depth--
					add(ruleLetter, position319)
				}
			l323:
				{
					position324, tokenIndex324, depth324 := position, tokenIndex, depth
					{
						position325 := position
						depth++
						{
							switch buffer[position] {
							case '$', '_':
								{
									position327, tokenIndex327, depth327 := position, tokenIndex, depth
									if buffer[position] != rune('_') {
										goto l328
									}
									position++
									goto l327
								l328:
									position, tokenIndex, depth = position327, tokenIndex327, depth327
									if buffer[position] != rune('$') {
										goto l324
									}
									position++
								}
							l327:
								break
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l324
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l324
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l324
								}
								position++
								break
							}
						}

						depth--
						add(ruleLetterOrDigit, position325)
					}
					goto l323
				l324:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
				}
				if !_rules[ruleSpacing]() {
					goto l317
				}
				depth--
				add(ruleIdentifier, position318)
			}
			return true
		l317:
			position, tokenIndex, depth = position317, tokenIndex317, depth317
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
			position347, tokenIndex347, depth347 := position, tokenIndex, depth
			{
				position348 := position
				depth++
				{
					switch buffer[position] {
					case 'I', 'i':
						{
							position350, tokenIndex350, depth350 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l351
							}
							position++
							goto l350
						l351:
							position, tokenIndex, depth = position350, tokenIndex350, depth350
							if buffer[position] != rune('I') {
								goto l347
							}
							position++
						}
					l350:
						{
							position352, tokenIndex352, depth352 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l353
							}
							position++
							goto l352
						l353:
							position, tokenIndex, depth = position352, tokenIndex352, depth352
							if buffer[position] != rune('N') {
								goto l347
							}
							position++
						}
					l352:
						{
							position354, tokenIndex354, depth354 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l355
							}
							position++
							goto l354
						l355:
							position, tokenIndex, depth = position354, tokenIndex354, depth354
							if buffer[position] != rune('T') {
								goto l347
							}
							position++
						}
					l354:
						break
					case 'F', 'f':
						{
							position356, tokenIndex356, depth356 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l357
							}
							position++
							goto l356
						l357:
							position, tokenIndex, depth = position356, tokenIndex356, depth356
							if buffer[position] != rune('F') {
								goto l347
							}
							position++
						}
					l356:
						{
							position358, tokenIndex358, depth358 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l359
							}
							position++
							goto l358
						l359:
							position, tokenIndex, depth = position358, tokenIndex358, depth358
							if buffer[position] != rune('L') {
								goto l347
							}
							position++
						}
					l358:
						{
							position360, tokenIndex360, depth360 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l361
							}
							position++
							goto l360
						l361:
							position, tokenIndex, depth = position360, tokenIndex360, depth360
							if buffer[position] != rune('O') {
								goto l347
							}
							position++
						}
					l360:
						{
							position362, tokenIndex362, depth362 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l363
							}
							position++
							goto l362
						l363:
							position, tokenIndex, depth = position362, tokenIndex362, depth362
							if buffer[position] != rune('A') {
								goto l347
							}
							position++
						}
					l362:
						{
							position364, tokenIndex364, depth364 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l365
							}
							position++
							goto l364
						l365:
							position, tokenIndex, depth = position364, tokenIndex364, depth364
							if buffer[position] != rune('T') {
								goto l347
							}
							position++
						}
					l364:
						break
					case 'B', 'b':
						{
							position366, tokenIndex366, depth366 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l367
							}
							position++
							goto l366
						l367:
							position, tokenIndex, depth = position366, tokenIndex366, depth366
							if buffer[position] != rune('B') {
								goto l347
							}
							position++
						}
					l366:
						{
							position368, tokenIndex368, depth368 := position, tokenIndex, depth
							if buffer[position] != rune('y') {
								goto l369
							}
							position++
							goto l368
						l369:
							position, tokenIndex, depth = position368, tokenIndex368, depth368
							if buffer[position] != rune('Y') {
								goto l347
							}
							position++
						}
					l368:
						{
							position370, tokenIndex370, depth370 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l371
							}
							position++
							goto l370
						l371:
							position, tokenIndex, depth = position370, tokenIndex370, depth370
							if buffer[position] != rune('T') {
								goto l347
							}
							position++
						}
					l370:
						{
							position372, tokenIndex372, depth372 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l373
							}
							position++
							goto l372
						l373:
							position, tokenIndex, depth = position372, tokenIndex372, depth372
							if buffer[position] != rune('E') {
								goto l347
							}
							position++
						}
					l372:
						break
					case 'W', 'w':
						{
							position374, tokenIndex374, depth374 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l375
							}
							position++
							goto l374
						l375:
							position, tokenIndex, depth = position374, tokenIndex374, depth374
							if buffer[position] != rune('W') {
								goto l347
							}
							position++
						}
					l374:
						{
							position376, tokenIndex376, depth376 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l377
							}
							position++
							goto l376
						l377:
							position, tokenIndex, depth = position376, tokenIndex376, depth376
							if buffer[position] != rune('O') {
								goto l347
							}
							position++
						}
					l376:
						{
							position378, tokenIndex378, depth378 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l379
							}
							position++
							goto l378
						l379:
							position, tokenIndex, depth = position378, tokenIndex378, depth378
							if buffer[position] != rune('R') {
								goto l347
							}
							position++
						}
					l378:
						{
							position380, tokenIndex380, depth380 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l381
							}
							position++
							goto l380
						l381:
							position, tokenIndex, depth = position380, tokenIndex380, depth380
							if buffer[position] != rune('D') {
								goto l347
							}
							position++
						}
					l380:
						break
					default:
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
								goto l347
							}
							position++
						}
					l382:
						{
							position384, tokenIndex384, depth384 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l385
							}
							position++
							goto l384
						l385:
							position, tokenIndex, depth = position384, tokenIndex384, depth384
							if buffer[position] != rune('W') {
								goto l347
							}
							position++
						}
					l384:
						{
							position386, tokenIndex386, depth386 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l387
							}
							position++
							goto l386
						l387:
							position, tokenIndex, depth = position386, tokenIndex386, depth386
							if buffer[position] != rune('O') {
								goto l347
							}
							position++
						}
					l386:
						{
							position388, tokenIndex388, depth388 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l389
							}
							position++
							goto l388
						l389:
							position, tokenIndex, depth = position388, tokenIndex388, depth388
							if buffer[position] != rune('R') {
								goto l347
							}
							position++
						}
					l388:
						{
							position390, tokenIndex390, depth390 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l391
							}
							position++
							goto l390
						l391:
							position, tokenIndex, depth = position390, tokenIndex390, depth390
							if buffer[position] != rune('D') {
								goto l347
							}
							position++
						}
					l390:
						break
					}
				}

				if !_rules[ruleSpace]() {
					goto l347
				}
				depth--
				add(ruleDATA_TYPE, position348)
			}
			return true
		l347:
			position, tokenIndex, depth = position347, tokenIndex347, depth347
			return false
		},
		/* 29 LBRK <- <('[' Spacing)> */
		func() bool {
			position392, tokenIndex392, depth392 := position, tokenIndex, depth
			{
				position393 := position
				depth++
				if buffer[position] != rune('[') {
					goto l392
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l392
				}
				depth--
				add(ruleLBRK, position393)
			}
			return true
		l392:
			position, tokenIndex, depth = position392, tokenIndex392, depth392
			return false
		},
		/* 30 RBRK <- <(']' Spacing)> */
		func() bool {
			position394, tokenIndex394, depth394 := position, tokenIndex, depth
			{
				position395 := position
				depth++
				if buffer[position] != rune(']') {
					goto l394
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l394
				}
				depth--
				add(ruleRBRK, position395)
			}
			return true
		l394:
			position, tokenIndex, depth = position394, tokenIndex394, depth394
			return false
		},
		/* 31 COMMA <- <(',' Spacing)> */
		func() bool {
			position396, tokenIndex396, depth396 := position, tokenIndex, depth
			{
				position397 := position
				depth++
				if buffer[position] != rune(',') {
					goto l396
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l396
				}
				depth--
				add(ruleCOMMA, position397)
			}
			return true
		l396:
			position, tokenIndex, depth = position396, tokenIndex396, depth396
			return false
		},
		/* 32 SEMICOLON <- <';'> */
		nil,
		/* 33 COLON <- <(':' Spacing)> */
		nil,
		/* 34 NL <- <'\n'> */
		func() bool {
			position400, tokenIndex400, depth400 := position, tokenIndex, depth
			{
				position401 := position
				depth++
				if buffer[position] != rune('\n') {
					goto l400
				}
				position++
				depth--
				add(ruleNL, position401)
			}
			return true
		l400:
			position, tokenIndex, depth = position400, tokenIndex400, depth400
			return false
		},
		/* 35 EOT <- <!.> */
		nil,
		/* 36 Literal <- <((FloatLiteral / ((&('"') StringLiteral) | (&('\'') CharLiteral) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') IntegerLiteral))) Spacing)> */
		nil,
		/* 37 IntegerLiteral <- <(HexNumeral / BinaryNumeral / OctalNumeral / DecimalNumeral)> */
		func() bool {
			position404, tokenIndex404, depth404 := position, tokenIndex, depth
			{
				position405 := position
				depth++
				{
					position406, tokenIndex406, depth406 := position, tokenIndex, depth
					if !_rules[ruleHexNumeral]() {
						goto l407
					}
					goto l406
				l407:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					{
						position409 := position
						depth++
						{
							position410, tokenIndex410, depth410 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l411
							}
							position++
							if buffer[position] != rune('b') {
								goto l411
							}
							position++
							goto l410
						l411:
							position, tokenIndex, depth = position410, tokenIndex410, depth410
							if buffer[position] != rune('0') {
								goto l408
							}
							position++
							if buffer[position] != rune('B') {
								goto l408
							}
							position++
						}
					l410:
						{
							position412, tokenIndex412, depth412 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l413
							}
							position++
							goto l412
						l413:
							position, tokenIndex, depth = position412, tokenIndex412, depth412
							if buffer[position] != rune('1') {
								goto l408
							}
							position++
						}
					l412:
					l414:
						{
							position415, tokenIndex415, depth415 := position, tokenIndex, depth
						l416:
							{
								position417, tokenIndex417, depth417 := position, tokenIndex, depth
								if buffer[position] != rune('_') {
									goto l417
								}
								position++
								goto l416
							l417:
								position, tokenIndex, depth = position417, tokenIndex417, depth417
							}
							{
								position418, tokenIndex418, depth418 := position, tokenIndex, depth
								if buffer[position] != rune('0') {
									goto l419
								}
								position++
								goto l418
							l419:
								position, tokenIndex, depth = position418, tokenIndex418, depth418
								if buffer[position] != rune('1') {
									goto l415
								}
								position++
							}
						l418:
							goto l414
						l415:
							position, tokenIndex, depth = position415, tokenIndex415, depth415
						}
						depth--
						add(ruleBinaryNumeral, position409)
					}
					goto l406
				l408:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					{
						position421 := position
						depth++
						if buffer[position] != rune('0') {
							goto l420
						}
						position++
					l424:
						{
							position425, tokenIndex425, depth425 := position, tokenIndex, depth
							if buffer[position] != rune('_') {
								goto l425
							}
							position++
							goto l424
						l425:
							position, tokenIndex, depth = position425, tokenIndex425, depth425
						}
						if c := buffer[position]; c < rune('0') || c > rune('7') {
							goto l420
						}
						position++
					l422:
						{
							position423, tokenIndex423, depth423 := position, tokenIndex, depth
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
								goto l423
							}
							position++
							goto l422
						l423:
							position, tokenIndex, depth = position423, tokenIndex423, depth423
						}
						depth--
						add(ruleOctalNumeral, position421)
					}
					goto l406
				l420:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					{
						position428 := position
						depth++
						{
							position429, tokenIndex429, depth429 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l430
							}
							position++
							goto l429
						l430:
							position, tokenIndex, depth = position429, tokenIndex429, depth429
							if c := buffer[position]; c < rune('1') || c > rune('9') {
								goto l404
							}
							position++
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
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l432
								}
								position++
								goto l431
							l432:
								position, tokenIndex, depth = position432, tokenIndex432, depth432
							}
						}
					l429:
						depth--
						add(ruleDecimalNumeral, position428)
					}
				}
			l406:
				depth--
				add(ruleIntegerLiteral, position405)
			}
			return true
		l404:
			position, tokenIndex, depth = position404, tokenIndex404, depth404
			return false
		},
		/* 38 DecimalNumeral <- <('0' / ([1-9] ('_'* [0-9])*))> */
		nil,
		/* 39 HexNumeral <- <((('0' 'x') / ('0' 'X')) HexDigits)> */
		func() bool {
			position436, tokenIndex436, depth436 := position, tokenIndex, depth
			{
				position437 := position
				depth++
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if buffer[position] != rune('0') {
						goto l439
					}
					position++
					if buffer[position] != rune('x') {
						goto l439
					}
					position++
					goto l438
				l439:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if buffer[position] != rune('0') {
						goto l436
					}
					position++
					if buffer[position] != rune('X') {
						goto l436
					}
					position++
				}
			l438:
				if !_rules[ruleHexDigits]() {
					goto l436
				}
				depth--
				add(ruleHexNumeral, position437)
			}
			return true
		l436:
			position, tokenIndex, depth = position436, tokenIndex436, depth436
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
			position444, tokenIndex444, depth444 := position, tokenIndex, depth
			{
				position445 := position
				depth++
				{
					position446, tokenIndex446, depth446 := position, tokenIndex, depth
					if buffer[position] != rune('e') {
						goto l447
					}
					position++
					goto l446
				l447:
					position, tokenIndex, depth = position446, tokenIndex446, depth446
					if buffer[position] != rune('E') {
						goto l444
					}
					position++
				}
			l446:
				{
					position448, tokenIndex448, depth448 := position, tokenIndex, depth
					{
						position450, tokenIndex450, depth450 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l451
						}
						position++
						goto l450
					l451:
						position, tokenIndex, depth = position450, tokenIndex450, depth450
						if buffer[position] != rune('-') {
							goto l448
						}
						position++
					}
				l450:
					goto l449
				l448:
					position, tokenIndex, depth = position448, tokenIndex448, depth448
				}
			l449:
				if !_rules[ruleDigits]() {
					goto l444
				}
				depth--
				add(ruleExponent, position445)
			}
			return true
		l444:
			position, tokenIndex, depth = position444, tokenIndex444, depth444
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
			position455, tokenIndex455, depth455 := position, tokenIndex, depth
			{
				position456 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l455
				}
				position++
			l457:
				{
					position458, tokenIndex458, depth458 := position, tokenIndex, depth
				l459:
					{
						position460, tokenIndex460, depth460 := position, tokenIndex, depth
						if buffer[position] != rune('_') {
							goto l460
						}
						position++
						goto l459
					l460:
						position, tokenIndex, depth = position460, tokenIndex460, depth460
					}
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l458
					}
					position++
					goto l457
				l458:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
				}
				depth--
				add(ruleDigits, position456)
			}
			return true
		l455:
			position, tokenIndex, depth = position455, tokenIndex455, depth455
			return false
		},
		/* 49 HexDigits <- <(HexDigit ('_'* HexDigit)*)> */
		func() bool {
			position461, tokenIndex461, depth461 := position, tokenIndex, depth
			{
				position462 := position
				depth++
				if !_rules[ruleHexDigit]() {
					goto l461
				}
			l463:
				{
					position464, tokenIndex464, depth464 := position, tokenIndex, depth
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
					if !_rules[ruleHexDigit]() {
						goto l464
					}
					goto l463
				l464:
					position, tokenIndex, depth = position464, tokenIndex464, depth464
				}
				depth--
				add(ruleHexDigits, position462)
			}
			return true
		l461:
			position, tokenIndex, depth = position461, tokenIndex461, depth461
			return false
		},
		/* 50 HexDigit <- <((&('A' | 'B' | 'C' | 'D' | 'E' | 'F') [A-F]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]))> */
		func() bool {
			position467, tokenIndex467, depth467 := position, tokenIndex, depth
			{
				position468 := position
				depth++
				{
					switch buffer[position] {
					case 'A', 'B', 'C', 'D', 'E', 'F':
						if c := buffer[position]; c < rune('A') || c > rune('F') {
							goto l467
						}
						position++
						break
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l467
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l467
						}
						position++
						break
					}
				}

				depth--
				add(ruleHexDigit, position468)
			}
			return true
		l467:
			position, tokenIndex, depth = position467, tokenIndex467, depth467
			return false
		},
		/* 51 CharLiteral <- <('\'' (Escape / (!('\'' / '\\') .)) '\'')> */
		nil,
		/* 52 StringLiteral <- <('"' (Escape / (!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('"') '"')) .))* '"')> */
		nil,
		/* 53 Escape <- <('\\' ((&('u') UnicodeEscape) | (&('\\') '\\') | (&('\'') '\'') | (&('"') '"') | (&('r') 'r') | (&('f') 'f') | (&('n') 'n') | (&('t') 't') | (&('b') 'b') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7') OctalEscape)))> */
		func() bool {
			position472, tokenIndex472, depth472 := position, tokenIndex, depth
			{
				position473 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l472
				}
				position++
				{
					switch buffer[position] {
					case 'u':
						{
							position475 := position
							depth++
							if buffer[position] != rune('u') {
								goto l472
							}
							position++
						l476:
							{
								position477, tokenIndex477, depth477 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l477
								}
								position++
								goto l476
							l477:
								position, tokenIndex, depth = position477, tokenIndex477, depth477
							}
							if !_rules[ruleHexDigit]() {
								goto l472
							}
							if !_rules[ruleHexDigit]() {
								goto l472
							}
							if !_rules[ruleHexDigit]() {
								goto l472
							}
							if !_rules[ruleHexDigit]() {
								goto l472
							}
							depth--
							add(ruleUnicodeEscape, position475)
						}
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l472
						}
						position++
						break
					case '\'':
						if buffer[position] != rune('\'') {
							goto l472
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l472
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l472
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l472
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l472
						}
						position++
						break
					case 't':
						if buffer[position] != rune('t') {
							goto l472
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l472
						}
						position++
						break
					default:
						{
							position478 := position
							depth++
							{
								position479, tokenIndex479, depth479 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('3') {
									goto l480
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l480
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l480
								}
								position++
								goto l479
							l480:
								position, tokenIndex, depth = position479, tokenIndex479, depth479
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l481
								}
								position++
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l481
								}
								position++
								goto l479
							l481:
								position, tokenIndex, depth = position479, tokenIndex479, depth479
								if c := buffer[position]; c < rune('0') || c > rune('7') {
									goto l472
								}
								position++
							}
						l479:
							depth--
							add(ruleOctalEscape, position478)
						}
						break
					}
				}

				depth--
				add(ruleEscape, position473)
			}
			return true
		l472:
			position, tokenIndex, depth = position472, tokenIndex472, depth472
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
