package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

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
	BOOL    = "BOOL"
)
