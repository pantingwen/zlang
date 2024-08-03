package lexer

import (
	"z/token"
)

var preToken token.Token

type Lexer struct {
	input       string // 输入的字符串
	position    int    // 已经读取的字符的位置
	readPostion int    // 准备读取的字符的位置
	ch          byte   // 已经读取的字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPostion >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPostion]
	}
	l.position = l.readPostion
	l.readPostion += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			for { // single line anntation
				l.readChar()
				ch := l.ch
				if ch == 10 {
					return l.NextToken()
				}
			}
		} else if l.peekChar() == '*' { // multiple lines annotation begin with /* end with */
			l.readChar()
			for {
				l.readChar()
				ch := l.ch
				if ch == '*' && l.peekChar() == '/' {
					l.readChar() // use readChar twice lose */ char
					l.readChar()
					return l.NextToken()
				}
				if ch == 0 { // find */ until the last of file
					l.readChar()
					return l.NextToken()
				}
			}
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			tok.Type = token.LE
			tok.Literal = l.input[l.position : l.position+2]
			l.readChar()
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			tok.Type = token.GE
			tok.Literal = l.input[l.position : l.position+2]
			l.readChar()
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '\n': // replace \n with ;
		if preToken.Literal != ";" && preToken.Literal != "{" && preToken.Literal != "}" {
			tok.Type = token.SEMICOLON
			tok.Literal = ";"
		} else {
			l.readChar()
			return l.NextToken()
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIndentifier()
			tok.Type = token.LookIndent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILIEGAL, l.ch)
		}
	}

	l.readChar()
	preToken = tok
	return tok
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIndentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func newToken(TokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: TokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhiteSpace() {
	for l.isWhiteSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.readPostion > len(l.input) {
		return 0
	} else {
		return l.input[l.readPostion]
	}
}
