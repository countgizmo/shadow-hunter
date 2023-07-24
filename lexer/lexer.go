package lexer

import (
	"unicode"

	"ziggytwister.com/shadow-hunter/token"
)

type Lexer struct {
	input           []rune
	currentPosition int
	nextPosition    int
	ch              rune
}

func New(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextPosition]
	}

	l.currentPosition = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LCURLY, l.ch)
	case '}':
		tok = newToken(token.RCURLY, l.ch)
	case '[':
		tok = newToken(token.LSQBRACKET, l.ch)
	case ']':
		tok = newToken(token.RSQBRACKET, l.ch)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isKeyword(l.ch) {
			tok.Literal = l.readKeyword()
			tok.Type = token.KEYWORD
			return tok
		} else if unicode.IsDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else if isString(l.ch) {
			tok.Literal = l.readString()
			tok.Type = token.STRING
			return tok
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isKeyword(ch rune) bool {
	return ch == ':'
}

func isString(ch rune) bool {
	return ch == '"'
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func (l *Lexer) readIdentifier() string {
	startPosition := l.currentPosition

	for unicode.IsLetter(l.ch) {
		l.readChar()
	}

	return string(l.input[startPosition:l.currentPosition])
}

func (l *Lexer) readKeyword() string {
	startPosition := l.currentPosition

	l.readChar() //reading the ':'
	for unicode.IsLetter(l.ch) {
		l.readChar()
	}

	return string(l.input[startPosition:l.currentPosition])
}

func (l *Lexer) readNumber() string {
	startPosition := l.currentPosition

	for unicode.IsDigit(l.ch) {
		l.readChar()
	}

	return string(l.input[startPosition:l.currentPosition])
}

func (l *Lexer) readString() string {
	l.readChar() //reading left side '"'
	startPosition := l.currentPosition

	for l.ch != '"' {
		l.readChar()
	}

	result := string(l.input[startPosition:l.currentPosition])
	l.readChar() //reading right side '"'
	return result
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' || l.ch == ',' {
		l.readChar()
	}
}
