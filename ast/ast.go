package ast

import (
	"bytes"

	"ziggytwister.com/shadow-hunter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (e *EDN) String() string {
	var out bytes.Buffer

	for _, s := range e.Elements {
		out.WriteString(s.String())
	}

	return out.String()
}

//
// Vectors
//

type VectorElement struct {
	Token    token.Token
	Elements []Element
}

func (ve *VectorElement) elementNode()         {}
func (ve *VectorElement) TokenLiteral() string { return "vector" }
func (ve *VectorElement) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, e := range ve.Elements {
		out.WriteString(e.String())
		if i < len(ve.Elements)-1 {
			out.WriteString(" ")
		}
	}
	out.WriteString("]")
	return out.String()
}

//
// Maps
//

type MapElement struct {
	Token  token.Token
	Keys   []Element
	Values []Element
}

func (me *MapElement) elementNode()         {}
func (me *MapElement) TokenLiteral() string { return me.Token.Literal }
func (me *MapElement) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for i, key := range me.Keys {
		out.WriteString(key.String())
		out.WriteString(" ")
		out.WriteString(me.Values[i].String())
		if i < len(me.Keys)-1 {
			out.WriteString("\n")
		}
	}
	out.WriteString("}")
	return out.String()
}

//
// Keywords
//

type KeywordElement struct {
	Token token.Token
	Value string
}

func (k *KeywordElement) elementNode()         {}
func (k *KeywordElement) TokenLiteral() string { return k.Token.Literal }
func (k *KeywordElement) String() string {
	var out bytes.Buffer
	out.WriteString(k.Value)
	return out.String()
}

//
// Strings
//

type StringElement struct {
	Token token.Token
	Value string
}

func (k *StringElement) elementNode()         {}
func (k *StringElement) TokenLiteral() string { return k.Token.Literal }
func (k *StringElement) String() string {
	var out bytes.Buffer
	out.WriteString(k.Value)
	return out.String()
}

//
// Ints
//

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) elementNode()         {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

//
// Bools
//

type BoolLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BoolLiteral) elementNode()         {}
func (bl *BoolLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BoolLiteral) String() string       { return bl.Token.Literal }
