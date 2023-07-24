package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	IDENT   = "IDENT"

	LPAREN     = "("
	RPAREN     = ")"
	LCURLY     = "{"
	RCURLY     = "}"
	LSQBRACKET = "["
	RSQBRACKET = "]"

	// Identifiers and friends
	KEYWORD = "KEYWORD"
	INT     = "INT"
	STRING  = "STRING"

	// Reserved
	TRUE  = "TRUE"
	FALSE = "FALSE"
)

var reserved = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := reserved[ident]; ok {
		return tok
	}

	return IDENT
}
