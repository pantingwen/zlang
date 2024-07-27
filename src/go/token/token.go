package token

const (
	ILIEGAL   = "ILIEGAL"
	EOF       = "EOF"
	IDENT     = "IDENT"
	INT       = "INT"
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// keyword
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	STRING   = "STRING"
	IMPORT   = "IMPORT"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	LE = "<="
	GT = ">"
	GE = ">="

	EQ     = "=="
	NOT_EQ = "!="

	LBRACKET = "["
	RBRACKET = "]"

	COLON = ":"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"import": IMPORT,
}

func LookIndent(indent string) TokenType {
	token, ok := keywords[indent]
	if ok {
		return token
	}
	return IDENT
}
