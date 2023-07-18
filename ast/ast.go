package ast

import "ziggytwister.com/shadow-hunter/token"

type Node interface {
	TokenLiteral() string
}

type Element interface {
	Node
	elementNode()
}

type EDN struct {
	Elements []Element
}

func (e *EDN) TokenLiteral() string {
	if len(e.Elements) > 0 {
		return e.Elements[0].TokenLiteral()
	} else {
		return ""
	}
}

type MapElement struct {
	Token  token.Token
	Keys   []Element
	Values []Element
}

func (me *MapElement) elementNode()         {}
func (me *MapElement) TokenLiteral() string { return me.Token.Literal }

type KeywordElement struct {
	Token token.Token
	Value string
}

func (k *KeywordElement) elementNode()         {}
func (k *KeywordElement) TokenLiteral() string { return k.Token.Literal }
