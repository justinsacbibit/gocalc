package gocalc

type lexer interface {
	token() *token
	peekToken() *token
}

type state int

const (
	stErr state = iota
	stStart
	stWhitespace
	stId
	stInt
	stZero
	stZeroB
	stBinInt
	stZeroX
	stHexInt
	stOctInt
	stFloat
	stLparen
	stRparen
	stComma
	stLogicalNot
	stNotEqual
	stBitwiseNot
	stStar
	stPercent
	stSlash
	stPlus
	stMinus
	stLeftShift
	stRightShift
	stLessThan
	stLessOrEqual
	stGreaterThan
	stGreaterOrEqual
	stEqual
	stBitwiseAnd
	stLogicalAnd
	stBitwiseXor
	stBitwiseOr
	stLogicalOr
	stT
	stTr
	stTru
	stTrue
	stF
	stFa
	stFal
	stFals
	stFalse
)

const whitespace = "\t\n\r "
const letters = "abcdefABCDEFghijklmnopqrstuvwxyzGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"
const hexDigits = "0123456789abcdefABCDEF"
const binaryDigits = "01"
const octalDigits = "01234567"
const oneToNine = "123456789"

const stateCount uint8 = 44
const maxTransitionCount uint16 = 256

var trans = [stateCount][maxTransitionCount]state{}

var stateTokens = [...]tokenType{
	tokenError,          // stErr
	tokenError,          // stStart
	tokenWhitespace,     // stWhitespace
	tokenIdentifier,     // stId
	tokenInt,            // stInt
	tokenInt,            // stZero
	tokenError,          // stZeroB
	tokenInt,            // stBinInt
	tokenError,          // stZeroX
	tokenInt,            // stHexInt
	tokenInt,            // stOctInt
	tokenFloat,          // stFloat
	tokenLeftParen,      // stLparen
	tokenRightParen,     // stRparen
	tokenComma,          // stComma
	tokenLogicalNot,     // stLogicalNot
	tokenNotEqual,       // stNotEqual
	tokenBitwiseNot,     // stBitwiseNot
	tokenStar,           // stStar
	tokenPercent,        // stPercent
	tokenSlash,          // stSlash
	tokenPlus,           // stPlus
	tokenMinus,          // stMinus
	tokenLeftShift,      // stLeftShift
	tokenRightShift,     // stRightShift
	tokenLessThan,       // stLessThan
	tokenLessOrEqual,    // stLessOrEqual
	tokenGreaterThan,    // stGreaterThan
	tokenGreaterOrEqual, // stGreaterOrEqual
	tokenEqual,          // stEqual
	tokenBitwiseAnd,     // stBitwiseAnd
	tokenLogicalAnd,     // stLogicalAnd
	tokenBitwiseXor,     // stBitwiseXor
	tokenBitwiseOr,      // stBitwiseOr
	tokenLogicalOr,      // stLogicalOr
	tokenIdentifier,     // stT
	tokenIdentifier,     // stTr
	tokenIdentifier,     // stTru
	tokenTrue,           // stTrue
	tokenIdentifier,     // stF
	tokenIdentifier,     // stFa
	tokenIdentifier,     // stFal
	tokenIdentifier,     // stFals
	tokenFalse,          // stFalse
}

func setTrans(current state, transition string, next state) {
	for _, r := range transition {
		trans[current][r] = next
	}
}

func init() {
	setTrans(stStart, whitespace, stWhitespace)
	setTrans(stWhitespace, whitespace, stWhitespace)

	// Single character tokens
	setTrans(stStart, "(", stLparen)
	setTrans(stStart, ")", stRparen)
	setTrans(stStart, ",", stComma)
	setTrans(stStart, "~", stBitwiseNot)
	setTrans(stStart, "*", stStar)
	setTrans(stStart, "%", stPercent)
	setTrans(stStart, "/", stSlash)
	setTrans(stStart, "+", stPlus)
	setTrans(stStart, "-", stMinus)
	setTrans(stStart, "=", stEqual)

	setTrans(stStart, letters, stId)
	setTrans(stId, letters+digits, stId)

	// Int
	setTrans(stStart, oneToNine, stInt)
	setTrans(stInt, digits, stInt)

	// Zero
	setTrans(stStart, "0", stZero)
	setTrans(stZero, ".", stFloat)
	// Binary
	setTrans(stZero, "b", stZeroB)
	setTrans(stZeroB, binaryDigits, stBinInt)
	setTrans(stBinInt, binaryDigits, stBinInt)
	// Hex
	setTrans(stZero, "xX", stZeroX)
	setTrans(stZeroX, hexDigits, stHexInt)
	setTrans(stHexInt, hexDigits, stHexInt)
	// Octal
	setTrans(stZero, octalDigits, stOctInt)
	setTrans(stOctInt, octalDigits, stOctInt)

	// Float
	setTrans(stInt, ".", stFloat)
	setTrans(stFloat, digits, stFloat)

	// Comparisons
	setTrans(stStart, "!", stLogicalNot)
	setTrans(stLogicalNot, "=", stNotEqual)
	setTrans(stStart, "<", stLessThan)
	setTrans(stLessThan, "=", stLessOrEqual)
	setTrans(stStart, ">", stGreaterThan)
	setTrans(stGreaterThan, "=", stGreaterOrEqual)
	setTrans(stStart, "&", stBitwiseAnd)
	setTrans(stBitwiseAnd, "&", stLogicalAnd)
	setTrans(stStart, "^", stBitwiseXor)
	setTrans(stStart, "|", stBitwiseOr)
	setTrans(stBitwiseOr, "|", stLogicalOr)

	// Shift
	setTrans(stLessThan, "<", stLeftShift)
	setTrans(stGreaterThan, ">", stRightShift)

	// Boolean literals
	setTrans(stStart, "t", stT)
	setTrans(stT, "r", stTr)
	setTrans(stTr, "u", stTru)
	setTrans(stTru, "e", stTrue)
	setTrans(stTrue, letters+digits, stId)

	setTrans(stStart, "f", stF)
	setTrans(stF, "a", stFa)
	setTrans(stFa, "l", stFal)
	setTrans(stFal, "s", stFals)
	setTrans(stFals, "e", stFalse)
	setTrans(stFalse, letters+digits, stId)
}

func newLexer(input string) lexer {
	return &gocalcLexer{
		input:  input,
		tokens: queue{},
		state:  stStart,
	}
}

func (l *gocalcLexer) push() {
	if len(l.input) == 0 {
		l.emit(tokenEOF)
		return
	}

	for len(l.tokens) < 1 {
		nextState := stErr

		if l.pos < len(l.input) {
			nextState = trans[l.state][l.input[l.pos]]
		}

		if nextState == stErr {
			curTokenType := stateTokens[l.state]
			if curTokenType == tokenError {
				l.emit(tokenError)
				return
			}

			if curTokenType != tokenWhitespace {
				l.emit(curTokenType)
			}
			l.start = l.pos
			l.state = stStart

			if l.start >= len(l.input) {
				l.emit(tokenEOF)
			}
		} else {
			l.state = nextState
			l.pos++
		}
	}
}

func (l *gocalcLexer) token() *token {
	if len(l.tokens) < 1 {
		l.push()
	}
	return l.tokens.pop()
}

func (l *gocalcLexer) peekToken() *token {
	if len(l.tokens) < 1 {
		l.push()
	}
	return l.tokens.first()
}

const eof = -1

type gocalcLexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens queue
	state  state
}

func (l *gocalcLexer) emit(t tokenType) {
	l.tokens.push(&token{
		typ: t,
		val: l.input[l.start:l.pos],
	})
	l.start = l.pos
}
