package lexer

import (
	"testing"

	"ziggytwister.com/shadow-hunter/token"
)

// The first goal is to tokenize this string
// {:tag :ret, :val "2", :form "(+ 1 1)", :ns "cljs.user", :ms 15}

func TestNextToken(t *testing.T) {
	input := `{:tag :ret, :val "2", :form "(+ 1 1)", :ns "cljs.user", :ms 15}`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LCURLY, "{"},
		{token.KEYWORD, ":tag"},
		{token.KEYWORD, ":ret"},
		{token.KEYWORD, ":val"},
		{token.STRING, "2"},
		{token.KEYWORD, ":form"},
		{token.STRING, "(+ 1 1)"},
		{token.KEYWORD, ":ns"},
		{token.STRING, "cljs.user"},
		{token.KEYWORD, ":ms"},
		{token.INT, "15"},
		{token.RCURLY, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextTokenWithNestedMaps(t *testing.T) {
	input := `{1 {:name "Jack"
                :done false}
	           2 {:name "Blob"}}`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LCURLY, "{"},
		{token.INT, "1"},
		{token.LCURLY, "{"},
		{token.KEYWORD, ":name"},
		{token.STRING, "Jack"},
		{token.KEYWORD, ":done"},
		{token.FALSE, "false"},
		{token.RCURLY, "}"},
		{token.INT, "2"},
		{token.LCURLY, "{"},
		{token.KEYWORD, ":name"},
		{token.STRING, "Blob"},
		{token.RCURLY, "}"},
		{token.RCURLY, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
