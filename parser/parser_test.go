package parser

import (
	"testing"

	"ziggytwister.com/shadow-hunter/lexer"
)

func TestKeywordElement(t *testing.T) {
	input := ":name"

	l := lexer.New(input)
	p := New(l)

	edn := p.ParseEDN()

	if edn == nil {
		t.Fatal("ParseEDN() returned nil")
	}

	if len(edn.Elements) != 1 {
		t.Fatalf("EDN expected to contain 1 element go %d",
			len(edn.Elements))
	}

	if edn.Elements[0].TokenLiteral() != ":name" {
		t.Fatalf("Expected first element to be :name got %v",
			edn.Elements[0].TokenLiteral())
	}
}

func TestVectorElement(t *testing.T) {
	input := "[:name :age]"

	l := lexer.New(input)
	p := New(l)

	edn := p.ParseEDN()

	if edn == nil {
		t.Fatal("ParseEDN() returned nil")
	}

	if len(edn.Elements) != 2 {
		t.Fatalf("EDN expected to contain 2 element go %d",
			len(edn.Elements))
	}

	if edn.Elements[0].TokenLiteral() != ":name" {
		t.Fatalf("Expected first element to be :name got %v",
			edn.Elements[0].TokenLiteral())
	}

	if edn.Elements[1].TokenLiteral() != ":age" {
		t.Fatalf("Expected first element to be :age got %v",
			edn.Elements[0].TokenLiteral())
	}
}

//func TestMapElement(t *testing.T) {
//	input := `
//{:name "Jack"}
//`
//
//	l := lexer.New(input)
//	p := New(l)
//
//	edn := p.ParseEDN()
//	if edn == nil {
//		t.Fatal("ParseEDN() returned nil")
//	}
//	if len(edn.Elements) != 1 {
//		t.Fatalf("EDN expected to contain 1 element go %d",
//			len(edn.Elements))
//	}
//
//	tests := []struct {
//		expectedKey   string
//		expectedValue string
//	}{
//		{":name", "Jack"},
//	}
//
//	for i, tt := range tests {
//		element := edn.Elements[i]
//
//		if !testMapElement(t, element, tt.expectedKey, tt.expectedValue) {
//			return
//		}
//	}
//}
//
//func testMapElement(t *testing.T, e ast.Element, key string, value string) bool {
//	if e.TokenLiteral() != "{" {
//		t.Errorf("I don't know what I'm doing yet")
//		return false
//	}
//
//	return true
//}
