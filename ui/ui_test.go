package ui

import (
	"testing"

	"ziggytwister.com/shadow-hunter/ast"
	"ziggytwister.com/shadow-hunter/token"
)

func TestGetCurrentDataSlice(t *testing.T) {
	edn := &ast.EDN{
		Elements: []ast.Element{
			&ast.MapElement{
				Token: token.Token{Type: token.LCURLY, Literal: "{"},
				Keys: []ast.Element{
					&ast.KeywordElement{
						Token: token.Token{Type: token.KEYWORD, Literal: ":name"},
						Value: ":name",
					},
					&ast.KeywordElement{
						Token: token.Token{Type: token.KEYWORD, Literal: ":extra"},
						Value: ":extra",
					},
				},
				Values: []ast.Element{
					&ast.StringElement{
						Token: token.Token{Type: token.STRING, Literal: "Jack"},
						Value: "Jack",
					},
					&ast.MapElement{
						Token: token.Token{Type: token.LCURLY, Literal: "{"},
						Keys: []ast.Element{
							&ast.KeywordElement{
								Token: token.Token{Type: token.KEYWORD, Literal: ":hobby"},
								Value: ":hobby",
							},
							&ast.KeywordElement{
								Token: token.Token{Type: token.KEYWORD, Literal: ":age"},
								Value: ":age",
							},
						},
						Values: []ast.Element{
							&ast.StringElement{
								Token: token.Token{Type: token.KEYWORD, Literal: "Painting"},
								Value: "Painting",
							},
							&ast.IntegerLiteral{
								Token: token.Token{Type: token.KEYWORD, Literal: "33"},
								Value: 33,
							},
						},
					},
				},
			},
		},
	}

	m := mainModel{
		edn:            edn,
		currentPathIdx: 0,
		path:           []int{0, 1, 1},
	}
	data := m.getCurrentDataSlice()

	mapElement, ok := data.(*ast.MapElement)

	if !ok {
		t.Fatalf("Expected first element to be a map go %T", data)
	}

	if element := mapElement.Keys[0].TokenLiteral(); element != ":name" {
		t.Fatalf("Expected first key of first map to be :name got %v", element)
	}

	if element := mapElement.Values[0].TokenLiteral(); element != "Jack" {
		t.Fatalf("Expected first value of first map to be Jack got %v", element)
	}

	if element := mapElement.Keys[1].TokenLiteral(); element != ":extra" {
		t.Fatalf("Expected second key of first map to be :extra got %v", element)
	}

	m.currentPathIdx = 1
	data = m.getCurrentDataSlice()

	mapElement, ok = data.(*ast.MapElement)

	if !ok {
		t.Fatalf("Expected second element to be a map go %T", data)
	}

	if element := mapElement.Keys[0].TokenLiteral(); element != ":hobby" {
		t.Fatalf("Expected first key of second map to be :hobby got %v", element)
	}

	if element := mapElement.Values[0].TokenLiteral(); element != "Painting" {
		t.Fatalf("Expected first value of second map to be Jack got %v", element)
	}

	if element := mapElement.Keys[1].TokenLiteral(); element != ":age" {
		t.Fatalf("Expected second key of second map to be :extra got %v", element)
	}

	if element := mapElement.Values[1].TokenLiteral(); element != "33" {
		t.Fatalf("Expected second value of second map to be Jack got %v", element)
	}

	m.currentPathIdx = 2
	data = m.getCurrentDataSlice()

	integerElement, ok := data.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Expected third element to be a number go %T", data)
	}

	if integerElement.Value != 33 {
		t.Fatalf("Expected integer value to be 33 got %d", integerElement.Value)
	}
}
